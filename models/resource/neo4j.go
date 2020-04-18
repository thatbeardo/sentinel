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
func (repo *neo4jRepository) Update(oldResource Element, newResource Element) (Element, error) {
	oldParentID, newParentID := extractOldAndNewParentIds(oldResource, newResource)
	result, err := repo.session.Run(`
		MATCH(child:Resource{id:$child_id})
		SET child.name=$name
		SET child.source_id=$source_id
		WITH child
		MATCH(nparent:Resource{id:$new_parent_id})
		WITH child,nparent
		MATCH (child)-[r:OWNED_BY]->(parent:Resource{id:old_parent_id})
		WITH child,parent,r,nparent
		DELETE r
		WITH child,nparent
		CREATE(child)-[:OWNED_BY]->(nparent)
		RETURN child.id`,
		map[string]interface{}{
			"child_id":      oldResource.ID,
			"name":          newResource.Attributes.Name,
			"source_id":     newResource.Attributes.SourceID,
			"new_parent_id": newParentID,
			"old_parent_id": oldParentID,
		})
	if err != nil {
		return Element{}, models.ErrDatabase
	}
	var id string
	for result.Next() {
		id = fmt.Sprint(result.Record().GetByIndex(0))
	}
	parentID := determineUpdatedParent(oldParentID, newParentID)
	return constructResourceResponse(oldResource.Attributes, id, parentID), nil
}

func decodeParentID(response interface{}) string {
	if response == nil {
		return ""
	}
	return fmt.Sprint(response)
}

func determineUpdatedParent(oldID string, newID string) string {
	if oldID == "" && newID == "" {
		return ""
	}
	if oldID == "" && newID != "" {
		return newID
	}
	if oldID != "" && newID == "" {
		return oldID
	}
	return newID
}

func extractOldAndNewParentIds(oldResource Element, newResource Element) (string, string) {
	var oldParentID string
	var newParentID string
	if newResource.Relationships != nil {
		newParentID = newResource.Relationships.Parent.Data.ID
	}
	if oldResource.Relationships != nil {
		oldParentID = oldResource.Relationships.Parent.Data.ID
	}
	return oldParentID, newParentID
}
