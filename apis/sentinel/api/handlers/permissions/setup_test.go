package permissions_test

import (
	handler "github.com/bithippie/guard-my-app/apis/sentinel/api/handlers"
	"github.com/bithippie/guard-my-app/apis/sentinel/api/handlers/permissions"
	mocks "github.com/bithippie/guard-my-app/apis/sentinel/mocks/authorization"
	permission "github.com/bithippie/guard-my-app/apis/sentinel/models/permission/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/permission/service"
	"github.com/bithippie/guard-my-app/apis/sentinel/testutil"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const noErrors = `{"data":{"type":"permission","attributes":{"name":"resource:read","permitted":"allow"}}}`
const typeFieldAbsent = `{"data":{"attributes":{"name":"resource:read","permitted":"allow"}}}`
const nameFieldAbsent = `{"data":{"type":"permission","attributes":{"permitted":"allow"}}}`
const permittedFieldAbsent = `{"data":{"type":"permission","attributes":{"name":"resource:read"}}}`

type mockService struct {
	GetAllPermissionsForcontextResponse             permission.Output
	GetAllPermissionsForcontextWithResourceResponse permission.Output
	CreateResponse                                  permission.OutputDetails
	UpdateResponse                                  permission.OutputDetails
	IsPermissionIdempotentResponse                  bool
	Err                                             error
}

func (m mockService) GetAllPermissionsForcontext(contextID string) (permission.Output, error) {
	return m.GetAllPermissionsForcontextResponse, m.Err
}

func (m mockService) GetAllPermissionsForcontextWithResource(contextID string, resourceID string) (permission.Output, error) {
	return m.GetAllPermissionsForcontextWithResourceResponse, m.Err
}

func (m mockService) IsPermissionIdempotent(input *permission.Input, contextID, resourceID string) (bool, error) {
	return m.IsPermissionIdempotentResponse, m.Err
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
	testutil.RemoveMiddleware()
	r := gin.Default()
	r.Use(cors.Default())
	r.NoRoute(handler.NoRoute)
	group := r.Group("/v1")
	permissions.Routes(group, s, mocks.AuthorizationService{})
	return r
}
