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
	db.AutoMigrate(&code.Code{})
	db.AutoMigrate(&project.Project{})
	db.AutoMigrate(&schedule.Schedule{})
	db.AutoMigrate(&code.DockerCredentials{})

	codeHandler := &code.Handler{
		Storer: &code.GormStorer{
			DB: db,
		},
	}

	// API
	e := echo.New()
	e.Use(token.Checker("foo"))
	e.Use(middleware.Logger())
	// Codes
	e.POST("/2/projects/:project/codes", codeHandler.Create)
	e.GET("/2/projects/:project/codes", codeHandler.Find)
	e.GET("/2/projects/:project/codes/:code_id", codeHandler.Get)
	e.DELETE("/2/projects/:project/codes/:code_id", codeHandler.Delete)
	e.POST("/2/projects/:project/credentials", codeHandler.Credentials)

	// Schedules
	e.POST("/2/projects/:project/schedules", CreateSchedule)
	e.GET("/2/projects/:project/schedules", GetSchedules)
	e.GET("/2/projects/:project/schedules/:schedule_id", GetSchedule)
	e.POST("/2/projects/:project/schedules/:schedule_id/cancel", CancelSchedule)

	// Tasks
	e.GET("/2/projects/:project/tasks", GetTasks)
	e.GET("/2/projects/:project/tasks/:task_id", GetTask)
	e.POST("/2/projects/:project/tasks", QueueTasks)
	e.POST("/2/projects/:project/tasks/:task_id/cancel", CancelTask)

	log.Fatal(e.Start(":8080"))
}

var notImplemented = struct {
	Message string `json:"message"`
}{
	"Not implemented",
}

func CancelTask(c echo.Context) error {
	c.JSON(http.StatusNotImplemented, notImplemented)
	return nil
}

func QueueTasks(c echo.Context) error {
	c.JSON(http.StatusNotImplemented, notImplemented)
	return nil
}

func GetTask(c echo.Context) error {
	c.JSON(http.StatusNotImplemented, notImplemented)
	return nil
}

func GetTasks(c echo.Context) error {
	c.JSON(http.StatusNotImplemented, notImplemented)
	return nil
}

func CancelSchedule(c echo.Context) error {
	c.JSON(http.StatusNotImplemented, notImplemented)
	return nil
}

func GetSchedule(c echo.Context) error {
	c.JSON(http.StatusNotImplemented, notImplemented)
	return nil
}

func GetSchedules(c echo.Context) error {
	c.JSON(http.StatusNotImplemented, notImplemented)
	return nil
}

func CreateSchedule(c echo.Context) error {
	c.JSON(http.StatusNotImplemented, notImplemented)
	return nil
}

func DeleteCode(c echo.Context) error {
	c.JSON(http.StatusNotImplemented, notImplemented)
	return nil
}

func GetCode(c echo.Context) error {
	c.JSON(http.StatusNotImplemented, notImplemented)
	return nil
}

func GetCodes(c echo.Context) error {
	c.JSON(http.StatusNotImplemented, notImplemented)
	return nil
}

func CreateCode(c echo.Context) error {
	c.JSON(http.StatusNotImplemented, notImplemented)
	return nil
}
