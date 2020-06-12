package authorization

import (
	permission "github.com/bithippie/guard-my-app/apis/sentinel/models/permission/dto"
	resource "github.com/bithippie/guard-my-app/apis/sentinel/models/resource/dto"
)

// Output is the response schema for the authorization endpoint
type Output struct {
	Data []Details `json:"data" binding:"required"`
}

// Details provides information about entities this principal has access to
type Details struct {
	ID            string              `json:"id"`
	Type          string              `json:"type" enums:"resource"`
	Attributes    resource.Attributes `json:"attributes"`
	Relationships Relationships       `json:"relationships"`
}

// Relationships represents permissions details for this principal
type Relationships struct {
	Permissions Permissions `json:"permissions"`
}

// Permissions represent all permissions access to this principal
type Permissions struct {
	Data []permission.Attributes `json:"data"`
}
