package session

import (
	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	grant "github.com/bithippie/guard-my-app/apis/sentinel/models/grant/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/injection"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/neo4j"
)

// Session interface defines methods needed to communicate/execute queries and a cleanup function when everything is done
type Session interface {
	Execute(statement string, parameters map[string]interface{}) (response grant.Output, err error)
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

type edge struct {
	ID        string `mapstructure:"id"`
	WithGrant bool   `mapstructure:"with_grant"`
}

type resource struct {
	ID       string `mapstructure:"id"`
	Name     string `mapstructure:"name"`
	SourceID string `mapstructure:"source_id"`
}

type context struct {
	ID   string `mapstructure:"id"`
	Name string `mapstructure:"name"`
}

// Execute runs the statement passed as a query and populates the data parameter with result
func (n session) Execute(statement string, parameters map[string]interface{}) (output grant.Output, err error) {
	results, err := n.session.Run(statement, parameters)
	grants := []grant.Details{}

	if err != nil {
		return grant.Output{}, models.ErrDatabase
	}
	if len(results) == 0 {
		output = grant.Output{Data: grants}
		return
	}

	for _, result := range results {
		var grant edge
		err = injection.EdgeDecoder(result, "grant", &grant)
		if err != nil {
			return
		}

		var context context
		err = injection.NodeDecoder(result, "context", &context)
		if err != nil {
			return
		}

		var principal resource
		err = injection.NodeDecoder(result, "principal", &principal)
		if err != nil {
			return
		}

		relationships := generateRelationships(context.ID, principal.ID)
		grants = append(grants, generateGrant(grant.ID, grant.WithGrant, relationships))
	}

	output = grant.Output{Data: grants}
	return
}

func generateRelationships(contextID, principalID string) grant.Relationships {
	// TODO:
	// Generate relationship only if ID is non-empty
	context := grant.Relationship{
		Data: grant.Data{
			Type: "context",
			ID:   contextID,
		},
	}

	principal := grant.Relationship{
		Data: grant.Data{
			Type: "resource",
			ID:   principalID,
		},
	}

	return grant.Relationships{
		Principal: &principal,
		Context:   &context,
	}
}

func generateGrant(id string, withGrant bool, relationships grant.Relationships) grant.Details {
	return grant.Details{
		InputDetails: grant.InputDetails{
			Type: "grant",
			Attributes: &grant.Attributes{
				WithGrant: withGrant,
			},
		},
		Relationships: relationships,
		ID:            id,
	}
}
