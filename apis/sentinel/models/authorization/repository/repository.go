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
	GetAuthorizationForPrincipal(principalID, contextID string, input authorization.Input) (authorization.Output, error)
	IsTargetOwnedByClient(clientID, tenant, targetID string) bool
	IsContextOwnedByClient(clientID, tenant, contextID string) bool
	IsPermissionOwnedByTenant(clientID, tenant, permissionID string) bool
}

type repository struct {
	session session.Session
}

// GetAuthorizationForPrincipal calls underlying session and returns output data to the service
func (repo repository) GetAuthorizationForPrincipal(principalID, contextID string, input authorization.Input) (output authorization.Output, err error) {
	return repo.session.Execute(
		fmt.Sprintf(`
			MATCH(target:Resource)-[:OWNED_BY*0..%s]->(ancestors:Resource)
			%s
			MATCH (principal:Resource{id: $principal_id})
			MATCH (context:Context{id:%s})-[permission:PERMISSION{%s}]->(target)
			%s
			MATCH (principal)<-[grant:GRANTED_TO]-(context)
			RETURN {target:target, permissions: COLLECT(permission)}`,
			strconv.Itoa(input.Depth),
			generateQueryFilter(input.Targets, "target", "id"),
			generateContextFilter(contextID),
			generatePermittedFilter(input),
			generateQueryFilter(input.Permissions, "permission", "name"),
		),
		map[string]interface{}{
			"principal_id": principalID,
			"context_id":   contextID,
		},
	)
}

func (repo repository) IsTargetOwnedByClient(clientID, tenant, targetID string) bool {
	results, err := repo.session.Execute(fmt.Sprintf(`
		MATCH (target:Resource{id: $target_id})
		%s`, generateOwnershipInspectionQuery()),
		map[string]interface{}{
			"client_id": clientID,
			"tenant":    tenant,
			"target_id": targetID,
		},
	)
	fmt.Println(results)
	return err == nil && len(results.Data) > 0
}

func (repo repository) IsContextOwnedByClient(clientID, tenant, contextID string) bool {
	results, err := repo.session.Execute(fmt.Sprintf(`
		MATCH (context:Context{id: $context_id})-[:GRANTED_TO]->(target:Resource)
		%s`, generateOwnershipInspectionQuery()),
		map[string]interface{}{
			"client_id":  clientID,
			"tenant":     tenant,
			"context_id": contextID,
		},
	)
	return err == nil && len(results.Data) > 0
}

func (repo repository) IsPermissionOwnedByTenant(clientID, tenant, permissionID string) bool {
	results, err := repo.session.Execute(fmt.Sprintf(`
		MATCH (context:Context)-[:PERMISSION{id: $permission_id}]->(target:Resource)
		%s`, generateOwnershipInspectionQuery()),
		map[string]interface{}{
			"client_id":     clientID,
			"tenant":        tenant,
			"permission_id": permissionID,
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

func generateOwnershipInspectionQuery() string {
	return `-[:OWNED_BY*0..]->(hub:Resource)
		<-[tenantPermission:PERMISSION {name:"sentinel:read", permitted:"allow"}]-(tenantContext:Context)
		-[:GRANTED_TO]->(tenant:Resource {source_id:$tenant})
		<-[clientPermission:PERMISSION {name:"sentinel:read", permitted:"allow"}]-(clientContext:Context)
		-[:GRANTED_TO]->(client:Resource {source_id:$client_id})
		RETURN {target: target, permissions: COLLECT(clientPermission)}`
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

func generateContextFilter(contextID string) (context string) {
	if contextID == "" {
		return "principal.context_id"
	}
	return "$context_id"
}
