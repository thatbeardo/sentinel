package grants_test

import (
	handler "github.com/bithippie/guard-my-app/apis/sentinel/api/handlers"
	"github.com/bithippie/guard-my-app/apis/sentinel/api/handlers/grants"
	mocks "github.com/bithippie/guard-my-app/apis/sentinel/mocks/authorization"
	grant "github.com/bithippie/guard-my-app/apis/sentinel/models/grant/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/grant/service"
	"github.com/bithippie/guard-my-app/apis/sentinel/testutil"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const noErrors = `{"data":{"type":"grant","attributes":{"with_grant":true}}}`
const typeFieldAbsent = `{"data":{"attributes":{"with_grant":true}}}`
const withGrantFieldAbsent = `{"data":{"type":"grant","attributes":{}}}`

type mockService struct {
	CreateResponse                            grant.OutputDetails
	GetPrincipalAndcontextForResourceResponse grant.Output
	GrantExistsResponse                       bool
	Err                                       error
}

func (m mockService) Create(*grant.Input, string, string) (grant.OutputDetails, error) {
	return m.CreateResponse, m.Err
}

func (m mockService) GetPrincipalAndcontextForResource(id string) (grant.Output, error) {
	return m.GetPrincipalAndcontextForResourceResponse, m.Err
}

func (m mockService) GrantExists(contextID, principalID string) (bool, error) {
	return m.GrantExistsResponse, m.Err
}

func setupRouter(s service.Service) *gin.Engine {
	testutil.RemoveMiddleware()
	r := gin.Default()
	r.Use(cors.Default())
	r.NoRoute(handler.NoRoute)
	group := r.Group("/v1")
	grants.Routes(group, s, mocks.AuthorizationService{})
	return r
}
