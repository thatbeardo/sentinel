package injection

import "github.com/bithippie/guard-my-app/apis/sentinel/api/claims"

// Reset takes all injected variables to their original values
func Reset() {
	ExtractClaims = extractClaimsDefault
}

// ExtractClaims reads the context and parses out the desired value from it
var ExtractClaims = extractClaimsDefault
var extractClaimsDefault = claims.Extract
