package testdata

import (
	grant "github.com/bithippie/guard-my-app/apis/sentinel/models/grant/dto"
)

// Output represent mock data being sent back to the caller
var Output = grant.Output{
	Data: []grant.Details{Details},
}

// Input represents the data sent by the user as input
var Input = &grant.Input{
	Data: grant.InputDetails{
		Type: "grant",
		Attributes: &grant.Attributes{
			WithGrant: true,
		},
	},
}

// Details is a mock response of a single grant
var Details = grant.Details{
	InputDetails: grant.InputDetails{
		Type: "grant",
		Attributes: &grant.Attributes{
			WithGrant: true,
		},
	},
	Relationships: relationships,
	ID:            "test-grant-id",
}

// OutputDetails is a mock response of a single grant
var OutputDetails = grant.OutputDetails{
	Data: Details,
}

var relationships = grant.Relationships{
	Context: &grant.Relationship{
		Data: grant.Data{
			Type: "context",
			ID:   "context-id",
		},
	},
	Principal: &grant.Relationship{
		Data: grant.Data{
			Type: "resource",
			ID:   "resource-id",
		},
	},
}
