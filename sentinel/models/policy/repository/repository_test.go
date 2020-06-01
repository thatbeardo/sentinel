package repository_test

import (
	"testing"

	models "github.com/bithippie/guard-my-app/sentinel/models"
	"github.com/bithippie/guard-my-app/sentinel/models/policy/outputs"
	"github.com/bithippie/guard-my-app/sentinel/models/policy/repository"
	"github.com/bithippie/guard-my-app/sentinel/models/policy/testdata"
	"github.com/stretchr/testify/assert"
)

var createStatement = `
		CREATE ( policy:Policy{ name:$name, id:randomUUID() })
		RETURN { policy: policy }`

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

func TestCreate_SessionReturnsResponse_NoErrors(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Response,
		ExpectedStatement: createStatement,
		ExpectedParameter: map[string]interface{}{
			"name": "test-policy",
		},
		t: t,
	}
	repository := repository.New(session)
	response, err := repository.Create(testdata.Payload)
	assert.Equal(t, testdata.Response.Data[0], response)
	assert.Nil(t, err)
}

func TestCreate_SessionReturnsEmptyResponse_DatabaseErrorReturned(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   outputs.Response{Data: []outputs.Policy{}},
		ExpectedStatement: createStatement,
		ExpectedParameter: map[string]interface{}{
			"name": "test-policy",
		},
		t: t,
	}
	repository := repository.New(session)
	_, err := repository.Create(testdata.Payload)
	assert.Equal(t, models.ErrDatabase, err)
}
