package outputs

import "github.com/bithippie/guard-my-app/sentinel/models/permission/inputs"

// Response represents the data that has to be sent back to the use
type Response struct {
	Data []Permission
}

// Permission represents whether a policy has access to a resource
type Permission struct {
	inputs.PermissionDetails
	ID string `json:"id"`
}
