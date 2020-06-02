package session_test

import (
	"errors"
	"testing"

	models "github.com/bithippie/guard-my-app/sentinel/models"
	"github.com/bithippie/guard-my-app/sentinel/models/grant/outputs"
	"github.com/bithippie/guard-my-app/sentinel/models/grant/session"
	"github.com/bithippie/guard-my-app/sentinel/models/grant/testdata"
	"github.com/bithippie/guard-my-app/sentinel/models/injection"
	"github.com/stretchr/testify/assert"
)

type mockRelationship struct {
	id      int64
	startId int64
	endId   int64
	relType string
	props   map[string]interface{}
}

func (m mockRelationship) Id() int64 {
	return m.id
}

func (m mockRelationship) StartId() int64 {
	return m.startId
}

func (m mockRelationship) EndId() int64 {
	return m.endId
}

func (m mockRelationship) Type() string {
	return m.relType
}

func (m mockRelationship) Props() map[string]interface{} {
	return m.props
}

type mockNeo4jSession struct {
	RunResponse []map[string]interface{}
	RunErr      error
}

func (m mockNeo4jSession) Run(statement string, parameters map[string]interface{}) ([]map[string]interface{}, error) {
	return m.RunResponse, m.RunErr
}

func TestExecute_RunReturnsError_ReturnDatabaseError(t *testing.T) {
	session := session.NewNeo4jSession(mockNeo4jSession{
		RunErr: errors.New("Database error"),
	})

	_, err := session.Execute(`cypher-query`, map[string]interface{}{})
	assert.Equal(t, models.ErrDatabase, err)
}

func TestExecute_DatabaseReturnsNoGrants_EmptyResourcesArrayReturned(t *testing.T) {
	session := session.NewNeo4jSession(mockNeo4jSession{
		RunResponse: []map[string]interface{}{},
	})

	response, err := session.Execute(`cypher-query`, map[string]interface{}{})
	assert.Equal(t, outputs.Response{Data: []outputs.Grant{}}, response)
	assert.Nil(t, err)
}

func TestExecute_DecodeFails_ReturnDatabaseError(t *testing.T) {
	defer injection.Reset()
	var errDecoding = errors.New("some-decoder-error")

	injection.MapDecoder = func(interface{}, interface{}) error { return errDecoding }
	session := session.NewNeo4jSession(mockNeo4jSession{
		RunResponse: generateValidResultMap(),
	})

	_, err := session.Execute(`cypher-query`, map[string]interface{}{})
	assert.Equal(t, errDecoding, err)
}

func TestExecute_DecodeFailsDuringCast_ReturnDatabaseError(t *testing.T) {
	results := generateValidResultMap()
	results[0]["grant"] = "invalid-field"
	session := session.NewNeo4jSession(mockNeo4jSession{
		RunResponse: results,
	})

	_, err := session.Execute(`cypher-query`, map[string]interface{}{})
	assert.Equal(t, models.ErrDatabase, err)
}

func TestExecute_NoErrorsFromDB_ReturnResponse(t *testing.T) {
	session := session.NewNeo4jSession(mockNeo4jSession{
		RunResponse: generateValidResultMap(),
	})

	response, err := session.Execute(`cypher-query`, map[string]interface{}{})

	assert.Equal(t, testdata.Response, response)
	assert.Nil(t, err)
}

func generateValidResultMap() []map[string]interface{} {
	result := map[string]interface{}{
		"grant": mockRelationship{
			id:      1,
			relType: "Grant",
			props: map[string]interface{}{
				"id":         "test-grant-id",
				"with_grant": true,
			},
		},
	}
	return []map[string]interface{}{result}
}
