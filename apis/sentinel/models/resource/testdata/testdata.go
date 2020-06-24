package testdata

import resource "github.com/bithippie/guard-my-app/apis/sentinel/models/resource/dto"

var attributes = &resource.Attributes{
	Name:     "test-resource",
	SourceID: "test-source-id",
}

var policy = resource.Data{
	Type: "policy",
	ID:   "policy-id",
}

var parent = &resource.Parent{
	Data: &resource.Data{
		Type: "resource",
		ID:   "parent-id",
	},
}

var policies = &resource.Policies{
	Data: []resource.Data{policy},
}

var detailsWithoutPolicies = resource.Details{
	ID:         "test-id",
	Type:       "resource",
	Attributes: attributes,
	Relationships: &resource.Relationships{
		Parent: parent,
	},
}

var detailsWithoutParent = resource.Details{
	ID:         "test-id",
	Type:       "resource",
	Attributes: attributes,
	Relationships: &resource.Relationships{
		Policies: policies,
	},
}

var detailsWithPolicies = resource.Details{
	ID:         "test-id",
	Type:       "resource",
	Attributes: attributes,
	Relationships: &resource.Relationships{
		Parent:   parent,
		Policies: policies,
	},
}

// Output represent a condition when all fields are present in the payload being sent back
var Output = resource.Output{
	Data: []resource.Details{detailsWithPolicies},
}

// OutputWithoutPolicies represents mock data for a response without policies
var OutputWithoutPolicies = resource.Output{
	Data: []resource.Details{detailsWithoutPolicies},
}

// OutputWithoutParent represents mock data for a response without parent
var OutputWithoutParent = resource.Output{
	Data: []resource.Details{detailsWithoutParent},
}

// EmptyOutput represents a condition when no resources are present in the database
var EmptyOutput = resource.Output{
	Data: []resource.Details{},
}

// Input is a mock for the data being sent by the user
var Input = &resource.Input{
	Data: &resource.InputDetails{
		Type:       "resource",
		Attributes: attributes,
		Relationships: &resource.InputRelationships{
			Parent: parent,
		},
	},
}

// InputWithoutParent is a mock for the data being sent by the user
var InputWithoutParent = &resource.Input{
	Data: &resource.InputDetails{
		Type:          "resource",
		Attributes:    attributes,
		Relationships: nil,
	},
}

// InputWithoutRelationship is a mock for the data being sent by the user
var InputWithoutRelationship = &resource.Input{
	Data: &resource.InputDetails{
		Type:       "resource",
		Attributes: attributes,
	},
}

// ModificationResult represents the response sent when a resource is Updated/Created
var ModificationResult = resource.OutputDetails{
	Data: detailsWithPolicies,
}
