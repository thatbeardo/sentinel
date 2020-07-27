package context

// Input is the payload that a POST endpoint expects.
type Input struct {
	Data *InputDetails `json:"data" binding:"required,dive"`
}

// InputDetails is the paylaod sent when creating a new resource
type InputDetails struct {
	Type       string      `json:"type" enums:"context" binding:"required"`
	Attributes *Attributes `json:"attributes" binding:"required,dive"`
}

// Attributes defines the name of the context
type Attributes struct {
	Name string `json:"name" binding:"required"`
}
