package resource

// Service recieves commands from handlers and forwards them to the repository
type Service interface {
	Get() (Response, error)
	Create(*Input) (Response, error)
}

type service struct {
	repository Repository
}

// NewService creates a service instance with the repository passed
func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (service *service) Get() (Response, error) {
	return service.repository.Get()
}

func (service *service) Create(resource *Input) (Response, error) {
	return service.repository.Create(resource)
}
