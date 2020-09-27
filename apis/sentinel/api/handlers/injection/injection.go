package injection

import (
	"encoding/json"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/middleware"
	"github.com/bithippie/guard-my-app/apis/sentinel/http"
)

// Reset takes all injected variables to their original values
func Reset() {
	VerifyResourceOwnership = verifyResourceOwnershipDefault
	ValidateNewResource = validateNewResourceDefault
	VerifyContextOwnership = verifyContextOwnershipDefault
	VerifyPermissionOwnership = verifyPermissionOwnershipDefault
	VerifyGrantExistence = verifyGrantExistenceDefault
	VerifyPermissionIdempotence = verifyPermissionIdempotenceDefault
	Marshal = marshalDefault
	Unmarshal = unmarshalDefault
	Post = postDefault
}

// VerifyResourceOwnership checks if the resource belongs to the correct tenant
var VerifyResourceOwnership = verifyResourceOwnershipDefault
var verifyResourceOwnershipDefault = middleware.VerifyResourceOwnership

// ValidateNewResource validates the resource sent in POST payload
var ValidateNewResource = validateNewResourceDefault
var validateNewResourceDefault = middleware.ValidateNewResource

// VerifyContextOwnership checks if the context being updated belongs to the current tenant
var VerifyContextOwnership = verifyContextOwnershipDefault
var verifyContextOwnershipDefault = middleware.VerifyContextOwnership

// VerifyPermissionOwnership checks if the permission belongs to this tenant
var VerifyPermissionOwnership = verifyPermissionOwnershipDefault
var verifyPermissionOwnershipDefault = middleware.VerifyRelationshipOwnership

// VerifyGrantExistence checks if a duplicate grant is being requested
var VerifyGrantExistence = verifyGrantExistenceDefault
var verifyGrantExistenceDefault = middleware.VerifyGrantExistence

// VerifyPermissionIdempotence checks if a duplicate permission is being created
var VerifyPermissionIdempotence = verifyPermissionIdempotenceDefault
var verifyPermissionIdempotenceDefault = middleware.VerifyPermissionIdempotence

// Marshal returns the JSON encoding of the parameter passed
var Marshal = marshalDefault
var marshalDefault = json.Marshal

// Unmarshal encodes the byte stream into json
var Unmarshal = unmarshalDefault
var unmarshalDefault = json.Unmarshal

// Post makes a POST request to an endpoint
var Post = postDefault
var postDefault = http.Post
