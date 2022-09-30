package storer

import (
	"github.com/loafoe/ferrite/types"
)

type Project interface {
	Create(project types.Project) (*types.Project, error)
	FindByID(id string) (*types.Project, error)
	FindLatest() (*types.Project, error)
}
