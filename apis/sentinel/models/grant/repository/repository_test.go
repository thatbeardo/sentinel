package repository_test

import (
	"testing"

	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	grant "github.com/bithippie/guard-my-app/apis/sentinel/models/grant/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/grant/repository"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/grant/testdata"
	"github.com/stretchr/testify/assert"
)

var getAllContextsAndPrincipalsStatement = `
		MATCH (context:Context)-[:PERMISSION]->(principal: Resource {id: $principalID})
		OPTIONAL MATCH (context)-[grant:GRANTED_TO]->(principal:Resource)
		RETURN {grant: grant, context:Context, principal:principal}`

var createStatement = `
		MATCH (context:Context), (principal: Resource)
		WHERE context.id = $contextID AND principal.id = $principalID
		CREATE (context)-[grant:GRANTED_TO {with_grant: $withGrant, id: randomUUID()}]->(principal)
		RETURN {grant: grant, context:Context, principal: principal}`

type mockSession struct {
	ExecuteResponse   grant.Output
	ExecuteErr        error
	ExpectedStatement string
	ExpectedParameter map[string]interface{}
	t                 *testing.T
}

func (m mockSession) Execute(statement string, parameters map[string]interface{}) (grant.Output, error) {
	assert.Equal(m.t, m.ExpectedStatement, statement)
	assert.Equal(m.t, m.ExpectedParameter, parameters)
	return m.ExecuteResponse, m.ExecuteErr
}

func TestCreate_SessionReturnsResponse_ReturnDatabaseErr(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: createStatement,
		ExpectedParameter: map[string]interface{}{
			"withGrant":   true,
			"contextID":   "test-context-id",
			"principalID": "test-principal-id",
		},
		t: t,
	}
	repository := repository.New(session)
	response, err := repository.Create(testdata.Input, "test-context-id", "test-principal-id")
	assert.Equal(t, testdata.OutputDetails, response)
	assert.Nil(t, err)
}

func TestCreate_SessionReturnsEmptyResponse_DatabaseErrorReturned(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   grant.Output{Data: []grant.Details{}},
		ExpectedStatement: createStatement,
		ExpectedParameter: map[string]interface{}{
			"withGrant":   true,
			"contextID":   "test-context-id",
			"principalID": "test-principal-id",
		},
		t: t,
	}
	repository := repository.New(session)
	_, err := repository.Create(testdata.Input, "test-context-id", "test-principal-id")
	assert.Equal(t, models.ErrNotFound, err)
}

func TestGetAllContextsAndPrincipals_SessionReturnsResponse_ReturnDatabaseErr(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: getAllContextsAndPrincipalsStatement,
		ExpectedParameter: map[string]interface{}{
			"principalID": "test-principal-id",
		},
		t: t,
	}
	repository := repository.New(session)
	response, err := repository.GetPrincipalAndcontextForResource("test-principal-id")
	assert.Equal(t, testdata.Output, response)
	assert.Nil(t, err)
}

func TestGetAllContextsAndPrincipals_SessionReturnsEmptyResponse_DatabaseErrorReturned(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   grant.Output{Data: []grant.Details{}},
		ExpectedStatement: getAllContextsAndPrincipalsStatement,
		ExpectedParameter: map[string]interface{}{
			"principalID": "test-principal-id",
		},
		t: t,
	}
	repository := repository.New(session)
	result, err := repository.GetPrincipalAndcontextForResource("test-principal-id")
	assert.Equal(t, grant.Output{Data: []grant.Details{}}, result)
	assert.Nil(t, err)
}
