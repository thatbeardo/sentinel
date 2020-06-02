package repository_test

import (
	"testing"

	models "github.com/bithippie/guard-my-app/sentinel/models"
	"github.com/bithippie/guard-my-app/sentinel/models/permission/outputs"
	"github.com/bithippie/guard-my-app/sentinel/models/permission/repository"
	"github.com/bithippie/guard-my-app/sentinel/models/permission/testdata"
	"github.com/stretchr/testify/assert"
)

var createStatement = `
		MATCH (policy: Policy), (target: Resource)
		WHERE policy.id = $policyID AND target.id = $targetID
		CREATE (policy)-[r:PERMISSION {name: $name, permitted: $permitted, id: randomUUID()}]->(target)
		RETURN {permission: r}`

type mockSession struct {
	ExecuteResponse   outputs.Response
	ExecuteErr        error
	ExpectedStatement string
	ExpectedParameter map[string]interface{}
	t                 *testing.T
}

func (m mockSession) Execute(statement string, parameters map[string]interface{}) (outputs.Response, error) {
	assert.Equal(m.t, m.ExpectedStatement, statement)
	assert.Equal(m.t, m.ExpectedParameter, parameters)
	return m.ExecuteResponse, m.ExecuteErr
}

func TestCreate_SessionReturnsResponse_ReturnDatabaseErr(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Response,
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
	response, err := repository.Create(testdata.Payload, "test-policy-id", "test-target-id")
	assert.Equal(t, testdata.Response.Data[0], response)
	assert.Nil(t, err)
}

func TestCreate_SessionReturnsEmptyResponse_DatabaseErrorReturned(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   outputs.Response{Data: []outputs.Permission{}},
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
	_, err := repository.Create(testdata.Payload, "test-policy-id", "test-target-id")
	assert.Equal(t, models.ErrDatabase, err)
}
