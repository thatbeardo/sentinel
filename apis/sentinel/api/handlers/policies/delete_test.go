package policies_test

import (
	"net/http"
	"testing"

	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	"github.com/bithippie/guard-my-app/apis/sentinel/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDelete_PolicyDeleted_ReturnStatusOK(t *testing.T) {
	mockService := mockService{}
	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "DELETE", "/v1/policies/test-id", "")
	defer cleanup()

	assert.Equal(t, http.StatusNoContent, response.StatusCode)
}

func TestDeleteResourcesServiceError(t *testing.T) {
	mockService := mockService{
		Err: models.ErrDatabase,
	}

	router := setupRouter(mockService)
	response, cleanup := testutil.PerformRequest(router, "DELETE", "/v1/policies/test-id", "")
	defer cleanup()

	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}
