package types

import (
	"github.com/philips-software/go-hsdp-api/iron"
)

type DockerCredentials struct {
	ProjectID string `gorm:"foreignKey:ProjectRefer"`
	iron.DockerCredentials
}
