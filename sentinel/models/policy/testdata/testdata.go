package testdata

import (
	policy "github.com/bithippie/guard-my-app/sentinel/models/policy/dto"
)

// Output represents the payload being sent back to the user
var Output = policy.Output{
	Data: []policy.Details{policies},
}

// OutputDetails consists data pertaining to one policy
var OutputDetails = policy.OutputDetails{
	Data: policies,
}

var policies = policy.Details{
	InputDetails: policy.InputDetails{
		Type: "policy",
		Attributes: &policy.Attributes{
			Name: "test-policy",
		},
	},
	ID: "test-id",
	Relationships: &policy.Relationships{
		Principals:      principalRelationship,
		TargetResources: targetRelationship,
	},
}

var principalRelationship = &policy.Relationship{
	Data: []policy.Resource{principal},
}

var targetRelationship = &policy.Relationship{
	Data: []policy.Resource{target},
}

var principal = policy.Resource{
	Type: "resource",
	ID:   "principal-resource-id",
}

var target = policy.Resource{
	Type: "resource",
	ID:   "target-resource-id",
}

// Input represents the data sent by the customer to the POST endpoint
var Input = &policy.Input{
	Data: &policy.InputDetails{
		Type: "policy",
		Attributes: &policy.Attributes{
			Name: "test-policy",
		},
	},
}
