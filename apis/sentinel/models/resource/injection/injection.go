package injection

import (
	"github.com/bithippie/guard-my-app/apis/sentinel/api/claims"
	"github.com/bithippie/guard-my-app/apis/sentinel/api/tenant"
)

// Reset takes all injected variables to their original values
func Reset() {
	ExtractClaims = extractClaimsDefault
}

// ExtractClaims reads the context and parses out the desired value from it
var ExtractClaims = extractClaimsDefault
var extractClaimsDefault = claims.Extract

// ExtractTenant reads the context and returns the tenant set
var ExtractTenant = extractTenantDefault
var extractTenantDefault = tenant.Extract
