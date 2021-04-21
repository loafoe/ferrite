package code

import (
	"gorm.io/gorm"
)

type GormStorer struct {
	DB *gorm.DB
}

func (c *GormStorer) Create(code Code) (*Code, error) {
	tx := c.DB.Create(code)
	err := tx.Commit().Error
	if err != nil {
		return nil, err
	}
	// TODO: retrieve from DB
	return &code, nil
}

func (c *GormStorer) Delete(code Code) error {
	panic("implement me")
}

func (c *GormStorer) Update(code Code) error {
	panic("implement me")
}

func (c *GormStorer) FindByID(id string) (Code, error) {
	panic("implement me")
}

func (c *GormStorer) SaveCredentials(creds DockerCredentials) error {
	tx := c.DB.Create(creds)
	return tx.Error
}
