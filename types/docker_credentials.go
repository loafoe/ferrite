package types

import (
	"github.com/dip-software/go-dip-api/iron"
)

type DockerCredentials struct {
	ProjectID string `gorm:"foreignKey:ProjectRefer"`
	iron.DockerCredentials
}
