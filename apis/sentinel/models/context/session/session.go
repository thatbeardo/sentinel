package session

import (
	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	context "github.com/bithippie/guard-my-app/apis/sentinel/models/context/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/injection"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/neo4j"
)

// Session interface defines methods needed to communicate/execute queries and a cleanup function when everything is done
type Session interface {
	Execute(statement string, parameters map[string]interface{}) (response context.Output, err error)
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

type node struct {
	Name string `mapstructure:"name"`
	ID   string `mapstructure:"id"`
}

type resource struct {
	ID string `mapstructure:"id"`
}

// Execute runs the statement passed as a query and populates the data parameter with result
func (n session) Execute(statement string, parameters map[string]interface{}) (response context.Output, err error) {
	results, err := n.session.Run(statement, parameters)
	contexts := []context.Details{}

	if err != nil {
		return context.Output{}, models.ErrDatabase
	}
	if len(results) == 0 {
		response = context.Output{Data: contexts}
		return
	}

	for _, result := range results {
		var data node
		var principals *context.Relationship
		var targets *context.Relationship

		err = injection.NodeDecoder(result, "context", &data)
		if err != nil {
			return
		}

		principals, err = decodeResources(result, "principals")
		if err != nil {
			return
		}

		targets, err = decodeResources(result, "targets")
		if err != nil {
			return
		}

		contexts = append(contexts, generatecontext(data.Name, data.ID, principals, targets))
	}

	response = context.Output{Data: contexts}
	return
}

func decodeResources(results map[string]interface{}, field string) (resources *context.Relationship, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = models.ErrDatabase
		}
	}()

	var data = []context.Resource{}
	resources = &context.Relationship{Data: data}
	if results[field] != nil {
		for _, node := range results[field].([]interface{}) {
			var resource = resource{}
			injection.MapDecoder(node.(neo4j.Node).Props(), &resource)
			data = append(data, context.Resource{Type: "resource", ID: resource.ID})
		}
		resources.Data = data
	}
	return
}

func generatecontext(name string, id string, principals, targetResources *context.Relationship) context.Details {
	return context.Details{
		InputDetails: context.InputDetails{
			Type: "context",
			Attributes: &context.Attributes{
				Name: name,
			},
		},
		ID: id,
		Relationships: &context.Relationships{
			Principals:      principals,
			TargetResources: targetResources,
		},
	}
}
