package resource

// Resource represents an entity created by the user.
type Resource struct {
	Name     string `json:"name"`
	SourceID string `json:"source_id"`
}
