package repository

import (
	models "github.com/bithippie/guard-my-app/sentinel/models"
	"github.com/bithippie/guard-my-app/sentinel/models/permission/inputs"
	"github.com/bithippie/guard-my-app/sentinel/models/permission/outputs"
	"github.com/bithippie/guard-my-app/sentinel/models/permission/session"
)

// Repository exposes wrapper methods around the underlying session
type Repository interface {
	Create(input *inputs.Payload, policyID string, targetID string) (outputs.Permission, error)
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

func (repo *repository) Create(input *inputs.Payload, policyID string, targetID string) (response outputs.Permission, err error) {
	permission, err := repo.session.Execute(`
		MATCH (policy: Policy), (target: Resource)
		WHERE policy.id = $policyID AND target.ID = $targetID
		CREATE (policy)-[r:PERMISSION {name: $name, permitted: $permitted}]->(target)
		RETURN {permission: r}`,
		map[string]interface{}{
			"name":      input.Data.Attributes.Name,
			"permitted": input.Data.Attributes.Permitted,
			"policyID":  policyID,
			"targetID":  targetID,
		})

	if len(permission.Data) == 0 {
		return outputs.Permission{}, models.ErrDatabase
	}
	response = permission.Data[0]
	return
}
