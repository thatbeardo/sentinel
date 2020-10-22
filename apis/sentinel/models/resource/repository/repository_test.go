package repository_test

import (
	"errors"
	"testing"

	context "github.com/bithippie/guard-my-app/apis/sentinel/mocks/context"
	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	contextTestData "github.com/bithippie/guard-my-app/apis/sentinel/models/context/testdata"
	resource "github.com/bithippie/guard-my-app/apis/sentinel/models/resource/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/repository"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/testdata"
	"github.com/stretchr/testify/assert"
)

var errTest = errors.New("test-error")
var errNotFound = errors.New("Data not found")

var getStatement = `
		MATCH(client:Resource{source_id:$client_id})
		<-[:GRANTED_TO]-(:Context)-[:PERMISSION{name:"sentinel:read"}]->
		(tenant:Resource{source_id: $tenant})<-[:GRANTED_TO]-(:Context)-[:PERMISSION{name:"sentinel:read"}]->(hub:Resource)
		<-[:OWNED_BY*0..]-(child:Resource)
		OPTIONAL MATCH (child)-[:OWNED_BY]->(parent: Resource)
		OPTIONAL MATCH (context:Context)-[:GRANTED_TO]->(child)
		RETURN {child: child, parent: parent, context: COLLECT(context)}`

var getByIDStatement = `
		MATCH(child:Resource{id: $id})
		OPTIONAL MATCH (child)-[:OWNED_BY]->(parent: Resource)
		OPTIONAL MATCH (context:Context)-[:GRANTED_TO]->(child)
		RETURN {child: child, parent: parent, context: COLLECT(context)}`

var getResourcesHub = `
		MATCH(client:Resource{source_id: $client_id})<-[:GRANTED_TO]-(:Context)
		-[:PERMISSION{name:"sentinel:read", permitted: "allow"}]->
		(tenant:Resource{source_id: $tenant})<-[:GRANTED_TO]-(:Context)
		-[:PERMISSION{name:"sentinel:read", permitted: "allow"}]->(hub:Resource)
		RETURN {child: hub}`

var createWithoutParentStatement = `
		MATCH(tenant:Resource{source_id: $tenant_id})<-[:GRANTED_TO]-(context:Context)
		WITH context
		CREATE (context)-[:PERMISSION]->(child:Resource{name:$name, source_id:$source_id, context_id:$context_id, id:randomUUID()})
		RETURN {child:child}`

var createWithParentStatement = `
		MATCH (parent:Resource{id:$parent_id})
		WITH parent
		CREATE (child:Resource{name:$name, source_id:$source_id, context_id:$context_id, id:randomUUID()})-[:OWNED_BY]->(parent)
		RETURN {child:child, parent:parent}`

var createTenantStatement = `
		CREATE (child:Resource{name:$name, source_id:$source_id, context_id:$context_id, id: randomUUID()})
		RETURN {child:child}`

var updateStatementNewParent = `
		MATCH(child:Resource{id:$child_id})
		SET child.name=$name SET child.context_id=$context_id SET child.source_id=$source_id 
		WITH child
		OPTIONAL MATCH (child)-[old_relationship:OWNED_BY]->(old_parent:Resource)
		DETACH DELETE old_relationship
		WITH child

		OPTIONAL MATCH (new_parent:Resource{id:$new_parent_id})
		CREATE(child)-[:OWNED_BY]->(new_parent)
		RETURN {child: child, parent: new_parent}`

var updateStatementOldParent = `
		MATCH(child:Resource{id:$child_id})
		SET child.name=$name SET child.context_id=$context_id SET child.source_id=$source_id 
		WITH child
		OPTIONAL MATCH (child)-[old_relationship:OWNED_BY]->(old_parent:Resource)RETURN {child: child, parent: old_parent}`

var updateStatementOldParentContextAndNameUpdateAbsent = `
		MATCH(child:Resource{id:$child_id})
		SET child.source_id=$source_id 
		WITH child
		OPTIONAL MATCH (child)-[old_relationship:OWNED_BY]->(old_parent:Resource)RETURN {child: child, parent: old_parent}`

var addContextStatement = `
		MATCH(principal:Resource{id: $principalID})
		CREATE (principal)<-[:GRANTED_TO]-(context:Context{ name:$name, id:randomUUID() })
		RETURN {context:context, principals: COLLECT(principal)}`

var getAassociatedContextsStatement = `
		MATCH (resource:Resource{id: $id}) 
		WITH resource
		MATCH(context)-[:GRANTED_TO]->(resource)
		WITH context
		OPTIONAL MATCH(context)-[:PERMISSION]->(target:Resource)
		OPTIONAL MATCH(context)-[:GRANTED_TO]->(principal:Resource)
		RETURN {context:context, principals:COLLECT(principal), targets:COLLECT(target)}`

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
		ExpectedParameter: map[string]interface{}{
			"tenant":    "test-tenant",
			"client_id": "client-id",
		},
		t: t,
	}
	repository := repository.New(session, context.Session{})
	response, err := repository.Get("client-id", "test-tenant")

	assert.Equal(t, testdata.Output, response)
	assert.Nil(t, err)
}

func TestGet_SessionReturnsError_ReturnsError(t *testing.T) {

	session := mockSession{
		ExecuteErr:        errTest,
		ExpectedStatement: getStatement,
		ExpectedParameter: map[string]interface{}{
			"tenant":    "test-tenant",
			"client_id": "client-id",
		},
		t: t,
	}
	repository := repository.New(session, context.Session{})

	_, err := repository.Get("client-id", "test-tenant")
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
	repository := repository.New(session, context.Session{})
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
	repository := repository.New(session, context.Session{})
	_, err := repository.GetByID("test-id")
	assert.Equal(t, errNotFound, err)
}

func TestCreate_CreateTenantResourceEmptyResponse_NoErrors(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   resource.Output{Data: []resource.Details{}},
		ExpectedStatement: createTenantStatement,
		ExpectedParameter: map[string]interface{}{
			"name":       "test-resource",
			"source_id":  "test-source-id",
			"context_id": "test-context-id",
		},
		t: t,
	}

	repository := repository.New(session, context.Session{})
	_, err := repository.CreateMetaResource(testdata.InputWithoutParent)
	assert.Equal(t, models.ErrNotFound, err)
}

func TestCreateTenant_SessionReturnsResponse_NoErrors(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: createTenantStatement,
		ExpectedParameter: map[string]interface{}{
			"name":       "test-resource",
			"source_id":  "test-source-id",
			"context_id": "test-context-id",
		},
		t: t,
	}

	repository := repository.New(session, context.Session{})
	response, err := repository.CreateMetaResource(testdata.InputWithoutParent)
	assert.Equal(t, testdata.ModificationResult, response)
	assert.Nil(t, err)
}

func TestAttachResourceToExistingParent_SessionReturnsError_ReturnsError(t *testing.T) {
	session := mockSession{
		ExecuteErr:        errNotFound,
		ExpectedStatement: createWithParentStatement,
		ExpectedParameter: map[string]interface{}{
			"name":       "test-resource",
			"source_id":  "test-source-id",
			"parent_id":  "parent-id",
			"context_id": "test-context-id",
		},
		t: t,
	}
	repository := repository.New(session, context.Session{})

	_, err := repository.Create(testdata.Input)
	assert.Equal(t, errNotFound, err)
}

func TestAttachResourceToExistingParent_SessionReturnsData_ReturnsError(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: createWithParentStatement,
		ExpectedParameter: map[string]interface{}{
			"name":       "test-resource",
			"source_id":  "test-source-id",
			"parent_id":  "parent-id",
			"context_id": "test-context-id",
		},
		t: t,
	}
	repository := repository.New(session, context.Session{})

	response, _ := repository.Create(testdata.Input)
	assert.Equal(t, testdata.ModificationResult, response)
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
			"context_id":    "test-context-id",
		},
		t: t,
	}
	repository := repository.New(session, context.Session{})
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
			"context_id":    "test-context-id",
		},
		t: t,
	}
	repository := repository.New(session, context.Session{})
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
			"context_id":    "test-context-id",
		},
		t: t,
	}
	repository := repository.New(session, context.Session{})
	_, err := repository.Update(testdata.OutputWithoutParent.Data[0], testdata.InputWithoutRelationship)
	assert.Equal(t, errNotFound, err)
}

func TestUpdate_InputWithout_ReturnsError(t *testing.T) {
	session := mockSession{
		ExecuteErr:        errNotFound,
		ExpectedStatement: updateStatementOldParentContextAndNameUpdateAbsent,
		ExpectedParameter: map[string]interface{}{
			"child_id":      "test-id",
			"name":          "",
			"source_id":     "test-source-id",
			"new_parent_id": "",
			"context_id":    "",
		},
		t: t,
	}
	repository := repository.New(session, context.Session{})
	_, err := repository.Update(testdata.OutputWithoutParent.Data[0], testdata.InputWithoutContext)
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
	repository := repository.New(session, context.Session{})
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
	repository := repository.New(session, context.Session{})
	err := repository.Delete("test-id")
	assert.Equal(t, errNotFound, err)
}

func TestAssociatecontext_SessionReturnsError_ReturnError(t *testing.T) {
	mockContextSession := context.Session{
		ExecuteErr:        errNotFound,
		ExecuteResponse:   contextTestData.Output,
		ExpectedStatement: addContextStatement,
		ExpectedParameter: map[string]interface{}{
			"principalID": "test-id",
			"name":        "test-context",
		},
		T: t,
	}

	repository := repository.New(mockSession{}, mockContextSession)
	_, err := repository.AddContext("test-id", contextTestData.Input)
	assert.Equal(t, errNotFound, err)
}

func TestAssociatecontext_SessionEmptyResults_ReturnError(t *testing.T) {
	mockContextSession := context.Session{
		ExpectedStatement: addContextStatement,
		ExpectedParameter: map[string]interface{}{
			"principalID": "test-id",
			"name":        "test-context",
		},
		T: t,
	}

	repository := repository.New(mockSession{}, mockContextSession)
	_, err := repository.AddContext("test-id", contextTestData.Input)
	assert.Equal(t, models.ErrDatabase, err)
}

func TestGetAllAssociatedContexts_SessionReturnsError_ReturnError(t *testing.T) {
	mockContextSession := context.Session{
		ExecuteErr:        models.ErrDatabase,
		ExpectedStatement: getAassociatedContextsStatement,
		ExpectedParameter: map[string]interface{}{
			"id": "test-id",
		},
		T: t,
	}

	repository := repository.New(mockSession{}, mockContextSession)
	_, err := repository.GetAllContexts("test-id")
	assert.Equal(t, models.ErrDatabase, err)
}

func TestGetAllAssociatedContexts_SessionReturnsData_ReturnData(t *testing.T) {
	mockContextSession := context.Session{
		ExecuteResponse:   contextTestData.Output,
		ExpectedStatement: getAassociatedContextsStatement,
		ExpectedParameter: map[string]interface{}{
			"id": "test-id",
		},
		T: t,
	}

	repository := repository.New(mockSession{}, mockContextSession)
	response, err := repository.GetAllContexts("test-id")
	assert.Nil(t, err)
	assert.Equal(t, contextTestData.Output, response)
}

func TestGetResourcesHub_SessionReturnsEmptyDataArray_ReturnError(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   resource.Output{Data: []resource.Details{}},
		ExpectedStatement: getResourcesHub,
		ExpectedParameter: map[string]interface{}{
			"client_id": "client-id",
			"tenant":    "tenant",
		},
		t: t,
	}

	repository := repository.New(session, context.Session{})
	_, err := repository.GetResourcesHub("client-id", "tenant")
	assert.Equal(t, models.ErrNotFound, err)
}

func TestGetResourcesHub_SessionReturnsError_ReturnError(t *testing.T) {
	session := mockSession{
		ExecuteErr:        models.ErrDatabase,
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: getResourcesHub,
		ExpectedParameter: map[string]interface{}{
			"client_id": "client-id",
			"tenant":    "tenant",
		},
		t: t,
	}

	repository := repository.New(session, context.Session{})
	_, err := repository.GetResourcesHub("client-id", "tenant")
	assert.Equal(t, models.ErrDatabase, err)
}

func TestGetResourcesHub_SessionReturnsData_ReturnData(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: getResourcesHub,
		ExpectedParameter: map[string]interface{}{
			"client_id": "client-id",
			"tenant":    "tenant",
		},
		t: t,
	}

	repository := repository.New(session, context.Session{})
	response, _ := repository.GetResourcesHub("client-id", "tenant")
	assert.Equal(t, testdata.OutputDetails, response)
}
