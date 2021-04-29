package storer

import (
	"github.com/philips-labs/ferrite/types"
	"io"
)

var (
	TaskNone = io.EOF
)

type Task interface {
	Create(schedule types.Task) (*types.Task, error)
	Delete(id string) error
	FindByID(id string) (*types.Task, error)
	FindByProjectID(id string) (*[]types.Task, error)
	Next() (*types.Task, error)
	SetStatus(id, status string) error
}
