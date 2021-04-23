package code

import (
	"ferrite/schedule"

	"github.com/philips-software/go-hsdp-api/iron"
)

type Code struct {
	iron.Code
	Schedules []schedule.Schedule `gorm:"foreignKey:CodeName;references:Name;constraint:OnDelete:CASCADE"`
}
