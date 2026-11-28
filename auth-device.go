package constantcontact

import (
	"context"

	"github.com/tekkamanendless/constantcontact/wrap"
)

type DeviceAuthorizeInput struct {
	_ string `api:"httppath:https://authz.constantcontact.com/oauth2/default/v1/device/authorize"`
	_ string `api:"httpmethod:POST"`

	ClientID     string `form:"client_id"`
	ResponseType string `form:"response_type"`
	Scope        string `form:"scope"` // Space-separated list of scopes.
	State        string `form:"state"`
}

type DeviceAuthorizeOutput struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresInSeconds        int    `json:"expires_in"`
	IntervalSeconds         int    `json:"interval"`
}

func (c *Client) DeviceAuthorize(ctx context.Context, input DeviceAuthorizeInput) (DeviceAuthorizeOutput, error) {
	var output DeviceAuthorizeOutput
	err := wrap.WrapRequest(ctx, c.restClient, input, &output)
	return output, err
}
