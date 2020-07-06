package mocks

import (
	"testing"

	context "github.com/bithippie/guard-my-app/apis/sentinel/models/context/dto"
	"github.com/stretchr/testify/assert"
)

// Session represents a mocked session for context and resource tags testing
type Session struct {
	ExecuteResponse   context.Output
	ExecuteErr        error
	ExpectedStatement string
	ExpectedParameter map[string]interface{}
	T                 *testing.T
}

// Execute is the method run by default from within the session
func (m Session) Execute(statement string, parameters map[string]interface{}) (context.Output, error) {
	assert.Equal(m.T, m.ExpectedStatement, statement)
	assert.Equal(m.T, m.ExpectedParameter, parameters)
	return m.ExecuteResponse, m.ExecuteErr
}
