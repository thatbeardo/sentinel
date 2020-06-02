package outputs

import "github.com/bithippie/guard-my-app/sentinel/models/policy/inputs"

// Response represents the final payload sent back to the user
type Response struct {
	Data []Policy `json:"data"`
}

// Policy represents data pertaining to one policy
type Policy struct {
	inputs.PolicyDetails
	ID            string         `json:"id"`
	Relationships *Relationships `json:"relationships,omitempty"`
}

// Relationships represent the relationships of an input payload
type Relationships struct {
	Principals      *PolicyDetails `json:"principals,omitempty" binding:"required,dive"`
	TargetResources *PolicyDetails `json:"target_resources,omitempty" binding:"required,dive"`
}

// PolicyDetails consists of array of resources associated with this policy
type PolicyDetails struct {
	Data []Resource `json:"data" binding:"required"`
}

// Resource denotes either the Principal or Target resource
type Resource struct {
	Type string `json:"type" binding:"required"`
	ID   string `json:"id" binding:"required"`
}
