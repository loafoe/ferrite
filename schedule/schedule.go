package schedule

import (
	"siderite-server/cluster"

	"github.com/philips-software/go-hsdp-api/iron"
)

type Schedule struct {
	iron.Schedule
	OnCluster cluster.Cluster `gorm:"foreignKey:Cluster"`
}
