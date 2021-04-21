package code

type Storer interface {
	Create(code Code) (*Code, error)
	Delete(code Code) error
	Update(code Code) error
	FindByID(id string) (Code, error)
	SaveCredentials(creds DockerCredentials) error
}
