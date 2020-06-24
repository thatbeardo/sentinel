package claims

import (
	"context"

	"github.com/dgrijalva/jwt-go"
)

// Extract parses the context and returns the claim that is being requested
func Extract(ctx context.Context, claim string) string {
	user := ctx.Value("user")
	mapClaims := user.(*jwt.Token).Claims.(jwt.MapClaims)
	if val, ok := mapClaims[claim]; ok {
		return val.(string)
	}
	return ""
}
