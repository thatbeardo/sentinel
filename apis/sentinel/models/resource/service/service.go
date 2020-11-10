package service

import (
	"context"
	"strings"

	"fmt"

	dto "github.com/bithippie/guard-my-app/apis/sentinel/models/context/dto"
	resource "github.com/bithippie/guard-my-app/apis/sentinel/models/resource/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/injection"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/repository"
)

// Service receives commands from handlers and forwards them to the repository
type Service interface {
	Get(context.Context) (resource.Output, error)
	GetByID(string) (resource.OutputDetails, error)

	Create(context.Context, *resource.Input) (resource.OutputDetails, error)

	Associatecontext(string, *dto.Input) (dto.OutputDetails, error)
	GetAllAssociatedContexts(string) (dto.Output, error)

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

func (service *service) Get(ctx context.Context) (resource.Output, error) {
	return service.repository.Get(injection.ExtractClaims(ctx, "azp"), injection.ExtractTenant(ctx))
}

func (service *service) GetByID(id string) (resource.OutputDetails, error) {
	return service.repository.GetByID(id)
}

func (service *service) Update(id string, input *resource.Input) (resource.OutputDetails, error) {
	child, err := service.repository.GetByID(id)
	if err != nil {
		return resource.OutputDetails{}, err
	}
	return service.repository.Update(child.Data, input)
}

func (service *service) Associatecontext(principalID string, input *dto.Input) (dto.OutputDetails, error) {
	return service.repository.AddContext(principalID, input)
}

func (service *service) GetAllAssociatedContexts(id string) (dto.Output, error) {
	return service.repository.GetAllContexts(id)
}

func (service *service) Create(ctx context.Context, input *resource.Input) (resource.OutputDetails, error) {
	scope := injection.ExtractClaims(ctx, "scope")
	fmt.Println("Scope" + scope)
	if strings.Contains(scope, "create:resource") {
		return service.repository.CreateMetaResource(input)
	}

	azp := injection.ExtractClaims(ctx, "azp")
	if input.Data.Relationships == nil {
		parent, _ := service.repository.GetResourcesHub(azp, injection.ExtractTenant(ctx))
		input.Data.Relationships = generateParentRelationship(parent.Data.ID)
	}

	return service.repository.Create(input)
}

func (service *service) Delete(id string) error {
	return service.repository.Delete(id)
}

func generateParentRelationship(id string) *resource.InputRelationships {
	return &resource.InputRelationships{
		Parent: &resource.Parent{
			Data: &resource.Data{
				Type: "resource",
				ID:   id,
			},
		},
	}
}
