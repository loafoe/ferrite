package postgres

import (
	"github.com/philips-labs/ferrite/storer"
	"github.com/philips-labs/ferrite/types"

	"gorm.io/gorm"
)

type TaskStorer struct {
	DB *gorm.DB
}

func (g *TaskStorer) Create(schedule types.Task) (*types.Task, error) {
	tx := g.DB.Create(schedule)
	if tx.Error != nil {
		return nil, tx.Error
	}
	createdSchedule, err := g.FindByID(schedule.ID)
	return createdSchedule, err
}

func (g *TaskStorer) Delete(id string) error {
	tx := g.DB.Delete(&types.Task{}, "id = ?", id)
	return tx.Error
}

func (g *TaskStorer) FindByID(id string) (*types.Task, error) {
	var task types.Task
	tx := g.DB.First(&task, "id = ?", id)
	return &task, tx.Error
}

func (g *TaskStorer) FindByProjectID(id string) (*[]types.Task, error) {
	var tasks []types.Task
	tx := g.DB.Find(&tasks, "project_id = ?", id)
	return &tasks, tx.Error
}

func (g *TaskStorer) FindByCodeName(codeName string) (*[]types.Task, error) {
	var schedules []types.Task
	tx := g.DB.Find(&schedules, "code_name = ?", codeName)
	return &schedules, tx.Error
}

func (g *TaskStorer) SetStatus(id, status string) error {
	tx := g.DB.Model(&types.Task{}).Where("id = ?", id).Update("status", status)
	return tx.Error
}

func (g *TaskStorer) Next() (*types.Task, error) {
	var task types.Task
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
		return nil, storer.TaskNone
	}
	return &task, tx.Error
}
