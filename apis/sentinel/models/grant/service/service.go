package service

import (
	grant "github.com/bithippie/guard-my-app/apis/sentinel/models/grant/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/grant/repository"
)

// Service receives commands from handlers and forwards them to the repository
type Service interface {
	Create(*grant.Input, string, string) (grant.OutputDetails, error)
	GetPrincipalAndPolicyForResource(string) (grant.Output, error)
}

type service struct {
	repository repository.Repository
}

// New creates a service instance with the repository passed
func New(repository repository.Repository) Service {
	return &service{repository: repository}
}

func (service *service) Create(input *grant.Input, policyID string, targetID string) (grant.OutputDetails, error) {
	// TODO:
	// Validate presence of both policy and resource before calling repository create method
	return service.repository.Create(input, policyID, targetID)
}

func (service *service) GetPrincipalAndPolicyForResource(id string) (grant.Output, error) {
	return service.repository.GetPrincipalAndPolicyForResource(id)
}
