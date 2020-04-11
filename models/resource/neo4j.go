package resource

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/neo4j"
	models "github.com/thatbeardo/go-sentinel/models"
)

// Neo4jSession is used to communicate with the underlying Graph database
type Neo4jSession interface {
	Run(cypher string, params map[string]interface{}, configurers ...func(*neo4j.TransactionConfig)) (neo4j.Result, error)
	Close() error
}

type neo4jRepository struct {
	session Neo4jSession
}

// NewNeo4jRepository is a factory method to create a neo4j repository
func NewNeo4jRepository(session Neo4jSession) Repository {

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
	result, err := repo.session.Run("MATCH(n:Resource) WHERE n.id = $id RETURN n.name, n.source_id", map[string]interface{}{
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
	return response, err
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

// Delete function deletes a node from the graph
func (repo *neo4jRepository) Delete(id string) error {
	_, err := repo.session.Run(`MATCH(n:Resource) WHERE n.id = $id DETACH DELETE n`,
		map[string]interface{}{
			"id": id,
		})
	if err != nil {
		return models.ErrDatabase
	}
	return nil
}
