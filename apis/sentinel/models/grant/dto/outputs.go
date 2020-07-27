package grant

// Output represents the data that has to be sent back to the use
type Output struct {
	Data []Details `json:"data"`
}

// Details represent data about one grant entity
type Details struct {
	InputDetails
	Relationships `json:"relationships"`
	ID            string `json:"id"`
}

// OutputDetails represents whether a context has access to a resource
type OutputDetails struct {
	Data Details `json:"data"`
}

// Data provides details about the context and Principal
type Data struct {
	Type string `json:"type" enums:"grant" binding:"required"`
	ID   string `json:"id" binding:"required"`
}

// Relationship has data pertaining to a context/Principal
type Relationship struct {
	Data Data `json:"data"`
}

// Relationships enlists details about the context and the resource that are attached with this grant
type Relationships struct {
	Context   *Relationship `json:"context" binding:"required"`
	Principal *Relationship `json:"principal" binding:"required"`
}
