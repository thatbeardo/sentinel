package repository

import (
	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	grant "github.com/bithippie/guard-my-app/apis/sentinel/models/grant/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/grant/session"
)

// Repository exposes wrapper methods around the underlying session
type Repository interface {
	Create(*grant.Input, string, string) (grant.OutputDetails, error)
	GetPrincipalAndcontextForResource(string) (grant.Output, error)
	GrantExists(contextID, principalID string) (bool, error)
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

func (repo *repository) GetPrincipalAndcontextForResource(principal string) (output grant.Output, err error) {
	return repo.session.Execute(`
		MATCH (context:Context)-[:PERMISSION]->(principal: Resource {id: $principalID})
		OPTIONAL MATCH (context)-[grant:GRANTED_TO]->(principal:Resource)
		RETURN {grant: grant, context:context, principal:principal}`,
		map[string]interface{}{
			"principalID": principal,
		})
}

func (repo *repository) Create(input *grant.Input, contextID, principalID string) (output grant.OutputDetails, err error) {
	result, err := repo.session.Execute(`
		MATCH (context:Context), (principal: Resource)
		WHERE context.id = $contextID AND principal.id = $principalID
		CREATE (context)-[grant:GRANTED_TO {with_grant: $withGrant, id: randomUUID()}]->(principal)
		RETURN {grant: grant, context:context, principal: principal}`,
		map[string]interface{}{
			"withGrant":   input.Data.Attributes.WithGrant,
			"contextID":   contextID,
			"principalID": principalID,
		})

	if len(result.Data) == 0 {
		return grant.OutputDetails{}, models.ErrNotFound
	}
	output = grant.OutputDetails{Data: result.Data[0]}
	return
}

func (repo *repository) GrantExists(contextID, principalID string) (bool, error) {
	result, err := repo.session.Execute(`
		MATCH (context:Context{id: $contextID})-[grant:GRANTED_TO]->(principal: Resource {id: $principalID})
		RETURN {grant: grant}`,
		map[string]interface{}{
			"contextID":   contextID,
			"principalID": principalID,
		})

	if err != nil {
		return false, err
	}

	return len(result.Data) > 0, err
}
