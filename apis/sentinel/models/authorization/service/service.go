package service

import (
	authorization "github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/repository"
)

// Service interface to expose methods that an authorization endpoint can expect
type Service interface {
	GetAuthorizationForPrincipal(string, authorization.Input) (output authorization.Output, err error)
	IsTargetOwnedByTenant(string, string) bool
	IsPolicyOwnedByTenant(string, string) bool
	IsPermissionOwnedByTenant(string, string) bool
}

type service struct {
	repository repository.Repository
}

func (s service) GetAuthorizationForPrincipal(principalID string, input authorization.Input) (output authorization.Output, err error) {
	return s.repository.GetAuthorizationForPrincipal(principalID, input)
}

func (s service) IsTargetOwnedByTenant(targetID string, tenantID string) bool {
	return s.repository.IsTargetOwnedByTenant(targetID, tenantID)
}

func (s service) IsPolicyOwnedByTenant(policyID string, tenantID string) bool {
	return s.repository.IsPolicyOwnedByTenant(policyID, tenantID)
}

func (s service) IsPermissionOwnedByTenant(permissionID, tenantID string) bool {
	return s.repository.IsPermissionOwnedByTenant(permissionID, tenantID)
}

// New is a singleton factory method to return service instances
func New(repository repository.Repository) Service {
	return service{
		repository: repository,
	}
}
