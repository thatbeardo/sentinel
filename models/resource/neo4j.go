package resource

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/neo4j"
	models "github.com/thatbeardo/go-sentinel/models"
)

type neo4jRepository struct {
	session neo4j.Session
}

// NewNeo4jRepository is a factory method to create a neo4j repository
func NewNeo4jRepository(session neo4j.Session) Repository {
	return &neo4jRepository{session}
}

// Get retrieves all the resources present in the graph
func (repo *neo4jRepository) Get() (Response, error) {
	result, err := repo.session.Run("MATCH(child:Resource) OPTIONAL MATCH (child: Resource)-[:OWNED_BY]->(parent: Resource) RETURN child.name, child.source_id, child.id, parent.id", map[string]interface{}{})
	if err != nil {
		return Response{}, models.ErrDatabase
	}
	var dtos []Element = []Element{}
	for result.Next() {
		resourceName := fmt.Sprint(result.Record().GetByIndex(0))
		resourceSourceID := fmt.Sprint(result.Record().GetByIndex(1))
		id := fmt.Sprint(result.Record().GetByIndex(2))
		parentID := decodeParentID(result.Record().GetByIndex(3))
		dtos = append(dtos, constructResourceResponse(&Resource{Name: resourceName, SourceID: resourceSourceID}, id, parentID))
	}
	return Response{Data: dtos}, nil
}

// GetByID function adds a resource node
func (repo *neo4jRepository) GetByID(id string) (Element, error) {
	result, err := repo.session.Run("MATCH(child:Resource) WHERE child.id = $id OPTIONAL MATCH (child: Resource)-[:OWNED_BY]->(parent: Resource) RETURN child.name, child.source_id, parent.id", map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return Element{}, models.ErrDatabase
	}
	var response Element
	for result.Next() {
		resourceName := fmt.Sprint(result.Record().GetByIndex(0))
		resourceSourceID := fmt.Sprint(result.Record().GetByIndex(1))
		parentID := decodeParentID(result.Record().GetByIndex(2))
		response = constructResourceResponse(&Resource{Name: resourceName, SourceID: resourceSourceID}, id, parentID)
	}
	if response.ID == "" {
		err = models.ErrNotFound
	}
	return response, err
}

// Create function adds a node to the graph
func (repo *neo4jRepository) Create(resource *Input) (Element, error) {
	var parentID string
	if resource.Data.Relationships != nil {
		parentID = resource.Data.Relationships.Parent.Data.ID
	}
	result, err := repo.session.Run(`
	CREATE(child:Resource{name:$name, source_id: $source_id})
	WITH child
	MATCH(parent:Resource{id:$parent_id})
	WITH child,parent
	CREATE(child)-[r:OWNED_BY]->(parent)
	return child`,
		map[string]interface{}{
			"name":      resource.Data.Attributes.Name,
			"source_id": resource.Data.Attributes.SourceID,
			"parent_id": parentID,
		})
	if err != nil {
		return Element{}, models.ErrDatabase
	}
	var id string
	for result.Next() {
		id = fmt.Sprint(result.Record().GetByIndex(0))
	}
	return constructResourceResponse(resource.Data.Attributes, id, ""), nil
}

// Delete function deletes a node from the graph
func (repo *neo4jRepository) Delete(id string) error {
	result, err := repo.session.Run(`MATCH (n:Resource { id: $id }) DETACH DELETE n`,
		map[string]interface{}{
			"id": id,
		})
	if err != nil {
		return models.ErrDatabase
	}
	result.Next()
	summary, err := result.Summary()
	if err != nil {
		return models.ErrDatabase
	}
	if summary.Counters().NodesDeleted() == 0 {
		return models.ErrNotFound
	}
	return nil
}

// Update function Edits the contents of a node
func (repo *neo4jRepository) Update(id string, resource *Input) (Element, error) {
	result, err := repo.session.Run("MATCH (n:Resource { id: $id }) SET n.name = $name, n.source_id = $source_id RETURN n.name, n.source_id",
		map[string]interface{}{
			"id":        id,
			"name":      resource.Data.Attributes.Name,
			"source_id": resource.Data.Attributes.SourceID,
		})
	if err != nil {
		return Element{}, models.ErrDatabase
	}
	for result.Next() {
		id = fmt.Sprint(result.Record().GetByIndex(0))
	}
	var parentID string = ""
	if resource.Data.Relationships != nil {
		parentID = resource.Data.Relationships.Parent.Data.ID
	}
	return constructResourceResponse(resource.Data.Attributes, id, parentID), nil
}

// UpdateOwnerRelationships is used to delete the old parent and attach a node to a new parent
func (repo *neo4jRepository) UpdateOwnership(id string, resource *Input) (Element, error) {
	result, err := repo.session.Run(`MATCH(child:Resource{id:$id})
	SET child.name = $name
	SET child.source_id = $source_id
	WITH child
	OPTIONAL MATCH(child)-[ownership:OWNED_BY]->(:Resource)
	DELETE ownership
	WITH child
	MATCH (parent:Resource{id: $parent_id})
	MERGE (child)-[r:OWNED_BY]->(parent)`,
		map[string]interface{}{
			"id":        id,
			"name":      resource.Data.Attributes.Name,
			"source_id": resource.Data.Attributes.SourceID,
			"parent_id": resource.Data.Relationships.Parent.Data.ID,
		})
	if err != nil {
		return Element{}, models.ErrDatabase
	}
	result.Next()
	summary, err := result.Summary()

	if err != nil || summary.Counters().RelationshipsCreated() != 1 {
		return Element{}, models.ErrDatabase
	}
	return constructResourceResponse(resource.Data.Attributes, id, resource.Data.Relationships.Parent.Data.ID), nil
}

func decodeParentID(response interface{}) string {
	if response == nil {
		return ""
	}
	return fmt.Sprint(response)
}
