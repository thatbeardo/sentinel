package tenant_test

import (
	"context"
	"testing"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/tenant"
	"github.com/stretchr/testify/assert"
)

func TestExtract_ContextWithValueSetPassedToExtract_ContextShouldReturnValue(t *testing.T) {
	context := context.Background()
	updatedContext := tenant.Add(context, "test-tenant")
	assert.Equal(t, "test-tenant", tenant.Extract(updatedContext))
}

func TestAdd_NewValuePassedToContext_ContextShouldBeUpdated(t *testing.T) {
	context := context.Background()
	updatedContext := tenant.Add(context, "test-tenant")
	assert.Equal(t, "test-tenant", tenant.Extract(updatedContext))
}

func TestExtract_NoValueSetOnContext_ReturnsBlankString(t *testing.T) {
	context := context.Background()
	assert.Equal(t, "", tenant.Extract(context))
}
