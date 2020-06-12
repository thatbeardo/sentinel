package session_test

import (
	"errors"
	"testing"

	"github.com/bithippie/guard-my-app/sentinel/mocks"
	models "github.com/bithippie/guard-my-app/sentinel/models"
	authorization "github.com/bithippie/guard-my-app/sentinel/models/authorization/dto"
	"github.com/bithippie/guard-my-app/sentinel/models/authorization/session"
	"github.com/bithippie/guard-my-app/sentinel/models/authorization/testdata"
	"github.com/bithippie/guard-my-app/sentinel/models/injection"
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

	_, err := session.Execute("cypher-query", testParameters, authorization.Input{})
	assert.Equal(t, err, models.ErrDatabase)
}

func TestExecute_DBReturnsZeroLengthResults_ReturnError(t *testing.T) {
	session := session.NewNeo4jSession(mockNeo4jSession{
		ExpectedStatement:  testStatement,
		ExpectedParameters: testParameters,
		RunResponse:        []map[string]interface{}{},
		t:                  t,
	})

	_, err := session.Execute("cypher-query", testParameters, authorization.Input{})
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
	_, err := session.Execute("cypher-query", testParameters, authorization.Input{
		Permissions: []string{"test-permission-id"},
	})
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

	_, err := session.Execute("cypher-query", testParameters, authorization.Input{})
	assert.Equal(t, models.ErrDatabase, err)
}

func TestExecute_DBReturnsCleanData_ReturnOutput(t *testing.T) {
	session := session.NewNeo4jSession(mockNeo4jSession{
		ExpectedStatement:  testStatement,
		ExpectedParameters: testParameters,
		RunResponse:        generateMockResponse(),
		t:                  t,
	})

	output, err := session.Execute("cypher-query", testParameters, authorization.Input{})
	assert.Equal(t, testdata.Output, output)
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
