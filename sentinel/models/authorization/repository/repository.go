package repository

import (
	"fmt"
	"strconv"
	"strings"

	authorization "github.com/bithippie/guard-my-app/sentinel/models/authorization/dto"
	"github.com/bithippie/guard-my-app/sentinel/models/authorization/session"
)

// Repository exposes wrapper methods around the underlying session
type Repository interface {
	GetAuthorizationForPrincipal(string, authorization.Input) (authorization.Output, error)
}

type repository struct {
	session session.Session
}

// GetAuthorizationForPrincipal calls underlying session and returns output data to the service
func (repo repository) GetAuthorizationForPrincipal(principalID string, input authorization.Input) (output authorization.Output, err error) {
	return repo.session.Execute(
		fmt.Sprintf(`
			MATCH(target:Resource)-[:OWNED_BY*0..%s]->(ancestors:Resource)
			%s
			MATCH (policy:Policy)-[permission:PERMISSION{%s}]->(target)
			%s
			MATCH (principal:Resource{id: $principalID})<-[grant:GRANTED_TO]-(policy)
			RETURN {target:target, permissions: COLLECT(permission)}`,
			strconv.Itoa(input.Depth),
			generateQueryFilter(input.Targets, "target", "id"),
			generatePermittedFilter(input),
			generateQueryFilter(input.Permissions, "permission", "name"),
		),
		map[string]interface{}{
			"principalID": principalID,
		},
		input,
	)
}

// New is a factory method to generate repository instances
func New(session session.Session) Repository {
	return repository{
		session: session,
	}
}

func generateQueryFilter(elements []string, property string, field string) (properties string) {
	if len(elements) != 0 {
		properties = fmt.Sprintf(`WHERE %s.%s IN ["%s"]`, property, field, strings.Join(elements, "\",\""))
	}
	return
}

func generatePermittedFilter(input authorization.Input) (properties string) {
	if !input.IncludeDenied {
		properties = `permitted:"allow"`
	}
	return
}
