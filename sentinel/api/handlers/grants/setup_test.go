package grants_test

import (
	handler "github.com/bithippie/guard-my-app/sentinel/api/handlers"
	"github.com/bithippie/guard-my-app/sentinel/api/handlers/grants"
	grant "github.com/bithippie/guard-my-app/sentinel/models/grant/dto"
	"github.com/bithippie/guard-my-app/sentinel/models/grant/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const noErrors = `{"data":{"type":"grant","attributes":{"with_grant":true}}}`
const typeFieldAbsent = `{"data":{"attributes":{"with_grant":true}}}`
const withGrantFieldAbsent = `{"data":{"type":"grant","attributes":{}}}`

type mockService struct {
	CreateResponse                           grant.OutputDetails
	GetPrincipalAndPolicyForResourceResponse grant.Output
	Err                                      error
}

func (m mockService) Create(*grant.Input, string, string) (grant.OutputDetails, error) {
	return m.CreateResponse, m.Err
}

func (m mockService) GetPrincipalAndPolicyForResource(id string) (grant.Output, error) {
	return m.GetPrincipalAndPolicyForResourceResponse, m.Err
}

func setupRouter(s service.Service) *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())
	r.NoRoute(handler.NoRoute)
	group := r.Group("/v1")
	grants.Routes(group, s)
	return r
}
