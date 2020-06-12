package repository_test

import (
	"testing"

	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	grant "github.com/bithippie/guard-my-app/apis/sentinel/models/grant/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/grant/repository"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/grant/testdata"
	"github.com/stretchr/testify/assert"
)

var getAllPoliciesAndPrincipalsStatement = `
		MATCH (policy:Policy)-[:PERMISSION]->(principal: Resource {id: $principalID})
		OPTIONAL MATCH (policy)-[grant:GRANTED_TO]->(principal:Resource)
		RETURN {grant: grant, policy:policy, principal:principal}`

var createStatement = `
		MATCH (policy: Policy), (principal: Resource)
		WHERE policy.id = $policyID AND principal.id = $principalID
		CREATE (policy)-[grant:GRANTED_TO {with_grant: $withGrant, id: randomUUID()}]->(principal)
		RETURN {grant: grant, policy: policy, principal: principal}`

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
			"policyID":    "test-policy-id",
			"principalID": "test-principal-id",
		},
		t: t,
	}
	repository := repository.New(session)
	response, err := repository.Create(testdata.Input, "test-policy-id", "test-principal-id")
	assert.Equal(t, testdata.OutputDetails, response)
	assert.Nil(t, err)
}

func TestCreate_SessionReturnsEmptyResponse_DatabaseErrorReturned(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   grant.Output{Data: []grant.Details{}},
		ExpectedStatement: createStatement,
		ExpectedParameter: map[string]interface{}{
			"withGrant":   true,
			"policyID":    "test-policy-id",
			"principalID": "test-principal-id",
		},
		t: t,
	}
	repository := repository.New(session)
	_, err := repository.Create(testdata.Input, "test-policy-id", "test-principal-id")
	assert.Equal(t, models.ErrNotFound, err)
}

func TestGetAllPoliciesAndPrincipals_SessionReturnsResponse_ReturnDatabaseErr(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: getAllPoliciesAndPrincipalsStatement,
		ExpectedParameter: map[string]interface{}{
			"principalID": "test-principal-id",
		},
		t: t,
	}
	repository := repository.New(session)
	response, err := repository.GetPrincipalAndPolicyForResource("test-principal-id")
	assert.Equal(t, testdata.Output, response)
	assert.Nil(t, err)
}

func TestGetAllPoliciesAndPrincipals_SessionReturnsEmptyResponse_DatabaseErrorReturned(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   grant.Output{Data: []grant.Details{}},
		ExpectedStatement: getAllPoliciesAndPrincipalsStatement,
		ExpectedParameter: map[string]interface{}{
			"principalID": "test-principal-id",
		},
		t: t,
	}
	repository := repository.New(session)
	result, err := repository.GetPrincipalAndPolicyForResource("test-principal-id")
	assert.Equal(t, grant.Output{Data: []grant.Details{}}, result)
	assert.Nil(t, err)
}
