package code

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/philips-software/go-hsdp-api/iron"
)

type Handler struct {
	Storer Storer
}

func (g *Handler) Create(c echo.Context) error {
	id := strings.Replace(uuid.New().String(), "-", "", -1)
	var code Code
	if err := c.Bind(&code); err != nil {
		return err
	}
	code.ID = id
	createdCode, err := g.Storer.Create(code)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, createdCode)
}

func (g *Handler) Delete(c echo.Context) error {
	panic("implement me")
}

func (g *Handler) Update(c echo.Context) error {
	panic("implement me")
}

func (g *Handler) Find(c echo.Context) error {
	panic("implement me")
}

func (g *Handler) Get(c echo.Context) error {
	panic("implement me")
}

func (g *Handler) Credentials(c echo.Context) error {
	var authRequest struct {
		Message string `json:"msg"`
	}
	if err := c.Bind(&authRequest); err != nil {
		return err
	}
	data, err := base64.StdEncoding.DecodeString(authRequest.Message)
	if err != nil {
		return err
	}
	var creds iron.DockerCredentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return err
	}
	if dbErr := g.Storer.SaveCredentials(DockerCredentials{
		DockerCredentials: creds,
	}); dbErr != nil {
		return dbErr
	}
	return c.JSON(http.StatusOK, creds)
}
