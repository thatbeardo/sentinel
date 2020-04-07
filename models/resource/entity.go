package resource

// Input is the payload that a POST endpoint expects.
type Input struct {
	Type       string   `json:"type" binding:"required"`
	ID         string   `json:"id"`
	Attributes Resource `json:"attributes" binding:"required"`
}

type relationshipsInput struct {
	Parent parent `json:"parent"`
}

// Response represents the final payload sent back to the user
type Response struct {
	Data []Dto `json:"data"`
}

// Dto consists of all details of a resource
type Dto struct {
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

func constructResourceResponse(resources []Resource, id []string) Response {
	var resourcesArray []Dto
	for index, resourceElement := range resources {
		relationships := generateResourceRelationship()
		payload := Dto{
			Type:          "resource",
			ID:            id[index],
			Attributes:    resourceElement,
			Relationships: relationships,
		}
		resourcesArray = append(resourcesArray, payload)
	}
	return Response{Data: resourcesArray}
}

func generateResourceRelationship() relationships {
	policy := identifier{Type: "policy", ID: "some-id"}
	policies := policies{Data: []identifier{policy}}

	parent := parent{Data: identifier{Type: "resource", ID: "some-id"}}
	relationships := relationships{Parent: parent, Policies: policies}
	return relationships
}
