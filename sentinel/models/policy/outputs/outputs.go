package outptus

import "github.com/bithippie/guard-my-app/sentinel/models/policy/inputs"

// Response represents the final payload sent back to the user
type Response struct {
	Data []Element `json:"data"`
}

// Element represents data pertaining to one policy
type Element struct {
	inputs.Element
	ID string `json:"id"`
}
