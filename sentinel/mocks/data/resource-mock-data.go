package data

import entity "github.com/bithippie/go-sentinel/sentinel/models/resource"

var policy = &entity.Identifier{
	Type: "policy",
	ID:   "policy-id",
}

var parent = &entity.Identifier{
	Type: "resource",
	ID:   "parent-id",
}

var parents = &entity.Parent{
	Data: parent,
}

var policies = &entity.Policies{
	Data: []*entity.Identifier{policy},
}

var inputRelationships = &entity.RelationshipsInput{
	Parent: parents,
}

var relationships = &entity.Relationships{
	Parent:   parents,
	Policies: policies,
}

var relationshipsWithoutParent = &entity.Relationships{
	Policies: policies,
}

var attributes = &entity.Resource{
	Name:     "test-resource",
	SourceID: "test-source-id",
}

var attributesWithDifferentName = &entity.Resource{
	Name:     "different-resource",
	SourceID: "different-source-id",
}

var parentAttributes = &entity.Resource{
	Name:     "parent-resource",
	SourceID: "parent-source-id",
}

var inputElement = &entity.InputElement{
	Type:          "resource",
	Attributes:    attributes,
	Relationships: inputRelationships,
}

var inputElementDifferentName = &entity.InputElement{
	Type:          "resource",
	Attributes:    attributesWithDifferentName,
	Relationships: inputRelationships,
}

var inputElementRelationshipsAbsent = &entity.InputElement{
	Type:       "resource",
	Attributes: attributes,
}

// Element represents response to be displayed to the user
var Element = entity.Element{
	Type:          "resource",
	Attributes:    attributes,
	Relationships: relationships,
	ID:            "test-id",
}

// ElementWithDifferentName represents response to be displayed to the user
var ElementWithDifferentName = entity.Element{
	Type:          "resource",
	Attributes:    attributesWithDifferentName,
	Relationships: relationships,
	ID:            "test-id",
}

// ElementWithoutParent represents standalone resource without any parents
var ElementWithoutParent = entity.Element{
	Type:          "resource",
	Attributes:    attributes,
	Relationships: relationshipsWithoutParent,
	ID:            "test-id",
}

// ElementRelationshipsAbsent represents response to be displayed to the user
var ElementRelationshipsAbsent = entity.Element{
	Type:       "resource",
	Attributes: attributes,
	ID:         "test-id",
}

// ElementPoliciesAbsent represents response to be displayed to the user
var ElementPoliciesAbsent = entity.Element{
	Type:       "resource",
	Attributes: attributes,
	Relationships: &entity.Relationships{
		Parent: parents,
	},
	ID: "test-id",
}

// ParentElement represents response to be displayed to the user
var ParentElement = entity.Element{
	Type:       "resource",
	Attributes: parentAttributes,
	ID:         "parent-id",
}

// Input represents the payload sent by the user to the service
var Input = &entity.Input{Data: inputElement}

// InputWithDifferentName represents an input with a different name
var InputWithDifferentName = &entity.Input{Data: inputElementDifferentName}

// InputRelationshipsAbsent represents the payload sent by the user to the service
var InputRelationshipsAbsent = &entity.Input{Data: inputElementRelationshipsAbsent}

// Response is the payload sent to the user
var Response = entity.Response{
	Data: []entity.Element{Element},
}

// ResponseWithoutPolicies has relationship without policies
var ResponseWithoutPolicies = entity.Response{
	Data: []entity.Element{ElementPoliciesAbsent},
}

// EmptyResponse denotes a case when no resources were found
var EmptyResponse = entity.Response{
	Data: []entity.Element{},
}

func generateResourceRelationship() entity.Relationships {
	policy := &entity.Identifier{Type: "policy", ID: "some-id"}
	policies := &entity.Policies{Data: []*entity.Identifier{policy}}

	parent := &entity.Parent{Data: &entity.Identifier{Type: "resource", ID: "some-id"}}
	relationships := entity.Relationships{Parent: parent, Policies: policies}
	return relationships
}
