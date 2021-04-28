package worker

import (
	"context"
	"ferrite/storer"
	"ferrite/types"
	"fmt"
	"io"
	"os"
	"time"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// Start starts a worker run
func Start(fs *storer.Ferrite) (chan bool, error) {
	done := make(chan bool)
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		fmt.Printf("starting worker..\n")
		for {
			err := fetchAndRunNextAvailableTask(fs)
			if err == storer.TaskNone {
				select {
				case <-ticker.C:
					continue
				case <-done:
					fmt.Printf("Received done signal. Exiting...\n")
					return
				}
			}
		}
	}()
	return done, nil
}

func fetchAndRunNextAvailableTask(fs *storer.Ferrite) error {
	t, err := fs.Task.Next()
	if err != nil {
		return err
	}
	fmt.Printf("new task: %s\n", t.ID)
	_ = fs.Task.SetStatus(t.ID, "running")
	if err := runTask(*t, fs); err != nil {
		fmt.Printf("error running task: %v\n", err)
		_ = fs.Task.SetStatus(t.ID, "error")
		return err
	}
	return fs.Task.SetStatus(t.ID, "done")
}

func runTask(t types.Task, fs *storer.Ferrite) error {
	ctx := context.Background()
	taskCode, err := fs.Code.FindByName(t.CodeName)
	if err != nil {
		return err
	}
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	out, err := cli.ImagePull(ctx, taskCode.Image, dockertypes.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()
	_, _ = io.Copy(os.Stdout, out)

	// Create volume
	vol, err := cli.VolumeCreate(ctx, volume.VolumeCreateBody{
		Driver: "local",
		Name:   t.ID,
	})
	if err != nil {
		return err
	}
	defer func() {
		_ = cli.VolumeRemove(ctx, t.ID, true)
	}()
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: taskCode.Image,
		Tty:   false,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeVolume,
				Source: vol.Name,
				Target: "/work",
			},
		},
	}, nil, nil, t.ID)
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, dockertypes.ContainerStartOptions{}); err != nil {
		return err
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}
	logs, err := cli.ContainerLogs(ctx, resp.ID, dockertypes.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return err
	}
	defer logs.Close()

	_, _ = stdcopy.StdCopy(os.Stdout, os.Stderr, logs)

	return cli.ContainerRemove(ctx, resp.ID, dockertypes.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
}
