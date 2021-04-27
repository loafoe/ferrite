package worker

import (
	"context"
	"ferrite/code"
	"ferrite/task"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// Start starts a worker run
func Start(storer task.Storer, codes code.Storer) (chan bool, error) {
	done := make(chan bool)
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		fmt.Printf("starting worker..\n")
		for {
			err := fetchAndRunNextAvailableTask(storer, codes)
			if err == task.None {
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

func fetchAndRunNextAvailableTask(storer task.Storer, codes code.Storer) error {
	t, err := storer.Next()
	if err != nil {
		return err
	}
	fmt.Printf("new task: %s\n", t.ID)
	_ = storer.SetStatus(t.ID, "running")
	if err := runTask(*t, codes); err != nil {
		fmt.Printf("error running task: %v\n", err)
		_ = storer.SetStatus(t.ID, "error")
		return err
	}
	return storer.SetStatus(t.ID, "done")
}

func runTask(t task.Task, codes code.Storer) error {
	ctx := context.Background()
	taskCode, err := codes.FindByName(t.CodeName)
	if err != nil {
		return err
	}
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	out, err := cli.ImagePull(ctx, taskCode.Image, types.ImagePullOptions{})
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
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: taskCode.Image,
		Tty:   false,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeVolume,
				Source: vol.Mountpoint,
				Target: "/work",
			},
		},
	}, nil, nil, t.ID)
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
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
	logs, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return err
	}
	defer logs.Close()

	_, _ = stdcopy.StdCopy(os.Stdout, os.Stderr, logs)

	return cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
}
