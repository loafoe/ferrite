package postgres

import (
	"github.com/philips-labs/ferrite/types"

	"gorm.io/gorm"
)

type ProjectStorer struct {
	DB *gorm.DB
}

func (g *ProjectStorer) Create(project types.Project) (*types.Project, error) {
	tx := g.DB.Create(project)
	if tx.Error != nil {
		return nil, tx.Error
	}
	createdProject, err := g.FindByID(project.ID)
	return createdProject, err
}

func (g *ProjectStorer) FindByID(id string) (*types.Project, error) {
	var project types.Project
	tx := g.DB.First(&project, "id = ?", id)
	return &project, tx.Error
}

func (g *ProjectStorer) FindLatest() (*types.Project, error) {
	var project types.Project
	tx := g.DB.First(&project).Order("created_at DESC")
	return &project, tx.Error
}
