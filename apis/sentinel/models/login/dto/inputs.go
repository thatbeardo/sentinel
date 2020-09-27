package login

// ClientCredentials are passed over to auth0 to fetch an access token
type ClientCredentials struct {
	ClientID     string `json:"client_id" binding:"required"`
	ClientSecret string `json:"client_secret" binding:"required"`
}
