package context

// Output represents the final payload sent back to the user
type Output struct {
	Data []Details `json:"data"`
}

// Details represents data pertaining to one context
type Details struct {
	InputDetails
	ID            string         `json:"id"`
	Relationships *Relationships `json:"relationships"`
}

// OutputDetails contains all data to uniquely identify a resource
type OutputDetails struct {
	Data Details `json:"data"`
}

// Relationships represent the relationships of an input payload
type Relationships struct {
	Principals      *Relationship `json:"principals" binding:"required,dive"`
	TargetResources *Relationship `json:"target_resources" binding:"required,dive"`
}

// Relationship consists of array of resources associated with this context
type Relationship struct {
	Data []Resource `json:"data" binding:"required,dive"`
}

// Resource denotes either the Principal or Target resource
type Resource struct {
	Type string `json:"type" enums:"resource"`
	ID   string `json:"id"`
}
