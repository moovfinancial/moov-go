package moov

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

func DefaultHttpClient() *http.Client {
	return &http.Client{
		Transport: http.DefaultTransport,
	}
}

// GetHTTPResponse performs an HTTP request and returns the response body or an error.
func GetHTTPResponse(c Client, method string, url string, data any, header map[string]string) ([]byte, int, error) {
	return CallHttp(context.Background(), c, method, url, data, header)
}

// GetHTTPResponse performs an HTTP request and returns the response body or an error.
func CallHttp(ctx context.Context, c Client, method string, url string, data any, header map[string]string) ([]byte, int, error) {

	reqBody, err := httpRequestBody(data)
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, 0, err
	}

	// by default send basic auth but allow the header to be overridden
	req.SetBasicAuth(c.Credentials.PublicKey, c.Credentials.SecretKey)

	for key, val := range header {
		req.Header.Set(key, val)
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	return body, resp.StatusCode, nil
}

func httpRequestBody(data any) (io.Reader, error) {
	if data != nil {
		payload, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		return bytes.NewBuffer(payload), nil
	}

	return nil, nil
}
