package repository

import (
	models "github.com/bithippie/guard-my-app/sentinel/models"
	"github.com/bithippie/guard-my-app/sentinel/models/grant/inputs"
	"github.com/bithippie/guard-my-app/sentinel/models/grant/outputs"
	"github.com/bithippie/guard-my-app/sentinel/models/grant/session"
)

// Repository exposes wrapper methods around the underlying session
type Repository interface {
	Create(input *inputs.Payload, policyID string, targetID string) (outputs.Grant, error)
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

func (repo *repository) Create(input *inputs.Payload, policyID string, targetID string) (response outputs.Grant, err error) {
	permission, err := repo.session.Execute(`
		MATCH (policy: Policy), (target: Resource)
		WHERE policy.id = $policyID AND target.id = $targetID
		CREATE (policy)-[r:GRANT {with_grant: $withGrant, id: randomUUID()}]->(target)
		RETURN {grant: r}`,
		map[string]interface{}{
			"withGrant": input.Data.Attributes.WithGrant,
			"policyID":  policyID,
			"targetID":  targetID,
		})

	if len(permission.Data) == 0 {
		return outputs.Grant{}, models.ErrDatabase
	}
	response = permission.Data[0]
	return
}
