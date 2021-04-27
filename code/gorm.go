package code

import (
	"gorm.io/gorm"
)

type GormStorer struct {
	DB *gorm.DB
}

func (c *GormStorer) Create(code Code) (*Code, error) {
	tx := c.DB.Create(code)
	if tx.Error != nil {
		return nil, tx.Error
	}
	createdCode, err := c.FindByID(code.ID)
	return createdCode, err
}

func (c *GormStorer) Delete(id string) error {
	tx := c.DB.Delete(&Code{}, "id = ?", id)
	return tx.Error
}

func (c *GormStorer) Update(code Code) error {
	panic("implement me")
}

func (c *GormStorer) FindByID(id string) (*Code, error) {
	var code Code
	tx := c.DB.First(&code, "id = ?", id)
	return &code, tx.Error
}

func (c *GormStorer) FindByName(name string) (*Code, error) {
	var code Code
	tx := c.DB.First(&code, "name = ?", name)
	return &code, tx.Error
}

func (c *GormStorer) SaveCredentials(creds DockerCredentials) error {
	tx := c.DB.Create(creds)
	return tx.Error
}

func (c *GormStorer) FindByProjectID(id string) (*[]Code, error) {
	var codes []Code
	tx := c.DB.Find(&codes, "project_id = ?", id)
	return &codes, tx.Error
}
