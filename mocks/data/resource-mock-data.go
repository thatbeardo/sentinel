package data

import "github.com/thatbeardo/go-sentinel/models/resource"

var policy = &resource.Identifier{
	Type: "policy",
	ID:   "policy-id",
}

var parent = &resource.Identifier{
	Type: "resource",
	ID:   "parent-id",
}

var parents = &resource.Parent{
	Data: parent,
}

var policies = &resource.Policies{
	Data: []*resource.Identifier{policy},
}

var inputRelationships = &resource.RelationshipsInput{
	Parent: parents,
}

var relationships = &resource.Relationships{
	Parent:   parents,
	Policies: policies,
}

var relationshipsWithoutParent = &resource.Relationships{
	Policies: policies,
}

var attributes = &resource.Resource{
	Name:     "test-resource",
	SourceID: "test-source-id",
}

var attributesWithDifferentName = &resource.Resource{
	Name:     "different-resource",
	SourceID: "different-source-id",
}

var parentAttributes = &resource.Resource{
	Name:     "parent-resource",
	SourceID: "parent-source-id",
}

var inputElement = &resource.InputElement{
	Type:          "resource",
	Attributes:    attributes,
	Relationships: inputRelationships,
}

var inputElementDifferentName = &resource.InputElement{
	Type:          "resource",
	Attributes:    attributesWithDifferentName,
	Relationships: inputRelationships,
}

var inputElementRelationshipsAbsent = &resource.InputElement{
	Type:       "resource",
	Attributes: attributes,
}

// Element represents response to be displayed to the user
var Element = resource.Element{
	Type:          "resource",
	Attributes:    attributes,
	Relationships: relationships,
	ID:            "test-id",
}

// ElementWithDifferentName represents response to be displayed to the user
var ElementWithDifferentName = resource.Element{
	Type:          "resource",
	Attributes:    attributesWithDifferentName,
	Relationships: relationships,
	ID:            "test-id",
}

// ElementWithoutParent represents standalone resource without any parents
var ElementWithoutParent = resource.Element{
	Type:          "resource",
	Attributes:    attributes,
	Relationships: relationshipsWithoutParent,
	ID:            "test-id",
}

// ElementRelationshipsAbsent represents response to be displayed to the user
var ElementRelationshipsAbsent = resource.Element{
	Type:       "resource",
	Attributes: attributes,
	ID:         "test-id",
}

// ParentElement represents response to be displayed to the user
var ParentElement = resource.Element{
	Type:       "resource",
	Attributes: parentAttributes,
	ID:         "parent-id",
}

// Input represente the payload sent by the user to the service
var Input = &resource.Input{Data: inputElement}

// InputWithDifferentName represents an input with a different name
var InputWithDifferentName = &resource.Input{Data: inputElementDifferentName}

// InputRelationshipsAbsent represente the payload sent by the user to the service
var InputRelationshipsAbsent = &resource.Input{Data: inputElementRelationshipsAbsent}

// Response is the payload sent to the user
var Response = resource.Response{
	Data: []resource.Element{Element},
}

// EmptyResponse denotes a case when no resources were found
var EmptyResponse = resource.Response{
	Data: []resource.Element{},
}

func generateResourceRelationship() resource.Relationships {
	policy := &resource.Identifier{Type: "policy", ID: "some-id"}
	policies := &resource.Policies{Data: []*resource.Identifier{policy}}

	parent := &resource.Parent{Data: &resource.Identifier{Type: "resource", ID: "some-id"}}
	relationships := resource.Relationships{Parent: parent, Policies: policies}
	return relationships
}
