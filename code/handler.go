package code

import (
	"encoding/base64"
	"encoding/json"
	"ferrite/project"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/philips-software/go-hsdp-api/iron"
)

type Handler struct {
	Storer        Storer
	ProjectStorer project.Storer
}

type codeResponse struct {
	Message string `json:"msg"`
}

func (g *Handler) Create(c echo.Context) error {
	var code Code
	projectID := c.Param("project")
	p, err := g.ProjectStorer.FindByID(projectID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, codeResponse{"invalid or unknown project"})
	}
	if err := c.Bind(&code); err != nil {
		return c.JSON(http.StatusBadRequest, codeResponse{err.Error()})
	}
	code.ProjectID = p.ID
	if code.ID != "" { // Update
		if err := g.Storer.Update(code); err != nil {
			return c.JSON(http.StatusBadRequest, codeResponse{err.Error()})
		}
		return c.JSON(http.StatusOK, code)
	}
	id := strings.Replace(uuid.New().String(), "-", "", -1)
	code.ID = id
	createdCode, err := g.Storer.Create(code)
	if err != nil {
		return c.JSON(http.StatusBadRequest, codeResponse{err.Error()})
	}
	return c.JSON(http.StatusCreated, createdCode)
}

func (g *Handler) Delete(c echo.Context) error {
	projectID := c.Param("project")
	codeID := c.Param("code")
	p, err := g.ProjectStorer.FindByID(projectID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, codeResponse{"invalid or unknown project"})
	}
	code, err := g.Storer.FindByID(codeID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, codeResponse{"invalid or unknown code"})
	}
	if code.ProjectID != p.ID {
		return c.JSON(http.StatusBadRequest, codeResponse{"invalid request"})
	}
	if err := g.Storer.Delete(code.ID); err != nil {
		return c.JSON(http.StatusBadRequest, codeResponse{err.Error()})
	}
	return c.JSON(http.StatusOK, codeResponse{"Deleted"})
}

func (g *Handler) Update(c echo.Context) error {
	panic("implement me")
}

func (g *Handler) Find(c echo.Context) error {
	projectID := c.Param("project")
	p, err := g.ProjectStorer.FindByID(projectID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, codeResponse{"invalid or unknown project"})
	}
	codes, err := g.Storer.FindByProjectID(p.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, codeResponse{err.Error()})
	}
	return c.JSON(http.StatusOK, codes)
}

func (g *Handler) Get(c echo.Context) error {
	var code Code
	projectID := c.Param("project")
	if err := c.Bind(&code); err != nil {
		return c.JSON(http.StatusBadRequest, codeResponse{err.Error()})
	}
	code.ProjectID = projectID
	foundCode, err := g.Storer.FindByID(code.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, codeResponse{err.Error()})
	}
	return c.JSON(http.StatusOK, foundCode)
}

func (g *Handler) Credentials(c echo.Context) error {
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
	if dbErr := g.Storer.SaveCredentials(DockerCredentials{
		DockerCredentials: creds,
	}); dbErr != nil {
		return c.JSON(http.StatusBadRequest, codeResponse{dbErr.Error()})
	}
	return c.JSON(http.StatusOK, creds)
}
