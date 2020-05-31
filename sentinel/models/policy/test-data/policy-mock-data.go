package testdata

import (
	"github.com/bithippie/guard-my-app/sentinel/models/policy/inputs"
	"github.com/bithippie/guard-my-app/sentinel/models/policy/outputs"
)

// Response represents the payload being sent back to the user
var Response = outputs.Response{
	Data: []outputs.Policy{policies},
}

var policies = outputs.Policy{
	PolicyDetails: inputs.PolicyDetails{
		Type: "policy",
		Attributes: &inputs.Attributes{
			Name: "test-policy",
		},
	},
	ID:            "test-id",
	Relationships: &outputs.Relationships{},
}

// Payload represents the data sent by the customer to the POST endpoint
var Payload = inputs.Payload{
	Data: &inputs.PolicyDetails{
		Type: "policy",
		Attributes: &inputs.Attributes{
			Name: "test-policy",
		},
	},
}
