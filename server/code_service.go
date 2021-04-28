package server

import (
	"encoding/base64"
	"encoding/json"
	"ferrite/storer"
	"ferrite/types"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/philips-software/go-hsdp-api/iron"
)

type CodeService struct {
	Storer *storer.Ferrite
}

type codeResponse struct {
	Message string `json:"msg"`
}

func (g *CodeService) Create(c echo.Context) error {
	var code types.Code
	projectID := c.Param("project")
	p, err := g.Storer.Code.FindByID(projectID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, codeResponse{"invalid or unknown project"})
	}
	data := c.FormValue("data")
	err = json.Unmarshal([]byte(data), &code)
	if err != nil {
		return c.JSON(http.StatusBadRequest, codeResponse{err.Error()})
	}
	code.ProjectID = p.ID
	if code.ID != "" { // Update
		if err := g.Storer.Code.Update(code); err != nil {
			return c.JSON(http.StatusBadRequest, codeResponse{err.Error()})
		}
		return c.JSON(http.StatusOK, code)
	}
	id := strings.Replace(uuid.New().String(), "-", "", -1)
	code.ID = id
	now := time.Now()
	code.CreatedAt = &now
	createdCode, err := g.Storer.Code.Create(code)
	if err != nil {
		return c.JSON(http.StatusBadRequest, codeResponse{err.Error()})
	}
	return c.JSON(http.StatusCreated, createdCode)
}

func (g *CodeService) Delete(c echo.Context) error {
	projectID := c.Param("project")
	codeID := c.Param("code")
	p, err := g.Storer.Project.FindByID(projectID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, codeResponse{"invalid or unknown project"})
	}
	code, err := g.Storer.Code.FindByID(codeID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, codeResponse{"invalid or unknown code"})
	}
	if code.ProjectID != p.ID {
		return c.JSON(http.StatusBadRequest, codeResponse{"invalid request"})
	}
	if err := g.Storer.Code.Delete(code.ID); err != nil {
		return c.JSON(http.StatusBadRequest, codeResponse{err.Error()})
	}
	return c.JSON(http.StatusOK, codeResponse{"Deleted"})
}

func (g *CodeService) Update(c echo.Context) error {
	panic("implement me")
}

func (g *CodeService) Find(c echo.Context) error {
	projectID := c.Param("project")
	p, err := g.Storer.Project.FindByID(projectID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, codeResponse{"invalid or unknown project"})
	}
	codes, err := g.Storer.Code.FindByProjectID(p.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, codeResponse{err.Error()})
	}
	return c.JSON(http.StatusOK, codes)
}

func (g *CodeService) Get(c echo.Context) error {
	projectID := c.Param("project")
	_, err := g.Storer.Project.FindByID(projectID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, codeResponse{"invalid or unknown project"})
	}
	codeID := c.Param("code")
	foundCode, err := g.Storer.Code.FindByID(codeID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, codeResponse{err.Error()})
	}
	return c.JSON(http.StatusOK, foundCode)
}

func (g *CodeService) Credentials(c echo.Context) error {
	var authRequest struct {
		Message string `json:"msg"`
	}
	if err := c.Bind(&authRequest); err != nil {
		return c.JSON(http.StatusBadRequest, codeResponse{err.Error()})
	}
	data, err := base64.StdEncoding.DecodeString(authRequest.Message)
	if err != nil {
		return c.JSON(http.StatusBadRequest, codeResponse{err.Error()})
	}
	var creds iron.DockerCredentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return c.JSON(http.StatusBadRequest, codeResponse{err.Error()})
	}
	if dbErr := g.Storer.Code.SaveCredentials(types.DockerCredentials{
		DockerCredentials: creds,
	}); dbErr != nil {
		return c.JSON(http.StatusBadRequest, codeResponse{dbErr.Error()})
	}
	return c.JSON(http.StatusOK, creds)
}
