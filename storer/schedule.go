package storer

import "github.com/loafoe/ferrite/types"

type Schedule interface {
	Create(schedule types.Schedule) (*types.Schedule, error)
	Delete(id string) error
	FindByID(id string) (*types.Schedule, error)
	FindByProjectID(id string) (*[]types.Schedule, error)
	FindByCodeName(codeName string) (*[]types.Schedule, error)
}
