package task

import (
	"ferrite/project"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Storer        Storer
	ProjectStorer project.Storer
}

type taskResponse struct {
	Message string `json:"msg"`
}

func (g *Handler) Create(c echo.Context) error {
	var createTasks struct {
		Tasks []Task `json:"tasks"`
	}
	projectID := c.Param("project")
	p, err := g.ProjectStorer.FindByID(projectID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, taskResponse{"invalid or unknown project"})
	}
	if err := c.Bind(&createTasks); err != nil {
		return c.JSON(http.StatusBadRequest, taskResponse{err.Error()})
	}
	var createdTasks struct {
		Tasks []Task `json:"tasks"`
	}
	for _, task := range createTasks.Tasks {
		task.ProjectID = p.ID
		if task.ID != "" {
			return c.JSON(http.StatusBadRequest, taskResponse{"cannot update task"})
		}
		id := strings.Replace(uuid.New().String(), "-", "", -1)
		task.ID = id
		now := time.Now()
		task.CreatedAt = &now
		task.UpdatedAt = &now
		task.StartTime = &now
		task.EndTime = &now
		createdTask, err := g.Storer.Create(task)
		if err != nil {
			return c.JSON(http.StatusBadRequest, taskResponse{err.Error()})
		}
		createdTasks.Tasks = append(createdTasks.Tasks, *createdTask)
	}
	return c.JSON(http.StatusOK, createdTasks)
}

func (g *Handler) Delete(c echo.Context) error {
	projectID := c.Param("project")
	taskID := c.Param("task")
	p, err := g.ProjectStorer.FindByID(projectID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, taskResponse{"invalid or unknown project"})
	}
	task, err := g.Storer.FindByID(taskID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, taskResponse{"invalid or unknown task"})
	}
	if task.ProjectID != p.ID {
		return c.JSON(http.StatusBadRequest, taskResponse{"invalid request"})
	}
	if err := g.Storer.Delete(task.ID); err != nil {
		return c.JSON(http.StatusBadRequest, taskResponse{err.Error()})
	}
	return c.JSON(http.StatusOK, taskResponse{"Cancelled"})
}

func (g *Handler) Find(c echo.Context) error {
	projectID := c.Param("project")
	p, err := g.ProjectStorer.FindByID(projectID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, taskResponse{"invalid or unknown project"})
	}
	schedules, err := g.Storer.FindByProjectID(p.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, taskResponse{err.Error()})
	}
	return c.JSON(http.StatusOK, schedules)
}

func (g *Handler) Get(c echo.Context) error {
	var task Task
	projectID := c.Param("project")
	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusBadRequest, taskResponse{err.Error()})
	}
	task.ProjectID = projectID
	foundTask, err := g.Storer.FindByID(task.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, taskResponse{err.Error()})
	}
	return c.JSON(http.StatusOK, foundTask)
}
