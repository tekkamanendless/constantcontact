package constantcontact

import "github.com/tekkamanendless/restapiclient"

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
	}
	if config.AccessToken != "" {
		restClientOptions = append(restClientOptions, restapiclient.OptionHeader("Authorization", "Bearer "+config.AccessToken))
	}

	return &Client{
		config:     config,
		restClient: restapiclient.New(config.BaseURL, restClientOptions...),
	}
}
