package service

import (
	grant "github.com/bithippie/guard-my-app/apis/sentinel/models/grant/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/grant/repository"
)

// Service receives commands from handlers and forwards them to the repository
type Service interface {
	Create(*grant.Input, string, string) (grant.OutputDetails, error)
	GetPrincipalAndcontextForResource(string) (grant.Output, error)
	GrantExists(contextID, principalID string) (bool, error)
}

type service struct {
	repository repository.Repository
}

// New creates a service instance with the repository passed
func New(repository repository.Repository) Service {
	return &service{repository: repository}
}

func (service *service) Create(input *grant.Input, contextID string, targetID string) (grant.OutputDetails, error) {
	// TODO:
	// Validate presence of both context and resource before calling repository create method
	return service.repository.Create(input, contextID, targetID)
}

func (service *service) GetPrincipalAndcontextForResource(id string) (grant.Output, error) {
	return service.repository.GetPrincipalAndcontextForResource(id)
}

func (service *service) GrantExists(contextID, principalID string) (bool, error) {
	return service.repository.GrantExists(contextID, principalID)
}
