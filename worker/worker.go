package worker

import (
	"ferrite/task"
	"fmt"
	"time"
)

// Start starts a worker run
func Start(storer task.Storer) (chan bool, error) {
	done := make(chan bool)
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		fmt.Printf("starting worker..\n")
		for {
			select {
			case <-ticker.C:
				err := fetchAndRunNextAvailableTask(storer)
				if err != nil {
					fmt.Printf("ticker: %v\n", err)
				}
			case <-done:
				fmt.Printf("Received done signal. Exiting...\n")
				return
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
	fmt.Printf("Found new task: %s\n", t.ID)
	_ = storer.SetStatus(t.ID, "running")
	fmt.Printf("TODO: should run here...\n")
	return storer.SetStatus(t.ID, "done")
}
