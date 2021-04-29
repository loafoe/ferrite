package storer

import "github.com/philips-labs/ferrite/types"

type Cluster interface {
	Create(cluster types.Cluster) (*types.Cluster, error)
	Delete(id string) error
	FindByID(id string) (*types.Cluster, error)
	FindLatest() (*types.Cluster, error)
}
