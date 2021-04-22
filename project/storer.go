package project

type Storer interface {
	Create(project Project) (*Project, error)
	FindByID(id string) (*Project, error)
}
