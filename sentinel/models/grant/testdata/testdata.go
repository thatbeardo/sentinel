package testdata

import (
	"github.com/bithippie/guard-my-app/sentinel/models/grant/inputs"
	"github.com/bithippie/guard-my-app/sentinel/models/grant/outputs"
)

// Response represent mock data being sent back to the caller
var Response = outputs.Response{
	Data: []outputs.Grant{grant},
}

// Payload represents the data sent by the user as input
var Payload = &inputs.Payload{
	Data: inputs.GrantDetails{
		Type: "grant",
		Attributes: &inputs.Attributes{
			WithGrant: true,
		},
	},
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
