package injection

import (
	"encoding/json"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/authentication"
	"github.com/bithippie/guard-my-app/apis/sentinel/api/claims"
	"github.com/bithippie/guard-my-app/apis/sentinel/api/tenant"
)

// Reset takes all injected variables to their original values
func Reset() {
	VerifyAccessToken = verifyAccessTokenDefault
	ExtractClaim = extractClaimDefault
	AddTenantToContext = addTenantToContextDefault
	Unmarshal = unmarshalDefault
}

// VerifyAccessToken verifies the callers identity using the digital signature
var VerifyAccessToken = verifyAccessTokenDefault
var verifyAccessTokenDefault = authentication.CheckJwt

// ExtractClaim parses context and returns the desired claim
var ExtractClaim = extractClaimDefault
var extractClaimDefault = claims.Extract

// AddTenantToContext updates the existing context with tenant information
var AddTenantToContext = addTenantToContextDefault
var addTenantToContextDefault = tenant.Add

// Unmarshal decodes a byte stream into the target interface
var Unmarshal = unmarshalDefault
var unmarshalDefault = json.Unmarshal
