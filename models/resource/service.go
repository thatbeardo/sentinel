package resource

// Service recieves commands from handlers and forwards them to the repository
type Service interface {
	Get() (Response, error)
	GetByID(string) (Element, error)
	Create(*Input) (Element, error)
	Update(string, *Input) (Element, error)
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

func (service *service) Update(id string, resource *Input) (Element, error) {
	if resource.Data.Relationships != nil {
		var err error
		if _, err = service.repository.GetByID(id); err != nil {
			return Element{}, err
		}

		// var parent Element
		_, err = service.repository.GetByID(resource.Data.Relationships.Parent.Data.ID)
		if err != nil {
			return Element{}, err
		}

		response, err := service.repository.UpdateOwnership(id, resource)
		if err != nil {
			return Element{}, err
		}

		return response, nil
	}
	return service.repository.Update(id, resource)
}

func (service *service) Create(resource *Input) (Element, error) {
	if resource.Data.Relationships != nil {
		parent, err := service.repository.GetByID(resource.Data.Relationships.Parent.Data.ID)
		if err != nil {
			return Element{}, err
		}

		var child Element
		child, err = service.repository.Create(resource)
		if err != nil {
			return Element{}, err
		}

		err = service.repository.CreateEdge(child.ID, parent.ID)
		if err != nil {
			return Element{}, err
		}
		return child, nil
	}
	return service.repository.Create(resource)
}

func (service *service) Delete(id string) error {
	return service.repository.Delete(id)
}
