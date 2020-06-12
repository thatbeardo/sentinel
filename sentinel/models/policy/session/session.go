package session

import (
	models "github.com/bithippie/guard-my-app/sentinel/models"
	"github.com/bithippie/guard-my-app/sentinel/models/injection"
	"github.com/bithippie/guard-my-app/sentinel/models/neo4j"
	policy "github.com/bithippie/guard-my-app/sentinel/models/policy/dto"
)

// Session interface defines methods needed to communicate/execute queries and a cleanup function when everything is done
type Session interface {
	Execute(statement string, parameters map[string]interface{}) (response policy.Output, err error)
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
func (n session) Execute(statement string, parameters map[string]interface{}) (response policy.Output, err error) {
	results, err := n.session.Run(statement, parameters)
	policies := []policy.Details{}

	if err != nil {
		return policy.Output{}, models.ErrDatabase
	}
	if len(results) == 0 {
		response = policy.Output{Data: policies}
		return
	}

	for _, result := range results {
		var data node
		var principals *policy.Relationship
		var targets *policy.Relationship

		err = injection.NodeDecoder(result, "policy", &data)
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

		policies = append(policies, generatePolicy(data.Name, data.ID, principals, targets))
	}

	response = policy.Output{Data: policies}
	return
}

func decodeResources(results map[string]interface{}, field string) (resources *policy.Relationship, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = models.ErrDatabase
		}
	}()

	var data = []policy.Resource{}
	if results[field] != nil {
		for _, node := range results[field].([]interface{}) {
			var resource = resource{}
			injection.MapDecoder(node.(neo4j.Node).Props(), &resource)
			data = append(data, policy.Resource{Type: "resource", ID: resource.ID})
		}
		resources = &policy.Relationship{Data: data}
	}
	return
}

func generatePolicy(name string, id string, principals, targetResources *policy.Relationship) policy.Details {
	return policy.Details{
		InputDetails: policy.InputDetails{
			Type: "policy",
			Attributes: &policy.Attributes{
				Name: name,
			},
		},
		ID: id,
		Relationships: &policy.Relationships{
			Principals:      principals,
			TargetResources: targetResources,
		},
	}
}
