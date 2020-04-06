package resourcestest

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thatbeardo/go-sentinel/models/resource"
	"github.com/thatbeardo/go-sentinel/server"
	"github.com/thatbeardo/go-sentinel/testutil"
	views "github.com/thatbeardo/go-sentinel/views/dto/resources"
)

func TestGetResourcesOk(t *testing.T) {

	expectedBody := views.WrapGetResources([]*resource.Resource{{Name: "Harshil", SourceID: "Mavain"}})

	router := server.SetupRouter(testutil.GetMockService())
	w := testutil.PerformRequest(router, "GET", "/resources/")

	assert.Equal(t, http.StatusOK, w.Code)

	var response views.ResourceResponse
	err := json.Unmarshal([]byte(w.Body.String()), &response)

	assert.NotEmpty(t, response.Data)
	assert.Nil(t, err)
	assert.Equal(t, expectedBody, response)
}
