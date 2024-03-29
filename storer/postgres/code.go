package postgres

import (
	"github.com/loafoe/ferrite/types"

	"gorm.io/gorm"
)

type CodeStorer struct {
	DB *gorm.DB
}

func (c *CodeStorer) Create(code types.Code) (*types.Code, error) {
	tx := c.DB.Create(code)
	if tx.Error != nil {
		return nil, tx.Error
	}
	createdCode, err := c.FindByID(code.ID)
	return createdCode, err
}

func (c *CodeStorer) Delete(id string) error {
	tx := c.DB.Delete(&types.Code{}, "id = ?", id)
	return tx.Error
}

func (c *CodeStorer) Update(code types.Code) error {
	var foundCode types.Code
	tx := c.DB.First(&foundCode, "id = ?", code.ID)
	if tx.Error != nil {
		return tx.Error
	}
	// Only update image for now
	tx = c.DB.Model(&types.Code{}).Where("id = ?", code.ID).Update("image", code.Image)
	return tx.Error
}

func (c *CodeStorer) FindByID(id string) (*types.Code, error) {
	var code types.Code
	tx := c.DB.First(&code, "id = ?", id)
	return &code, tx.Error
}

func (c *CodeStorer) FindByName(name string) (*types.Code, error) {
	var code types.Code
	tx := c.DB.First(&code, "name = ?", name)
	return &code, tx.Error
}

func (c *CodeStorer) SaveCredentials(creds types.DockerCredentials) error {
	tx := c.DB.Create(creds)
	return tx.Error
}

func (c *CodeStorer) FindByProjectID(id string) (*[]types.Code, error) {
	var codes []types.Code
	tx := c.DB.Find(&codes, "project_id = ?", id)
	return &codes, tx.Error
}
