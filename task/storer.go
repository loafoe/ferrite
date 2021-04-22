package task

type Storer interface {
	Create(schedule Task) (*Task, error)
	Delete(id string) error
	FindByID(id string) (*Task, error)
	FindByProjectID(id string) (*[]Task, error)
}
