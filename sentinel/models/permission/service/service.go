package service

import (
	"github.com/bithippie/guard-my-app/sentinel/models/permission/inputs"
	"github.com/bithippie/guard-my-app/sentinel/models/permission/outputs"
	"github.com/bithippie/guard-my-app/sentinel/models/permission/repository"
)

// Service receives commands from handlers and forwards them to the repository
type Service interface {
	Create(*inputs.Payload, string, string) (outputs.Permission, error)
}

type service struct {
	repository repository.Repository
}

// NewService creates a service instance with the repository passed
func NewService(repository repository.Repository) Service {
	return &service{repository: repository}
}

func (service *service) Create(payload *inputs.Payload, policyID string, targetID string) (outputs.Permission, error) {
	// TODO:
	// Validate presence of both policy and resource before calling repository create method
	return service.repository.Create(payload, policyID, targetID)
}
