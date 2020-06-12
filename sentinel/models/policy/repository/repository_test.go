package repository_test

import (
	"errors"
	"testing"

	models "github.com/bithippie/guard-my-app/sentinel/models"
	policy "github.com/bithippie/guard-my-app/sentinel/models/policy/dto"
	"github.com/bithippie/guard-my-app/sentinel/models/policy/repository"
	"github.com/bithippie/guard-my-app/sentinel/models/policy/testdata"
	"github.com/stretchr/testify/assert"
)

var errTest = errors.New("test-error")

var getStatement = `
		MATCH (policy:Policy) 
		OPTIONAL MATCH(policy)-[:GRANTED_TO]->(principal:Resource)
		OPTIONAL MATCH(policy)-[:PERMISSION]->(target:Resource)
		RETURN {policy:policy, principals:COLLECT(principal), targets:COLLECT(target)}`

var getByIDStatement = `
		MATCH (policy:Policy) 
		WHERE policy.id = $id
		OPTIONAL MATCH(policy)-[:GRANTED_TO]->(principal:Resource)
		OPTIONAL MATCH(policy)-[:PERMISSION]->(target:Resource)
		RETURN {policy:policy, principals:COLLECT(principal), targets:COLLECT(target)}`

var createStatement = `
		CREATE ( policy:Policy{ name:$name, id:randomUUID() })
		RETURN { policy: policy }`

var updateStatement = `
		MATCH (policy:Policy)
		WHERE policy.id = $id
		SET policy.name = $name
		RETURN { policy: policy }`

var deleteStatement = `
		MATCH (policy:Policy)
		WHERE policy.id = $id
		DETACH DELETE policy`

type mockSession struct {
	ExecuteResponse   policy.Output
	ExecuteErr        error
	ExpectedStatement string
	ExpectedParameter map[string]interface{}
	t                 *testing.T
}

func (m mockSession) Execute(statement string, parameters map[string]interface{}) (policy.Output, error) {
	assert.Equal(m.t, m.ExpectedStatement, statement)
	assert.Equal(m.t, m.ExpectedParameter, parameters)
	return m.ExecuteResponse, m.ExecuteErr
}

func TestGet_SessionReturnsError_ErrorReturned(t *testing.T) {
	session := mockSession{
		ExecuteErr:        errTest,
		ExpectedStatement: getStatement,
		ExpectedParameter: map[string]interface{}{},
		t:                 t,
	}
	repository := repository.New(session)
	_, err := repository.Get()
	assert.Equal(t, errTest, err)
}

func TestGet_SessionReturnsResponse_NoErrorsReturned(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: getStatement,
		ExpectedParameter: map[string]interface{}{},
		t:                 t,
	}
	repository := repository.New(session)
	response, err := repository.Get()
	assert.Equal(t, testdata.Output, response)
	assert.Nil(t, err)
}

func TestGetByID_SessionReturnsError_ErrorReturned(t *testing.T) {
	session := mockSession{
		ExecuteErr:        errTest,
		ExpectedStatement: getByIDStatement,
		ExpectedParameter: map[string]interface{}{
			"id": "test-id",
		},
		t: t,
	}
	repository := repository.New(session)
	_, err := repository.GetByID("test-id")
	assert.Equal(t, models.ErrNotFound, err)
}

func TestGetByID_SessionReturnsResponse_NoErrorsReturned(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: getByIDStatement,
		ExpectedParameter: map[string]interface{}{
			"id": "test-id",
		},
		t: t,
	}
	repository := repository.New(session)
	response, err := repository.GetByID("test-id")
	assert.Equal(t, testdata.OutputDetails, response)
	assert.Nil(t, err)
}

func TestCreate_SessionReturnsError_ErrorReturned(t *testing.T) {
	session := mockSession{
		ExecuteErr:        errors.New("test-error"),
		ExpectedStatement: createStatement,
		ExpectedParameter: map[string]interface{}{
			"name": "test-policy",
		},
		t: t,
	}
	repository := repository.New(session)
	response, err := repository.Create(testdata.Input)
	assert.Equal(t, policy.OutputDetails{}, response)
	assert.NotNil(t, err)
}

func TestCreate_SessionReturnsEmptyResponse_DatabaseErrorReturned(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   policy.Output{Data: []policy.Details{}},
		ExpectedStatement: createStatement,
		ExpectedParameter: map[string]interface{}{
			"name": "test-policy",
		},
		t: t,
	}
	repository := repository.New(session)
	_, err := repository.Create(testdata.Input)
	assert.Equal(t, models.ErrDatabase, err)
}

func TestCreate_SessionReturnsResponse_NoErrorReturned(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: createStatement,
		ExpectedParameter: map[string]interface{}{
			"name": "test-policy",
		},
		t: t,
	}
	repository := repository.New(session)
	response, err := repository.Create(testdata.Input)
	assert.Nil(t, err)
	assert.Equal(t, testdata.OutputDetails, response)
}

func TestUpdate_SessionReturnsError_ErrorReturned(t *testing.T) {
	session := mockSession{
		ExecuteErr:        errors.New("test-error"),
		ExpectedStatement: updateStatement,
		ExpectedParameter: map[string]interface{}{
			"name": "test-policy",
			"id":   "test-id",
		},
		t: t,
	}
	repository := repository.New(session)
	response, err := repository.Update("test-id", testdata.Input)
	assert.Equal(t, policy.OutputDetails{}, response)
	assert.NotNil(t, err)
}

func TestUpdate_SessionReturnsEmptyResponse_DatabaseErrorReturned(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   policy.Output{Data: []policy.Details{}},
		ExpectedStatement: updateStatement,
		ExpectedParameter: map[string]interface{}{
			"name": "test-policy",
			"id":   "test-id",
		},
		t: t,
	}
	repository := repository.New(session)
	_, err := repository.Update("test-id", testdata.Input)
	assert.Equal(t, models.ErrDatabase, err)
}

func TestUpdate_SessionReturnsResponse_NoErrorReturned(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: updateStatement,
		ExpectedParameter: map[string]interface{}{
			"name": "test-policy",
			"id":   "test-id",
		},
		t: t,
	}
	repository := repository.New(session)
	response, err := repository.Update("test-id", testdata.Input)
	assert.Nil(t, err)
	assert.Equal(t, testdata.OutputDetails, response)
}

func TestDelete_SessionReturnsResponse_ReturnsNoErrors(t *testing.T) {
	session := mockSession{
		ExpectedStatement: deleteStatement,
		ExpectedParameter: map[string]interface{}{
			"id": "test-id",
		},
		t: t,
	}
	repository := repository.New(session)
	err := repository.Delete("test-id")
	assert.Nil(t, err)
}

func TestDelete_SessionReturnsError_ReturnsNoErrors(t *testing.T) {
	errTest := errors.New("test-error")
	session := mockSession{
		ExecuteErr:        errTest,
		ExpectedStatement: deleteStatement,
		ExpectedParameter: map[string]interface{}{
			"id": "test-id",
		},
		t: t,
	}
	repository := repository.New(session)
	err := repository.Delete("test-id")
	assert.Equal(t, errTest, err)
}
