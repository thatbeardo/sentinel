package session

import (
	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/injection"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/neo4j"
	permission "github.com/bithippie/guard-my-app/apis/sentinel/models/permission/dto"
)

// Session interface defines methods needed to communicate/execute queries and a cleanup function when everything is done
type Session interface {
	Execute(string, map[string]interface{}) (permission.Output, error)
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
	Permitted string `mapstructure:"permitted"`
}

// Execute runs the statement passed as a query and populates the data parameter with result
func (n session) Execute(statement string, parameters map[string]interface{}) (output permission.Output, err error) {
	results, err := n.session.Run(statement, parameters)
	permissions := []permission.Details{}

	if err != nil {
		return permission.Output{}, models.ErrDatabase
	}
	if len(results) == 0 {
		output = permission.Output{Data: permissions}
		return
	}

	for _, result := range results {
		var permission relationship
		err = injection.EdgeDecoder(result, "permission", &permission)
		if err != nil {
			return
		}
		permissions = append(permissions, generatePermission(permission.ID, permission.Name, permission.Permitted))
	}

	output = permission.Output{Data: permissions}
	return
}

func generatePermission(id string, name string, permitted string) permission.Details {
	return permission.Details{
		InputDetails: permission.InputDetails{
			Type: "permission",
			Attributes: &permission.Attributes{
				Name:      name,
				Permitted: permitted,
			},
		},
		ID: id,
	}
}
