package constantcontact

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/tekkamanendless/httperror"
	"github.com/tekkamanendless/restapiclient"
)

type Client struct {
	config     Config
	restClient *restapiclient.Client
}

type Config struct {
	BaseURL      string
	AccessToken  string
	RefreshToken string
}

type Option func(*Config)

func WithAccessToken(accessToken string) Option {
	return func(c *Config) {
		c.AccessToken = accessToken
	}
}

func WithBaseURL(baseURL string) Option {
	return func(c *Config) {
		c.BaseURL = baseURL
	}
}

func WithRefreshToken(refreshToken string) Option {
	return func(c *Config) {
		c.RefreshToken = refreshToken
	}
}

func NewClient(opts ...Option) *Client {
	config := Config{
		BaseURL: "https://api.cc.email",
	}
	for _, opt := range opts {
		opt(&config)
	}

	restClientOptions := []restapiclient.Option{
		restapiclient.OptionHeader("Accept", "application/json"),
		restapiclient.OptionErrorHandler(func(response *http.Response) error {
			httpErr := httperror.ErrorFromStatus(response.StatusCode)

			body, err := io.ReadAll(response.Body)
			if err != nil {
				return nil // Use the original error.
			}

			var errorOutput ErrorOutput
			err = json.Unmarshal(body, &errorOutput)
			if err != nil {
				return nil // Use the original error.
			}

			var errs []error
			for _, e := range errorOutput {
				errs = append(errs, fmt.Errorf("%s: %s", e.ErrorKey, e.ErrorMessage))
			}

			if len(errs) == 0 {
				return nil // Use the original error.
			}

			return fmt.Errorf("%w: %w", httpErr, errors.Join(errs...))
		}),
	}
	if config.AccessToken != "" {
		restClientOptions = append(restClientOptions, restapiclient.OptionHeader("Authorization", "Bearer "+config.AccessToken))
	}

	return &Client{
		config:     config,
		restClient: restapiclient.New(config.BaseURL, restClientOptions...),
	}
}
