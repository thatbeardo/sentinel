package repository

import (
	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	grant "github.com/bithippie/guard-my-app/apis/sentinel/models/grant/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/grant/session"
)

// Repository exposes wrapper methods around the underlying session
type Repository interface {
	Create(*grant.Input, string, string) (grant.OutputDetails, error)
	GetPrincipalAndPolicyForResource(string) (grant.Output, error)
}

type repository struct {
	session session.Session
}

// New is a factory method to generate repository instances
func New(session session.Session) Repository {
	return &repository{
		session: session,
	}
}

func (repo *repository) GetPrincipalAndPolicyForResource(principal string) (output grant.Output, err error) {
	return repo.session.Execute(`
		MATCH (policy:Policy)-[:PERMISSION]->(principal: Resource {id: $principalID})
		OPTIONAL MATCH (policy)-[grant:GRANTED_TO]->(principal:Resource)
		RETURN {grant: grant, policy:policy, principal:principal}`,
		map[string]interface{}{
			"principalID": principal,
		})
}

func (repo *repository) Create(input *grant.Input, policyID string, principalID string) (output grant.OutputDetails, err error) {
	result, err := repo.session.Execute(`
		MATCH (policy: Policy), (principal: Resource)
		WHERE policy.id = $policyID AND principal.id = $principalID
		CREATE (policy)-[grant:GRANTED_TO {with_grant: $withGrant, id: randomUUID()}]->(principal)
		RETURN {grant: grant, policy: policy, principal: principal}`,
		map[string]interface{}{
			"withGrant":   input.Data.Attributes.WithGrant,
			"policyID":    policyID,
			"principalID": principalID,
		})

	if len(result.Data) == 0 {
		return grant.OutputDetails{}, models.ErrNotFound
	}
	output = grant.OutputDetails{Data: result.Data[0]}
	return
}
