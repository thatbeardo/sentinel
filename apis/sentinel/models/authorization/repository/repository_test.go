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

type mockSession struct {
	ExecuteResponse   authorization.Output
	ExecuteErr        error
	ExpectedStatement string
	ExpectedParameter map[string]interface{}
	t                 *testing.T
}

func (m mockSession) Execute(statement string, parameters map[string]interface{}, input authorization.Input) (authorization.Output, error) {
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
