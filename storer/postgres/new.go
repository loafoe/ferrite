package postgres

import (
	"database/sql"

	"github.com/loafoe/ferrite/storer"
	"github.com/loafoe/ferrite/types"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(conn *sql.DB) (*storer.Ferrite, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: conn,
	}))
	if err != nil {
		return nil, err
	}
	// Auto Migrate
	_ = db.AutoMigrate(&types.Cluster{})
	_ = db.AutoMigrate(&types.Code{})
	_ = db.AutoMigrate(&types.Project{})
	_ = db.AutoMigrate(&types.Schedule{})
	_ = db.AutoMigrate(&types.DockerCredentials{})
	_ = db.AutoMigrate(&types.Task{})

	return &storer.Ferrite{
		Code: &CodeStorer{
			DB: db,
		},
		Project: &ProjectStorer{
			DB: db,
		},
		Task: &TaskStorer{
			DB: db,
		},
		Cluster: &ClusterStorer{
			DB: db,
		},
		Schedule: &ScheduleStorer{
			DB: db,
		},
	}, nil
}
