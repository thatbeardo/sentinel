package repository_test

import (
	"errors"
	"testing"

	authorization "github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/repository"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/testdata"
	"github.com/stretchr/testify/assert"
)

var depth4TargetsPresentPermittedAbsentPermissionsPresent = `
			MATCH(target:Resource)-[:OWNED_BY*0..4]->(ancestors:Resource)
			WHERE target.id IN ["ghi","jkl"]
			MATCH (principal:Resource{id: $principal_id})
			MATCH (context:Context{id:principal.context_id})-[permission:PERMISSION{permitted:"allow"}]->(target)
			WHERE permission.name IN ["abc","def"]
			MATCH (principal)<-[grant:GRANTED_TO]-(context)
			RETURN {target:target, permissions: COLLECT(permission)}`

var depth4TargetsAbsentPermittedAbsentPermissionsPresent = `
			MATCH(target:Resource)-[:OWNED_BY*0..4]->(ancestors:Resource)
			
			MATCH (principal:Resource{id: $principal_id})
			MATCH (context:Context{id:principal.context_id})-[permission:PERMISSION{permitted:"allow"}]->(target)
			WHERE permission.name IN ["abc","def"]
			MATCH (principal)<-[grant:GRANTED_TO]-(context)
			RETURN {target:target, permissions: COLLECT(permission)}`

var depth4TargetsAbsentPermittedAbsentPermissionNamesAbsent = `
			MATCH(target:Resource)-[:OWNED_BY*0..4]->(ancestors:Resource)
			
			MATCH (principal:Resource{id: $principal_id})
			MATCH (context:Context{id:$context_id})-[permission:PERMISSION{}]->(target)
			
			MATCH (principal)<-[grant:GRANTED_TO]-(context)
			RETURN {target:target, permissions: COLLECT(permission)}`

var resourceOwnershipStatement = `
		MATCH (target:Resource{id: $target_id})
		-[:OWNED_BY*0..]->(hub:Resource)
		<-[tenantPermission:PERMISSION {name:"sentinel:read", permitted:"allow"}]-(tenantContext:Context)
		-[:GRANTED_TO]->(tenant:Resource {source_id:$tenant})
		<-[clientPermission:PERMISSION {name:"sentinel:read", permitted:"allow"}]-(clientContext:Context)
		-[:GRANTED_TO]->(client:Resource {source_id:$client_id})
		RETURN {target: target, permissions: COLLECT(clientPermission)}`

var contextOwnershipStatement = `
		MATCH (context:Context{id: $context_id})-[:GRANTED_TO]->(target:Resource)
		-[:OWNED_BY*0..]->(hub:Resource)
		<-[tenantPermission:PERMISSION {name:"sentinel:read", permitted:"allow"}]-(tenantContext:Context)
		-[:GRANTED_TO]->(tenant:Resource {source_id:$tenant})
		<-[clientPermission:PERMISSION {name:"sentinel:read", permitted:"allow"}]-(clientContext:Context)
		-[:GRANTED_TO]->(client:Resource {source_id:$client_id})
		RETURN {target: target, permissions: COLLECT(clientPermission)}`

var permissionOwnershipStatement = `
		MATCH (context:Context)-[:PERMISSION{id: $permission_id}]->(target:Resource)
		-[:OWNED_BY*0..]->(hub:Resource)
		<-[tenantPermission:PERMISSION {name:"sentinel:read", permitted:"allow"}]-(tenantContext:Context)
		-[:GRANTED_TO]->(tenant:Resource {source_id:$tenant})
		<-[clientPermission:PERMISSION {name:"sentinel:read", permitted:"allow"}]-(clientContext:Context)
		-[:GRANTED_TO]->(client:Resource {source_id:$client_id})
		RETURN {target: target, permissions: COLLECT(clientPermission)}`

type mockSession struct {
	ExecuteResponse   authorization.Output
	ExecuteErr        error
	ExpectedStatement string
	ExpectedParameter map[string]interface{}
	t                 *testing.T
}

func (m mockSession) Execute(statement string, parameters map[string]interface{}) (authorization.Output, error) {
	assert.Equal(m.t, m.ExpectedStatement, statement)
	assert.Equal(m.t, m.ExpectedParameter, parameters)
	return m.ExecuteResponse, m.ExecuteErr
}

func TestGetAuthorizationTargetsMentioned_SessionReturnsData_ReturnData(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: depth4TargetsPresentPermittedAbsentPermissionsPresent,
		ExpectedParameter: map[string]interface{}{
			"principal_id": "test-principal-id",
			"context_id":   "",
		},
		t: t,
	}

	repository := repository.New(session)
	response, err := repository.GetAuthorizationForPrincipal(
		"test-principal-id",
		"",
		authorization.Input{Depth: 4,
			IncludeDenied: false,
			Permissions:   []string{"abc", "def"},
			Targets:       []string{"ghi", "jkl"},
		},
	)

	assert.Equal(t, testdata.Output, response)
	assert.Nil(t, err)
}

func TestGetAuthorizationFor_SessionReturnsData_ReturnData(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: depth4TargetsAbsentPermittedAbsentPermissionsPresent,
		ExpectedParameter: map[string]interface{}{
			"principal_id": "test-principal-id",
			"context_id":   "",
		},
		t: t,
	}

	repository := repository.New(session)
	response, err := repository.GetAuthorizationForPrincipal(
		"test-principal-id",
		"",
		authorization.Input{Depth: 4, IncludeDenied: false, Permissions: []string{"abc", "def"}},
	)

	assert.Equal(t, testdata.Output, response)
	assert.Nil(t, err)
}

func TestGetAuthorizationForNonPermittedPrincipals_SessionReturnsData_ReturnData(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: depth4TargetsAbsentPermittedAbsentPermissionNamesAbsent,
		ExpectedParameter: map[string]interface{}{
			"principal_id": "test-principal-id",
			"context_id":   "test-context-id",
		},
		t: t,
	}

	repository := repository.New(session)
	response, err := repository.GetAuthorizationForPrincipal(
		"test-principal-id",
		"test-context-id",
		authorization.Input{Depth: 4, IncludeDenied: true},
	)

	assert.Equal(t, testdata.Output, response)
	assert.Nil(t, err)
}

func TestGetAuthorizationForPrincipal_SessionReturnsError_ErrorReturned(t *testing.T) {
	session := mockSession{
		ExecuteErr:        errors.New("some-test-error"),
		ExpectedStatement: depth4TargetsAbsentPermittedAbsentPermissionNamesAbsent,
		ExpectedParameter: map[string]interface{}{
			"principal_id": "test-principal-id",
			"context_id":   "test-context-id",
		},
		t: t,
	}

	repository := repository.New(session)
	_, err := repository.GetAuthorizationForPrincipal(
		"test-principal-id",
		"test-context-id",
		authorization.Input{Depth: 4, IncludeDenied: true},
	)

	assert.Equal(t, errors.New("some-test-error"), err)
}

func TestIsTargetOwnedByClient_RepositoryReturnsErr_ReturnFalse(t *testing.T) {
	session := mockSession{
		ExecuteErr:        errors.New("some-test-error"),
		ExpectedStatement: resourceOwnershipStatement,
		ExpectedParameter: map[string]interface{}{
			"client_id": "client-id",
			"tenant":    "tenant",
			"target_id": "target-id",
		},
		t: t,
	}

	repository := repository.New(session)
	result := repository.IsTargetOwnedByClient("client-id", "tenant", "target-id")

	assert.False(t, result)
}

func TestIsTargetOwnedByClient_RepositoryReturnsEmptyData_ReturnFalse(t *testing.T) {
	output := testdata.Output
	output.Data = []authorization.Details{}
	session := mockSession{
		ExecuteResponse:   output,
		ExpectedStatement: resourceOwnershipStatement,
		ExpectedParameter: map[string]interface{}{
			"client_id": "client-id",
			"tenant":    "tenant",
			"target_id": "target-id",
		},
		t: t,
	}

	repository := repository.New(session)
	result := repository.IsTargetOwnedByClient("client-id", "tenant", "target-id")

	assert.False(t, result)
}

func TestIsTargetOwnedByClient_RepositoryReturnsData_ReturnTrue(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: resourceOwnershipStatement,
		ExpectedParameter: map[string]interface{}{
			"client_id": "client-id",
			"tenant":    "tenant",
			"target_id": "target-id",
		},
		t: t,
	}

	repository := repository.New(session)
	result := repository.IsTargetOwnedByClient("client-id", "tenant", "target-id")

	assert.True(t, result)
}

func TestIsContextOwnedByTenant_SessionReturnsEmptyResponse_ReturnFalse(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   authorization.Output{Data: []authorization.Details{}},
		ExpectedStatement: contextOwnershipStatement,
		ExpectedParameter: map[string]interface{}{
			"client_id":  "client-id",
			"tenant":     "tenant",
			"context_id": "context-id",
		},
		t: t,
	}

	repository := repository.New(session)
	result := repository.IsContextOwnedByClient("client-id", "tenant", "context-id")

	assert.False(t, result)
}

func TestIscontextOwnedByTenant_SessionReturnsError_ReturnFalse(t *testing.T) {
	session := mockSession{
		ExecuteErr:        errors.New("some-test-error"),
		ExpectedStatement: contextOwnershipStatement,
		ExpectedParameter: map[string]interface{}{
			"client_id":  "client-id",
			"tenant":     "tenant",
			"context_id": "context-id",
		},
		t: t,
	}

	repository := repository.New(session)
	result := repository.IsContextOwnedByClient("client-id", "tenant", "context-id")

	assert.False(t, result)
}

func TestIsContextOwnedByTenant_SessionReturnsData_ReturnTrue(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: contextOwnershipStatement,
		ExpectedParameter: map[string]interface{}{
			"client_id":  "client-id",
			"tenant":     "tenant",
			"context_id": "context-id",
		},
		t: t,
	}

	repository := repository.New(session)
	result := repository.IsContextOwnedByClient("client-id", "tenant", "context-id")

	assert.True(t, result)
}

func TestIsPermissionOwnedByTenant_SessionReturnsEmptyResponse_ReturnFalse(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   authorization.Output{Data: []authorization.Details{}},
		ExpectedStatement: permissionOwnershipStatement,
		ExpectedParameter: map[string]interface{}{
			"client_id":     "client-id",
			"tenant":        "tenant",
			"permission_id": "permission-id",
		},
		t: t,
	}

	repository := repository.New(session)
	result := repository.IsPermissionOwnedByTenant("client-id", "tenant", "permission-id")

	assert.False(t, result)
}

func TestIsPermissionOwnedByTenant_SessionReturnsError_ReturnFalse(t *testing.T) {
	session := mockSession{
		ExecuteErr:        errors.New("some-test-error"),
		ExpectedStatement: permissionOwnershipStatement,
		ExpectedParameter: map[string]interface{}{
			"client_id":     "client-id",
			"tenant":        "tenant",
			"permission_id": "permission-id",
		},
		t: t,
	}

	repository := repository.New(session)
	result := repository.IsPermissionOwnedByTenant("client-id", "tenant", "permission-id")

	assert.False(t, result)
}

func TestIsPermissionOwnedByTenant_SessionReturnsData_ReturnTrue(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: permissionOwnershipStatement,
		ExpectedParameter: map[string]interface{}{
			"client_id":     "client-id",
			"tenant":        "tenant",
			"permission_id": "permission-id",
		},
		t: t,
	}

	repository := repository.New(session)
	result := repository.IsPermissionOwnedByTenant("client-id", "tenant", "permission-id")

	assert.True(t, result)
}
