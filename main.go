package main

import (
	"ferrite/cluster"
	"ferrite/code"
	"ferrite/project"
	"ferrite/schedule"
	"ferrite/task"
	"ferrite/token"
	"fmt"
	"log"
	"os"

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
	_ = db.AutoMigrate(&cluster.Cluster{})
	_ = db.AutoMigrate(&code.Code{})
	_ = db.AutoMigrate(&project.Project{})
	_ = db.AutoMigrate(&schedule.Schedule{})
	_ = db.AutoMigrate(&code.DockerCredentials{})
	_ = db.AutoMigrate(&task.Task{})

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

	scheduleHandler := &schedule.Handler{
		Storer: &schedule.GormStorer{
			DB: db,
		},
		ProjectStorer: &project.GormStorer{
			DB: db,
		},
	}

	taskHandler := &task.Handler{
		Storer: &task.GormStorer{
			DB: db,
		},
		ProjectStorer: &project.GormStorer{
			DB: db,
		},
	}

	clusterHandler := &cluster.Handler{
		Storer: &cluster.GormStorer{
			DB: db,
		},
	}

	// API
	e := echo.New()
	e.Use(token.Checker(os.Getenv("TOKEN")))
	e.Use(middleware.Logger())

	// Clusters
	e.POST("/2/clusters", clusterHandler.Create)
	e.GET("/2/clusters/:cluster", clusterHandler.Get)

	// Projects
	e.POST("/2/projects", projectHandler.Create)
	e.GET("/2/projects/:project", projectHandler.Get)

	// Codes
	e.POST("/2/projects/:project/codes", codeHandler.Create)
	e.GET("/2/projects/:project/codes", codeHandler.Find)
	e.GET("/2/projects/:project/codes/:code", codeHandler.Get)
	e.DELETE("/2/projects/:project/codes/:code", codeHandler.Delete)
	e.POST("/2/projects/:project/credentials", codeHandler.Credentials)

	// Schedules
	e.POST("/2/projects/:project/schedules", scheduleHandler.Create)
	e.GET("/2/projects/:project/schedules", scheduleHandler.Find)
	e.GET("/2/projects/:project/schedules/:schedule", scheduleHandler.Get)
	e.POST("/2/projects/:project/schedules/:schedule/cancel", scheduleHandler.Delete)

	// Tasks
	e.GET("/2/projects/:project/tasks", taskHandler.Find)
	e.GET("/2/projects/:project/tasks/:task", taskHandler.Get)
	e.POST("/2/projects/:project/tasks", taskHandler.Create)
	e.POST("/2/projects/:project/tasks/:task/cancel", taskHandler.Delete)

	log.Fatal(e.Start(":8080"))
}
