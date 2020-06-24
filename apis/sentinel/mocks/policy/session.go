package mocks

import (
	"testing"

	policy "github.com/bithippie/guard-my-app/apis/sentinel/models/policy/dto"
	"github.com/stretchr/testify/assert"
)

// Session represents a mocked session for policy and resource tags testing
type Session struct {
	ExecuteResponse   policy.Output
	ExecuteErr        error
	ExpectedStatement string
	ExpectedParameter map[string]interface{}
	T                 *testing.T
}

// Execute is the method run by default from within the session
func (m Session) Execute(statement string, parameters map[string]interface{}) (policy.Output, error) {
	assert.Equal(m.T, m.ExpectedStatement, statement)
	assert.Equal(m.T, m.ExpectedParameter, parameters)
	return m.ExecuteResponse, m.ExecuteErr
}
