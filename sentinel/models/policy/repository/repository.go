package repository

import (
	models "github.com/bithippie/guard-my-app/sentinel/models"
	policy "github.com/bithippie/guard-my-app/sentinel/models/policy/dto"
	"github.com/bithippie/guard-my-app/sentinel/models/policy/session"
)

// Repository exposes wrapper methods around the underlying session
type Repository interface {
	Get() (policy.Output, error)
	GetByID(string) (policy.OutputDetails, error)
	Create(*policy.Input) (policy.OutputDetails, error)
	Update(string, *policy.Input) (policy.OutputDetails, error)
	Delete(string) error
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

func (repo *repository) Get() (output policy.Output, err error) {
	return repo.session.Execute(`
		MATCH (policy:Policy) 
		OPTIONAL MATCH(policy)-[:GRANTED_TO]->(principal:Resource)
		OPTIONAL MATCH(policy)-[:PERMISSION]->(target:Resource)
		RETURN {policy:policy, principals:COLLECT(principal), targets:COLLECT(target)}`,
		map[string]interface{}{})
}

func (repo *repository) GetByID(id string) (output policy.OutputDetails, err error) {
	result, err := repo.session.Execute(`
		MATCH (policy:Policy) 
		WHERE policy.id = $id
		OPTIONAL MATCH(policy)-[:GRANTED_TO]->(principal:Resource)
		OPTIONAL MATCH(policy)-[:PERMISSION]->(target:Resource)
		RETURN {policy:policy, principals:COLLECT(principal), targets:COLLECT(target)}`,
		map[string]interface{}{
			"id": id,
		})
	if len(result.Data) == 0 {
		return output, models.ErrNotFound
	}
	output = policy.OutputDetails{Data: result.Data[0]}
	return
}

func (repo *repository) Create(input *policy.Input) (output policy.OutputDetails, err error) {
	result, err := repo.session.Execute(`
		CREATE ( policy:Policy{ name:$name, id:randomUUID() })
		RETURN { policy: policy }`,
		map[string]interface{}{
			"name": input.Data.Attributes.Name,
		})

	if len(result.Data) == 0 {
		return policy.OutputDetails{}, models.ErrDatabase
	}
	output = policy.OutputDetails{Data: result.Data[0]}
	return
}

func (repo *repository) Update(id string, input *policy.Input) (output policy.OutputDetails, err error) {
	result, err := repo.session.Execute(`
		MATCH (policy:Policy)
		WHERE policy.id = $id
		SET policy.name = $name
		RETURN { policy: policy }`,
		map[string]interface{}{
			"id":   id,
			"name": input.Data.Attributes.Name,
		})

	if len(result.Data) == 0 {
		return policy.OutputDetails{}, models.ErrDatabase
	}
	output = policy.OutputDetails{Data: result.Data[0]}
	return
}

func (repo *repository) Delete(id string) (err error) {
	_, err = repo.session.Execute(`
		MATCH (policy:Policy)
		WHERE policy.id = $id
		DETACH DELETE policy`,
		map[string]interface{}{
			"id": id,
		})

	return
}
