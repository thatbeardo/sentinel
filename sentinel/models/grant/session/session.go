package session

import (
	models "github.com/bithippie/guard-my-app/sentinel/models"
	"github.com/bithippie/guard-my-app/sentinel/models/grant/inputs"
	"github.com/bithippie/guard-my-app/sentinel/models/grant/outputs"
	"github.com/bithippie/guard-my-app/sentinel/models/injection"
	"github.com/bithippie/guard-my-app/sentinel/models/neo4j"
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

type relationship struct {
	ID        string `mapstructure:"id"`
	Name      string `mapstructure:"name"`
	WithGrant bool   `mapstructure:"with_grant"`
}

// Execute runs the statement passed as a query and populates the data parameter with result
func (n session) Execute(statement string, parameters map[string]interface{}) (response outputs.Response, err error) {
	results, err := n.session.Run(statement, parameters)
	grants := []outputs.Grant{}

	if err != nil {
		return outputs.Response{}, models.ErrDatabase
	}
	if len(results) == 0 {
		response = outputs.Response{Data: grants}
		return
	}

	var grant relationship
	for _, result := range results {
		err = decodeField(result, "grant", &grant)
		if err != nil {
			return
		}
		grants = append(grants, generateGrant(grant.ID, grant.Name, grant.WithGrant))
	}

	response = outputs.Response{Data: grants}
	return
}

func decodeField(results map[string]interface{}, field string, target interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = models.ErrDatabase
		}
	}()
	if results[field] != nil {
		node := results[field].(neo4j.Relationship)
		err = injection.MapDecoder(node.Props(), &target)
	}
	return
}

func generateGrant(id string, name string, withGrant bool) (permission outputs.Grant) {
	return outputs.Grant{
		GrantDetails: inputs.GrantDetails{
			Type: "grant",
			Attributes: &inputs.Attributes{
				WithGrant: withGrant,
			},
		},
		ID: id,
	}
}
