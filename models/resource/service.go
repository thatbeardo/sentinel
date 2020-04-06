package resource

// Service recieves commands from handlers and forwards them to the repository
type Service interface {
	Get() ([]*Resource, error)
	Create(*Resource) (*Resource, error)
}

type service struct {
	repository Repository
}

// NewService creates a service instance with the repository passed
func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (service *service) Get() ([]*Resource, error) {
	return service.repository.Get()
}

func (service *service) Create(resource *Resource) (*Resource, error) {
	return service.Create(resource)
}
