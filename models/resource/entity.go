package resource

// Input is the payload that a POST endpoint expects.
type Input struct {
	Data *InputElement `json:"data" binding:"required,dive"`
}

// InputElement is the paylaod sent when creating a new resource
type InputElement struct {
	Type          string              `json:"type" binding:"required"`
	Attributes    *Resource           `json:"attributes" binding:"required,dive"`
	Relationships *RelationshipsInput `json:"relationships,omitempty"`
}

// RelationshipsInput represent the relationships of an input payload
type RelationshipsInput struct {
	Parent *Parent `json:"parent,omitempty" binding:"required,dive"`
}

// Response represents the final payload sent back to the user
type Response struct {
	Data []Element `json:"data"`
}

// Element consists of all details of a resource
type Element struct {
	Type          string         `json:"type" binding:"required"`
	ID            string         `json:"id"`
	Attributes    *Resource      `json:"attributes" binding:"required,dive"`
	Relationships *Relationships `json:"relationships"`
}

// Relationships provides details about a given resource like policies and parent
type Relationships struct {
	Parent   *Parent   `json:"parent"`
	Policies *Policies `json:"policies"`
}

// Identifier helps acts as a reference to any given entity
type Identifier struct {
	Type string `json:"type,omitempty" enums:"policy, resource, grant, permission" binding:"required"`
	ID   string `json:"id,omitempty" binding:"required"`
}

// Parent is the parent resource of this resource
type Parent struct {
	Data *Identifier `json:"data,omitempty" binding:"required,dive"`
}

// Policies defined list of policies applicable to this resource
type Policies struct {
	Data []*Identifier `json:"data"`
}

// Resource represents an entity created by the user.
type Resource struct {
	Name     string `json:"name"`
	SourceID string `json:"source_id" binding:"required"`
}

func constructResourceResponse(resource *Resource, id ...string) Element {
	relationships := generateResourceRelationship(id)
	return Element{
		Type:          "resource",
		ID:            id[0],
		Attributes:    resource,
		Relationships: relationships,
	}
}

func generateResourceRelationship(ids []string) *Relationships {
	policy := &Identifier{Type: "policy", ID: "policy-id"}
	policies := &Policies{Data: []*Identifier{policy}}
	var parent *Parent

	if len(ids) == 2 {
		parent = &Parent{Data: &Identifier{Type: "resource", ID: "parent-id"}}
	}
	relationships := &Relationships{Parent: parent, Policies: policies}
	return relationships
}
