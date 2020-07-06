package service

import (
	context "github.com/bithippie/guard-my-app/apis/sentinel/models/context/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/context/repository"
)

// Service receives commands from handlers and forwards them to the repository
type Service interface {
	GetByID(string) (context.OutputDetails, error)
	Update(string, *context.Input) (context.OutputDetails, error)
	Delete(string) error
}

type service struct {
	repository repository.Repository
}

// New creates a service instance with the repository passed
func New(repository repository.Repository) Service {
	return &service{repository: repository}
}

func (service *service) GetByID(id string) (context.OutputDetails, error) {
	return service.repository.GetByID(id)
}

func (service *service) Update(id string, input *context.Input) (context context.OutputDetails, err error) {
	_, err = service.repository.GetByID(id)

	if err != nil {
		return
	}

	return service.repository.Update(id, input)
}

func (service *service) Delete(id string) error {
	return service.repository.Delete(id)
}
