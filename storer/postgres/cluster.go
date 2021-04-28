package postgres

import (
	"ferrite/types"

	"gorm.io/gorm"
)

type ClusterStorer struct {
	DB *gorm.DB
}

func (c *ClusterStorer) Create(cluster types.Cluster) (*types.Cluster, error) {
	tx := c.DB.Create(cluster)
	if tx.Error != nil {
		return nil, tx.Error
	}
	createdCluster, err := c.FindByID(cluster.ID)
	return createdCluster, err
}

func (c *ClusterStorer) Delete(id string) error {
	tx := c.DB.Delete(&types.Cluster{}, "id = ?", id)
	return tx.Error
}

func (c *ClusterStorer) FindByID(id string) (*types.Cluster, error) {
	var cluster types.Cluster
	tx := c.DB.First(&cluster, "id = ?", id)
	return &cluster, tx.Error
}
