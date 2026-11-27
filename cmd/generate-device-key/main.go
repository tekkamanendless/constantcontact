package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/tekkamanendless/constantcontact"
	"github.com/tekkamanendless/restapiclient"
)

func main() {
	ctx := context.Background()

	clientID := os.Getenv("CONSTANT_CONTACT_CLIENT_ID")
	fmt.Printf("CONSTANT_CONTACT_CLIENT_ID=%s\n", clientID)

	constantContactClient := constantcontact.NewClient()

	type DeviceAuthorizeOutput struct {
		DeviceCode              string `json:"device_code"`
		UserCode                string `json:"user_code"`
		VerificationURI         string `json:"verification_uri"`
		VerificationURIComplete string `json:"verification_uri_complete"`
		ExpiresInSeconds        int    `json:"expires_in"`
		IntervalSeconds         int    `json:"interval"`
	}

	var deviceAuthorizeOutput DeviceAuthorizeOutput
	{
		values := url.Values{}
		values.Add("client_id", clientID)
		values.Add("response_type", "code")
		values.Add("scope", "contact_data campaign_data offline_access")
		values.Add("state", "some-random-state")

		err := constantContactClient.Do(ctx, http.MethodPost, "https://authz.constantcontact.com/oauth2/default/v1/device/authorize", values.Encode(), &deviceAuthorizeOutput, restapiclient.OptionHeader("Content-Type", "application/x-www-form-urlencoded"))
		if err != nil {
			panic(err)
		}

		fmt.Printf("Output: %+v\n", deviceAuthorizeOutput)

		fmt.Printf("Open the following URL in your browser: %s\n", deviceAuthorizeOutput.VerificationURIComplete)
		fmt.Printf("Enter the following code: %s\n", deviceAuthorizeOutput.UserCode)
		fmt.Printf("Press enter to continue...\n")
		os.Stdin.Read(make([]byte, 1))
	}

	for {
		time.Sleep(time.Duration(deviceAuthorizeOutput.IntervalSeconds) * time.Second)

		{
			input := constantcontact.PostOauth2TokenInput{
				ClientID:   clientID,
				GrantType:  constantcontact.GrantTypeDeviceCode,
				DeviceCode: deviceAuthorizeOutput.DeviceCode,
			}
			var output constantcontact.PostOauth2TokenOutput
			output, err := constantContactClient.PostOauth2Token(ctx, input)
			if err != nil {
				if strings.Contains(err.Error(), "authorization_pending") {
					continue
				}
				if strings.Contains(err.Error(), "expired_token") {
					fmt.Printf("You waited too long before entering the user code. Please try again.\n")
					os.Exit(1)
				}
				panic(err)
			}

			fmt.Printf("Output: %+v\n", output)

			fmt.Printf("\n")
			fmt.Printf("CONSTANT_CONTACT_CLIENT_ID=%s\n", clientID)
			fmt.Printf("CONSTANT_CONTACT_REFRESH_TOKEN=%s\n", output.RefreshToken)
		}
	}
}
