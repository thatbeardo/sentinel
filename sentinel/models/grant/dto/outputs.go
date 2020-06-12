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

// OutputDetails represents whether a policy has access to a resource
type OutputDetails struct {
	Data Details `json:"data"`
}

// Data provides details about the Policy and Principal
type Data struct {
	Type string `json:"type" binding:"required"`
	ID   string `json:"id" binding:"required"`
}

// Relationship has data pertaining to a Policy/Principal
type Relationship struct {
	Data Data `json:"data"`
}

// Relationships enlists details about the Policy and the resource that are attached with this grant
type Relationships struct {
	Policy    *Relationship `json:"policy" binding:"required"`
	Principal *Relationship `json:"principal" binding:"required"`
}
