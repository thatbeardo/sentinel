package service

import (
	entity "github.com/thatbeardo/go-sentinel/models/resource"
	"github.com/thatbeardo/go-sentinel/models/resource/repository"
)

// Service recieves commands from handlers and forwards them to the repository
type Service interface {
	Get() (entity.Response, error)
	GetByID(string) (entity.Element, error)
	Create(*entity.Input) (entity.Element, error)
	Update(string, *entity.Input) (entity.Element, error)
	Delete(string) error
}

type service struct {
	repository repository.Repository
}

// NewService creates a service instance with the repository passed
func NewService(repository repository.Repository) Service {
	return &service{repository: repository}
}

func (service *service) Get() (entity.Response, error) {
	return service.repository.Get()
}

func (service *service) GetByID(id string) (entity.Element, error) {
	return service.repository.GetByID(id)
}

func (service *service) Update(id string, resource *entity.Input) (entity.Element, error) {
	child, err := service.repository.GetByID(id)
	if err != nil {
		return entity.Element{}, err
	}

	if resource.Data.Relationships != nil {
		_, err = service.repository.GetByID(resource.Data.Relationships.Parent.Data.ID)
		if err != nil {
			return entity.Element{}, err
		}
	}

	return service.repository.Update(child, resource)
}

func (service *service) Create(resource *entity.Input) (entity.Element, error) {
	if resource.Data.Relationships != nil {
		_, err := service.repository.GetByID(resource.Data.Relationships.Parent.Data.ID)
		if err != nil {
			return entity.Element{}, err
		}
	}
	return service.repository.Create(resource)
}

func (service *service) Delete(id string) error {
	return service.repository.Delete(id)
}
