package resource

// Output consists of a list of all resources to be returned
type Output struct {
	Data []Details `json:"data"`
}

// Details to uniquely identify a resource
type Details struct {
	ID            string         `json:"id"`
	Type          string         `json:"type,omitempty" enums:"resource"`
	Attributes    *Attributes    `json:"attributes"`
	Relationships *Relationships `json:"relationships,omitempty"`
}

// OutputDetails contains all data to uniquely identify a resource
type OutputDetails struct {
	Data Details `json:"data"`
}

// Policies represents all the policies that have a grant to this resource
type Policies struct {
	Data []Data `json:"data"`
}

// Relationships provides details about a given resource like policies and parent
type Relationships struct {
	Parent   *Parent   `json:"parent,omitempty"`
	Policies *Policies `json:"policies,omitempty"`
}
