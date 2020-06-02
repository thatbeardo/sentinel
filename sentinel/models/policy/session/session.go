package session

import (
	models "github.com/bithippie/guard-my-app/sentinel/models"
	"github.com/bithippie/guard-my-app/sentinel/models/injection"
	"github.com/bithippie/guard-my-app/sentinel/models/neo4j"
	"github.com/bithippie/guard-my-app/sentinel/models/policy/inputs"
	"github.com/bithippie/guard-my-app/sentinel/models/policy/outputs"
)

// Session interface defines methods needed to communicate/execute queries and a cleanup function when everything is done
type Session interface {
	Execute(statement string, parameters map[string]interface{}) (response outputs.Response, err error)
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

// Execute runs the statement passed as a query and populates the data parameter with result
func (n session) Execute(statement string, parameters map[string]interface{}) (response outputs.Response, err error) {
	results, err := n.session.Run(statement, parameters)
	policies := []outputs.Policy{}

	if err != nil {
		return outputs.Response{}, models.ErrDatabase
	}
	if len(results) == 0 {
		response = outputs.Response{Data: policies}
		return
	}

	var policy node
	for _, result := range results {
		err = decodeField(result, "policy", &policy)
		if err != nil {
			return
		}
		policies = append(policies, generatePolicy(policy.Name, policy.ID))
	}

	response = outputs.Response{Data: policies}
	return
}

func decodeField(results map[string]interface{}, field string, target interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = models.ErrDatabase
		}
	}()
	if results[field] != nil {
		node := results[field].(neo4j.Node)
		err = injection.MapDecoder(node.Props(), &target)
	}
	return
}

func generatePolicy(name string, id string) (policy outputs.Policy) {
	return outputs.Policy{
		PolicyDetails: inputs.PolicyDetails{
			Type: "policy",
			Attributes: &inputs.Attributes{
				Name: name,
			},
		},
		ID:            id,
		Relationships: &outputs.Relationships{},
	}
}
