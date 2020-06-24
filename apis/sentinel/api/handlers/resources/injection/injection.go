package injection

import "github.com/bithippie/guard-my-app/apis/sentinel/api/middleware"

// Reset takes all injected variables to their original values
func Reset() {
	VerifyResourceOwnership = verifyResourceOwnershipDefault
	ValidateNewResource = validateNewResourceDefault
}

// VerifyResourceOwnership is a middleware injected at runtime to check if the resource belongs to the correct tenant
var VerifyResourceOwnership = verifyResourceOwnershipDefault
var verifyResourceOwnershipDefault = middleware.VerifyResourceOwnership

// ValidateNewResource is responsible for checking the resource sent in POST payload
var ValidateNewResource = validateNewResourceDefault
var validateNewResourceDefault = middleware.ValidateNewResource
