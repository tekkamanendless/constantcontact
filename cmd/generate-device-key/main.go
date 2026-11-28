package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/tekkamanendless/constantcontact"
)

func main() {
	ctx := context.Background()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	clientID := os.Getenv("CONSTANT_CONTACT_CLIENT_ID")
	fmt.Printf("CONSTANT_CONTACT_CLIENT_ID=%s\n", clientID)

	constantContactClient := constantcontact.NewClient()

	deviceAuthorizeInput := constantcontact.DeviceAuthorizeInput{
		ClientID:     clientID,
		ResponseType: "code",
		Scope:        "contact_data campaign_data offline_access",
		State:        "some-random-state",
	}
	deviceAuthorizeOutput, err := constantContactClient.DeviceAuthorize(ctx, deviceAuthorizeInput)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Output: %+v\n", deviceAuthorizeOutput)

	fmt.Printf("Open the following URL in your browser: %s\n", deviceAuthorizeOutput.VerificationURIComplete)
	fmt.Printf("Enter the following code: %s\n", deviceAuthorizeOutput.UserCode)
	fmt.Printf("Press enter to continue...\n")
	os.Stdin.Read(make([]byte, 1))

	tries := 0
	for {
		if tries > 0 {
			time.Sleep(time.Duration(deviceAuthorizeOutput.IntervalSeconds) * time.Second)
		}
		tries++

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
		break
	}
}
