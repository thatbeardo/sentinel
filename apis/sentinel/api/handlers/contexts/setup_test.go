package contexts_test

import (
	handler "github.com/bithippie/guard-my-app/apis/sentinel/api/handlers"
	contexts "github.com/bithippie/guard-my-app/apis/sentinel/api/handlers/contexts"
	mocks "github.com/bithippie/guard-my-app/apis/sentinel/mocks/authorization"
	context "github.com/bithippie/guard-my-app/apis/sentinel/models/context/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/context/service"
	"github.com/bithippie/guard-my-app/apis/sentinel/testutil"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const noErrors = `{"data":{"type":"context","attributes":{"name":"valid-request"}}}`
const nameAbsentBadRequest = `{"data":{"type":"context","attributes":{}}}`

const typeAbsentBadRequest = `{"data":{"attributes":{"name":"valid-request"}}}`

type mockService struct {
	GetResponse     context.Output
	GetByIDResponse context.OutputDetails
	CreateResponse  context.OutputDetails
	UpdateResponse  context.OutputDetails
	Err             error
}

func (m mockService) Get() (context.Output, error) {
	return m.GetResponse, m.Err
}

func (m mockService) GetByID(string) (context.OutputDetails, error) {
	return m.GetByIDResponse, m.Err
}

func (m mockService) Create(*context.Input) (context.OutputDetails, error) {
	return m.CreateResponse, m.Err
}

func (m mockService) Update(string, *context.Input) (context.OutputDetails, error) {
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
	contexts.Routes(group, s, mocks.AuthorizationService{})
	return r
}
