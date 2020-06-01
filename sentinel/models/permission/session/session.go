package session

import (
	models "github.com/bithippie/guard-my-app/sentinel/models"
	"github.com/bithippie/guard-my-app/sentinel/models/injection"
	"github.com/bithippie/guard-my-app/sentinel/models/neo4j"
	"github.com/bithippie/guard-my-app/sentinel/models/permission/inputs"
	"github.com/bithippie/guard-my-app/sentinel/models/permission/outputs"
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
	ID        string `mapstructure:"id"`
	Name      string `mapstructure:"name"`
	Permitted string `mapstructure:"permitted"`
}

// Execute runs the statement passed as a query and populates the data parameter with result
func (n session) Execute(statement string, parameters map[string]interface{}) (response outputs.Response, err error) {
	results, err := n.session.Run(statement, parameters)
	permissions := []outputs.Permission{}

	if err != nil {
		return outputs.Response{}, models.ErrDatabase
	}
	if len(results) == 0 {
		response = outputs.Response{Data: permissions}
		return
	}

	var permission node
	for _, result := range results {
		err = decodeField(result, "permission", &permission)
		if err != nil {
			return
		}
		permissions = append(permissions, generatePermission(permission.ID, permission.Name, permission.Permitted))
	}

	response = outputs.Response{Data: permissions}
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

func generatePermission(id string, name string, permitted string) (permission outputs.Permission) {
	return outputs.Permission{
		PermissionDetails: inputs.PermissionDetails{
			Type: "permission",
			Attributes: &inputs.Attributes{
				Name:      name,
				Permitted: permitted,
			},
		},
		ID: id,
	}
}
