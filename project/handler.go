package project

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Storer Storer
}

type projectResponse struct {
	Message string `json:"msg"`
}

func (h *Handler) Create(c echo.Context) error {
	var project Project
	if err := c.Bind(&project); err != nil {
		return c.JSON(http.StatusBadRequest, projectResponse{err.Error()})
	}
	id := strings.Replace(uuid.New().String(), "-", "", -1)
	project.ID = id
	createdProject, err := h.Storer.Create(project)
	if err != nil {
		return c.JSON(http.StatusBadRequest, projectResponse{err.Error()})
	}
	return c.JSON(http.StatusCreated, createdProject)
}

func (h *Handler) Get(c echo.Context) error {
	projectID := c.Param("project")
	foundProject, err := h.Storer.FindByID(projectID)
	if err != nil {
		return c.JSON(http.StatusNotFound, projectResponse{err.Error()})
	}
	return c.JSON(http.StatusOK, foundProject)
}
