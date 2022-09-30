package main

import (
	"fmt"
	"log"
	"os"

	"github.com/loafoe/ferrite/server"
	pg "github.com/loafoe/ferrite/storer/postgres"
	"github.com/loafoe/ferrite/token"
	"github.com/loafoe/ferrite/worker"

	"github.com/cloudfoundry-community/gautocloud"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/philips-software/gautocloud-connectors/hsdp"
)

func main() {
	// Database
	var pgclient *hsdp.PostgreSQLClient

	err := gautocloud.Inject(&pgclient)
	if err != nil {
		fmt.Printf("error discovering database: %v\n", err)
		return
	}

	fs, err := pg.New(pgclient.DB)
	if err != nil {
		fmt.Printf("error configuring postgres storer: %v\n", err)
		return
	}

	// Check if we should go in worker mode
	if len(os.Args) > 1 && os.Args[1] == "worker" {
		_, err := worker.Start(fs)
		if err != nil {
			fmt.Printf("error starting worker: %v\n", err)
		}
		select {}
		return
	}

	svc, err := server.New(fs)
	if err != nil {
		fmt.Printf("error initializing server: %v\n", err)
		return
	}

	// API
	e := echo.New()
	e.Use(token.Checker(os.Getenv("TOKEN")))
	e.Use(middleware.Logger())

	// Bootstrap
	e.POST("/bootstrap", svc.Bootstrap.Bootstrap)

	// Clusters
	e.POST("/2/clusters", svc.Cluster.Create)
	e.GET("/2/clusters/:cluster", svc.Cluster.Get)

	// Projects
	e.POST("/2/projects", svc.Project.Create)
	e.GET("/2/projects/:project", svc.Project.Get)

	// Codes
	e.POST("/2/projects/:project/codes", svc.Code.Create)
	e.GET("/2/projects/:project/codes", svc.Code.Find)
	e.GET("/2/projects/:project/codes/:code", svc.Code.Get)
	e.DELETE("/2/projects/:project/codes/:code", svc.Code.Delete)
	e.POST("/2/projects/:project/credentials", svc.Code.Credentials)

	// Schedules
	e.POST("/2/projects/:project/schedules", svc.Schedule.Create)
	e.GET("/2/projects/:project/schedules", svc.Schedule.Find)
	e.GET("/2/projects/:project/schedules/:schedule", svc.Schedule.Get)
	e.POST("/2/projects/:project/schedules/:schedule/cancel", svc.Schedule.Delete)

	// Tasks
	e.GET("/2/projects/:project/tasks", svc.Task.Find)
	e.GET("/2/projects/:project/tasks/:task", svc.Task.Get)
	e.POST("/2/projects/:project/tasks", svc.Task.Create)
	e.POST("/2/projects/:project/tasks/:task/cancel", svc.Task.Delete)

	log.Fatal(e.Start(":8080"))
}
