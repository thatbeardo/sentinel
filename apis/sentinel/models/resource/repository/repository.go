package repository

import (
	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	resource "github.com/bithippie/guard-my-app/apis/sentinel/models/resource/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/session"
)

// Repository is used by the service to communicate with the underlying database
type Repository interface {
	Get() (resource.Output, error)
	GetByID(string) (resource.OutputDetails, error)
	Create(*resource.Input) (resource.OutputDetails, error)
	Update(resource.Details, *resource.Input) (resource.OutputDetails, error)
	Delete(string) error
}

type repository struct {
	session session.Session
}

// Get retrieves all the resources present in the graph
func (repo *repository) Get() (resource.Output, error) {
	return repo.session.Execute(`
		MATCH(child:Resource)
		OPTIONAL MATCH (child: Resource)-[:OWNED_BY]->(parent: Resource)
		OPTIONAL MATCH (policy: Policy)-[:GRANTED_TO]->(child: Resource)
		RETURN {child: child, parent: parent, policy: COLLECT(policy)}`,
		map[string]interface{}{})
}

// GetByID function adds a resource node
func (repo *repository) GetByID(id string) (resource.OutputDetails, error) {
	results, err := repo.session.Execute(`
		MATCH(child:Resource)
		WHERE child.id = $id
		OPTIONAL MATCH (child: Resource)-[:OWNED_BY]->(parent: Resource)
		OPTIONAL MATCH (policy: Policy)-[:GRANTED_TO]->(child: Resource)
		RETURN {child: child, parent: parent, policy: COLLECT(policy)}`,
		map[string]interface{}{
			"id": id,
		})
	if len(results.Data) == 0 {
		return resource.OutputDetails{}, models.ErrNotFound
	}
	return resource.OutputDetails{Data: results.Data[0]}, err
}

// Create function adds a node to the graph
func (repo *repository) Create(input *resource.Input) (resource.OutputDetails, error) {
	var parentID string
	if input.Data.Relationships != nil {
		parentID = input.Data.Relationships.Parent.Data.ID
	}
	results, err := repo.session.Execute(`
		CREATE(child:Resource{name:$name, source_id: $source_id, id: randomUUID()})
		WITH child
		OPTIONAL MATCH(parent:Resource{id:$parent_id})
		WITH child,parent
		FOREACH (o IN CASE WHEN parent IS NOT NULL THEN [parent] ELSE [] END | CREATE (child)-[:OWNED_BY]->(parent))
		RETURN {child: child, parent: parent}`,
		map[string]interface{}{
			"name":      input.Data.Attributes.Name,
			"source_id": input.Data.Attributes.SourceID,
			"parent_id": parentID,
		})
	if len(results.Data) == 0 {
		return resource.OutputDetails{}, models.ErrNotFound
	}
	return resource.OutputDetails{Data: results.Data[0]}, err
}

// Update function Edits the contents of a node
func (repo *repository) Update(oldResource resource.Details, newResource *resource.Input) (resource.OutputDetails, error) {
	newParentID := extractParentID(newResource)
	statement := generateUpdateStatement(newParentID)
	results, err := repo.session.Execute(statement,
		map[string]interface{}{
			"child_id":      oldResource.ID,
			"name":          newResource.Data.Attributes.Name,
			"source_id":     newResource.Data.Attributes.SourceID,
			"new_parent_id": newParentID,
		})
	if len(results.Data) == 0 {
		return resource.OutputDetails{}, models.ErrNotFound
	}
	return resource.OutputDetails{Data: results.Data[0]}, err
}

// Delete function deletes a node from the graph
func (repo *repository) Delete(id string) error {
	_, err := repo.session.Execute(`
		MATCH (n:Resource { id: $id }) DETACH DELETE n`,
		map[string]interface{}{
			"id": id,
		})
	return err
}

// New is a factory method to generate repository instances
func New(session session.Session) Repository {
	return &repository{
		session: session,
	}
}

func extractParentID(newResource *resource.Input) string {
	var parentID string
	if newResource.Data.Relationships != nil {
		parentID = newResource.Data.Relationships.Parent.Data.ID
	}
	return parentID
}

func generateUpdateStatement(newParentID string) (statement string) {
	statement = `
		MATCH(child:Resource{id:$child_id})
		SET child.name=$name
		SET child.source_id=$source_id
		WITH child
		OPTIONAL MATCH (child)-[old_relationship:OWNED_BY]->(old_parent:Resource)`

	if newParentID == "" {
		statement += `RETURN {child: child, parent: old_parent}`
	} else {
		statement += `
		DETACH DELETE old_relationship
		WITH child

		OPTIONAL MATCH (new_parent:Resource{id:$new_parent_id})
		CREATE(child)-[:OWNED_BY]->(new_parent)
		RETURN {child: child, parent: new_parent}`
	}
	return
}
