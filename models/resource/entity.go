package resource

// Input is the payload that a POST endpoint expects.
type Input struct {
	Data resourceInput `json:"data" bindings:"required"`
}

type resourceInput struct {
	Type          string             `json:"type" bindings:"required"`
	Attributes    Resource           `json:"attributes" binding:"required"`
	Relationships relationshipsInput `json:"relationships"`
}

type relationshipsInput struct {
	Parent parent `json:"parent"`
}

// Response represents the final payload sent back to the user
type Response struct {
	Data []Element `json:"data"`
}

// Element consists of all details of a resource
type Element struct {
	Type          string        `json:"type" binding:"required"`
	ID            string        `json:"id"`
	Attributes    Resource      `json:"attributes" binding:"required"`
	Relationships relationships `json:"relationships"`
}

type relationships struct {
	Parent   parent   `json:"parent"`
	Policies policies `json:"policies"`
}

type identifier struct {
	Type string `json:"type" enums:"policy, resource, grant, permission"`
	ID   string `json:"id"`
}

type parent struct {
	Data identifier `json:"data"`
}

type policies struct {
	Data []identifier `json:"data"`
}

// Resource represents an entity created by the user.
type Resource struct {
	Name     string `json:"name"`
	SourceID string `json:"source_id" binding:"required"`
}

func constructResourceResponse(resource Resource, id string) Element {
	relationships := generateResourceRelationship()
	return Element{
		Type:          "resource",
		ID:            id,
		Attributes:    resource,
		Relationships: relationships,
	}
}

func generateResourceRelationship() relationships {
	policy := identifier{Type: "policy", ID: "some-id"}
	policies := policies{Data: []identifier{policy}}

	parent := parent{Data: identifier{Type: "resource", ID: "some-id"}}
	relationships := relationships{Parent: parent, Policies: policies}
	return relationships
}
