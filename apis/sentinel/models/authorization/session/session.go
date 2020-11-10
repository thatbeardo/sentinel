package session

import (
	"fmt"

	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	authorization "github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/injection"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/neo4j"
	permission "github.com/bithippie/guard-my-app/apis/sentinel/models/permission/dto"
	resource "github.com/bithippie/guard-my-app/apis/sentinel/models/resource/dto"
)

// Session interface defines methods needed to communicate/execute queries and a cleanup function when everything is done
type Session interface {
	Execute(string, map[string]interface{}) (authorization.Output, error)
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

// Execute runs the statement passed as a query and populates the data parameter with result
func (n session) Execute(statement string, parameters map[string]interface{}) (output authorization.Output, err error) {
	fmt.Println(statement)
	results, err := n.session.Run(statement, parameters)
	if err != nil {
		return
	}

	if len(results) == 0 {
		err = models.ErrNotFound
		return
	}

	details, err := generateOutputData(results)
	if err != nil {
		err = models.ErrDatabase
		return
	}

	output = authorization.Output{Data: details}
	return
}

func decodeEdges(results map[string]interface{}, field string) (permissions []permission.Attributes, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = models.ErrDatabase
		}
	}()

	if results[field] != nil {
		for _, node := range results[field].([]interface{}) {
			var egde edge
			injection.MapDecoder(node.(neo4j.Relationship).Props(), &egde)
			permissions = append(permissions, generatePermission(egde.Name, egde.Permitted))
		}
	}
	return
}

func generateOutputData(results []map[string]interface{}) (details []authorization.Details, err error) {
	entitlements := make(map[target]map[string]entitlement)
	for _, result := range results {
		var target target
		err = injection.NodeDecoder(result, "target", &target)
		if err != nil {
			return
		}

		var permissions []permission.Attributes
		permissions, err = decodeEdges(result, "permissions")
		if err != nil {
			return
		}

		var length int64 
		if len, ok := result["length"]; ok {
			length = len.(int64)
		}	
		updateEntitlementsMap(entitlements, target, permissions, length)
	}
	for t, p := range entitlements {
		details = append(details, generateDetail(t, p))
	}
	fmt.Println(entitlements)
	return
}

func updateEntitlementsMap(input map[target]map[string]entitlement, target target, permissions []permission.Attributes, length int64) map[target]map[string]entitlement  {
	if _, ok := input[target]; !ok {
		input[target] = make(map[string]entitlement)
	}
	for _, permission := range permissions {
		e := input[target]
		if permissionDetails, ok := e[permission.Name]; ok {
			if length < permissionDetails.Length {
				e[permission.Name] = entitlement{Permission: permission, Length: length}
			}
		} else {
			e[permission.Name] = entitlement{Permission: permission, Length: length}
		}
	}
	return input
}

func generateDetail(target target, entitlements map[string]entitlement) authorization.Details {
	if len(entitlements) == 0 {
		return authorization.Details{}
	}
	var permissions []permission.Attributes
	for _, entitlement := range entitlements {
		permissions = append(permissions, entitlement.Permission)
	}

	return authorization.Details{
		ID:   target.ID,
		Type: "resource",
		Attributes: resource.Attributes{
			Name:     target.Name,
			SourceID: target.SourceID,
		},
		Relationships: authorization.Relationships{
			Permissions: authorization.Permissions{
				Data: permissions,
			},
		},
	}
}

func generatePermission(name, permitted string) permission.Attributes {
	return permission.Attributes{
		Name:      name,
		Permitted: permitted,
	}
}

type target struct {
	Name     string `mapstructure:"name"`
	SourceID string `mapstructure:"source_id"`
	ID       string `mapstructure:"id"`
}

type edge struct {
	Name      string `mapstructure:"name"`
	Permitted string `mapstructure:"permitted"`
}

type entitlement struct {
	Permission permission.Attributes
	Length int64 
}