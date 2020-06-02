package testdata

import (
	"github.com/bithippie/guard-my-app/sentinel/models/grant/inputs"
	"github.com/bithippie/guard-my-app/sentinel/models/grant/outputs"
)

// Response represent mock data being sent back to the caller
var Response = outputs.Response{
	Data: []outputs.Grant{grant},
}

var grant = outputs.Grant{
	GrantDetails: inputs.GrantDetails{
		Type: "grant",
		Attributes: &inputs.Attributes{
			WithGrant: true,
		},
	},
	ID: "test-grant-id",
}
