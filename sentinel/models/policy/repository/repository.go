package repository

import (
	models "github.com/bithippie/guard-my-app/sentinel/models"
	"github.com/bithippie/guard-my-app/sentinel/models/policy/inputs"
	"github.com/bithippie/guard-my-app/sentinel/models/policy/outputs"
	"github.com/bithippie/guard-my-app/sentinel/models/policy/session"
)

// Repository exposes wrapper methods around the underlying session
type Repository interface {
	Create(*inputs.Payload) (outputs.Policy, error)
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

func (repo *repository) Create(input *inputs.Payload) (response outputs.Policy, err error) {
	policy, err := repo.session.Execute(`
		CREATE ( policy:Policy{ name:$name, id:randomUUID() })
		RETURN { policy: policy }`,
		map[string]interface{}{
			"name": input.Data.Attributes.Name,
		})

	if len(policy.Data) == 0 {
		return outputs.Policy{}, models.ErrDatabase
	}
	response = policy.Data[0]
	return
}
