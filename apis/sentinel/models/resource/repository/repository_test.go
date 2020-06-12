package repository_test

import (
	"errors"
	"testing"

	resource "github.com/bithippie/guard-my-app/apis/sentinel/models/resource/dto"
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

var createStatement = `
		CREATE(child:Resource{name:$name, source_id: $source_id, id: randomUUID()})
		WITH child
		OPTIONAL MATCH(parent:Resource{id:$parent_id})
		WITH child,parent
		FOREACH (o IN CASE WHEN parent IS NOT NULL THEN [parent] ELSE [] END | CREATE (child)-[:OWNED_BY]->(parent))
		RETURN {child: child, parent: parent}`

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
	repository := repository.New(session)
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
	repository := repository.New(session)
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
	repository := repository.New(session)
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
	repository := repository.New(session)
	_, err := repository.GetByID("test-id")
	assert.Equal(t, errNotFound, err)
}

func TestCreate_SessionReturnsResponse_NoErrors(t *testing.T) {
	session := mockSession{
		ExecuteResponse:   testdata.Output,
		ExpectedStatement: createStatement,
		ExpectedParameter: map[string]interface{}{
			"name":      "test-resource",
			"source_id": "test-source-id",
			"parent_id": "parent-id",
		},
		t: t,
	}
	repository := repository.New(session)
	response, err := repository.Create(testdata.Input)
	assert.Equal(t, testdata.ModificationResult, response)
	assert.Nil(t, err)
}

func TestCreate_SessionReturnsError_ReturnsError(t *testing.T) {
	session := mockSession{
		ExecuteErr:        errNotFound,
		ExpectedStatement: createStatement,
		ExpectedParameter: map[string]interface{}{
			"name":      "test-resource",
			"source_id": "test-source-id",
			"parent_id": "parent-id",
		},
		t: t,
	}
	repository := repository.New(session)
	_, err := repository.Create(testdata.Input)
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
	repository := repository.New(session)
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
	repository := repository.New(session)
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
	repository := repository.New(session)
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
	repository := repository.New(session)
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
	repository := repository.New(session)
	err := repository.Delete("test-id")
	assert.Equal(t, errNotFound, err)
}
