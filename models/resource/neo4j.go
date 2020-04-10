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
	result, err := repo.session.Run("MATCH(n:Resource) RETURN n.name, n.source_id, n.id", map[string]interface{}{})
	if err != nil {
		return Response{}, models.ErrDatabase
	}
	var dtos []Element = []Element{}
	for result.Next() {
		resourceName := fmt.Sprint(result.Record().GetByIndex(0))
		resourceSourceID := fmt.Sprint(result.Record().GetByIndex(1))
		id := fmt.Sprint(result.Record().GetByIndex(2))
		dtos = append(dtos, constructResourceResponse(Resource{Name: resourceName, SourceID: resourceSourceID}, id))
	}
	return Response{Data: dtos}, nil
}

// GetByID function adds a resource node
func (repo *neo4jRepository) GetByID(id string) (Element, error) {
	result, err := repo.session.Run("MATCH(n:Resource) WHERE n.id = $id", map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return Element{}, models.ErrDatabase
	}
	var response Element
	for result.Next() {
		resourceName := fmt.Sprint(result.Record().GetByIndex(0))
		resourceSourceID := fmt.Sprint(result.Record().GetByIndex(1))
		response = constructResourceResponse(Resource{Name: resourceName, SourceID: resourceSourceID}, id)
	}
	if response.ID == "" {
		err = models.ErrNotFound
	}
	return response, nil
}

// Create function adds a node to the graph
func (repo *neo4jRepository) Create(resource *Input) (Element, error) {
	result, err := repo.session.Run("CREATE (n:Resource { name: $name, source_id: $source_id, id: randomUUID() }) RETURN n.id",
		map[string]interface{}{
			"name":      resource.Data.Attributes.Name,
			"source_id": resource.Data.Attributes.SourceID,
		})
	if err != nil {
		return Element{}, models.ErrDatabase
	}
	var id string
	for result.Next() {
		id = fmt.Sprint(result.Record().GetByIndex(0))
	}
	return constructResourceResponse(resource.Data.Attributes, id), nil
}
