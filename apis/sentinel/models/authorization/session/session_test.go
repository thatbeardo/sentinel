package session_test

import (
	"errors"
	"testing"

	"github.com/bithippie/guard-my-app/apis/sentinel/mocks"
	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	authorization "github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/session"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/testdata"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/injection"
	"github.com/stretchr/testify/assert"
)

type mockNeo4jSession struct {
	RunResponse        []map[string]interface{}
	RunErr             error
	ExpectedStatement  string
	ExpectedParameters map[string]interface{}
	t                  *testing.T
}

func (m mockNeo4jSession) Run(statement string, parameters map[string]interface{}) ([]map[string]interface{}, error) {
	assert.Equal(m.t, m.ExpectedStatement, statement)
	assert.Equal(m.t, m.ExpectedParameters, parameters)
	return m.RunResponse, m.RunErr
}

var testParameters = map[string]interface{}{"test-parameter": "test-parameter"}
var testStatement = "cypher-query"

func TestExecute_DBReturnsError_ReturnError(t *testing.T) {
	session := session.NewNeo4jSession(mockNeo4jSession{
		ExpectedStatement:  testStatement,
		ExpectedParameters: testParameters,
		RunErr:             models.ErrDatabase,
		t:                  t,
	})

	_, err := session.Execute("cypher-query", testParameters)
	assert.Equal(t, err, models.ErrDatabase)
}

func TestExecute_DBReturnsZeroLengthResults_ReturnError(t *testing.T) {
	session := session.NewNeo4jSession(mockNeo4jSession{
		ExpectedStatement:  testStatement,
		ExpectedParameters: testParameters,
		RunResponse:        []map[string]interface{}{},
		t:                  t,
	})

	_, err := session.Execute("cypher-query", testParameters)
	assert.Equal(t, err, models.ErrNotFound)
}

func TestExecute_DecodingTargetFails_ReturnError(t *testing.T) {
	session := session.NewNeo4jSession(mockNeo4jSession{
		ExpectedStatement:  testStatement,
		ExpectedParameters: testParameters,
		RunResponse:        generateMockResponse(),
		t:                  t,
	})

	defer injection.Reset()
	injection.NodeDecoder = func(map[string]interface{}, string, interface{}) error {
		return errors.New("some-decoding-error")
	}
	_, err := session.Execute("cypher-query", testParameters)
	assert.Equal(t, err, models.ErrDatabase)
}

func TestExecute_DecodingPermissionsFails_ReturnError(t *testing.T) {
	response := generateMockResponse()
	response[0]["permissions"] = "invalid-entry"
	session := session.NewNeo4jSession(mockNeo4jSession{
		ExpectedStatement:  testStatement,
		ExpectedParameters: testParameters,
		RunResponse:        response,
		t:                  t,
	})

	_, err := session.Execute("cypher-query", testParameters)
	assert.Equal(t, models.ErrDatabase, err)
}

func TestExecute_DBReturnsCleanData_ReturnOutput(t *testing.T) {
	session := session.NewNeo4jSession(mockNeo4jSession{
		ExpectedStatement:  testStatement,
		ExpectedParameters: testParameters,
		RunResponse:        generateMockResponse(),
		t:                  t,
	})

	output, err := session.Execute("cypher-query", testParameters)
	assert.Equal(t, testdata.Output, output)
	assert.Nil(t, err)
}

func TestExecute_DbReturnsDataWithMultiplePermissions_ReturnShortestPath(t *testing.T) {
	response := generateMultiLengthPermissions()
	expectedResponse := testdata.Output

	session := session.NewNeo4jSession(mockNeo4jSession{
		ExpectedStatement:  testStatement,
		ExpectedParameters: testParameters,
		RunResponse:        response,
		t:                  t,
	})

	output, err := session.Execute("cypher-query", testParameters)
	assert.Equal(t, expectedResponse, output)
	assert.Nil(t, err)
}

func TestExecute_DBReturnsCleanDataNoPermissions_ReturnOutput(t *testing.T) {
	response := generateMockResponse()
	expectedResponse := testdata.Output
	expectedResponse.Data[0] = authorization.Details{}

	response[0]["permissions"] = []interface{}{}
	session := session.NewNeo4jSession(mockNeo4jSession{
		ExpectedStatement:  testStatement,
		ExpectedParameters: testParameters,
		RunResponse:        response,
		t:                  t,
	})

	output, err := session.Execute("cypher-query", testParameters)
	assert.Equal(t, expectedResponse, output)
	assert.Nil(t, err)
}


func generateMockResponse() []map[string]interface{} {
	result := map[string]interface{}{
		"permissions": []interface{}{
			mocks.NewRelationship(1, 0, 0, "Permission",
				map[string]interface{}{
					"permitted": "allow",
					"name":      "test-permission",
				}),
		},
		"target": mocks.NewNode(1,
			[]string{"Resource"},
			map[string]interface{}{
				"id":        "test-target-id",
				"name":      "test-target",
				"source_id": "test-target-source-id",
			}),
	}
	return []map[string]interface{}{result}
}

func generateMultiLengthPermissions() []map[string]interface{} {
	shortestLengthResult := map[string]interface{}{
		"length": int64(2),
		"permissions": []interface{}{
			mocks.NewRelationship(1, 0, 0, "Permission",
				map[string]interface{}{
					"permitted": "allow",
					"name":      "test-permission",
				}),
		},
		"target": mocks.NewNode(1,
			[]string{"Resource"},
			map[string]interface{}{
				"id":        "test-target-id",
				"name":      "test-target",
				"source_id": "test-target-source-id",
			}),
	}

	longerLengthResult := map[string]interface{}{
		"length": int64(5),
		"permissions": []interface{}{
			mocks.NewRelationship(1, 0, 0, "Permission",
				map[string]interface{}{
					"permitted": "deny",
					"name":      "test-permission",
				}),
		},
		"target": mocks.NewNode(1,
			[]string{"Resource"},
			map[string]interface{}{
				"id":        "test-target-id",
				"name":      "test-target",
				"source_id": "test-target-source-id",
			}),
	}

	return []map[string]interface{}{longerLengthResult, shortestLengthResult}
}
