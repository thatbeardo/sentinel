package outputs

import "github.com/bithippie/guard-my-app/sentinel/models/grant/inputs"

// Response represents the data that has to be sent back to the use
type Response struct {
	Data []Grant
}

// Grant represents whether a policy has access to a resource
type Grant struct {
	inputs.GrantDetails
	ID string `json:"id"`
}

// Relationships enlists details about the Policy and the resource that are attached with this grant
type Relationships struct {
	Policy    *GrantDetails `json:"policy" binding:"required"`
	Principal *GrantDetails `json:"principal" binding:"required"`
}

// GrantDetails provides details about the Policy and Principal
type GrantDetails struct {
	Type string `json:"type" binding:"required"`
	ID   string `json:"id" binding:"required"`
}
