package resources_test

import (
	"context"

	handler "github.com/bithippie/guard-my-app/apis/sentinel/api/handlers"
	"github.com/bithippie/guard-my-app/apis/sentinel/api/handlers/resources"

	mocks "github.com/bithippie/guard-my-app/apis/sentinel/mocks/authorization"

	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	contextDto "github.com/bithippie/guard-my-app/apis/sentinel/models/context/dto"
	resource "github.com/bithippie/guard-my-app/apis/sentinel/models/resource/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/service"
	"github.com/bithippie/guard-my-app/apis/sentinel/testutil"
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
	AssociateResponse                contextDto.OutputDetails
	GetAllAssociatedContextsResponse contextDto.Output

	GetErr       error
	GetByIDErr   error
	CreateErr    error
	UpdateErr    error
	DeleteErr    error
	AssociateErr error
}

func (m mockService) Get(context.Context) (resource.Output, error) {
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

func (m mockService) Associatecontext(string, *contextDto.Input) (contextDto.OutputDetails, error) {
	return m.AssociateResponse, m.AssociateErr
}

func (m mockService) GetAllAssociatedContexts(string) (contextDto.Output, error) {
	return m.GetAllAssociatedContextsResponse, m.AssociateErr
}

func (m mockService) Delete(string) error {
	return m.DeleteErr
}

func setupRouter(s service.Service) *gin.Engine {

	testutil.RemoveMiddleware()
	r := gin.Default()
	r.Use(cors.Default())
	r.NoRoute(handler.NoRoute)
	group := r.Group("/v1")
	resources.Routes(group, s, mocks.AuthorizationService{})
	return r
}
