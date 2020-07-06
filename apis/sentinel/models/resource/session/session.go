package session

import (
	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/injection"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/neo4j"
	resource "github.com/bithippie/guard-my-app/apis/sentinel/models/resource/dto"
)

// Session interface defines methods needed to communicate/execute queries and a cleanup function when everything is done
type Session interface {
	Execute(statement string, parameters map[string]interface{}) (response resource.Output, err error)
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
	Name      string `mapstructure:"name"`
	SourceID  string `mapstructure:"source_id"`
	ID        string `mapstructure:"id"`
	ContextID string `mapstructure:"context_id"`
}

type contextNode struct {
	ID   string `mapstructure:"id"`
	Name string `mapstructure:"name"`
}

// Execute runs the statement passed as a query and populates the data parameter with result
func (n session) Execute(statement string, parameters map[string]interface{}) (response resource.Output, err error) {
	resultMap, err := n.session.Run(statement, parameters)
	resources := []resource.Details{}

	if err != nil {
		return resource.Output{}, models.ErrDatabase
	}
	if len(resultMap) == 0 {
		response = resource.Output{Data: resources}
		return
	}

	var childResource, parentResource resourceNode
	for _, result := range resultMap {

		err = injection.NodeDecoder(result, "child", &childResource)
		if err != nil {
			return
		}

		err = injection.NodeDecoder(result, "parent", &parentResource)
		if err != nil {
			return
		}

		var contexts *resource.Contexts
		contexts, err = decodeContexts(result, "context")
		if err != nil {
			return
		}

		parent := generateParentResource(parentResource.ID)
		resources = append(resources, generateDetails(childResource, parent, contexts))
	}

	response = resource.Output{Data: resources}
	return
}

func decodeContexts(results map[string]interface{}, field string) (contexts *resource.Contexts, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = models.ErrDatabase
		}
	}()

	var data = []resource.Data{}
	if results[field] != nil {
		for _, node := range results[field].([]interface{}) {
			var contextNode = contextNode{}
			injection.MapDecoder(node.(neo4j.Node).Props(), &contextNode)
			data = append(data, resource.Data{Type: "context", ID: contextNode.ID})
		}
		contexts = &resource.Contexts{Data: data}
	}
	return
}

func generateParentResource(parentResourceID string) (parent *resource.Parent) {
	if parentResourceID != "" {
		parent = &resource.Parent{
			Data: &resource.Data{
				Type: "resource",
				ID:   parentResourceID,
			},
		}
	}
	return
}

func generateDetails(childResource resourceNode, parent *resource.Parent, contexts *resource.Contexts) resource.Details {
	return resource.Details{
		Type: "resource",
		ID:   childResource.ID,
		Attributes: &resource.Attributes{
			Name:      childResource.Name,
			SourceID:  childResource.SourceID,
			ContextID: childResource.ContextID,
		},
		Relationships: &resource.Relationships{
			Parent:   parent,
			Contexts: contexts,
		},
	}
}
