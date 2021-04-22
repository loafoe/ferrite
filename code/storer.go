package code

type Storer interface {
	Create(code Code) (*Code, error)
	Delete(id string) error
	Update(code Code) error
	FindByID(id string) (*Code, error)
	FindByProjectID(id string) (*[]Code, error)
	SaveCredentials(creds DockerCredentials) error
}
