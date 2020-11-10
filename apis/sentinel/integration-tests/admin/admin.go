package admin

import (
	"encoding/json"
	"github.com/bithippie/guard-my-app/apis/sentinel/integration-tests/service"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/context/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/login/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/permission/dto"
	"net/http"
	"strings"
)

const url = `http://localhost:8080`

// ScaffoldNewClient provisions resources and permissions for a new sentinel client
func ScaffoldNewClient(clientName, environmentName string) {
	accessToken := generateAccessToken()
	resourceID := createClientResource(clientName, "client-source", accessToken)
	contextID := createClientContext("client-context", resourceID, accessToken)
	tenantID := createTenantResource("tenant", environmentName, accessToken)
	createTenantPermissions(contextID, tenantID, accessToken)
	
	tenantContextID := createTenantContext("tenant-context", tenantID, accessToken)
	hubID := createResourceHub("hub", "hub-source", accessToken)
	createHubPermissions(tenantContextID, hubID, accessToken)
}

func generateAccessToken() (accessToken string) {
	body := strings.NewReader(loginPayload)
	request := service.PrepareRequest(
		"POST",
		url,
		"/v1/login",
		"dummy-token",
		body,
	)
	response := service.ExecuteRequest(request)
	defer response.Body.Close()

	var credentials login.BearerToken
	data := service.ParseResponse(response, http.StatusOK)
	json.Unmarshal(data, &credentials)
	return credentials.AccessToken
}

func createClientResource(name, sourceID, accessToken string) (resourceID string) {
	return createResource(name, sourceID, accessToken)
}

func createClientContext(name, resourceID, accessToken string) (contextID string) {
	return createContext(name, resourceID, accessToken)
}

func createTenantResource(name, sourceID, accessToken string) (tenantID string) {
	return createResource(name, sourceID, accessToken)
}

func createTenantPermissions(contextID, tenantID, accessToken string) (permissionID string) {
	return createPermission(contextID, tenantID, accessToken)
}

func createTenantContext(name, tenantID, accessToken string) (tenantContext string) {
	return createContext(name, tenantID, accessToken)
}

func createResourceHub(name, sourceID, accessToken string) (resourceID string) {
	return createResource(name, sourceID, accessToken)
}

func createHubPermissions(contextID, hubID, accessToken string) (permissionID string) {
	return createPermission(contextID, hubID, accessToken)
}

func createResource(name, sourceID, accessToken string) (resourceID string) {
	body := strings.NewReader(createResourcePayload(name, sourceID))
	request := service.PrepareRequest(
		"POST",
		url,
		"/v1/resources?claimant=integration",
		accessToken,
		body,
	)
	response := service.ExecuteRequest(request)
	defer response.Body.Close()

	var resource resource.OutputDetails
	data := service.ParseResponse(response, http.StatusAccepted)
	json.Unmarshal(data, &resource)
	return resource.Data.ID
}

func createContext(contextName, resourceID, accessToken string) (tenantContext string) {
	body := strings.NewReader(createContextPayload(contextName))
	request := service.PrepareRequest(
		"POST",
		url,
		"/v1/resources/"+resourceID+"/contexts?claimant=integration",
		accessToken,
		body)
	response := service.ExecuteRequest(request)
	defer response.Body.Close()

	var context context.OutputDetails
	data := service.ParseResponse(response, http.StatusAccepted)
	json.Unmarshal(data, &context)
	return context.Data.ID
}

func createPermission(contextID, resourceID, accessToken string) (permissionID string) {
	body := strings.NewReader(createPermissionPayload)
	request := service.PrepareRequest(
		"PUT",
		url,
		"/v1/permissions/"+contextID+"/resources/"+resourceID+"?claimant=integration",
		accessToken,
		body)
	response := service.ExecuteRequest(request)
	defer response.Body.Close()

	var permission permission.OutputDetails
	data := service.ParseResponse(response, http.StatusAccepted)
	json.Unmarshal(data, &permission)
	return permission.Data.ID
}