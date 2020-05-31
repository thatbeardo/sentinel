package inputs

// Payload is the payload that a POST endpoint expects.
type Payload struct {
	Data *Element `json:"data" binding:"required,dive"`
}

// Element is the paylaod sent when creating a new resource
type Element struct {
	Type          string         `json:"type" binding:"required"`
	Attributes    *Attributes    `json:"attributes" binding:"required,dive"`
	Relationships *Relationships `json:"relationships,omitempty"`
}

// Attributes defines the name of the policy
type Attributes struct {
	Name string `json:"name" binding:"required"`
}

// Relationships represent the relationships of an input payload
type Relationships struct {
	Principals      *PolicyDetails `json:"principals,omitempty" binding:"required,dive"`
	TargetResources *PolicyDetails `json:"target_resources,omitempty" binding:"required,dive"`
}

// PolicyDetails consists of array of resources associated with this policy
type PolicyDetails struct {
	Data []Resource `json:"data" binding:"required"`
}

// Resource denotes either the Principal or Target resource
type Resource struct {
	Type string `json:"type" binding:"required"`
	ID   string `json:"id" binding:"required"`
}

// // Response represents the final payload sent back to the user
// type Response struct {
// 	Data []Element `json:"data"`
// }

// // Element consists of all details of a resource
// type Element struct {
// 	Type          string         `json:"type" binding:"required"`
// 	ID            string         `json:"id"`
// 	Attributes    *Resource      `json:"attributes" binding:"required,dive"`
// 	Relationships *Relationships `json:"relationships"`
// }

// // Relationships provides details about a given resource like policies and parent
// type Relationships struct {
// 	Parent   *Parent   `json:"parent,omitempty"`
// 	Policies *Policies `json:"policies"`
// }

// // Identifier helps acts as a reference to any given entity
// type Identifier struct {
// 	Type string `json:"type,omitempty" enums:"policy, resource, grant, permission" binding:"required"`
// 	ID   string `json:"id,omitempty" binding:"required"`
// }

// // Parent is the parent resource of this resource
// type Parent struct {
// 	Data *Identifier `json:"data,omitempty" binding:"required,dive"`
// }

// // Policies defined list of policies applicable to this resource
// type Policies struct {
// 	Data []*Identifier `json:"data"`
// }

// // Resource represents an entity created by the user.
// type Resource struct {
// 	Name     string `json:"name" mapstructure:"name"`
// 	SourceID string `json:"source_id" binding:"required" mapstructure:"source_id"`
// }

// func constructResourceResponse(resource *Resource, id ...string) Element {
// 	relationships := generateResourceRelationship(id[1])
// 	return Element{
// 		Type:          "resource",
// 		ID:            id[0],
// 		Attributes:    resource,
// 		Relationships: relationships,
// 	}
// }

// func generateResourceRelationship(id string) *Relationships {
// 	policy := &Identifier{Type: "policy", ID: "policy-id"}
// 	policies := &Policies{Data: []*Identifier{policy}}
// 	var parent *Parent

// 	if id != "" {
// 		parent = &Parent{Data: &Identifier{Type: "resource", ID: id}}
// 	}
// 	relationships := &Relationships{Parent: parent, Policies: policies}
// 	return relationships
// }
