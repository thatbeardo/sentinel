package inputs

// Payload represent the data sent from the user to create a permission
type Payload struct {
	Data []PermissionDetails `json:"data" binding:"required"`
}

// PermissionDetails represents whether a policy has access to a resource
type PermissionDetails struct {
	Type       string      `json:"type" binding:"required"`
	Attributes *Attributes `json:"attributes" binding:"required,dive"`
}

// Attributes represents name and whether if a policy is permitted to access a resource
type Attributes struct {
	Name      string `json:"name" binding:"required"`
	Permitted string `json:"permitted" binding:"required"`
}
