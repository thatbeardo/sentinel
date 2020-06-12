package authorizations_test

import (
	"testing"

	handler "github.com/bithippie/guard-my-app/sentinel/api/handlers"
	authorizations "github.com/bithippie/guard-my-app/sentinel/api/handlers/authorization"
	authorization "github.com/bithippie/guard-my-app/sentinel/models/authorization/dto"
	"github.com/bithippie/guard-my-app/sentinel/models/authorization/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const noErrors = `{"data":{"type":"grant","attributes":{"with_grant":true}}}`
const typeFieldAbsent = `{"data":{"attributes":{"with_grant":true}}}`
const withGrantFieldAbsent = `{"data":{"type":"grant","attributes":{}}}`

type mockService struct {
	GetAuthorizationForPrincipalResponse authorization.Output
	ExpectedInput                        authorization.Input
	Err                                  error
	t                                    *testing.T
}

func (m mockService) GetAuthorizationForPrincipal(principalID string, input authorization.Input) (output authorization.Output, err error) {
	assert.Equal(m.t, m.ExpectedInput, input)
	return m.GetAuthorizationForPrincipalResponse, m.Err
}

func setupRouter(s service.Service) *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())
	r.NoRoute(handler.NoRoute)
	group := r.Group("/v1")
	authorizations.Routes(group, s)
	return r
}
