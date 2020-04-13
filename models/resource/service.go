package resource

// Service recieves commands from handlers and forwards them to the repository
type Service interface {
	Get() (Response, error)
	Create(*Input) (Element, error)
	GetByID(string) (Element, error)
	Delete(string) error
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

func (service *service) GetByID(id string) (Element, error) {
	return service.repository.GetByID(id)
}

func (service *service) Create(resource *Input) (Element, error) {
	if resource.Data.Relationships != nil {
		_, err := service.repository.GetByID(resource.Data.Relationships.Parent.Data.ID)
		if err != nil {
			return Element{}, err
		}
		return service.repository.Create(resource)
	}
	return service.repository.Create(resource)
}

func (service *service) Delete(id string) error {
	return service.repository.Delete(id)
}
