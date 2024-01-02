package moov

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) ScopedAccessToken(accountID string, scope []string) (ClientCredentialsGrantToAccessTokenResponse, error) {
	token := ClientCredentialsGrantToAccessTokenResponse{}
	renderedScope := make([]string, len(scope))
	for i := 0; i < len(scope); i++ {
		renderedScope[i] = strings.Replace(scope[i], "{accountID}", accountID, 1)
	}
	params := url.Values{}
	params.Add("grant_type", "client_credentials")
	params.Add("scope", strings.Join(renderedScope, " "))
	req, err := http.NewRequest("POST", "https://api.moov.io/oauth2/token?"+params.Encode(), nil)
	if err != nil {
		err := errors.Join(err, ErrURL)
		return token, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.Credentials.PublicKey, c.Credentials.SecretKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// *url.Error
		return token, err
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return token, err
	}
	if err := json.Unmarshal(resBody, &token); err != nil {
		return token, err
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
		err := errors.Join(err, ErrURL)
		return token, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	req.SetBasicAuth(c.Credentials.PublicKey, c.Credentials.SecretKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// *url.Error
		return token, err
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return token, err
	}
	if err := json.Unmarshal(resBody, &token); err != nil { // Parse []byte to go struct pointer
		return token, err
	}
	return token, nil
}
