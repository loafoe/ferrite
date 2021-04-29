package server

import (
	"github.com/philips-labs/ferrite/storer"
	"github.com/philips-labs/ferrite/types"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TaskService struct {
	Storer *storer.Ferrite
}

type taskResponse struct {
	Message string `json:"msg"`
}

func (g *TaskService) Create(c echo.Context) error {
	var createTasks struct {
		Tasks []types.Task `json:"tasks"`
	}
	projectID := c.Param("project")
	p, err := g.Storer.Project.FindByID(projectID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, taskResponse{"invalid or unknown project"})
	}
	if err := c.Bind(&createTasks); err != nil {
		return c.JSON(http.StatusBadRequest, taskResponse{err.Error()})
	}
	var createdTasks struct {
		Tasks []types.Task `json:"tasks"`
	}
	for _, task := range createTasks.Tasks {
		task.ProjectID = p.ID
		if task.ID != "" {
			return c.JSON(http.StatusBadRequest, taskResponse{"cannot update task"})
		}
		id := strings.Replace(uuid.New().String(), "-", "", -1)
		task.ID = id
		now := time.Now()
		task.Status = "new"
		task.CreatedAt = &now
		task.UpdatedAt = &now
		task.StartTime = &now
		task.EndTime = &now
		createdTask, err := g.Storer.Task.Create(task)
		if err != nil {
			return c.JSON(http.StatusBadRequest, taskResponse{err.Error()})
		}
		createdTasks.Tasks = append(createdTasks.Tasks, *createdTask)
	}
	return c.JSON(http.StatusOK, createdTasks)
}

func (g *TaskService) Delete(c echo.Context) error {
	projectID := c.Param("project")
	taskID := c.Param("task")
	p, err := g.Storer.Project.FindByID(projectID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, taskResponse{"invalid or unknown project"})
	}
	task, err := g.Storer.Task.FindByID(taskID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, taskResponse{"invalid or unknown task"})
	}
	if task.ProjectID != p.ID {
		return c.JSON(http.StatusBadRequest, taskResponse{"invalid request"})
	}
	if err := g.Storer.Task.Delete(task.ID); err != nil {
		return c.JSON(http.StatusBadRequest, taskResponse{err.Error()})
	}
	return c.JSON(http.StatusOK, taskResponse{"Cancelled"})
}

func (g *TaskService) Find(c echo.Context) error {
	projectID := c.Param("project")
	p, err := g.Storer.Project.FindByID(projectID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, taskResponse{"invalid or unknown project"})
	}
	tasks, err := g.Storer.Task.FindByProjectID(p.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, taskResponse{err.Error()})
	}
	var taskResponse struct {
		Tasks []types.Task `json:"tasks"`
	}
	for _, task := range *tasks {
		taskResponse.Tasks = append(taskResponse.Tasks, task)
	}
	return c.JSON(http.StatusOK, taskResponse)
}

func (g *TaskService) Get(c echo.Context) error {
	var task types.Task
	projectID := c.Param("project")
	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusBadRequest, taskResponse{err.Error()})
	}
	task.ProjectID = projectID
	foundTask, err := g.Storer.Task.FindByID(task.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, taskResponse{err.Error()})
	}
	return c.JSON(http.StatusOK, foundTask)
}
