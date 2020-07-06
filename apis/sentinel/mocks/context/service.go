package mocks

import context "github.com/bithippie/guard-my-app/apis/sentinel/models/context/dto"

// Service represents the mock context service used by resources tag and context tags
type Service struct {
	GetResponse     context.Output
	GetByIDResponse context.OutputDetails
	CreateResponse  context.OutputDetails
	UpdateResponse  context.OutputDetails
	Err             error
}

// Get method is used to fetch all resources
func (m Service) Get() (context.Output, error) {
	return m.GetResponse, m.Err
}

// GetByID fetches a particular resource by ID
func (m Service) GetByID(string) (context.OutputDetails, error) {
	return m.GetByIDResponse, m.Err
}

// Create is used to add a new resource to the database
func (m Service) Create(*context.Input) (context.OutputDetails, error) {
	return m.CreateResponse, m.Err
}

// Update is used to edit fields on an existing resource
func (m Service) Update(string, *context.Input) (context.OutputDetails, error) {
	return m.UpdateResponse, m.Err
}

// Delete method removes a reosource from the database
func (m Service) Delete(string) error {
	return m.Err
}
