package testdata

import (
	"github.com/bithippie/guard-my-app/sentinel/models/permission/inputs"
	"github.com/bithippie/guard-my-app/sentinel/models/permission/outputs"
)

// Response denotes mocked data that is sent back in response body
var Response = outputs.Response{
	Data: []outputs.Permission{permission},
}

var permission = outputs.Permission{
	ID: "test-id",
	PermissionDetails: inputs.PermissionDetails{
		Type: "permission",
		Attributes: &inputs.Attributes{
			Name:      "test-permission",
			Permitted: "allow",
		},
	},
}
