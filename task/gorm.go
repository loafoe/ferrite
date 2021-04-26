package task

import (
	"gorm.io/gorm"
)

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

func (g *GormStorer) SetStatus(id, status string) error {
	tx := g.DB.Model(&Task{}).Where("id = ?", id).Update("status", status)
	return tx.Error
}

func (g *GormStorer) Next() (*Task, error) {
	var task Task
	tx := g.DB.Raw(`
WITH task AS (
  	SELECT
		id
	FROM
		tasks
	WHERE
		status = 'new'
	LIMIT 1
	FOR UPDATE SKIP LOCKED
) 
UPDATE
	tasks
SET
	status = 'pending'
FROM 
	task
WHERE
	tasks.id = task.id
RETURNING
	tasks.*;
`).Scan(&task)
	if task.ID == "" {
		return nil, None
	}
	return &task, tx.Error
}
