package resource

import (
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	models "github.com/thatbeardo/go-sentinel/models"
)

// Session interface defines methods needed to communicate/execute queries and a cleanup function when everything is done
type Session interface {
	Execute(statement string, parameters map[string]interface{}) (response Response, err error)
}

type neo4jSession struct {
	session neo4j.Session
}

type resourceNode struct {
	Name     string `mapstructure:"name"`
	SourceID string `mapstructure:"source_id"`
	ID       string `mapstructure:"id"`
}

// NewNeo4jSession is a factory method to create Neo4J session instances
func NewNeo4jSession(session neo4j.Session) Session {
	return neo4jSession{
		session: session,
	}
}

// Execute runs the statement passed as a query and opulates the data parameter with result
func (n neo4jSession) Execute(statement string, parameters map[string]interface{}) (response Response, err error) {
	result, err := n.session.Run(statement, parameters)
	if err != nil {
		return Response{}, models.ErrDatabase
	}

	resources := []Element{}
	for result.Next() {
		test := result.Record().GetByIndex(0)
		resultMap := test.(map[string]interface{})
		var childResource, parentResource resourceNode

		err = decodeResource(resultMap, "child", &childResource)
		if err != nil {
			return
		}

		decodeResource(resultMap, "parent", &parentResource)
		if err != nil {
			return
		}

		parent := generateParentResource(parentResource.ID)
		resources = append(resources, generateElement(childResource, parent))
	}
	response = Response{Data: resources}
	return
}

func decodeResource(results map[string]interface{}, resourceKey string, target interface{}) (err error) {
	if results[resourceKey] != nil {
		node := results[resourceKey].(neo4j.Node)
		err = mapstructure.Decode(node.Props(), &target)
	}
	return
}

func generateParentResource(parentResourceID string) (parent *Parent) {
	if parentResourceID != "" {
		parent = &Parent{
			Data: &Identifier{
				Type: "resource",
				ID:   parentResourceID,
			},
		}
	}
	return
}

func generateElement(childResource resourceNode, parent *Parent) (element Element) {
	return Element{
		Type: "resource",
		ID:   childResource.ID,
		Attributes: &Resource{
			Name:     childResource.Name,
			SourceID: childResource.SourceID,
		},
		Relationships: &Relationships{
			Parent: parent,
		},
	}
}
