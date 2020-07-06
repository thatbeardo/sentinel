package service

import (
	permission "github.com/bithippie/guard-my-app/apis/sentinel/models/permission/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/permission/repository"
)

// Service receives commands from handlers and forwards them to the repository
type Service interface {
	GetAllPermissionsForcontext(string) (permission.Output, error)
	GetAllPermissionsForcontextWithResource(string, string) (permission.Output, error)
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

func (service *service) Create(input *permission.Input, contextID string, targetID string) (permission.OutputDetails, error) {
	// TODO:
	// Validate presence of both context and resource before calling repository create method
	return service.repository.Create(input, contextID, targetID)
}

func (service *service) GetAllPermissionsForcontext(contextID string) (permission.Output, error) {
	return service.repository.GetAllPermissionsForcontext(contextID)
}

func (service *service) GetAllPermissionsForcontextWithResource(contextID string, resourceID string) (permission.Output, error) {
	return service.repository.GetAllPermissionsForcontextWithResource(contextID, resourceID)
}

func (service *service) Update(id string, input *permission.Input) (permission.OutputDetails, error) {
	// TODO:
	// Validate presence of context before calling repository
	return service.repository.Update(id, input)
}

func (service *service) Delete(id string) error {
	// TODO:
	// Validate presence of context before calling repository
	return service.repository.Delete(id)
}
