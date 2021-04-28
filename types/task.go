package types

import (
	"github.com/philips-software/go-hsdp-api/iron"
)

// Task represents a task that needs to be scheduled and run on a cluster
// The task lifecycle states are
// `new` - New in the system
// `pending` - Pulled by a runner but not running yet
// `running` - Running on a worker
// `done` - Done, success
// `error` - Error during run
type Task struct {
	iron.Task
}
