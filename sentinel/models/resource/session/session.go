package session

import (
	models "github.com/bithippie/guard-my-app/sentinel/models"
	entity "github.com/bithippie/guard-my-app/sentinel/models/resource"
	"github.com/bithippie/guard-my-app/sentinel/models/resource/injection"
	"github.com/bithippie/guard-my-app/sentinel/models/resource/neo4j"
)

// Session interface defines methods needed to communicate/execute queries and a cleanup function when everything is done
type Session interface {
	Execute(statement string, parameters map[string]interface{}) (response entity.Response, err error)
}

type session struct {
	session neo4j.Runner
}

// NewNeo4jSession is a factory method to create Neo4J session instances
func NewNeo4jSession(neo4jsession neo4j.Runner) Session {
	return session{
		session: neo4jsession,
	}
}

type resourceNode struct {
	Name     string `mapstructure:"name"`
	SourceID string `mapstructure:"source_id"`
	ID       string `mapstructure:"id"`
}

// Execute runs the statement passed as a query and populates the data parameter with result
func (n session) Execute(statement string, parameters map[string]interface{}) (response entity.Response, err error) {
	resultMap, err := n.session.Run(statement, parameters)
	if err != nil {
		return entity.Response{}, models.ErrDatabase
	}

	resources := []entity.Element{}
	var childResource, parentResource resourceNode

	err = decodeResource(resultMap, "child", &childResource)
	if err != nil {
		return
	}

	err = decodeResource(resultMap, "parent", &parentResource)
	if err != nil {
		return
	}

	parent := generateParentResource(parentResource.ID)
	resources = append(resources, generateElement(childResource, parent))

	response = entity.Response{Data: resources}
	return
}

func decodeResource(results map[string]interface{}, resourceKey string, target interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = models.ErrDatabase
		}
	}()
	if results[resourceKey] != nil {
		node := results[resourceKey].(neo4j.Node)
		err = injection.MapDecoder(node.Props(), &target)
	}
	return
}

func generateParentResource(parentResourceID string) (parent *entity.Parent) {
	if parentResourceID != "" {
		parent = &entity.Parent{
			Data: &entity.Identifier{
				Type: "resource",
				ID:   parentResourceID,
			},
		}
	}
	return
}

func generateElement(childResource resourceNode, parent *entity.Parent) (element entity.Element) {
	return entity.Element{
		Type: "resource",
		ID:   childResource.ID,
		Attributes: &entity.Resource{
			Name:     childResource.Name,
			SourceID: childResource.SourceID,
		},
		Relationships: &entity.Relationships{
			Parent: parent,
		},
	}
}
