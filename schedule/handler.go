package schedule

import (
	"ferrite/project"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Storer        Storer
	ProjectStorer project.Storer
}

type scheduleResponse struct {
	Message string `json:"msg"`
}

func (g *Handler) Create(c echo.Context) error {
	var createSchedules struct {
		Schedules []Schedule `json:"schedules"`
	}
	projectID := c.Param("project")
	p, err := g.ProjectStorer.FindByID(projectID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, scheduleResponse{"invalid or unknown project"})
	}
	if err := c.Bind(&createSchedules); err != nil {
		return c.JSON(http.StatusBadRequest, scheduleResponse{err.Error()})
	}
	var createdSchedules []Schedule
	for _, schedule := range createSchedules.Schedules {
		schedule.ProjectID = p.ID
		if schedule.ID != "" {
			return c.JSON(http.StatusBadRequest, scheduleResponse{"cannot update schedule"})
		}
		id := strings.Replace(uuid.New().String(), "-", "", -1)
		schedule.ID = id
		createdSchedule, err := g.Storer.Create(schedule)
		if err != nil {
			return c.JSON(http.StatusBadRequest, scheduleResponse{err.Error()})
		}
		createdSchedules = append(createdSchedules, *createdSchedule)
	}
	return c.JSON(http.StatusCreated, createdSchedules)
}

func (g *Handler) Delete(c echo.Context) error {
	projectID := c.Param("project")
	scheduleID := c.Param("schedule")
	p, err := g.ProjectStorer.FindByID(projectID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, scheduleResponse{"invalid or unknown project"})
	}
	schedule, err := g.Storer.FindByID(scheduleID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, scheduleResponse{"invalid or unknown schedule"})
	}
	if schedule.ProjectID != p.ID {
		return c.JSON(http.StatusBadRequest, scheduleResponse{"invalid request"})
	}
	if err := g.Storer.Delete(schedule.ID); err != nil {
		return c.JSON(http.StatusBadRequest, scheduleResponse{err.Error()})
	}
	return c.JSON(http.StatusOK, scheduleResponse{"Cancelled"})
}

func (g *Handler) Find(c echo.Context) error {
	projectID := c.Param("project")
	p, err := g.ProjectStorer.FindByID(projectID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, scheduleResponse{"invalid or unknown project"})
	}
	schedules, err := g.Storer.FindByProjectID(p.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, scheduleResponse{err.Error()})
	}
	return c.JSON(http.StatusOK, schedules)
}

func (g *Handler) Get(c echo.Context) error {
	var schedule Schedule
	projectID := c.Param("project")
	if err := c.Bind(&schedule); err != nil {
		return c.JSON(http.StatusBadRequest, scheduleResponse{err.Error()})
	}
	schedule.ProjectID = projectID
	scheduleID := c.Param("schedule")
	foundSchedule, err := g.Storer.FindByID(scheduleID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, scheduleResponse{err.Error()})
	}
	return c.JSON(http.StatusOK, foundSchedule)
}
