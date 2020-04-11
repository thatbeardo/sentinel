package mocks

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// MockSession is used to run test cases with different behaviors.
type MockSession struct {
	SetResponse func() (neo4j.Result, error)
}

// Run executes the query passed to it
func (m MockSession) Run(cypher string, params map[string]interface{}, configurers ...func(*neo4j.TransactionConfig)) (neo4j.Result, error) {
	return m.SetResponse()
}

// Close terminates connection with the databse
func (m MockSession) Close() error {
	return nil
}

// GetMockResult returns
func GetMockResult() *Result {
	mockResult := &Result{}
	mockRecord := &Record{}
	mockRecord.On("GetByIndex", 0).Return("test-resource")
	mockRecord.On("GetByIndex", 1).Return("test-source-id")
	mockRecord.On("GetByIndex", 2).Return("test-id")
	mockResult.On("Next").Return(true).Once()
	mockResult.On("Next").Return(false).Once()
	mockResult.On("Record").Return(mockRecord)
	return mockResult
}
