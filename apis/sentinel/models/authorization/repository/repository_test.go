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
			MATCH (policy:Policy)-[permission:PERMISSION{permitted:"allow"}]->(target)
			WHERE permission.name IN ["abc","def"]
			MATCH (principal:Resource{id: $principalID})<-[grant:GRANTED_TO]-(policy)
			RETURN {target:target, permissions: COLLECT(permission)}`

var depth4TargetsAbsentPermittedAbsentPermissionsPresent = `
			MATCH(target:Resource)-[:OWNED_BY*0..4]->(ancestors:Resource)
			
			MATCH (policy:Policy)-[permission:PERMISSION{permitted:"allow"}]->(target)
			WHERE permission.name IN ["abc","def"]
			MATCH (principal:Resource{id: $principalID})<-[grant:GRANTED_TO]-(policy)
			RETURN {target:target, permissions: COLLECT(permission)}`

var depth4TargetsAbsentPermittedAbsentPermissionNamesAbsent = `
			MATCH(target:Resource)-[:OWNED_BY*0..4]->(ancestors:Resource)
			
			MATCH (policy:Policy)-[permission:PERMISSION{}]->(target)
			
			MATCH (principal:Resource{id: $principalID})<-[grant:GRANTED_TO]-(policy)
			RETURN {target:target, permissions: COLLECT(permission)}`

var ownershipStatement = `
		MATCH (target:Resource{id: $targetID})
		-[OWNED_BY*0..]->(ancestors:Resource)
		<-[permission:PERMISSION]-(policy:Policy)
		-[:GRANTED_TO]->(tenant:Resource{source_id:$tenantID})
		RETURN {target: target, permissions: COLLECT(permission)}`

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
			"principalID": "test-principal-id",
		},
		t: t,
	}

	repository := repository.New(session)
	response, err := repository.GetAuthorizationForPrincipal(
		"test-principal-id",
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
			"principalID": "test-principal-id",
		},
		t: t,
	}

	repository := repository.New(session)
	response, err := repository.GetAuthorizationForPrincipal(
		"test-principal-id",
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
			"principalID": "test-principal-id",
		},
		t: t,
	}

	repository := repository.New(session)
	response, err := repository.GetAuthorizationForPrincipal(
		"test-principal-id",
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
			"principalID": "test-principal-id",
		},
		t: t,
	}

	repository := repository.New(session)
	_, err := repository.GetAuthorizationForPrincipal(
		"test-principal-id",
		authorization.Input{Depth: 4, IncludeDenied: true},
	)

	assert.Equal(t, errors.New("some-test-error"), err)
}

func TestIsTargetOwnedByTenant_RepositoryReturnsErr_ReturnFalse(t *testing.T) {
	session := mockSession{
		ExecuteErr:        errors.New("some-test-error"),
		ExpectedStatement: ownershipStatement,
		ExpectedParameter: map[string]interface{}{
			"targetID": "target-id",
			"tenantID": "tenant-id",
		},
		t: t,
	}

	repository := repository.New(session)
	result := repository.IsTargetOwnedByTenant("target-id", "tenant-id")

	assert.False(t, result)
}

func TestIsTargetOwnedByTenant_RepositoryReturnsEmptyData_ReturnFalse(t *testing.T) {
	output := testdata.Output
	output.Data = []authorization.Details{}
	session := mockSession{
		ExecuteResponse:   output,
		ExpectedStatement: ownershipStatement,
		ExpectedParameter: map[string]interface{}{
			"targetID": "target-id",
			"tenantID": "tenant-id",
		},
		t: t,
	}

	repository := repository.New(session)
	result := repository.IsTargetOwnedByTenant("target-id", "tenant-id")

	assert.False(t, result)
}

func TestIsTargetOwnedByTenant_RepositoryReturnsData_ReturnTrue(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: ownershipStatement,
		ExpectedParameter: map[string]interface{}{
			"targetID": "target-id",
			"tenantID": "tenant-id",
		},
		t: t,
	}

	repository := repository.New(session)
	result := repository.IsTargetOwnedByTenant("target-id", "tenant-id")

	assert.True(t, result)
}
