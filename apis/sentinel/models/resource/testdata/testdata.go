package testdata

import resource "github.com/bithippie/guard-my-app/apis/sentinel/models/resource/dto"

var attributes = &resource.Attributes{
	Name:      "test-resource",
	SourceID:  "test-source-id",
	ContextID: "test-context-id",
}

var attributesWithoutContext = &resource.Attributes{
	SourceID: "test-source-id",
}

var context = resource.Data{
	Type: "context",
	ID:   "context-id",
}

var parent = &resource.Parent{
	Data: &resource.Data{
		Type: "resource",
		ID:   "parent-id",
	},
}

var contexts = &resource.Contexts{
	Data: []resource.Data{context},
}

var detailsWithoutContexts = resource.Details{
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
		Contexts: contexts,
	},
}

var detailsWithContexts = resource.Details{
	ID:         "test-id",
	Type:       "resource",
	Attributes: attributes,
	Relationships: &resource.Relationships{
		Parent:   parent,
		Contexts: contexts,
	},
}

// Output represent a condition when all fields are present in the payload being sent back
var Output = resource.Output{
	Data: []resource.Details{detailsWithContexts},
}

// OutputWithoutContexts represents mock data for a response without contexts
var OutputWithoutContexts = resource.Output{
	Data: []resource.Details{detailsWithoutContexts},
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

// InputWithoutContext is a payload where context is absent
var InputWithoutContext = &resource.Input{
	Data: &resource.InputDetails{
		Type:       "resource",
		Attributes: attributesWithoutContext,
	},
}

// ModificationResult represents the response sent when a resource is Updated/Created
var ModificationResult = resource.OutputDetails{
	Data: detailsWithContexts,
}

// OutputDetails is the same as ModificationResult just with a different name
var OutputDetails = resource.OutputDetails{
	Data: detailsWithContexts,
}

// OutputDetailsWithoutContexts is used when testing create function of resource service
var OutputDetailsWithoutContexts = resource.OutputDetails{
	Data: detailsWithoutContexts,
}

// ResourcesHubOutputDetails is the ancestor of all resources
var ResourcesHubOutputDetails = resource.OutputDetails{
	Data: resource.Details{
		ID:   "parent-id",
		Type: "resource",
	},
}
