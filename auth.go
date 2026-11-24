package constantcontact

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/tekkamanendless/restapiclient"
)

type Oauth2Input struct {
	ClientID     string
	RefreshToken string
	Scope        string // Space-separated list of scopes.
}

type Oauth2Output struct {
	TokenType        string `json:"token_type"`
	ExpiresInSeconds int    `json:"expires_in"`
	AccessToken      string `json:"access_token"`
	Scope            string `json:"scope"`
	RefreshToken     string `json:"refresh_token"`
}

func (c *Client) Oauth2(ctx context.Context, input Oauth2Input) (Oauth2Output, error) {
	var output Oauth2Output

	queryParams := url.Values{}
	queryParams.Add("refresh_token", input.RefreshToken)
	queryParams.Add("grant_type", "refresh_token")
	queryParams.Add("client_id", input.ClientID)
	queryParams.Add("scope", input.Scope)
	targetURL := "https://authz.constantcontact.com/oauth2/default/v1/token?" + queryParams.Encode()

	err := c.restClient.Do(ctx, http.MethodPost, targetURL, nil, &output, restapiclient.OptionHeader("Content-Type", "application/x-www-form-urlencoded"))
	if err != nil {
		return output, fmt.Errorf("failed to refresh token: %w", err)
	}

	c.config.AccessToken = output.AccessToken
	c.restClient = c.restClient.WithOptions(restapiclient.OptionHeader("Authorization", "Bearer "+output.AccessToken))

	return output, nil
}
