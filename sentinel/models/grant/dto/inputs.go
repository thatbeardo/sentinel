package grant

// Input represent the data sent from the user to create a permission
type Input struct {
	Data InputDetails `json:"data" binding:"required"`
}

// InputDetails represents whether a policy has access to a resource
type InputDetails struct {
	Type       string      `json:"type" binding:"required"`
	Attributes *Attributes `json:"attributes" binding:"required,dive"`
}

// Attributes represents name and whether if a policy is permitted to access a resource
type Attributes struct {
	WithGrant bool `json:"with_grant" binding:"required"`
}
