package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bithippie/guard-my-app/apis/sentinel/mocks"
	policy "github.com/bithippie/guard-my-app/apis/sentinel/mocks/policy"
	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	policyTestData "github.com/bithippie/guard-my-app/apis/sentinel/models/policy/testdata"
	resource "github.com/bithippie/guard-my-app/apis/sentinel/models/resource/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/injection"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/repository"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/testdata"
	"github.com/stretchr/testify/assert"
)

var errTest = errors.New("test-error")
var errNotFound = errors.New("Data not found")

var getStatement = `
		MATCH(child:Resource)
		OPTIONAL MATCH (child: Resource)-[:OWNED_BY]->(parent: Resource)
		OPTIONAL MATCH (policy: Policy)-[:GRANTED_TO]->(child: Resource)
		RETURN {child: child, parent: parent, policy: COLLECT(policy)}`

var getByIDStatement = `
		MATCH(child:Resource)
		WHERE child.id = $id
		OPTIONAL MATCH (child: Resource)-[:OWNED_BY]->(parent: Resource)
		OPTIONAL MATCH (policy: Policy)-[:GRANTED_TO]->(child: Resource)
		RETURN {child: child, parent: parent, policy: COLLECT(policy)}`

var createWithoutParentStatement = `
		MATCH(tenant:Resource{source_id: $tenant_id})<-[:GRANTED_TO]-(policy:Policy)
		WITH policy
		CREATE (policy)-[:PERMISSION]->(child:Resource{name:$name, source_id:$source_id, id:randomUUID()})
		RETURN {child:child}`

var createWithParentStatement = `
			MATCH (parent:Resource{id:$parent_id})
			WITH parent
			CREATE (child:Resource{name:$name, source_id:$source_id, id:randomUUID()})-[:OWNED_BY]->(parent)
			RETURN {child:child, parent:parent}`

var updateStatementNewParent = `
		MATCH(child:Resource{id:$child_id})
		SET child.name=$name
		SET child.source_id=$source_id
		WITH child
		OPTIONAL MATCH (child)-[old_relationship:OWNED_BY]->(old_parent:Resource)
		DETACH DELETE old_relationship
		WITH child

		OPTIONAL MATCH (new_parent:Resource{id:$new_parent_id})
		CREATE(child)-[:OWNED_BY]->(new_parent)
		RETURN {child: child, parent: new_parent}`

var updateStatementOldParent = `
		MATCH(child:Resource{id:$child_id})
		SET child.name=$name
		SET child.source_id=$source_id
		WITH child
		OPTIONAL MATCH (child)-[old_relationship:OWNED_BY]->(old_parent:Resource)RETURN {child: child, parent: old_parent}`

var associatePolicyStatement = `
		MATCH(principal:Resource{id: $principalID})
		CREATE (principal)<-[:GRANTED_TO]-(policy:Policy{ name:$name, id:randomUUID() })
		RETURN {policy: policy, principals: COLLECT(principal)}`

var getAassociatedPoliciesStatement = `
		MATCH (resource:Resource{id: $id}) 
		WITH resource
		MATCH(policy)-[:GRANTED_TO]->(resource)
		WITH policy
		OPTIONAL MATCH(policy)-[:PERMISSION]->(target:Resource)
		OPTIONAL MATCH(policy)-[:GRANTED_TO]->(principal:Resource)
		RETURN {policy:policy, principals:COLLECT(principal), targets:COLLECT(target)}`

var deleteStatement = `
		MATCH (n:Resource { id: $id }) DETACH DELETE n`

type mockSession struct {
	ExecuteResponse   resource.Output
	ExecuteErr        error
	ExpectedStatement string
	ExpectedParameter map[string]interface{}
	t                 *testing.T
}

func (m mockSession) Execute(statement string, parameters map[string]interface{}) (resource.Output, error) {
	assert.Equal(m.t, m.ExpectedStatement, statement)
	assert.Equal(m.t, m.ExpectedParameter, parameters)
	return m.ExecuteResponse, m.ExecuteErr
}

func TestGet_SessionReturnsResponse_NoErrors(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: getStatement,
		ExpectedParameter: map[string]interface{}{},
		t:                 t,
	}
	repository := repository.New(session, policy.Session{})
	response, err := repository.Get()
	assert.Equal(t, testdata.Output, response)
	assert.Nil(t, err)
}

func TestGet_SessionReturnsError_ReturnsError(t *testing.T) {
	session := mockSession{
		ExecuteErr:        errTest,
		ExpectedStatement: getStatement,
		ExpectedParameter: map[string]interface{}{},
		t:                 t,
	}
	repository := repository.New(session, policy.Session{})
	_, err := repository.Get()
	assert.Equal(t, errTest, err)
}

func TestGetByID_SessionReturnsResponse_NoErrors(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: getByIDStatement,
		ExpectedParameter: map[string]interface{}{
			"id": "test-id",
		},
		t: t,
	}
	repository := repository.New(session, policy.Session{})
	response, err := repository.GetByID("test-id")
	assert.Equal(t, testdata.ModificationResult, response)
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
	repository := repository.New(session, policy.Session{})
	_, err := repository.GetByID("test-id")
	assert.Equal(t, errNotFound, err)
}

func TestCreate_PayloadWithParentSessionReturnsResponse_NoErrors(t *testing.T) {
	defer injection.Reset()
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: createWithParentStatement,
		ExpectedParameter: map[string]interface{}{
			"name":      "test-resource",
			"source_id": "test-source-id",
			"parent_id": "parent-id",
			"tenant_id": "test-tenant",
		},
		t: t,
	}
	injection.ExtractClaims = func(ctx context.Context, claim string) string {
		return "test-tenant"
	}

	repository := repository.New(session, policy.Session{})
	response, err := repository.Create(mocks.Context{}, testdata.Input)
	assert.Equal(t, testdata.ModificationResult, response)
	assert.Nil(t, err)
}

func TestCreate_PayloadWithoutParentSessionReturnsResponse_NoErrors(t *testing.T) {
	defer injection.Reset()

	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: createWithoutParentStatement,
		ExpectedParameter: map[string]interface{}{
			"name":      "test-resource",
			"source_id": "test-source-id",
			"parent_id": "",
			"tenant_id": "test-tenant",
		},
		t: t,
	}
	injection.ExtractClaims = func(ctx context.Context, claim string) string {
		return "test-tenant"
	}

	repository := repository.New(session, policy.Session{})
	response, err := repository.Create(mocks.Context{}, testdata.InputWithoutParent)
	assert.Equal(t, testdata.ModificationResult, response)
	assert.Nil(t, err)
}

func TestCreate_SessionReturnsError_ReturnsError(t *testing.T) {
	defer injection.Reset()
	injection.ExtractClaims = func(ctx context.Context, claim string) string {
		return "test-tenant"
	}
	session := mockSession{
		ExecuteErr:        errNotFound,
		ExpectedStatement: createWithParentStatement,
		ExpectedParameter: map[string]interface{}{
			"name":      "test-resource",
			"source_id": "test-source-id",
			"parent_id": "parent-id",
			"tenant_id": "test-tenant",
		},
		t: t,
	}
	repository := repository.New(session, policy.Session{})
	_, err := repository.Create(mocks.Context{}, testdata.Input)
	assert.Equal(t, errNotFound, err)
}

func TestUpdate_NewParentProvidedSessionReturnsResponse_NoErrors(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: updateStatementNewParent,
		ExpectedParameter: map[string]interface{}{
			"child_id":      "test-id",
			"name":          "test-resource",
			"source_id":     "test-source-id",
			"new_parent_id": "parent-id",
		},
		t: t,
	}
	repository := repository.New(session, policy.Session{})
	response, err := repository.Update(testdata.Output.Data[0], testdata.Input)
	assert.Equal(t, testdata.ModificationResult, response)
	assert.Nil(t, err)
}

func TestUpdate_ParentAbsentSessionReturnsError_ReturnsError(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: updateStatementOldParent,
		ExpectedParameter: map[string]interface{}{
			"child_id":      "test-id",
			"name":          "test-resource",
			"source_id":     "test-source-id",
			"new_parent_id": "",
		},
		t: t,
	}
	repository := repository.New(session, policy.Session{})
	response, err := repository.Update(testdata.OutputWithoutParent.Data[0], testdata.InputWithoutRelationship)
	assert.Equal(t, testdata.ModificationResult, response)
	assert.Nil(t, err)
}

func TestUpdate_SessionReturnsError_ReturnsError(t *testing.T) {
	session := mockSession{
		ExecuteErr:        errNotFound,
		ExpectedStatement: updateStatementOldParent,
		ExpectedParameter: map[string]interface{}{
			"child_id":      "test-id",
			"name":          "test-resource",
			"source_id":     "test-source-id",
			"new_parent_id": "",
		},
		t: t,
	}
	repository := repository.New(session, policy.Session{})
	_, err := repository.Update(testdata.OutputWithoutParent.Data[0], testdata.InputWithoutRelationship)
	assert.Equal(t, errNotFound, err)
}

func TestDelete_SessionReturnsResponse_ReturnsNoErrors(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: deleteStatement,
		ExpectedParameter: map[string]interface{}{
			"id": "test-id",
		},
		t: t,
	}
	repository := repository.New(session, policy.Session{})
	err := repository.Delete("test-id")
	assert.Nil(t, err)
}

func TestDelete_SessionReturnsError_ReturnsNoErrors(t *testing.T) {
	session := mockSession{
		ExecuteErr:        errNotFound,
		ExpectedStatement: deleteStatement,
		ExpectedParameter: map[string]interface{}{
			"id": "test-id",
		},
		t: t,
	}
	repository := repository.New(session, policy.Session{})
	err := repository.Delete("test-id")
	assert.Equal(t, errNotFound, err)
}

func TestAssociatePolicy_SessionReturnsError_ReturnError(t *testing.T) {
	mockPolicySession := policy.Session{
		ExecuteErr:        errNotFound,
		ExecuteResponse:   policyTestData.Output,
		ExpectedStatement: associatePolicyStatement,
		ExpectedParameter: map[string]interface{}{
			"principalID": "test-id",
			"name":        "test-policy",
		},
		T: t,
	}

	repository := repository.New(mockSession{}, mockPolicySession)
	_, err := repository.AssociatePolicy("test-id", policyTestData.Input)
	assert.Equal(t, errNotFound, err)
}

func TestAssociatePolicy_SessionEmptyResults_ReturnError(t *testing.T) {
	mockPolicySession := policy.Session{
		ExpectedStatement: associatePolicyStatement,
		ExpectedParameter: map[string]interface{}{
			"principalID": "test-id",
			"name":        "test-policy",
		},
		T: t,
	}

	repository := repository.New(mockSession{}, mockPolicySession)
	_, err := repository.AssociatePolicy("test-id", policyTestData.Input)
	assert.Equal(t, models.ErrDatabase, err)
}

func TestGetAllAssociatedPolicies_SessionReturnsError_ReturnError(t *testing.T) {
	mockPolicySession := policy.Session{
		ExecuteErr:        models.ErrDatabase,
		ExpectedStatement: getAassociatedPoliciesStatement,
		ExpectedParameter: map[string]interface{}{
			"id": "test-id",
		},
		T: t,
	}

	repository := repository.New(mockSession{}, mockPolicySession)
	_, err := repository.GetAllAssociatedPolicies("test-id")
	assert.Equal(t, models.ErrDatabase, err)
}

func TestGetAllAssociatedPolicies_SessionReturnsData_ReturnData(t *testing.T) {
	mockPolicySession := policy.Session{
		ExecuteResponse:   policyTestData.Output,
		ExpectedStatement: getAassociatedPoliciesStatement,
		ExpectedParameter: map[string]interface{}{
			"id": "test-id",
		},
		T: t,
	}

	repository := repository.New(mockSession{}, mockPolicySession)
	response, err := repository.GetAllAssociatedPolicies("test-id")
	assert.Nil(t, err)
	assert.Equal(t, policyTestData.Output, response)
}
