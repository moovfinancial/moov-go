package moov

import (
	"cmp"
	"io"
	"net/http"
	"os"

	"golang.org/x/time/rate"
)

type Client struct {
	Credentials Credentials
	HttpClient  *http.Client
	rateLimiter *rate.Limiter

	decoder Decoder

	bearerToken string
	origin      string
	referer     string

	moovURLScheme string
}

const defaultMoovURLScheme = "https"

// NewClient returns a moov.Client with credentials read from environment variables.
func NewClient(configurables ...ClientConfigurable) (*Client, error) {
	// Default client configuration if no configurables were specificied
	client := &Client{
		Credentials:   CredentialsFromEnv(),
		HttpClient:    DefaultHttpClient(),
		moovURLScheme: cmp.Or(os.Getenv("MOOV_URL_SCHEME"), defaultMoovURLScheme),
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

// WithToken configures the client to authenticate outbound calls using the
// given bearer token instead of Basic auth. Sets Credentials.Token and
// revalidates. Intended for pass-through scenarios where a caller-supplied
// access token should authenticate every call made by this client.
// Bearer tokens may 401 unless at least one of the `Origin` or `Referer`
// headers is sent, so pass WithOrigin and/or WithReferer options for those.
// Any origin/referer already set on c carries over unless an option overrides it.
func (c *Client) WithBearerToken(t string, opts ...BearerTokenOption) *Client {
	client := &Client{
		rateLimiter:   c.rateLimiter,
		decoder:       c.decoder,
		moovURLScheme: c.moovURLScheme,
		Credentials:   c.Credentials,
		HttpClient:    c.HttpClient,
		bearerToken:   t,
		origin:        c.origin,
		referer:       c.referer,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// BearerTokenOption configures the headers sent alongside a bearer token.
type BearerTokenOption func(c *Client)

// WithOrigin sets the Origin header sent on every call made by the client.
func WithOrigin(origin string) BearerTokenOption {
	return func(c *Client) {
		c.origin = origin
	}
}

// WithReferer sets the Referer header sent on every call made by the client.
func WithReferer(referer string) BearerTokenOption {
	return func(c *Client) {
		c.referer = referer
	}
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

type Decoder func(r io.Reader, contentType string, item any) error

func WithDecoder(dec Decoder) ClientConfigurable {
	return func(c *Client) error {
		c.decoder = dec
		return nil
	}
}

func WithMoovURLScheme(scheme string) ClientConfigurable {
	return func(c *Client) error {
		if scheme == "" {
			return nil // no-op
		}
		c.moovURLScheme = scheme
		return nil
	}
}
