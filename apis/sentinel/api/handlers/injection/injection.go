package injection

import "github.com/bithippie/guard-my-app/apis/sentinel/api/middleware"

// Reset takes all injected variables to their original values
func Reset() {
	VerifyResourceOwnership = verifyResourceOwnershipDefault
	ValidateNewResource = validateNewResourceDefault
}

// VerifyResourceOwnership checks if the resource belongs to the correct tenant
var VerifyResourceOwnership = verifyResourceOwnershipDefault
var verifyResourceOwnershipDefault = middleware.VerifyResourceOwnership

// ValidateNewResource validates the resource sent in POST payload
var ValidateNewResource = validateNewResourceDefault
var validateNewResourceDefault = middleware.ValidateNewResource

// VerifyPolicyOwnership checks if the policy being updated belongs to the current tenant
var VerifyPolicyOwnership = verifyPolicyOwnershipDefault
var verifyPolicyOwnershipDefault = middleware.VerifyPolicyOwnership

// VerifyPermissionOwnership checks if the permission belongs to this tenant
var VerifyPermissionOwnership = verifyPermissionOwnershipDefault
var verifyPermissionOwnershipDefault = middleware.VerifyRelationshipOwnership
