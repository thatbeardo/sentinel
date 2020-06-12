package session

import (
	models "github.com/bithippie/guard-my-app/sentinel/models"
	"github.com/bithippie/guard-my-app/sentinel/models/injection"
	"github.com/bithippie/guard-my-app/sentinel/models/neo4j"
	resource "github.com/bithippie/guard-my-app/sentinel/models/resource/dto"
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
	Name     string `mapstructure:"name"`
	SourceID string `mapstructure:"source_id"`
	ID       string `mapstructure:"id"`
}

type policyNode struct {
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

		var policies *resource.Policies
		policies, err = decodePolicies(result, "policy")
		if err != nil {
			return
		}

		parent := generateParentResource(parentResource.ID)
		resources = append(resources, generateDetails(childResource, parent, policies))
	}

	response = resource.Output{Data: resources}
	return
}

func decodePolicies(results map[string]interface{}, field string) (policies *resource.Policies, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = models.ErrDatabase
		}
	}()

	var data = []resource.Data{}
	if results[field] != nil {
		for _, node := range results[field].([]interface{}) {
			var policyNode = policyNode{}
			injection.MapDecoder(node.(neo4j.Node).Props(), &policyNode)
			data = append(data, resource.Data{Type: "policy", ID: policyNode.ID})
		}
		policies = &resource.Policies{Data: data}
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

func generateDetails(childResource resourceNode, parent *resource.Parent, policies *resource.Policies) resource.Details {
	return resource.Details{
		Type: "resource",
		ID:   childResource.ID,
		Attributes: &resource.Attributes{
			Name:     childResource.Name,
			SourceID: childResource.SourceID,
		},
		Relationships: &resource.Relationships{
			Parent:   parent,
			Policies: policies,
		},
	}
}
