package cluster

type Storer interface {
	Create(cluster Cluster) (*Cluster, error)
	Delete(id string) error
	FindByID(id string) (*Cluster, error)
}
