package cluster

import "gorm.io/gorm"

type GormStorer struct {
	DB *gorm.DB
}

func (c *GormStorer) Create(cluster Cluster) (*Cluster, error) {
	tx := c.DB.Create(cluster)
	if tx.Error != nil {
		return nil, tx.Error
	}
	createdCluster, err := c.FindByID(cluster.ID)
	return createdCluster, err
}

func (c *GormStorer) Delete(id string) error {
	tx := c.DB.Delete(&Cluster{}, "id = ?", id)
	return tx.Error
}

func (c *GormStorer) FindByID(id string) (*Cluster, error) {
	var cluster Cluster
	tx := c.DB.First(&cluster, "id = ?", id)
	return &cluster, tx.Error
}
