package resourcestest

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thatbeardo/go-sentinel/pkg/resource"
	"github.com/thatbeardo/go-sentinel/server"
	"github.com/thatbeardo/go-sentinel/testutil"
	views "github.com/thatbeardo/go-sentinel/views/responses/resources"
)

func TestResourcesOk(t *testing.T) {

	expectedBody := views.WrapGetResource(resource.Resource{Name: "Harshil", SourceID: "Mavain"})

	router := server.SetupRouter(testutil.GetMockService())
	w := testutil.PerformRequest(router, "GET", "/resources/")

	assert.Equal(t, http.StatusOK, w.Code)

	var response views.ResourceResponse
	err := json.Unmarshal([]byte(w.Body.String()), &response)

	assert.NotEmpty(t, response.Data)
	assert.Nil(t, err)
	assert.Equal(t, expectedBody, response)
}
