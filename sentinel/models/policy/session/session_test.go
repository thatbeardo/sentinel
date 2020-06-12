package session_test

import (
	"errors"
	"testing"

	"github.com/bithippie/guard-my-app/sentinel/mocks"
	models "github.com/bithippie/guard-my-app/sentinel/models"
	"github.com/bithippie/guard-my-app/sentinel/models/injection"
	policy "github.com/bithippie/guard-my-app/sentinel/models/policy/dto"
	"github.com/bithippie/guard-my-app/sentinel/models/policy/session"
	"github.com/bithippie/guard-my-app/sentinel/models/policy/testdata"
	"github.com/stretchr/testify/assert"
)

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

func TestExecute_DatabaseReturnsNoPolicies_EmptyResourcesArrayReturned(t *testing.T) {
	session := session.NewNeo4jSession(mockNeo4jSession{
		RunResponse: []map[string]interface{}{},
	})

	response, err := session.Execute(`cypher-query`, map[string]interface{}{})
	assert.Equal(t, policy.Output{Data: []policy.Details{}}, response)
	assert.Nil(t, err)
}

func TestExecute_DecodeFails_ReturnDatabaseError(t *testing.T) {
	defer injection.Reset()
	var errDecoding = errors.New("some-decoder-error")

	injection.NodeDecoder = func(map[string]interface{}, string, interface{}) error { return errDecoding }
	session := session.NewNeo4jSession(mockNeo4jSession{
		RunResponse: generateValidResultMap(),
	})

	_, err := session.Execute(`cypher-query`, map[string]interface{}{})
	assert.Equal(t, errDecoding, err)
}

func TestExecute_DecodeFailsDuringCast_ReturnDatabaseError(t *testing.T) {
	results := generateValidResultMap()
	results[0]["policy"] = "invalid-field"
	session := session.NewNeo4jSession(mockNeo4jSession{
		RunResponse: results,
	})

	_, err := session.Execute(`cypher-query`, map[string]interface{}{})
	assert.Equal(t, models.ErrDatabase, err)
}

func TestExecute_PrincipalResourcesDecodeFails_ReturnDatabaseError(t *testing.T) {
	results := generateValidResultMap()
	results[0]["principals"] = "invalid-field"
	session := session.NewNeo4jSession(mockNeo4jSession{
		RunResponse: results,
	})

	_, err := session.Execute(`cypher-query`, map[string]interface{}{})
	assert.Equal(t, models.ErrDatabase, err)
}

func TestExecute_TargetResourcesDecodeFails_ReturnDatabaseError(t *testing.T) {
	results := generateValidResultMap()
	results[0]["targets"] = "invalid-field"
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

	assert.Equal(t, testdata.Output, response)
	assert.Nil(t, err)
}

func generateValidResultMap() []map[string]interface{} {
	result := map[string]interface{}{
		"policy": mocks.NewNode(1,
			[]string{"Policy"},
			map[string]interface{}{
				"id":   "test-id",
				"name": "test-policy",
			}),
		"principals": []interface{}{
			mocks.NewNode(1,
				[]string{"Resource"},
				map[string]interface{}{
					"id":   "principal-resource-id",
					"name": "test-principal",
				}),
		},
		"targets": []interface{}{
			mocks.NewNode(1,
				[]string{"Resource"},
				map[string]interface{}{
					"id":   "target-resource-id",
					"name": "test-target",
				}),
		},
	}
	return []map[string]interface{}{result}
}
