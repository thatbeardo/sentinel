package testdata

import (
	authorization "github.com/bithippie/guard-my-app/apis/sentinel/models/authorization/dto"
	permission "github.com/bithippie/guard-my-app/apis/sentinel/models/permission/dto"
	resource "github.com/bithippie/guard-my-app/apis/sentinel/models/resource/dto"
)

// Output is used as test data
var Output = authorization.Output{
	Data: []authorization.Details{details},
}

var details = authorization.Details{
	ID:            "test-target-id",
	Type:          "resource",
	Attributes:    attributes,
	Relationships: relationships,
}

var relationships = authorization.Relationships{
	Permissions: permissions,
}

var attributes = resource.Attributes{
	Name:     "test-target",
	SourceID: "test-target-source-id",
}

var permissions = authorization.Permissions{
	Data: []permission.Attributes{
		{
			Name:      "test-permission",
			Permitted: "allow",
		},
	},
}
