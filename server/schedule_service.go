package server

import (
	"ferrite/storer"
	"ferrite/types"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ScheduleService struct {
	Storer *storer.Ferrite
}

type scheduleResponse struct {
	Message string `json:"msg"`
}

func (g *ScheduleService) Create(c echo.Context) error {
	var createSchedules struct {
		Schedules []types.Schedule `json:"schedules"`
	}
	projectID := c.Param("project")
	p, err := g.Storer.Project.FindByID(projectID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, scheduleResponse{"invalid or unknown project"})
	}
	if err := c.Bind(&createSchedules); err != nil {
		return c.JSON(http.StatusBadRequest, scheduleResponse{err.Error()})
	}
	var createdSchedules struct {
		Schedules []types.Schedule `json:"schedules"`
	}
	for _, schedule := range createSchedules.Schedules {
		schedule.ProjectID = p.ID
		if schedule.ID != "" {
			return c.JSON(http.StatusBadRequest, scheduleResponse{"cannot update schedule"})
		}
		id := strings.Replace(uuid.New().String(), "-", "", -1)
		schedule.ID = id
		now := time.Now()
		schedule.CreatedAt = &now
		schedule.UpdatedAt = &now
		createdSchedule, err := g.Storer.Schedule.Create(schedule)
		if err != nil {
			return c.JSON(http.StatusBadRequest, scheduleResponse{err.Error()})
		}
		createdSchedules.Schedules = append(createSchedules.Schedules, *createdSchedule)
	}
	return c.JSON(http.StatusOK, createdSchedules)
}

func (g *ScheduleService) Delete(c echo.Context) error {
	projectID := c.Param("project")
	scheduleID := c.Param("schedule")
	p, err := g.Storer.Project.FindByID(projectID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, scheduleResponse{"invalid or unknown project"})
	}
	schedule, err := g.Storer.Schedule.FindByID(scheduleID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, scheduleResponse{"invalid or unknown schedule"})
	}
	if schedule.ProjectID != p.ID {
		return c.JSON(http.StatusBadRequest, scheduleResponse{"invalid request"})
	}
	if err := g.Storer.Schedule.Delete(schedule.ID); err != nil {
		return c.JSON(http.StatusBadRequest, scheduleResponse{err.Error()})
	}
	return c.JSON(http.StatusOK, scheduleResponse{"Cancelled"})
}

func (g *ScheduleService) Find(c echo.Context) error {
	projectID := c.Param("project")
	p, err := g.Storer.Project.FindByID(projectID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, scheduleResponse{"invalid or unknown project"})
	}
	schedules, err := g.Storer.Schedule.FindByProjectID(p.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, scheduleResponse{err.Error()})
	}
	var scheduleResponse struct {
		Schedules []types.Schedule `json:"schedules"`
	}
	for _, schedule := range *schedules {
		scheduleResponse.Schedules = append(scheduleResponse.Schedules, schedule)
	}
	return c.JSON(http.StatusOK, scheduleResponse)
}

func (g *ScheduleService) Get(c echo.Context) error {
	var schedule types.Schedule
	projectID := c.Param("project")
	if err := c.Bind(&schedule); err != nil {
		return c.JSON(http.StatusBadRequest, scheduleResponse{err.Error()})
	}
	schedule.ProjectID = projectID
	scheduleID := c.Param("schedule")
	foundSchedule, err := g.Storer.Schedule.FindByID(scheduleID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, scheduleResponse{err.Error()})
	}
	return c.JSON(http.StatusOK, foundSchedule)
}
