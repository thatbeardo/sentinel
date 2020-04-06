package views

import "github.com/thatbeardo/go-sentinel/models/resource"

// ResourceResponse encapsulates resource get end point response as per json schema
type ResourceResponse struct {
	Data []Resource `json:"data"`
}

// Resource denotes the dto
type Resource struct {
	Type          string             `json:"type"`
	ID            string             `json:"id"`
	Attributes    *resource.Resource `json:"attributes"`
	Relationships relationships      `json:"relationships"`
}

// WrapGetResources wraps the resource into json schema
func WrapGetResources(resources []*resource.Resource) ResourceResponse {
	resourcesArray := make([]Resource, len(resources))

	for _, resourceElement := range resources {
		policy := identifier{Type: "policy", ID: "some-id"}
		policies := policies{Data: []identifier{policy}}

		parent := parent{Data: identifier{Type: "resource", ID: "some-id"}}
		relationships := relationships{Parent: parent, Policies: policies}

		payload := Resource{
			Type:          "resource",
			ID:            "sample-id",
			Attributes:    resourceElement,
			Relationships: relationships,
		}
		resourcesArray = append(resourcesArray, payload)
	}

	return ResourceResponse{Data: resourcesArray}
}
