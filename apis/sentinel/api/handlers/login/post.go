package login

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/bithippie/guard-my-app/apis/sentinel/api/handlers/injection"
	"github.com/bithippie/guard-my-app/apis/sentinel/api/views"
	login "github.com/bithippie/guard-my-app/apis/sentinel/models/login/dto"
	"github.com/gin-gonic/gin"
)

// @Summary Fetch access token to make authenticated requests to the API
// @Tags Login
// @Description Allows the user to fetch an access_token and make authenticated requests to Sentinel
// @Accept  json
// @Produce  json
// @Param input body login.ClientCredentials true "ClientID and ClientSecret needed to authenticate a users identity"
// @Success 200 {object} login.BearerToken	"ok"
// @Success 500 {object} views.ErrView
// @Router /v1/login [post]
func getAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {

		var input login.ClientCredentials
		if err := c.ShouldBind(&input); err != nil {
			views.Wrap(err, c)
			return
		}

		requestBody, err := injection.Marshal(map[string]string{
			"client_id":     input.ClientID,
			"client_secret": input.ClientSecret,
			"audience":      "https://api.guardmy.app/",
			"grant_type":    "client_credentials",
		})
		if err != nil {
			views.Wrap(err, c)
			return
		}

		resp, err := injection.Post("https://bithippie.auth0.com/oauth/token", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			views.Wrap(err, c)
			return
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		var bearerToken login.BearerToken
		err = injection.Unmarshal(body, &bearerToken)
		if err != nil {
			views.Wrap(err, c)
			return
		}

		c.JSON(http.StatusOK, bearerToken)
	}
}
