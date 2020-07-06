package testdata

import (
	context "github.com/bithippie/guard-my-app/apis/sentinel/models/context/dto"
)

// Output represents the payload being sent back to the user
var Output = context.Output{
	Data: []context.Details{contexts},
}

// OutputDetails consists data pertaining to one context
var OutputDetails = context.OutputDetails{
	Data: contexts,
}

var contexts = context.Details{
	InputDetails: context.InputDetails{
		Type: "context",
		Attributes: &context.Attributes{
			Name: "test-context",
		},
	},
	ID: "test-id",
	Relationships: &context.Relationships{
		Principals:      principalRelationship,
		TargetResources: targetRelationship,
	},
}

var principalRelationship = &context.Relationship{
	Data: []context.Resource{principal},
}

var targetRelationship = &context.Relationship{
	Data: []context.Resource{target},
}

var principal = context.Resource{
	Type: "resource",
	ID:   "principal-resource-id",
}

var target = context.Resource{
	Type: "resource",
	ID:   "target-resource-id",
}

// Input represents the data sent by the customer to the POST endpoint
var Input = &context.Input{
	Data: &context.InputDetails{
		Type: "context",
		Attributes: &context.Attributes{
			Name: "test-context",
		},
	},
}
