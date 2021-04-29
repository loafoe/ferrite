package storer

import (
	"github.com/philips-labs/ferrite/types"
)

type Code interface {
	Create(code types.Code) (*types.Code, error)
	Delete(id string) error
	Update(code types.Code) error
	FindByID(id string) (*types.Code, error)
	FindByName(id string) (*types.Code, error)
	FindByProjectID(id string) (*[]types.Code, error)
	SaveCredentials(creds types.DockerCredentials) error
}
