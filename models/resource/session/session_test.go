package session_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thatbeardo/go-sentinel/mocks"
	"github.com/thatbeardo/go-sentinel/mocks/data"
	models "github.com/thatbeardo/go-sentinel/models"
	"github.com/thatbeardo/go-sentinel/models/resource/injection"
	"github.com/thatbeardo/go-sentinel/models/resource/session"
)

type mockNode struct {
	id     int64
	labels []string
	props  map[string]interface{}
}

func (m mockNode) Id() int64 {
	return m.id
}

func (m mockNode) Labels() []string {
	return m.labels
}

func (m mockNode) Props() map[string]interface{} {
	return m.props
}

func TestExecute_RunReturnsError_ReturnDatabaseError(t *testing.T) {
	mockSession := &mocks.Session{}
	session := session.NewNeo4jSession(mockSession)

	mockSession.On("Run", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(nil, errors.New("Database error"))
	_, err := session.Execute(`cypher-query`, map[string]interface{}{})

	assert.Equal(t, models.ErrDatabase, err)
}

func TestExecute_DecodeFails_ReturnDatabaseError(t *testing.T) {
	defer injection.Reset()
	var errDecoding = errors.New("some-decoder-error")

	injection.MapDecoder = func(interface{}, interface{}) error { return errDecoding }
	mockSession := &mocks.Session{}
	session := session.NewNeo4jSession(mockSession)

	resultMap := generateValidResultMap()
	mockSession.On("Run", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(returnResultWithData(resultMap), nil)
	_, err := session.Execute(`cypher-query`, map[string]interface{}{})

	assert.Equal(t, errDecoding, err)
}

func TestExecute_NoErrorsFromDB_ReturnResponse(t *testing.T) {
	mockSession := &mocks.Session{}
	session := session.NewNeo4jSession(mockSession)

	mockSession.On("Run", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(returnResultWithData(generateValidResultMap()), nil)
	response, err := session.Execute(`cypher-query`, map[string]interface{}{})

	assert.Equal(t, data.ResponseWithoutPolicies, response)
	assert.Nil(t, err)
}

func TestExecute_ParentDecodeFails_ReturnError(t *testing.T) {
	mockSession := &mocks.Session{}
	session := session.NewNeo4jSession(mockSession)

	resultMap := generateValidResultMap()
	resultMap["parent"] = "invalid entry"
	mockSession.On("Run", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).Return(returnResultWithData(resultMap), nil)
	_, err := session.Execute(`cypher-query`, map[string]interface{}{})

	assert.Equal(t, models.ErrDatabase, err)

}

func generateValidResultMap() map[string]interface{} {
	return map[string]interface{}{
		"child": mockNode{
			id:     1,
			labels: []string{"Resource"},
			props: map[string]interface{}{
				"id":        "test-id",
				"name":      "test-resource",
				"source_id": "test-source-id",
			},
		},
		"parent": mockNode{
			id:     2,
			labels: []string{"Resource"},
			props: map[string]interface{}{
				"id":        "parent-id",
				"name":      "parent",
				"source_id": "parent-source-id",
			},
		},
	}
}

func returnResultWithData(data interface{}) *mocks.Result {
	mockResult := &mocks.Result{}
	mockResult.On("Next").Return(true).Once()
	mockResult.On("Next").Return(false).Once()
	mockRecord := &mocks.Record{}
	mockRecord.On("GetByIndex", 0).Return(data)
	mockResult.On("Record").Return(mockRecord)
	return mockResult
}
