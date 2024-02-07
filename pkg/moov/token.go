package moov

import (
	"context"
	"net/http"
)

type AccessTokenRequest struct {
	GrantType string `json:"grant_type"`
	// If not specified in `Authorization: Basic` it can be specified here
	ClientId *string `json:"client_id,omitempty"`
	// If not specified in `Authorization: Basic` it can be specified here
	ClientSecret *string `json:"client_secret,omitempty"`
	// Scopes that were allowed. Space delimited list of scopes. Take a look at the `ApplicationScope` model to see whats available
	Scope string `json:"scope,omitempty"`
	// Allowed endpoints this token is usable for. Space delimited list of urls.
	Audience string `json:"audience,omitempty"`
	// String passed to the authorization server to gain access to the system
	RefreshToken string `json:"refresh_token,omitempty"`
}

func (atr *AccessTokenRequest) WithScopes(scopes ...ScopeBuilder) error {
	scp, err := buildScopes(scopes...)
	if err != nil {
		return err
	}

	if atr.Scope != "" {
		scp = scp + " " + atr.Scope
	}

	atr.Scope = scp
	return nil
}

// AccessTokenResponse Possible responses that can be made to the /token endpoints. The type of response returned here will mirror the type of the request
type AccessTokenResponse struct {
	// String passed to the authorization server to gain access to the system
	AccessToken string `json:"access_token,omitempty"`
	// String passed to the authorization server to gain access to the system
	RefreshToken string `json:"refresh_token,omitempty"`
	// Type of token returned. Opaque is not parse-able while JWT follows JWT spec.
	TokenType string `json:"token_type,omitempty"`
	// Number of seconds the token is valid for.
	ExpiresIn int32 `json:"expires_in,omitempty"`
	// A space-delimited list of [scopes](https://docs.moov.io/guides/developer-tools/api-keys/scopes/) that are allowed
	Scope string `json:"scope,omitempty"`
}

// Makes the call for creating the access tokens
func (c *Client) accessToken(ctx context.Context, tokenReq AccessTokenRequest, scopes ...ScopeBuilder) (*AccessTokenResponse, error) {
	if err := tokenReq.WithScopes(scopes...); err != nil {
		return nil, err
	}

	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, "/oauth2/token"),
		JsonBody(tokenReq))
	if err != nil {
		return nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		return UnmarshalObjectResponse[AccessTokenResponse](resp)
	default:
		return nil, resp.Error()
	}
}

func (c *Client) RefreshAccessToken(ctx context.Context, refreshToken string) (*AccessTokenResponse, error) {
	return c.accessToken(ctx, AccessTokenRequest{
		GrantType:    "refresh_token",
		ClientId:     &c.Credentials.PublicKey,
		ClientSecret: &c.Credentials.SecretKey,
		RefreshToken: refreshToken,
	})
}

// Creates an access token that gives access to the ping endpoint. This allows for testing access tokens from server side
// or browser communications to ensure everything works.
func (c *Client) PingAccessToken(ctx context.Context) (*AccessTokenResponse, error) {
	return c.accessToken(ctx, AccessTokenRequest{
		GrantType:    "client_credentials",
		ClientId:     &c.Credentials.PublicKey,
		ClientSecret: &c.Credentials.SecretKey,
	}, Scopes.Ping())
}

// Creates a access token that can only be used for creating a new account thats connected to the client.
// This allows for a browser to create the account directly against Moov without that account information needing to
// send sensitive information through the clients backend services.
func (c *Client) AccountCreationToken(ctx context.Context) (*AccessTokenResponse, error) {
	return c.accessToken(ctx, AccessTokenRequest{
		GrantType:    "client_credentials",
		ClientId:     &c.Credentials.PublicKey,
		ClientSecret: &c.Credentials.SecretKey,
	}, Scopes.AccountsWrite())
}

// Creates an access token to access a connected account.
// This allows for a browser to access an account directly against Moov servers without that needing to send sensitive
// information through the clients backend services.
func (c *Client) AccessToken(ctx context.Context, scopes ...ScopeBuilder) (*AccessTokenResponse, error) {
	return c.accessToken(ctx, AccessTokenRequest{
		GrantType:    "client_credentials",
		ClientId:     &c.Credentials.PublicKey,
		ClientSecret: &c.Credentials.SecretKey,
	}, scopes...)
}
