package testdata

import (
	permission "github.com/bithippie/guard-my-app/apis/sentinel/models/permission/dto"
)

// Input represent the input used to create a permission
var Input = &permission.Input{
	Data: permission.InputDetails{
		Type: "permission",
		Attributes: &permission.Attributes{
			Name:      "test-permission",
			Permitted: "allow",
		},
	},
}

// Output denotes mocked data that is sent back in response body
var Output = permission.Output{
	Data: []permission.Details{details},
}

// OutputDetails mocks data pertaining to one permission
var OutputDetails = permission.OutputDetails{
	Data: details,
}

var details = permission.Details{
	ID: "test-id",
	InputDetails: permission.InputDetails{
		Type:       "permission",
		Attributes: attributes,
	},
}

var attributes = &permission.Attributes{
	Name:      "test-permission",
	Permitted: "allow",
}
