package repository_test

import (
	"errors"
	"testing"

	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	permission "github.com/bithippie/guard-my-app/apis/sentinel/models/permission/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/permission/repository"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/permission/testdata"
	"github.com/stretchr/testify/assert"
)

var errTest = errors.New("some-test-error")

var getAllPermissionsForPolicyStatement = `
		MATCH (policy:Policy { id: $policyID } )-[permission:PERMISSION]->(resource:Resource)
		RETURN {permission: permission}`

var getAllPermissionsForPolicyWithResourceStatement = `
		MATCH (policy:Policy { id: $policyID } )-[permission:PERMISSION]->(resource:Resource { id: $resourceID })
		RETURN {permission: permission}`

var createStatement = `
		MATCH (policy: Policy), (target: Resource)
		WHERE policy.id = $policyID AND target.id = $targetID
		CREATE (policy)-[r:PERMISSION {name: $name, permitted: $permitted, id: randomUUID()}]->(target)
		RETURN {permission: r}`

var updateStatement = `
		MATCH(policy:Policy)-[permission:PERMISSION]->(resource:Resource)
		WHERE permission.id = $id
		SET permission.name = $name
		SET permission.permitted = $permitted
		RETURN {permission: permission}`

var deleteStatement = `
		MATCH(policy:Policy)-[permission:PERMISSION]->(resource:Resource)
		WHERE permission.id = $id
		DELETE permission`

type mockSession struct {
	ExecuteResponse   permission.Output
	ExecuteErr        error
	ExpectedStatement string
	ExpectedParameter map[string]interface{}
	t                 *testing.T
}

func (m mockSession) Execute(statement string, parameters map[string]interface{}) (permission.Output, error) {
	assert.Equal(m.t, m.ExpectedStatement, statement)
	assert.Equal(m.t, m.ExpectedParameter, parameters)
	return m.ExecuteResponse, m.ExecuteErr
}

func TestGetAllPermissionsForPolicy_SessionReturnsResponse_ReturnDatabaseErr(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: getAllPermissionsForPolicyStatement,
		ExpectedParameter: map[string]interface{}{
			"policyID": "test-policy-id",
		},
		t: t,
	}
	repository := repository.New(session)
	response, err := repository.GetAllPermissionsForPolicy("test-policy-id")
	assert.Equal(t, testdata.Output, response)
	assert.Nil(t, err)
}

func TestGetAllPermissionsForPolicy_SessionReturnsEmptyResponse_DatabaseErrorReturned(t *testing.T) {
	session := mockSession{
		ExecuteErr:        errTest,
		ExpectedStatement: getAllPermissionsForPolicyStatement,
		ExpectedParameter: map[string]interface{}{
			"policyID": "test-policy-id",
		},
		t: t,
	}
	repository := repository.New(session)
	_, err := repository.GetAllPermissionsForPolicy("test-policy-id")
	assert.Equal(t, errTest, err)
}

func TestGetAllPermissionsForPolicyWithResource_SessionReturnsResponse_ReturnDatabaseErr(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: getAllPermissionsForPolicyWithResourceStatement,
		ExpectedParameter: map[string]interface{}{
			"policyID":   "test-policy-id",
			"resourceID": "test-resource-id",
		},
		t: t,
	}
	repository := repository.New(session)
	response, err := repository.GetAllPermissionsForPolicyWithResource("test-policy-id", "test-resource-id")
	assert.Equal(t, testdata.Output, response)
	assert.Nil(t, err)
}

func TestGetAllPermissionsForPolicyWithResource_SessionReturnsEmptyResponse_DatabaseErrorReturned(t *testing.T) {
	session := mockSession{
		ExecuteErr:        errTest,
		ExpectedStatement: getAllPermissionsForPolicyWithResourceStatement,
		ExpectedParameter: map[string]interface{}{
			"policyID":   "test-policy-id",
			"resourceID": "test-resource-id",
		},
		t: t,
	}
	repository := repository.New(session)
	_, err := repository.GetAllPermissionsForPolicyWithResource("test-policy-id", "test-resource-id")
	assert.Equal(t, errTest, err)
}

func TestCreate_SessionReturnsResponse_ReturnDatabaseErr(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: createStatement,
		ExpectedParameter: map[string]interface{}{
			"name":      "test-permission",
			"permitted": "allow",
			"policyID":  "test-policy-id",
			"targetID":  "test-target-id",
		},
		t: t,
	}
	repository := repository.New(session)
	response, err := repository.Create(testdata.Input, "test-policy-id", "test-target-id")
	assert.Equal(t, testdata.OutputDetails, response)
	assert.Nil(t, err)
}

func TestCreate_SessionReturnsEmptyResponse_DatabaseErrorReturned(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   permission.Output{Data: []permission.Details{}},
		ExpectedStatement: createStatement,
		ExpectedParameter: map[string]interface{}{
			"name":      "test-permission",
			"permitted": "allow",
			"policyID":  "test-policy-id",
			"targetID":  "test-target-id",
		},
		t: t,
	}
	repository := repository.New(session)
	_, err := repository.Create(testdata.Input, "test-policy-id", "test-target-id")
	assert.Equal(t, models.ErrNotFound, err)
}

func TestUpdate_SessionReturnsResponse_ReturnDatabaseErr(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: updateStatement,
		ExpectedParameter: map[string]interface{}{
			"name":      "test-permission",
			"permitted": "allow",
			"id":        "test-permission-id",
		},
		t: t,
	}
	repository := repository.New(session)
	response, err := repository.Update("test-permission-id", testdata.Input)
	assert.Equal(t, testdata.OutputDetails, response)
	assert.Nil(t, err)
}

func TestUpdate_SessionReturnsEmptyResponse_DatabaseErrorReturned(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   permission.Output{Data: []permission.Details{}},
		ExpectedStatement: updateStatement,
		ExpectedParameter: map[string]interface{}{
			"name":      "test-permission",
			"permitted": "allow",
			"id":        "test-permission-id",
		},
		t: t,
	}
	repository := repository.New(session)
	_, err := repository.Update("test-permission-id", testdata.Input)
	assert.Equal(t, models.ErrNotFound, err)
}

func TestDelete_SessionReturnsResponse_ReturnDatabaseErr(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: deleteStatement,
		ExpectedParameter: map[string]interface{}{
			"id": "test-permission-id",
		},
		t: t,
	}
	repository := repository.New(session)
	err := repository.Delete("test-permission-id")
	assert.Nil(t, err)
}

func TestDelete_SessionReturnsErr_DatabaseErrorReturned(t *testing.T) {
	session := mockSession{
		ExecuteErr:        errTest,
		ExpectedStatement: deleteStatement,
		ExpectedParameter: map[string]interface{}{
			"id": "test-permission-id",
		},
		t: t,
	}
	repository := repository.New(session)
	err := repository.Delete("test-permission-id")
	assert.Equal(t, errTest, err)
}
