package moov

import (
	"net/http"
)

type Client struct {
	Credentials Credentials
	HttpClient  *http.Client
}

// NewClient returns a moov.Client with credentials read from environment variables.
func NewClient(configurables ...ClientConfigurable) (*Client, error) {
	// Default client configuration if no configurable were specified
	client := &Client{
		Credentials: CredentialsFromEnv(),
		HttpClient:  DefaultHttpClient(),
	}

	// Apply all the configurable functions to the client
	for _, configurable := range configurables {
		if err := configurable(client); err != nil {
			return nil, err
		}
	}

	// Lets make sure that whatever they passed in for the credentials is valid.
	if err := client.Credentials.Validate(); err != nil {
		return nil, err
	}

	return client, nil
}

type ClientConfigurable func(c *Client) error

func WithCredentials(credentials Credentials) ClientConfigurable {
	return func(c *Client) error {
		c.Credentials = credentials
		return c.Credentials.Validate()
	}
}

func WithHttpClient(client *http.Client) ClientConfigurable {
	return func(c *Client) error {
		c.HttpClient = client
		return nil
	}
}
