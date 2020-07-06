package repository_test

import (
	"errors"
	"testing"

	mocks "github.com/bithippie/guard-my-app/apis/sentinel/mocks/context"
	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	context "github.com/bithippie/guard-my-app/apis/sentinel/models/context/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/context/repository"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/context/testdata"
	"github.com/stretchr/testify/assert"
)

var errTest = errors.New("test-error")

var getStatement = `
		MATCH (context:Context) 
		OPTIONAL MATCH(context)-[:GRANTED_TO]->(principal:Resource)
		OPTIONAL MATCH(context)-[:PERMISSION]->(target:Resource)
		RETURN {context:Context, principals:COLLECT(principal), targets:COLLECT(target)}`

var getByIDStatement = `
		MATCH (context:Context) 
		WHERE context.id = $id
		OPTIONAL MATCH(context)-[:GRANTED_TO]->(principal:Resource)
		OPTIONAL MATCH(context)-[:PERMISSION]->(target:Resource)
		RETURN {context:Context, principals:COLLECT(principal), targets:COLLECT(target)}`

var createStatement = `
		CREATE ( context:Context{ name:$name, id:randomUUID() })
		RETURN { context:Context }`

var updateStatement = `
		MATCH (context:Context)
		WHERE context.id = $id
		SET context.name = $name
		RETURN { context:Context }`

var deleteStatement = `
		MATCH (context:Context)
		WHERE context.id = $id
		DETACH DELETE context`

func TestGetByID_SessionReturnsError_ErrorReturned(t *testing.T) {
	session := mocks.Session{
		ExecuteErr:        errTest,
		ExpectedStatement: getByIDStatement,
		ExpectedParameter: map[string]interface{}{
			"id": "test-id",
		},
		T: t,
	}
	repository := repository.New(session)
	_, err := repository.GetByID("test-id")
	assert.Equal(t, models.ErrNotFound, err)
}

func TestGetByID_SessionReturnsResponse_NoErrorsReturned(t *testing.T) {
	session := mocks.Session{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: getByIDStatement,
		ExpectedParameter: map[string]interface{}{
			"id": "test-id",
		},
		T: t,
	}
	repository := repository.New(session)
	response, err := repository.GetByID("test-id")
	assert.Equal(t, testdata.OutputDetails, response)
	assert.Nil(t, err)
}

func TestUpdate_SessionReturnsError_ErrorReturned(t *testing.T) {
	session := mocks.Session{
		ExecuteErr:        errors.New("test-error"),
		ExpectedStatement: updateStatement,
		ExpectedParameter: map[string]interface{}{
			"name": "test-context",
			"id":   "test-id",
		},
		T: t,
	}
	repository := repository.New(session)
	response, err := repository.Update("test-id", testdata.Input)
	assert.Equal(t, context.OutputDetails{}, response)
	assert.NotNil(t, err)
}

func TestUpdate_SessionReturnsEmptyResponse_DatabaseErrorReturned(t *testing.T) {
	session := mocks.Session{
		ExecuteResponse:   context.Output{Data: []context.Details{}},
		ExpectedStatement: updateStatement,
		ExpectedParameter: map[string]interface{}{
			"name": "test-context",
			"id":   "test-id",
		},
		T: t,
	}
	repository := repository.New(session)
	_, err := repository.Update("test-id", testdata.Input)
	assert.Equal(t, models.ErrDatabase, err)
}

func TestUpdate_SessionReturnsResponse_NoErrorReturned(t *testing.T) {
	session := mocks.Session{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: updateStatement,
		ExpectedParameter: map[string]interface{}{
			"name": "test-context",
			"id":   "test-id",
		},
		T: t,
	}
	repository := repository.New(session)
	response, err := repository.Update("test-id", testdata.Input)
	assert.Nil(t, err)
	assert.Equal(t, testdata.OutputDetails, response)
}

func TestDelete_SessionReturnsResponse_ReturnsNoErrors(t *testing.T) {
	session := mocks.Session{
		ExpectedStatement: deleteStatement,
		ExpectedParameter: map[string]interface{}{
			"id": "test-id",
		},
		T: t,
	}
	repository := repository.New(session)
	err := repository.Delete("test-id")
	assert.Nil(t, err)
}

func TestDelete_SessionReturnsError_ReturnsNoErrors(t *testing.T) {
	errTest := errors.New("test-error")
	session := mocks.Session{
		ExecuteErr:        errTest,
		ExpectedStatement: deleteStatement,
		ExpectedParameter: map[string]interface{}{
			"id": "test-id",
		},
		T: t,
	}
	repository := repository.New(session)
	err := repository.Delete("test-id")
	assert.Equal(t, errTest, err)
}
