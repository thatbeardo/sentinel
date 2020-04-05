package views

import "github.com/thatbeardo/go-sentinel/models/resource"

// ResourceResponse encapsulates resource get end point response as per json schema
type ResourceResponse struct {
	Data []data `json:"data"`
}

type data struct {
	Type          string            `json:"type"`
	ID            string            `json:"id"`
	Attributes    resource.Resource `json:"attributes"`
	Relationships relationships     `json:"relationships"`
}

type relationships struct {
	Parent   parent   `json:"parent"`
	Policies policies `json:"policies"`
}

type parent struct {
	Data genericData `json:"data"`
}

type policies struct {
	Data []genericData `json:"data"`
}

type genericData struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// WrapGetResource wraps the resource into json schema
func WrapGetResource(resource resource.Resource) ResourceResponse {
	genericDataElement := genericData{}
	dataArray := []genericData{genericDataElement}
	policies := policies{Data: dataArray}

	parent := parent{Data: genericDataElement}
	relationships := relationships{Parent: parent, Policies: policies}

	dataPayload := data{
		Type:          "resource",
		ID:            "sample-id",
		Attributes:    resource,
		Relationships: relationships,
	}

	return ResourceResponse{Data: []data{dataPayload}}
}
