package service

import (
	"github.com/bithippie/guard-my-app/sentinel/models/policy/inputs"
	"github.com/bithippie/guard-my-app/sentinel/models/policy/outputs"
	"github.com/bithippie/guard-my-app/sentinel/models/policy/repository"
)

// Service receives commands from handlers and forwards them to the repository
type Service interface {
	Create(*inputs.Payload) (outputs.Policy, error)
}

type service struct {
	repository repository.Repository
}

// NewService creates a service instance with the repository passed
func NewService(repository repository.Repository) Service {
	return &service{repository: repository}
}

func (service *service) Create(payload *inputs.Payload) (outputs.Policy, error) {
	return service.repository.Create(payload)
}
