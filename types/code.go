package types

import (
	"github.com/dip-software/go-dip-api/iron"
)

type Code struct {
	iron.Code
	Name      string     `gorm:"uniqueIndex" json:"name"`
	Schedules []Schedule `gorm:"foreignKey:CodeName;references:Name;constraint:OnDelete:CASCADE" json:"-"`
	Tasks     []Task     `gorm:"foreignKey:CodeName;references:Name;constraint:OnDelete:CASCADE" json:"-"`
}
