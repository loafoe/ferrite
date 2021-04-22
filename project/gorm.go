package project

import "gorm.io/gorm"

type GormStorer struct {
	DB *gorm.DB
}

func (g *GormStorer) Create(project Project) (*Project, error) {
	tx := g.DB.Create(project)
	if tx.Error != nil {
		return nil, tx.Error
	}
	createdProject, err := g.FindByID(project.ID)
	return createdProject, err
}

func (g *GormStorer) FindByID(id string) (*Project, error) {
	var project Project
	tx := g.DB.First(&project, "id = ?", id)
	return &project, tx.Error
}
