package repository

import (
	"fmt"
	"strconv"
	"strings"

	authorization "github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/session"
)

// Repository exposes wrapper methods around the underlying session
type Repository interface {
	GetAuthorizationForPrincipal(string, authorization.Input) (authorization.Output, error)
	IsTargetOwnedByTenant(string, string) bool
  IsPolicyOwnedByTenant(string, string) bool
	IsPermissionOwnedByTenant(string, string) bool
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
	)
}


func (repo repository) IsTargetOwnedByTenant(targetID, tenantID string) bool {
	results, err := repo.session.Execute(`
		MATCH (target:Resource{id: $targetID})
		-[OWNED_BY*0..]->(ancestors:Resource)
		<-[permission:PERMISSION]-(policy:Policy)
		-[:GRANTED_TO]->(tenant:Resource{source_id:$tenantID})
		RETURN {target: target, permissions: COLLECT(permission)}`,
		map[string]interface{}{
			"targetID": targetID,
			"tenantID": tenantID,
		},
	)
	return err == nil && len(results.Data) > 0
}


func (repo repository) IsPolicyOwnedByTenant(policyID, tenantID string) bool {
	results, err := repo.session.Execute(`
		MATCH (policy:Policy{id: $policyID})
		-[:GRANTED_TO]->(target:Resource)
		-[:OWNED_BY*0..]->(ancestors:Resource)
		<-[permission:PERMISSION]-(tenantPolicy:Policy)
		-[:GRANTED_TO]->(tenant:Resource{source_id: $tenantID})
		RETURN {target: target, permissions: COLLECT(permission)}`,
		map[string]interface{}{
			"policyID": policyID,
			"tenantID": tenantID,
		},
	)
	return err == nil && len(results.Data) > 0
}

func (repo repository) IsPermissionOwnedByTenant(permissionID, tenantID string) bool {
	results, err := repo.session.Execute(`
		MATCH (policy:Policy)-[:PERMISSION{id: $permissionID}]->(target:Resource)
		-[:OWNED_BY*0..]->(ancestors:Resource)
		<-[permission:PERMISSION]-(tenantPolicy:Policy)
		-[:GRANTED_TO]->(tenant:Resource{source_id: $tenantID})
		RETURN {target: target, permissions: COLLECT(permission)}`,
		map[string]interface{}{
			"permissionID": permissionID,
			"tenantID":     tenantID,
		},
	)
	return err == nil && len(results.Data) > 0
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
