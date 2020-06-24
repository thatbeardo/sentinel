package policies_test

import (
	handler "github.com/bithippie/guard-my-app/apis/sentinel/api/handlers"
	"github.com/bithippie/guard-my-app/apis/sentinel/api/handlers/policies"
	policy "github.com/bithippie/guard-my-app/apis/sentinel/models/policy/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/policy/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const noErrors = `{"data":{"type":"policy","attributes":{"name":"valid-request"}}}`
const nameAbsentBadRequest = `{"data":{"type":"policy","attributes":{}}}`

const typeAbsentBadRequest = `{"data":{"attributes":{"name":"valid-request"}}}`

type mockService struct {
	GetResponse     policy.Output
	GetByIDResponse policy.OutputDetails
	CreateResponse  policy.OutputDetails
	UpdateResponse  policy.OutputDetails
	Err             error
}

func (m mockService) Get() (policy.Output, error) {
	return m.GetResponse, m.Err
}

func (m mockService) GetByID(string) (policy.OutputDetails, error) {
	return m.GetByIDResponse, m.Err
}

func (m mockService) Create(*policy.Input) (policy.OutputDetails, error) {
	return m.CreateResponse, m.Err
}

func (m mockService) Update(string, *policy.Input) (policy.OutputDetails, error) {
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
	policies.Routes(group, s)
	return r
}
