package resources_test

import (
	"context"

	handler "github.com/bithippie/guard-my-app/apis/sentinel/api/handlers"
	"github.com/bithippie/guard-my-app/apis/sentinel/api/handlers/resources"
	"github.com/bithippie/guard-my-app/apis/sentinel/api/handlers/resources/injection"
	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	policy "github.com/bithippie/guard-my-app/apis/sentinel/models/policy/dto"
	resource "github.com/bithippie/guard-my-app/apis/sentinel/models/resource/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var noMiddleware = func(c *gin.Context) {}

type mockService struct {
	GetByIDErrCandidate              string
	GetResponse                      resource.Output
	GetByIDResponse                  resource.OutputDetails
	CreateResponse                   resource.OutputDetails
	UpdateResponse                   resource.OutputDetails
	AssociateResponse                policy.OutputDetails
	GetAllAssociatedPoliciesResponse policy.Output

	GetErr       error
	GetByIDErr   error
	CreateErr    error
	UpdateErr    error
	DeleteErr    error
	AssociateErr error
}

func (m mockService) Get() (resource.Output, error) {
	return m.GetResponse, m.GetErr
}

func (m mockService) GetByID(id string) (resource.OutputDetails, error) {
	if m.GetByIDErrCandidate == id {
		return m.GetByIDResponse, models.ErrNotFound
	}
	return m.GetByIDResponse, m.GetByIDErr
}

func (m mockService) Create(context.Context, *resource.Input) (resource.OutputDetails, error) {
	return m.CreateResponse, m.CreateErr
}

func (m mockService) Update(string, *resource.Input) (resource.OutputDetails, error) {
	return m.UpdateResponse, m.UpdateErr
}

func (m mockService) AssociatePolicy(string, *policy.Input) (policy.OutputDetails, error) {
	return m.AssociateResponse, m.AssociateErr
}

func (m mockService) GetAllAssociatedPolicies(string) (policy.Output, error) {
	return m.GetAllAssociatedPoliciesResponse, m.AssociateErr
}

func (m mockService) Delete(string) error {
	return m.DeleteErr
}

func setupRouter(s service.Service) *gin.Engine {
	injection.VerifyResourceOwnership = noMiddleware
	injection.ValidateNewResource = noMiddleware
	r := gin.Default()
	r.Use(cors.Default())
	r.NoRoute(handler.NoRoute)
	group := r.Group("/v1")
	resources.Routes(group, s)
	return r
}
