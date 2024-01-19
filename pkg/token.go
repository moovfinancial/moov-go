package moov

import (
	"context"
	"net/http"
	"strings"
)

const (
	SCOPE_ACCOUNTS_READ  = "/accounts.read"
	SCOPE_ACCOUNTS_WRITE = "/accounts.write"

	// TODO: Add the rest of the scopes https://docs.moov.io/guides/developer-tools/scopes/
	SCOPE_BANK_ACCOUNT_READ       = "/accounts/{accountID}/bank-accounts.read"
	SCOPE_BANK_ACCOUNT_WRITE      = "/accounts/{accountID}/bank-accounts.write"
	SCOPE_CARDS_READ              = "/accounts/{accountID}/cards.read"
	SCOPE_CARDS_WRITE             = "/accounts/{accountID}/cards.write"
	SCOPE_CAPABILITIES_READ       = "/accounts/{accountID}/capabilities.read"
	SCOPE_CAPABILITIES_WRITE      = "/accounts/{accountID}/capabilities.write"
	SCOPE_DOCUMENTS_READ          = "/accounts/{accountID}/documents.read"
	SCOPE_DOCUMENTS_WRITE         = "/accounts/{accountID}/documents.write"
	SCOPE_FILES_READ              = "/accounts/{accountID}/files.read"
	SCOPE_FILES_WRITE             = "/accounts/{accountID}/files.write"
	SCOPE_PAYMENT_METHODS_READ    = "/accounts/{accountID}/payment-methods.read"
	SCOPE_PAYMENT_METHODS_WRITE   = "/accounts/{accountID}/payment-methods.write"
	SCOPE_PROFILE_READ            = "/accounts/{accountID}/profile.read"
	SCOPE_PROFILE_WRITE           = "/accounts/{accountID}/profile.write"
	SCOPE_PREFILE_ENRICHMENT_READ = "/accounts/{accountID}/prefile-enrichment.read"
	SCOPE_REPRESENTATIVES_READ    = "/accounts/{accountID}/representatives.read"
	SCOPE_REPRESENTATIVES_WRITE   = "/accounts/{accountID}/representatives.write"
	SCOPE_TRANSFERS_READ          = "/accounts/{accountID}/transfers.read"
	SCOPE_TRANSFERS_WRITE         = "/accounts/{accountID}/transfers.write"
	SCOPE_WALLETS_READ            = "/accounts/{accountID}/wallets.read"
	SCOPE_FED_READ                = "/fed.read"
	SCOPE_PING_READ               = "/ping.read"
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
func (c *Client) AccessToken(ctx context.Context, tokenReq AccessTokenRequest) (*AccessTokenResponse, error) {
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
	return c.AccessToken(ctx, AccessTokenRequest{
		GrantType:    "refresh_token",
		ClientId:     &c.Credentials.PublicKey,
		ClientSecret: &c.Credentials.SecretKey,
		RefreshToken: refreshToken,
	})
}

// Creates an access token that gives access to the ping endpoint. This allows for testing access tokens from server side
// or browser communications to ensure everything works.
func (c *Client) PingAccessToken(ctx context.Context) (*AccessTokenResponse, error) {
	return c.AccessToken(ctx, AccessTokenRequest{
		GrantType:    "client_credentials",
		ClientId:     &c.Credentials.PublicKey,
		ClientSecret: &c.Credentials.SecretKey,
		Scope:        SCOPE_PING_READ,
	})
}

// Creates a access token that can only be used for creating a new account thats connected to the client.
// This allows for a browser to create the account directly against Moov without that account information needing to
// send sensitive information through the clients backend services.
func (c *Client) AccountCreationToken(ctx context.Context) (*AccessTokenResponse, error) {
	return c.AccessToken(ctx, AccessTokenRequest{
		GrantType:    "client_credentials",
		ClientId:     &c.Credentials.PublicKey,
		ClientSecret: &c.Credentials.SecretKey,
		Scope:        SCOPE_ACCOUNTS_WRITE,
	})
}

// Creates an access token to access a connected account.
// This allows for a browser to access an account directly against Moov servers without that needing to send sensitive
// information through the clients backend services.
func (c *Client) AccountAccessToken(ctx context.Context, accountID string, scopes []string) (*AccessTokenResponse, error) {
	renderedScopes := strings.Join(scopes, " ")
	renderedScopes = strings.ReplaceAll(renderedScopes, "{accountID}", accountID)

	return c.AccessToken(ctx, AccessTokenRequest{
		GrantType:    "client_credentials",
		ClientId:     &c.Credentials.PublicKey,
		ClientSecret: &c.Credentials.SecretKey,
		Scope:        renderedScopes,
	})
}
