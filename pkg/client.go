package moov

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

// New create4s a new Moov client with the appropriate secret key.
//func New(private)

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

func NewClient(creds Credentials) *Client {
	nc := &Client{
		Credentials: Credentials{
			AccountID: creds.AccountID,
			PublicKey: creds.PublicKey,
			SecretKey: creds.SecretKey,
			Domain:    creds.Domain,
		},
	}
	return nc
}

// BasicAuth calls
func (c Client) BasicAuthToken() (ClientCredentialsGrantToAccessTokenResponse, error) {
	token := ClientCredentialsGrantToAccessTokenResponse{}
	if c.Credentials.PublicKey == "" || c.Credentials.SecretKey == "" {
		// Make error for token's not set.
		return token, fmt.Errorf("API Keys are not set")
	}

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
