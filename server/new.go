package server

import (
	"github.com/loafoe/ferrite/storer"
)

type Ferrite struct {
	Code      CodeService
	Project   ProjectService
	Cluster   ClusterService
	Task      TaskService
	Schedule  ScheduleService
	Bootstrap BootstrapService
}

func New(fs *storer.Ferrite) (*Ferrite, error) {
	return &Ferrite{
		Code: CodeService{
			Storer: fs,
		},
		Project: ProjectService{
			Storer: fs,
		},
		Cluster: ClusterService{
			Storer: fs,
		},
		Task: TaskService{
			Storer: fs,
		},
		Schedule: ScheduleService{
			Storer: fs,
		},
		Bootstrap: BootstrapService{
			Storer: fs,
		},
	}, nil
}
