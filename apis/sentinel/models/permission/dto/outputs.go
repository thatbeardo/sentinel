package permission

// Output represents the data that has to be sent back to the use
type Output struct {
	Data []Details `json:"data"`
}

// OutputDetails consists of data pertaining to a single permission
type OutputDetails struct {
	Data Details `json:"data"`
}

// Details represents whether a context has access to a resource
type Details struct {
	InputDetails
	ID string `json:"id"`
}
