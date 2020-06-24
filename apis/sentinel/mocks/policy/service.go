package mocks

import policy "github.com/bithippie/guard-my-app/apis/sentinel/models/policy/dto"

// Service represents the mock policy service used by resources tag and policy tags
type Service struct {
	GetResponse     policy.Output
	GetByIDResponse policy.OutputDetails
	CreateResponse  policy.OutputDetails
	UpdateResponse  policy.OutputDetails
	Err             error
}

// Get method is used to fetch all resources
func (m Service) Get() (policy.Output, error) {
	return m.GetResponse, m.Err
}

// GetByID fetches a particular resource by ID
func (m Service) GetByID(string) (policy.OutputDetails, error) {
	return m.GetByIDResponse, m.Err
}

// Create is used to add a new resource to the database
func (m Service) Create(*policy.Input) (policy.OutputDetails, error) {
	return m.CreateResponse, m.Err
}

// Update is used to edit fields on an existing resource
func (m Service) Update(string, *policy.Input) (policy.OutputDetails, error) {
	return m.UpdateResponse, m.Err
}

// Delete method removes a reosource from the database
func (m Service) Delete(string) error {
	return m.Err
}
