package repository

import (
	models "github.com/thatbeardo/go-sentinel/models"
	entity "github.com/thatbeardo/go-sentinel/models/resource"
	"github.com/thatbeardo/go-sentinel/models/resource/session"
)

// Repository is used by the service to communicate with the underlying database
type Repository interface {
	Get() (entity.Response, error)
	GetByID(string) (entity.Element, error)
	Create(*entity.Input) (entity.Element, error)
	Update(entity.Element, *entity.Input) (entity.Element, error)
	Delete(string) error
}

type repository struct {
	session session.Session
}

// Get retrieves all the resources present in the graph
func (repo *repository) Get() (entity.Response, error) {
	return repo.session.Execute(`
		MATCH(child:Resource)
		OPTIONAL MATCH (child: Resource)-[:OWNED_BY]->(parent: Resource)
		RETURN {child: child, parent: parent}`,
		map[string]interface{}{})
}

// GetByID function adds a resource node
func (repo *repository) GetByID(id string) (entity.Element, error) {
	elements, err := repo.session.Execute(`
		MATCH(child:Resource)
		WHERE child.id = $id
		OPTIONAL MATCH (child: Resource)-[:OWNED_BY]->(parent: Resource)
		RETURN {child: child, parent: parent}`,
		map[string]interface{}{
			"id": id,
		})
	if len(elements.Data) == 0 {
		return entity.Element{}, models.ErrNotFound
	}
	return elements.Data[0], err
}

// Create function adds a node to the graph
func (repo *repository) Create(resource *entity.Input) (entity.Element, error) {
	var parentID string
	if resource.Data.Relationships != nil {
		parentID = resource.Data.Relationships.Parent.Data.ID
	}
	elements, err := repo.session.Execute(`
		CREATE(child:Resource{name:$name, source_id: $source_id, id: randomUUID()})
		WITH child
		OPTIONAL MATCH(parent:Resource{id:$parent_id})
		WITH child,parent
		FOREACH (o IN CASE WHEN parent IS NOT NULL THEN [parent] ELSE [] END | CREATE (child)-[:OWNED_BY]->(parent))
		RETURN {child: child, parent: parent}`,
		map[string]interface{}{
			"name":      resource.Data.Attributes.Name,
			"source_id": resource.Data.Attributes.SourceID,
			"parent_id": parentID,
		})
	if len(elements.Data) == 0 {
		return entity.Element{}, models.ErrNotFound
	}
	return elements.Data[0], err
}

// Update function Edits the contents of a node
func (repo *repository) Update(oldResource entity.Element, newResource *entity.Input) (entity.Element, error) {
	newParentID := extractParentID(newResource)
	statement := generateUpdateStatement(newParentID)
	elements, err := repo.session.Execute(statement,
		map[string]interface{}{
			"child_id":      oldResource.ID,
			"name":          newResource.Data.Attributes.Name,
			"source_id":     newResource.Data.Attributes.SourceID,
			"new_parent_id": newParentID,
		})
	if len(elements.Data) == 0 {
		return entity.Element{}, models.ErrNotFound
	}
	return elements.Data[0], err
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

func extractParentID(newResource *entity.Input) string {
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
