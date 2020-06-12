package permissions_test

import (
	handler "github.com/bithippie/guard-my-app/sentinel/api/handlers"
	"github.com/bithippie/guard-my-app/sentinel/api/handlers/permissions"
	permission "github.com/bithippie/guard-my-app/sentinel/models/permission/dto"
	"github.com/bithippie/guard-my-app/sentinel/models/permission/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const noErrors = `{"data":{"type":"permission","attributes":{"name":"resource:read","permitted":"allow"}}}`
const typeFieldAbsent = `{"data":{"attributes":{"name":"resource:read","permitted":"allow"}}}`
const nameFieldAbsent = `{"data":{"type":"permission","attributes":{"permitted":"allow"}}}`
const permittedFieldAbsent = `{"data":{"type":"permission","attributes":{"name":"resource:read"}}}`

type mockService struct {
	GetAllPermissionsForPolicyResponse             permission.Output
	GetAllPermissionsForPolicyWithResourceResponse permission.Output
	CreateResponse                                 permission.OutputDetails
	UpdateResponse                                 permission.OutputDetails
	Err                                            error
}

func (m mockService) GetAllPermissionsForPolicy(policyID string) (permission.Output, error) {
	return m.GetAllPermissionsForPolicyResponse, m.Err
}

func (m mockService) GetAllPermissionsForPolicyWithResource(policyID string, resourceID string) (permission.Output, error) {
	return m.GetAllPermissionsForPolicyWithResourceResponse, m.Err
}

func (m mockService) Create(*permission.Input, string, string) (permission.OutputDetails, error) {
	return m.CreateResponse, m.Err
}

func (m mockService) Update(string, *permission.Input) (permission.OutputDetails, error) {
	return m.UpdateResponse, m.Err
}

func (m mockService) Delete(string) error {
	return m.Err
}

func setupRouter(s service.Service) *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())
	r.NoRoute(handler.NoRoute)
	group := r.Group("/v1")
	permissions.Routes(group, s)
	return r
}
