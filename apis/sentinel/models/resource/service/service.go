package service

import (
	resource "github.com/bithippie/guard-my-app/apis/sentinel/models/resource/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/repository"
)

// Service receives commands from handlers and forwards them to the repository
type Service interface {
	Get() (resource.Output, error)
	GetByID(string) (resource.OutputDetails, error)
	Create(*resource.Input) (resource.OutputDetails, error)
	Update(string, *resource.Input) (resource.OutputDetails, error)
	Delete(string) error
}

type service struct {
	repository repository.Repository
}

// NewService creates a service instance with the repository passed
func NewService(repository repository.Repository) Service {
	return &service{repository: repository}
}

func (service *service) Get() (resource.Output, error) {
	return service.repository.Get()
}

func (service *service) GetByID(id string) (resource.OutputDetails, error) {
	return service.repository.GetByID(id)
}

func (service *service) Update(id string, input *resource.Input) (resource.OutputDetails, error) {
	child, err := service.repository.GetByID(id)
	if err != nil {
		return resource.OutputDetails{}, err
	}

	if input.Data.Relationships != nil {
		_, err = service.repository.GetByID(input.Data.Relationships.Parent.Data.ID)
		if err != nil {
			return resource.OutputDetails{}, err
		}
	}

	return service.repository.Update(child.Data, input)
}

func (service *service) Create(input *resource.Input) (resource.OutputDetails, error) {
	if input.Data.Relationships != nil {
		_, err := service.repository.GetByID(input.Data.Relationships.Parent.Data.ID)
		if err != nil {
			return resource.OutputDetails{}, err
		}
	}
	return service.repository.Create(input)
}

func (service *service) Delete(id string) error {
	return service.repository.Delete(id)
}
