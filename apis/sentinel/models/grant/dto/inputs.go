package grant

// Input represent the data sent from the user to create a permission
type Input struct {
	Data InputDetails `json:"data" binding:"required"`
}

// InputDetails represents whether a context has access to a resource
type InputDetails struct {
	Type       string      `json:"type" enums:"grant" binding:"required"`
	Attributes *Attributes `json:"attributes" binding:"required,dive"`
}

// Attributes represents name and whether if a context is permitted to access a resource
type Attributes struct {
	WithGrant bool `json:"with_grant" binding:"required"`
}
