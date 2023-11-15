package moov

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	// SCOPE_ACCOUNTS_READ allows a new Moov account to be created
	SCOPE_ACCOUNTS_READ = "/accounts.read"
	// SCOPE_ACCOUNTS_WRITE View all connected accounts
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

var (
	ErrAuthCreditionalsNotSet = errors.New("API Keys are not set or invalid credentials")
	ErrAuthNetwork            = errors.New("network error")
)

// New create4s a new Moov client with the appropriate secret key.

type Client struct {
	Credentials Credentials
}

type Credentials struct {
	// AccountID Facilitator account ID
	AccountID string `json:"accountID,omitempty"`
	// PubliocKey Public key value from API key
	PublicKey string `json:"publicKey,omitempty"`
	// SecretKey Secret key value from API key
	SecretKey string `json:"secretKey,omitempty"`
	// Domain One of the domains from API key
	Domain string `json:"domain,omitempty"`
}

type ClientCredentialsGrantToAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int32  `json:"expires_in"`
	Scope        string `json:"scope"`
}

func NewClient(creds Credentials) (*Client, error) {
	nc := &Client{
		Credentials: Credentials{
			AccountID: creds.AccountID,
			PublicKey: creds.PublicKey,
			SecretKey: creds.SecretKey,
			Domain:    creds.Domain,
		},
	}

	if creds.PublicKey == "" || creds.SecretKey == "" {
		// Make error for token's not set.
		return nc, ErrAuthCreditionalsNotSet
	}

	// Ping the server to make sure we have valid credentials
	err := nc.Ping()
	if err != nil {
		return nc, err
	}

	return nc, nil
}

func (c *Client) ScopedAccessToken(accountID string, scope []string) (ClientCredentialsGrantToAccessTokenResponse, error) {
	token := ClientCredentialsGrantToAccessTokenResponse{}
	renderedScope := make([]string, len(scope))
	for i := 0; i < len(scope); i++ {
		// TODO: check if string is null or empty
		renderedScope[i] = strings.Replace(scope[i], "{accountID}", accountID, 1)
	}
	params := url.Values{}
	params.Add("grant_type", "client_credentials")
	params.Add("scope", strings.Join(renderedScope, " "))
	req, err := http.NewRequest("POST", "https://api.moov.io/oauth2/token?"+params.Encode(), nil)
	log.Println("https://api.moov.io/oauth2/token?" + params.Encode())
	if err != nil {
		// Todo: return an error
		log.Fatal(err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	req.SetBasicAuth(c.Credentials.PublicKey, c.Credentials.SecretKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// Todo: return an error
		log.Fatal(err)
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		// Todo: return an error
		log.Fatal(err)
	}

	if err := json.Unmarshal(resBody, &token); err != nil { // Parse []byte to go struct pointer
		// Todo: return an error
		log.Fatal("Can not unmarshal JSON")
	}
	return token, nil
}

// SingleUseAccessToken gets a token that can be used to instantiate Moov.js client side
func (c Client) SingleUseAccessToken() (ClientCredentialsGrantToAccessTokenResponse, error) {
	token := ClientCredentialsGrantToAccessTokenResponse{}
	params := url.Values{}
	params.Add("grant_type", "client_credentials")
	params.Add("scope", "/accounts.write")
	req, err := http.NewRequest("POST", "https://api.moov.io/oauth2/token?"+params.Encode(), nil)
	if err != nil {
		// Todo: return an error
		log.Fatal(err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	req.SetBasicAuth(c.Credentials.PublicKey, c.Credentials.SecretKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// Todo: return an error
		log.Fatal(err)
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		// Todo: return an error
		log.Fatal(err)
	}

	if err := json.Unmarshal(resBody, &token); err != nil { // Parse []byte to go struct pointer
		// Todo: return an error
		log.Fatal("Can not unmarshal JSON")
	}
	return token, nil
}

// Ping calls the ping endpoint to make sure we have valid credentials
func (c Client) Ping() error {
	req, _ := http.NewRequest(http.MethodGet, "https://api.moov.io/ping", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.Credentials.PublicKey, c.Credentials.SecretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ErrAuthNetwork
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized:
		return ErrAuthCreditionalsNotSet
	}
	return nil
}