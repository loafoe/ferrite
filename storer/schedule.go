package storer

import "github.com/philips-labs/ferrite/types"

type Schedule interface {
	Create(schedule types.Schedule) (*types.Schedule, error)
	Delete(id string) error
	FindByID(id string) (*types.Schedule, error)
	FindByProjectID(id string) (*[]types.Schedule, error)
	FindByCodeName(codeName string) (*[]types.Schedule, error)
}
