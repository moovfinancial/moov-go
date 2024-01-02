package moov

import (
	"net/http"
)

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
		return ErrAuthCredentialsNotSet
	}
	return nil
}
