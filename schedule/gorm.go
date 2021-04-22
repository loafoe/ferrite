package schedule

import (
	"gorm.io/gorm"
)

type GormStorer struct {
	DB *gorm.DB
}

func (g *GormStorer) Create(schedule Schedule) (*Schedule, error) {
	tx := g.DB.Create(schedule)
	if tx.Error != nil {
		return nil, tx.Error
	}
	createdSchedule, err := g.FindByID(schedule.ID)
	return createdSchedule, err
}

func (g *GormStorer) Delete(id string) error {
	tx := g.DB.Delete(&Schedule{}, "id = ?", id)
	return tx.Error
}

func (g *GormStorer) FindByID(id string) (*Schedule, error) {
	var schedule Schedule
	tx := g.DB.First(&schedule, "id = ?", id)
	return &schedule, tx.Error
}

func (g *GormStorer) FindByProjectID(id string) (*[]Schedule, error) {
	var schedules []Schedule
	tx := g.DB.Find(&schedules, "project_id = ?", id)
	return &schedules, tx.Error
}

func (g *GormStorer) FindByCodeName(codeName string) (*[]Schedule, error) {
	var schedules []Schedule
	tx := g.DB.Find(&schedules, "code_name = ?", codeName)
	return &schedules, tx.Error
}
