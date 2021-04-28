package storer

import (
	"ferrite/types"
)

type Project interface {
	Create(project types.Project) (*types.Project, error)
	FindByID(id string) (*types.Project, error)
}
