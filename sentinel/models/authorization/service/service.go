package service

import (
	authorization "github.com/bithippie/guard-my-app/sentinel/models/authorization/dto"
	"github.com/bithippie/guard-my-app/sentinel/models/authorization/repository"
)

// Service interface to expose methods that an authorization endpoint can expect
type Service interface {
	GetAuthorizationForPrincipal(string, authorization.Input) (output authorization.Output, err error)
}

type service struct {
	repository repository.Repository
}

func (s service) GetAuthorizationForPrincipal(principalID string, input authorization.Input) (output authorization.Output, err error) {
	return s.repository.GetAuthorizationForPrincipal(principalID, input)
}

// New is a factory method to return service instances
func New(repository repository.Repository) Service {
	return service{
		repository: repository,
	}
}
