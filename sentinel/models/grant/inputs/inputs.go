package inputs

// Payload represent the data sent from the user to create a permission
type Payload struct {
	Data GrantDetails `json:"data" binding:"required"`
}

// GrantDetails represents whether a policy has access to a resource
type GrantDetails struct {
	Type       string      `json:"type" binding:"required"`
	Attributes *Attributes `json:"attributes" binding:"required,dive"`
}

// Attributes represents name and whether if a policy is permitted to access a resource
type Attributes struct {
	WithGrant bool `json:"with_grant" binding:"required"`
}
