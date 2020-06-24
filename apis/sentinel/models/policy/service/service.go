package service

import (
	policy "github.com/bithippie/guard-my-app/apis/sentinel/models/policy/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/policy/repository"
)

// Service receives commands from handlers and forwards them to the repository
type Service interface {
	GetByID(string) (policy.OutputDetails, error)
	Update(string, *policy.Input) (policy.OutputDetails, error)
	Delete(string) error
}

type service struct {
	repository repository.Repository
}

// New creates a service instance with the repository passed
func New(repository repository.Repository) Service {
	return &service{repository: repository}
}

func (service *service) GetByID(id string) (policy.OutputDetails, error) {
	return service.repository.GetByID(id)
}

func (service *service) Update(id string, input *policy.Input) (policy policy.OutputDetails, err error) {
	_, err = service.repository.GetByID(id)

	if err != nil {
		return
	}

	return service.repository.Update(id, input)
}

func (service *service) Delete(id string) error {
	return service.repository.Delete(id)
}
