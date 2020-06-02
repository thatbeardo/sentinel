package inputs

// Payload is the payload that a POST endpoint expects.
type Payload struct {
	Data *PolicyDetails `json:"data" binding:"required,dive"`
}

// PolicyDetails is the paylaod sent when creating a new resource
type PolicyDetails struct {
	Type       string      `json:"type" binding:"required"`
	Attributes *Attributes `json:"attributes" binding:"required,dive"`
}

// Attributes defines the name of the policy
type Attributes struct {
	Name string `json:"name" binding:"required"`
}
