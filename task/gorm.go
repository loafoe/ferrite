package task

import "gorm.io/gorm"

type GormStorer struct {
	DB *gorm.DB
}

func (g *GormStorer) Create(schedule Task) (*Task, error) {
	tx := g.DB.Create(schedule)
	if tx.Error != nil {
		return nil, tx.Error
	}
	createdSchedule, err := g.FindByID(schedule.ID)
	return createdSchedule, err
}

func (g *GormStorer) Delete(id string) error {
	tx := g.DB.Delete(&Task{}, "id = ?", id)
	return tx.Error
}

func (g *GormStorer) FindByID(id string) (*Task, error) {
	var task Task
	tx := g.DB.First(&task, "id = ?", id)
	return &task, tx.Error
}

func (g *GormStorer) FindByProjectID(id string) (*[]Task, error) {
	var tasks []Task
	tx := g.DB.Find(&tasks, "project_id = ?", id)
	return &tasks, tx.Error
}

func (g *GormStorer) FindByCodeName(codeName string) (*[]Task, error) {
	var schedules []Task
	tx := g.DB.Find(&schedules, "code_name = ?", codeName)
	return &schedules, tx.Error
}
