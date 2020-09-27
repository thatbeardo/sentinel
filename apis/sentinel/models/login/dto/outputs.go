package login

// BearerToken contains an access token to make authenticated requests to Sentinel
type BearerToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}
