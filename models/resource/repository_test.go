package resource_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thatbeardo/go-sentinel/mocks/data"
	"github.com/thatbeardo/go-sentinel/models/resource"
)

var errTest = errors.New("test-error")
var errNotFound = errors.New("Data not found")

var getStatement = `
		MATCH(child:Resource) 
		OPTIONAL MATCH (child: Resource)-[:OWNED_BY]->(parent: Resource) 
		RETURN {child: child, parent: parent}`

var getByIDStatement = `
		MATCH(child:Resource) 
		WHERE child.id = $id 
		OPTIONAL MATCH (child: Resource)-[:OWNED_BY]->(parent: Resource) 
		RETURN {child: child, parent: parent}`

type mockSession struct {
	ExecuteResponse   resource.Response
	ExecuteErr        error
	ExpectedStatement string
	ExpectedParameter map[string]interface{}
	t                 *testing.T
}

func (m mockSession) Execute(statement string, parameters map[string]interface{}) (resource.Response, error) {
	assert.Equal(m.t, m.ExpectedStatement, statement)
	assert.Equal(m.t, m.ExpectedParameter, parameters)
	return m.ExecuteResponse, m.ExecuteErr
}

func TestGet_SessionReturnsResponse_NoErrors(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   data.Response,
		ExpectedStatement: getStatement,
		ExpectedParameter: map[string]interface{}{},
		t:                 t,
	}
	repository := resource.New(session)
	response, err := repository.Get()
	assert.Equal(t, data.Response, response)
	assert.Nil(t, err)
}

func TestGet_SessionReturnsError_ReturnsError(t *testing.T) {
	session := mockSession{
		ExecuteErr:        errTest,
		ExpectedStatement: getStatement,
		ExpectedParameter: map[string]interface{}{},
		t:                 t,
	}
	repository := resource.New(session)
	_, err := repository.Get()
	assert.Equal(t, errTest, err)
}

func TestGetByID_SessionReturnsResponse_NoErrors(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   data.Response,
		ExpectedStatement: getByIDStatement,
		ExpectedParameter: map[string]interface{}{
			"id": "test-id",
		},
		t: t,
	}
	repository := resource.New(session)
	response, err := repository.GetByID("test-id")
	assert.Equal(t, data.Response.Data[0], response)
	assert.Nil(t, err)
}

func TestGetByID_SessionReturnsError_ReturnsError(t *testing.T) {
	session := mockSession{
		ExecuteErr:        errNotFound,
		ExpectedStatement: getByIDStatement,
		ExpectedParameter: map[string]interface{}{
			"id": "test-id",
		},
		t: t,
	}
	repository := resource.New(session)
	_, err := repository.GetByID("test-id")
	assert.Equal(t, errNotFound, err)
}
