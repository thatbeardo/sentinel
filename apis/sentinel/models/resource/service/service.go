package service

import (
	"context"

	policy "github.com/bithippie/guard-my-app/apis/sentinel/models/policy/dto"
	resource "github.com/bithippie/guard-my-app/apis/sentinel/models/resource/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/repository"
)

// Service receives commands from handlers and forwards them to the repository
type Service interface {
	Get() (resource.Output, error)
	GetByID(string) (resource.OutputDetails, error)
	Create(context.Context, *resource.Input) (resource.OutputDetails, error)
	AssociatePolicy(string, *policy.Input) (policy.OutputDetails, error)
	GetAllAssociatedPolicies(string) (policy.Output, error)
	Update(string, *resource.Input) (resource.OutputDetails, error)
	Delete(string) error
}

type service struct {
	repository repository.Repository
}

// New creates a service instance with the repository passed
func New(repository repository.Repository) Service {
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

func (service *service) AssociatePolicy(principalID string, input *policy.Input) (policy.OutputDetails, error) {
	return service.repository.AssociatePolicy(principalID, input)
}

func (service *service) GetAllAssociatedPolicies(id string) (policy.Output, error) {
	return service.repository.GetAllAssociatedPolicies(id)
}

func (service *service) Create(ctx context.Context, input *resource.Input) (resource.OutputDetails, error) {
	return service.repository.Create(ctx, input)
}

func (service *service) Delete(id string) error {
	return service.repository.Delete(id)
}
