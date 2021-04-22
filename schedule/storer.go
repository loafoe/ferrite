package schedule

type Storer interface {
	Create(schedule Schedule) (*Schedule, error)
	Delete(id string) error
	FindByID(id string) (*Schedule, error)
	FindByProjectID(id string) (*[]Schedule, error)
	FindByCodeName(codeName string) (*[]Schedule, error)
}
