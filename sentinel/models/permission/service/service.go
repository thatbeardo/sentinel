package service

import (
	permission "github.com/bithippie/guard-my-app/sentinel/models/permission/dto"
	"github.com/bithippie/guard-my-app/sentinel/models/permission/repository"
)

// Service receives commands from handlers and forwards them to the repository
type Service interface {
	GetAllPermissionsForPolicy(string) (permission.Output, error)
	GetAllPermissionsForPolicyWithResource(string, string) (permission.Output, error)
	Create(*permission.Input, string, string) (permission.OutputDetails, error)
	Update(string, *permission.Input) (permission.OutputDetails, error)
	Delete(string) error
}

type service struct {
	repository repository.Repository
}

// New creates a service instance with the repository passed
func New(repository repository.Repository) Service {
	return &service{repository: repository}
}

func (service *service) Create(input *permission.Input, policyID string, targetID string) (permission.OutputDetails, error) {
	// TODO:
	// Validate presence of both policy and resource before calling repository create method
	return service.repository.Create(input, policyID, targetID)
}

func (service *service) GetAllPermissionsForPolicy(policyID string) (permission.Output, error) {
	return service.repository.GetAllPermissionsForPolicy(policyID)
}

func (service *service) GetAllPermissionsForPolicyWithResource(policyID string, resourceID string) (permission.Output, error) {
	return service.repository.GetAllPermissionsForPolicyWithResource(policyID, resourceID)
}

func (service *service) Update(id string, input *permission.Input) (permission.OutputDetails, error) {
	// TODO:
	// Validate presence of policy before calling repository
	return service.repository.Update(id, input)
}

func (service *service) Delete(id string) error {
	// TODO:
	// Validate presence of policy before calling repository
	return service.repository.Delete(id)
}
