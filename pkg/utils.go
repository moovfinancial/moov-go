package moov

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

type Credentials struct {
	AccountID string `yaml:"account_id,omitempty"`
	PublicKey string `yaml:"public_key,omitempty"`
	SecretKey string `yaml:"secret_key,omitempty"`
	Domain    string `yaml:"domain,omitempty"`
}

func readConfig() (Credentials, error) {
	var cred Credentials

	// Read the config file
	data, err := os.ReadFile("../config.yaml")
	if err != nil {
		return cred, err
	}

	// Unmarshal the YAML data into the cred struct
	if err := yaml.Unmarshal(data, &cred); err != nil {
		return cred, err
	}

	return cred, nil
}

// GetHTTPResponse performs an HTTP request and returns the response body or an error.
func GetHTTPResponse(c Client, method string, url string, data any, header map[string]string) ([]byte, int, error) {
	// Make an HTTP request
	var req *http.Request
	if data != nil {
		payload, err := json.Marshal(data)
		if err != nil {
			return nil, 0, err
		}

		req, err = http.NewRequest(method, url, bytes.NewBuffer(payload))
		if err != nil {
			return nil, 0, err
		}
	} else {
		var err error
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, 0, err
		}
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for key, val := range header {
		req.Header.Set(key, val)
	}
	req.SetBasicAuth(c.Credentials.PublicKey, c.Credentials.SecretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	log.Printf("fetch data from url:\n%s\n response:\n%s", url, string(body))

	return body, resp.StatusCode, nil
}
