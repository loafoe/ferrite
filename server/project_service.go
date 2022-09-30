package server

import (
	"net/http"
	"strings"

	"github.com/loafoe/ferrite/storer"
	"github.com/loafoe/ferrite/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ProjectService struct {
	Storer *storer.Ferrite
}

type projectResponse struct {
	Message string `json:"msg"`
}

func (h *ProjectService) Create(c echo.Context) error {
	var project types.Project
	if err := c.Bind(&project); err != nil {
		return c.JSON(http.StatusBadRequest, projectResponse{err.Error()})
	}
	id := strings.Replace(uuid.New().String(), "-", "", -1)
	project.ID = id
	createdProject, err := h.Storer.Project.Create(project)
	if err != nil {
		return c.JSON(http.StatusBadRequest, projectResponse{err.Error()})
	}
	return c.JSON(http.StatusCreated, createdProject)
}

func (h *ProjectService) Get(c echo.Context) error {
	projectID := c.Param("project")
	foundProject, err := h.Storer.Project.FindByID(projectID)
	if err != nil {
		return c.JSON(http.StatusNotFound, projectResponse{err.Error()})
	}
	return c.JSON(http.StatusOK, foundProject)
}
