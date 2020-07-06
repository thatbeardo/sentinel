package resource

// Attributes defines data about the resource being created
type Attributes struct {
	Name      string `json:"name" mapstructure:"name"`
	SourceID  string `json:"source_id" binding:"required" mapstructure:"source_id"`
	ContextID string `json:"context_id" mapstructure:"context_id"`
}

// Data represents a resource/context this resource refers to
type Data struct {
	Type string `json:"type,omitempty" enums:"resource" binding:"required"`
	ID   string `json:"id,omitempty" binding:"required"`
}

// Parent represents a node that owns this resource
type Parent struct {
	Data *Data `json:"data,omitempty" binding:"required,dive"`
}

// InputRelationships denotes the parent that needs to be attached
type InputRelationships struct {
	Parent *Parent `json:"parent,omitempty" binding:"required,dive"`
}

// InputDetails represents the data sent by the user
type InputDetails struct {
	Type          string              `json:"type,omitempty" enums:"resource" binding:"required"`
	Attributes    *Attributes         `json:"attributes" binding:"required,dive"`
	Relationships *InputRelationships `json:"relationships,omitempty"`
}

// Input is the data sent by the user to the system
type Input struct {
	Data *InputDetails `json:"data" binding:"required"`
}
