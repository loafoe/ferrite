package main

import (
	"fmt"
	"log"
	"net/http"
	"siderite-server/code"
	"siderite-server/project"
	"siderite-server/schedule"
	"siderite-server/token"

	"github.com/cloudfoundry-community/gautocloud"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/philips-software/gautocloud-connectors/hsdp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Database
	var svc *hsdp.PostgreSQLClient

	err := gautocloud.Inject(&svc)
	if err != nil {
		fmt.Printf("error discovering database: %v\n", err)
		return
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: svc.DB,
	}))
	if err != nil {
		fmt.Printf("error configuring gorm: %v\n", err)
		return
	}
	// Auto Migrate
	_ = db.AutoMigrate(&code.Code{})
	_ = db.AutoMigrate(&project.Project{})
	_ = db.AutoMigrate(&schedule.Schedule{})
	_ = db.AutoMigrate(&code.DockerCredentials{})

	codeHandler := &code.Handler{
		Storer: &code.GormStorer{
			DB: db,
		},
		ProjectStorer: &project.GormStorer{
			DB: db,
		},
	}

	projectHandler := &project.Handler{
		Storer: &project.GormStorer{
			DB: db,
		},
	}

	// API
	e := echo.New()
	e.Use(token.Checker("foo"))
	e.Use(middleware.Logger())

	// Projects
	e.POST("/2/projects", projectHandler.Create)

	// Codes
	e.POST("/2/projects/:project/codes", codeHandler.Create)
	e.GET("/2/projects/:project/codes", codeHandler.Find)
	e.GET("/2/projects/:project/codes/:code_id", codeHandler.Get)
	e.DELETE("/2/projects/:project/codes/:code_id", codeHandler.Delete)
	e.POST("/2/projects/:project/credentials", codeHandler.Credentials)

	// Schedules
	e.POST("/2/projects/:project/schedules", NotImplemented)
	e.GET("/2/projects/:project/schedules", NotImplemented)
	e.GET("/2/projects/:project/schedules/:schedule_id", NotImplemented)
	e.POST("/2/projects/:project/schedules/:schedule_id/cancel", NotImplemented)

	// Tasks
	e.GET("/2/projects/:project/tasks", NotImplemented)
	e.GET("/2/projects/:project/tasks/:task_id", NotImplemented)
	e.POST("/2/projects/:project/tasks", NotImplemented)
	e.POST("/2/projects/:project/tasks/:task_id/cancel", NotImplemented)

	log.Fatal(e.Start(":8080"))
}

var notImplemented = struct {
	Message string `json:"message"`
}{
	"Not implemented",
}

func NotImplemented(c echo.Context) error {
	c.JSON(http.StatusNotImplemented, notImplemented)
	return nil
}
