package storer

import (
	"github.com/philips-labs/ferrite/types"
)

type Project interface {
	Create(project types.Project) (*types.Project, error)
	FindByID(id string) (*types.Project, error)
	FindLatest() (*types.Project, error)
}
