package resources_test

import (
	handler "github.com/bithippie/guard-my-app/apis/sentinel/api/handlers"
	"github.com/bithippie/guard-my-app/apis/sentinel/api/handlers/resources"
	models "github.com/bithippie/guard-my-app/apis/sentinel/models"
	resource "github.com/bithippie/guard-my-app/apis/sentinel/models/resource/dto"
	"github.com/bithippie/guard-my-app/apis/sentinel/models/resource/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type mockService struct {
	GetByIDErrCandidate string
	GetResponse         resource.Output
	GetByIDResponse     resource.OutputDetails
	CreateResponse      resource.OutputDetails
	UpdateResponse      resource.OutputDetails

	GetErr     error
	GetByIDErr error
	CreateErr  error
	UpdateErr  error
	DeleteErr  error
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

func (m mockService) Create(*resource.Input) (resource.OutputDetails, error) {
	return m.CreateResponse, m.CreateErr
}

func (m mockService) Update(string, *resource.Input) (resource.OutputDetails, error) {
	return m.UpdateResponse, m.UpdateErr
}

func (m mockService) Delete(string) error {
	return m.DeleteErr
}

func setupRouter(s service.Service) *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())
	r.NoRoute(handler.NoRoute)
	group := r.Group("/v1")
	resources.Routes(group, s)
	return r
}
