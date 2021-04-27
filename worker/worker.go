package worker

import (
	"context"
	"ferrite/task"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// Start starts a worker run
func Start(storer task.Storer) (chan bool, error) {
	done := make(chan bool)
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		fmt.Printf("starting worker..\n")
		for {
			err := fetchAndRunNextAvailableTask(storer)
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

func fetchAndRunNextAvailableTask(storer task.Storer) error {
	t, err := storer.Next()
	if err != nil {
		return err
	}
	fmt.Printf("new task: %s\n", t.ID)
	_ = storer.SetStatus(t.ID, "running")
	if err := runTask(*t); err != nil {
		fmt.Printf("error running task: %v\n", err)
		_ = storer.SetStatus(t.ID, "error")
		return err
	}
	return storer.SetStatus(t.ID, "done")
}

func runTask(t task.Task) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	out, err := cli.ImagePull(ctx, t.CodeName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()
	_, _ = io.Copy(os.Stdout, out)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: t.CodeName,
		Tty:   false,
	}, nil, nil, nil, t.ID)
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

	return nil
}
