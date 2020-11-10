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
			MATCH path=(principal:Resource { id: $principal_id })
				<-[:GRANTED_TO]-
					(context:Context{id: %s})
						-[permission:PERMISSION]->
						(ancestors:Resource)<-[:OWNED_BY*0..%s]-
					(target:Resource)
				%s
			RETURN {target: target, permissions: COLLECT(permission), length: length(path)}`,
			generateContextFilter(contextID),
			strconv.Itoa(input.Depth),
			generateQueryFilter(input),
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

func generateQueryFilter(input authorization.Input) (filter string) {
	targetsFilter := generateTargetFilter(input)
	permissionsFilter := generatePermissionFilter(input)

	if len(targetsFilter) == 0 && len(permissionsFilter) == 0 {
		return
	} else if len(targetsFilter) != 0 && len(permissionsFilter) == 0 { 
		filter = fmt.Sprintf(`WHERE %s `, targetsFilter)
	} else if len(targetsFilter) == 0 && len(permissionsFilter) != 0 { 
		filter = fmt.Sprintf(`WHERE %s `, permissionsFilter)
	} else {
		filter = fmt.Sprintf(`WHERE %s AND %s `,  targetsFilter, permissionsFilter)
	}

	return
}

func generatePermissionFilter(input authorization.Input) (filter string) {
	if !input.IncludeDenied {
		filter += `permission.permitted IN ["allow"] `
	}
	if len(input.Permissions) != 0 {
		if len(filter) != 0 {
			filter += `AND ` 
		}
 		filter += fmt.Sprintf(`permission.name IN ["%s"] `, strings.Join(input.Permissions, "\",\""))
	}
	return
}


func generateTargetFilter(input authorization.Input) (properties string) {
	if len(input.Targets) != 0 {
		properties = fmt.Sprintf(`targets.id IN ["%s"]`, strings.Join(input.Targets, "\",\""))
	}
	return
}

func generateContextFilter(contextID string) (context string) {
	if contextID == "" {
		return "principal.context_id"
	}
	return "$context_id"
}
