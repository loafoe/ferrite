package types

import (
	"github.com/philips-software/go-hsdp-api/iron"
)

type Code struct {
	iron.Code
	Name      string     `gorm:"uniqueIndex" json:"name"`
	Schedules []Schedule `gorm:"foreignKey:CodeName;references:Name;constraint:OnDelete:CASCADE" json:"-"`
	Tasks     []Task     `gorm:"foreignKey:CodeName;references:Name;constraint:OnDelete:CASCADE" json:"-"`
}
