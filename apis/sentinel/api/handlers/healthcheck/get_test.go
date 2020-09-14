package healthcheck_test

import (
	"net/http"
	"testing"

	"github.com/bithippie/guard-my-app/apis/sentinel/testutil"
)

func TestHealthCheckEndpoint_PingHealthCheck_ReturnsStatsOK(t *testing.T) {
	router := setupRouter()
	response, cleanup := testutil.PerformRequest(router, "GET", "/v1/healthcheck", "")
	defer cleanup()

	testutil.ValidateResponse(t, response, "Sentinel online", http.StatusOK)
}
