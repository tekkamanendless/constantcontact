package constantcontact

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/tekkamanendless/restapiclient"
)

type GrantType string

const (
	GrantTypeRefreshToken GrantType = "refresh_token"
	GrantTypeDeviceCode   GrantType = "urn:ietf:params:oauth:grant-type:device_code"
)

type PostOauth2TokenInput struct {
	_ string `api:"httppath:https://authz.constantcontact.com/oauth2/default/v1/token"`
	_ string `api:"httpmethod:POST"`

	ClientID     string    `form:"client_id"`
	DeviceCode   string    `form:"device_code,omitempty"`
	GrantType    GrantType `form:"grant_type"`
	RefreshToken string    `form:"refresh_token,omitempty"`
	Scope        string    `form:"scope,omitempty"` // Space-separated list of scopes.
}

type PostOauth2TokenOutput struct {
	TokenType        string `json:"token_type"`
	ExpiresInSeconds int    `json:"expires_in"`
	AccessToken      string `json:"access_token"`
	Scope            string `json:"scope"` // Space-separated list of scopes.
	RefreshToken     string `json:"refresh_token"`
}

// PostOauth2Token requests an Oauth2 access token.
func (c *Client) PostOauth2Token(ctx context.Context, input PostOauth2TokenInput) (PostOauth2TokenOutput, error) {
	var output PostOauth2TokenOutput

	queryParams := url.Values{}
	queryParams.Add("client_id", input.ClientID)
	if input.DeviceCode != "" {
		queryParams.Add("device_code", input.DeviceCode)
	}
	queryParams.Add("grant_type", string(input.GrantType))
	if input.RefreshToken != "" {
		queryParams.Add("refresh_token", input.RefreshToken)
	}
	if input.Scope != "" {
		queryParams.Add("scope", input.Scope)
	}
	targetURL := "https://authz.constantcontact.com/oauth2/default/v1/token?" + queryParams.Encode()

	err := c.restClient.Do(ctx, http.MethodPost, targetURL, nil, &output, restapiclient.OptionHeader("Content-Type", "application/x-www-form-urlencoded"))
	if err != nil {
		return output, fmt.Errorf("failed to refresh token: %w", err)
	}

	return output, nil
}

// AuthenticateWithOauth2 authenticates the client with the OAuth2 token.
//
// If successful, the client will be authenticated and the access token will be set for all subsequent requests.
func (c *Client) AuthenticateWithOauth2(ctx context.Context, input PostOauth2TokenInput) error {
	output, err := c.PostOauth2Token(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to authenticate with oauth2: %w", err)
	}

	c.config.AccessToken = output.AccessToken
	c.restClient = c.restClient.WithOptions(restapiclient.OptionHeader("Authorization", "Bearer "+output.AccessToken))

	return nil
}
