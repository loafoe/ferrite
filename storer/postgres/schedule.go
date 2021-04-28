package postgres

import (
	"ferrite/types"

	"gorm.io/gorm"
)

type ScheduleStorer struct {
	DB *gorm.DB
}

func (g *ScheduleStorer) Create(schedule types.Schedule) (*types.Schedule, error) {
	tx := g.DB.Create(schedule)
	if tx.Error != nil {
		return nil, tx.Error
	}
	createdSchedule, err := g.FindByID(schedule.ID)
	return createdSchedule, err
}

func (g *ScheduleStorer) Delete(id string) error {
	tx := g.DB.Delete(&types.Schedule{}, "id = ?", id)
	return tx.Error
}

func (g *ScheduleStorer) FindByID(id string) (*types.Schedule, error) {
	var schedule types.Schedule
	tx := g.DB.First(&schedule, "id = ?", id)
	return &schedule, tx.Error
}

func (g *ScheduleStorer) FindByProjectID(id string) (*[]types.Schedule, error) {
	var schedules []types.Schedule
	tx := g.DB.Find(&schedules, "project_id = ?", id)
	return &schedules, tx.Error
}

func (g *ScheduleStorer) FindByCodeName(codeName string) (*[]types.Schedule, error) {
	var schedules []types.Schedule
	tx := g.DB.Find(&schedules, "code_name = ?", codeName)
	return &schedules, tx.Error
}
