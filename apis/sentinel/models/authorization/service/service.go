package service

import (
	"sync"

	authorization "github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/repository"
)

// Service interface to expose methods that an authorization endpoint can expect
type Service interface {
	GetAuthorizationForPrincipal(string, authorization.Input) (output authorization.Output, err error)
	IsTargetOwnedByTenant(targetID string, tenantID string) bool
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

var instance service
var once sync.Once

// New is a singleton factory method to return service instances
func New(repository repository.Repository) Service {
	once.Do(func() {
		instance = service{
			repository: repository,
		}
	})
	return instance
}

// NewWithoutRepository is a method to access the already populated instance field
func NewWithoutRepository() Service {
	return instance
}

// SetRepository is a helper method to inject the repository that needs to be assigned
func SetRepository(repository repository.Repository) {
	instance.repository = repository
}
