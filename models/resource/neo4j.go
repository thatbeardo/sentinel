package resource

import "github.com/neo4j/neo4j-go-driver/neo4j"

type neo4jRepository struct {
	session neo4j.Session
}

// NewNeo4jRepository is a factory method to create a neo4j repository
func NewNeo4jRepository(session neo4j.Session) Repository {
	return &neo4jRepository{session}
}

var resourceArray = []*Resource{}

// Get function adds a resource node
func (repo *neo4jRepository) Get() ([]*Resource, error) {
	return resourceArray, nil
}

func (repo *neo4jRepository) Create(resource *Resource) (*Resource, error) {
	return nil, nil
}
