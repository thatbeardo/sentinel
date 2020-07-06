package repository

import (
	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	context "github.com/bithippie/guard-my-app/apis/sentinel/models/context/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/context/session"
)

// Repository exposes wrapper methods around the underlying session
type Repository interface {
	GetByID(string) (context.OutputDetails, error)
	Update(string, *context.Input) (context.OutputDetails, error)
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

func (repo *repository) GetByID(id string) (output context.OutputDetails, err error) {
	result, err := repo.session.Execute(`
		MATCH (context:Context) 
		WHERE context.id = $id
		OPTIONAL MATCH(context)-[:GRANTED_TO]->(principal:Resource)
		OPTIONAL MATCH(context)-[:PERMISSION]->(target:Resource)
		RETURN {context:Context, principals:COLLECT(principal), targets:COLLECT(target)}`,
		map[string]interface{}{
			"id": id,
		})
	if len(result.Data) == 0 {
		return output, models.ErrNotFound
	}
	output = context.OutputDetails{Data: result.Data[0]}
	return
}

func (repo *repository) Update(id string, input *context.Input) (output context.OutputDetails, err error) {
	result, err := repo.session.Execute(`
		MATCH (context:Context)
		WHERE context.id = $id
		SET context.name = $name
		RETURN { context:Context }`,
		map[string]interface{}{
			"id":   id,
			"name": input.Data.Attributes.Name,
		})

	if len(result.Data) == 0 {
		return context.OutputDetails{}, models.ErrDatabase
	}
	output = context.OutputDetails{Data: result.Data[0]}
	return
}

func (repo *repository) Delete(id string) (err error) {
	_, err = repo.session.Execute(`
		MATCH (context:Context)
		WHERE context.id = $id
		DETACH DELETE context`,
		map[string]interface{}{
			"id": id,
		})

	return
}
