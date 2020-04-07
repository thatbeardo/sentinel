package resource

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type neo4jRepository struct {
	session neo4j.Session
}

// NewNeo4jRepository is a factory method to create a neo4j repository
func NewNeo4jRepository(session neo4j.Session) Repository {
	return &neo4jRepository{session}
}

var resourceArray = []*Resource{&Resource{Name: "Harshil", SourceID: "Mavani"}}

// Get function adds a resource node
func (repo *neo4jRepository) Get() (Response, error) {
	result, err := repo.session.Run("MATCH(n:Resource) RETURN n.name, n.source_id, n.id", map[string]interface{}{})
	if err != nil {
		return Response{}, err
	}

	var resources []Resource = []Resource{}
	var ids []string
	for result.Next() {
		resourceName := fmt.Sprint(result.Record().GetByIndex(0))
		resourceSourceID := fmt.Sprint(result.Record().GetByIndex(1))
		resources = append(resources, Resource{Name: resourceName, SourceID: resourceSourceID})
		ids = append(ids, fmt.Sprint(result.Record().GetByIndex(2)))
	}
	return constructResourceResponse(resources, ids), nil
}

func (repo *neo4jRepository) Create(resource *Input) (Response, error) {
	result, err := repo.session.Run("CREATE (n:Resource { name: $name, source_id: $source_id, id: randomUUID() }) RETURN n.id",
		map[string]interface{}{
			"name":      resource.Data.Attributes.Name,
			"source_id": resource.Data.Attributes.SourceID,
		})
	if err != nil {
		return Response{Data: []Dto{}}, err
	}
	var id string
	for result.Next() {
		id = fmt.Sprint(result.Record().GetByIndex(0))
	}
	return constructResourceResponse([]Resource{resource.Data.Attributes}, []string{id}), nil
}
