package service

import (
	"context"

	authorization "github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/repository"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/injection"
)

// Service interface to expose methods that an authorization endpoint can expect
type Service interface {
	GetAuthorizationForPrincipal(principalID, contextID string, input authorization.Input) (output authorization.Output, err error)
	IsTargetOwnedByClient(ctx context.Context, targetID string) bool
	IsContextOwnedByClient(ctx context.Context, contextID string) bool
	IsPermissionOwnedByTenant(ctx context.Context, permissionID string) bool
}

type service struct {
	repository repository.Repository
}

func (s service) GetAuthorizationForPrincipal(principalID, contextID string, input authorization.Input) (output authorization.Output, err error) {
	return s.repository.GetAuthorizationForPrincipal(principalID, contextID, input)
}

func (s service) IsTargetOwnedByClient(ctx context.Context, targetID string) bool {
	return s.repository.IsTargetOwnedByClient(injection.ExtractClaims(ctx, "azp"), injection.ExtractTenant(ctx), targetID)
}

func (s service) IsContextOwnedByClient(ctx context.Context, contextID string) bool {
	return s.repository.IsContextOwnedByClient(injection.ExtractClaims(ctx, "azp"), injection.ExtractTenant(ctx), contextID)
}

func (s service) IsPermissionOwnedByTenant(ctx context.Context, permissionID string) bool {
	return s.repository.IsPermissionOwnedByTenant(injection.ExtractClaims(ctx, "azp"), injection.ExtractTenant(ctx), permissionID)
}

// New is a singleton factory method to return service instances
func New(repository repository.Repository) Service {
	return service{
		repository: repository,
	}
}
