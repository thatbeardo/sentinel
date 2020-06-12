package authorization

// Input represent an object encapsulating all query paraneters for the authorization endpoint
type Input struct {
	Targets       []string
	Permissions   []string
	Depth         int
	IncludeDenied bool
}
