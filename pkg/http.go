package moov

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func DefaultHttpClient() *http.Client {
	return &http.Client{
		Transport: http.DefaultTransport,
	}
}

// GetHTTPResponse performs an HTTP request and returns the response body or an error.
func (c *Client) GetHTTPResponse(method string, url string, data any, header map[string]string) ([]byte, int, error) {
	reqBody, err := httpRequestBody(data)
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequestWithContext(context.Background(), method, url, reqBody)
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

func (c *Client) CallHttp(ctx context.Context, endpoint EndpointArg, args ...callArg) (CallResponse, error) {
	call, err := newCall(endpoint, args...)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://%s%s", c.Credentials.Host, call.path)

	req, err := http.NewRequestWithContext(ctx, call.method, url, call.body)
	if err != nil {
		return nil, err
	}

	qry := req.URL.Query()
	for k, v := range call.params {
		qry.Add(k, v)
	}
	req.URL.RawQuery = qry.Encode()

	for k, v := range call.headers {
		req.Header.Add(k, v)
	}

	if call.token != nil {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", *call.token))
	} else {
		req.SetBasicAuth(c.Credentials.PublicKey, c.Credentials.SecretKey)
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	return &httpCallResponse{
		resp: resp,
		body: body,
	}, nil
}

var _ CallResponse = &httpCallResponse{}

type httpCallResponse struct {
	resp *http.Response
	body []byte
}

func (r *httpCallResponse) Status() CallStatus {
	switch r.resp.StatusCode {
	case http.StatusOK, http.StatusNoContent:
		return StatusCompleted
	case http.StatusCreated:
		return StatusStarted

	case http.StatusBadRequest:
		return StatusBadRequest
	case http.StatusConflict:
		return StatusStateConflict
	case http.StatusUnprocessableEntity:
		return StatusFailedValidation

	case http.StatusNotFound:
		return StatusNotFound
	case http.StatusUnauthorized:
		return StatusUnauthenticated
	case http.StatusForbidden:
		return StatusUnauthorized

	case http.StatusTooManyRequests:
		return StatusRateLimited

	default:
		return StatusServerError
	}
}

func (r *httpCallResponse) Unmarshal(item any) error {
	ct := strings.ToLower(r.resp.Header.Get("content-type"))

	if _, ok := item.(*string); ok {
		item = string(r.body)
		return nil
	}

	if _, ok := item.([]byte); ok {
		item = r.body
		return nil
	}

	if strings.Contains(ct, "json") {
		// content type checking here...
		return json.Unmarshal(r.body, item)
	}

	return fmt.Errorf("unknown content-type: %s", ct)
}

func (r *httpCallResponse) Error() error {
	switch r.Status() {
	case StatusCompleted, StatusStarted:
		return nil
	default:
		return &httpCallError{
			status:     r.Status(),
			requestId:  r.resp.Header.Get("X-Request-ID"),
			statusCode: r.resp.StatusCode,
		}
	}
}

var _ HttpCallError = &httpCallError{}

type HttpCallError interface {
	error
	Status() CallStatus
	RequestId() string
	StatusCode() int
}

type httpCallError struct {
	status     CallStatus
	requestId  string
	statusCode int
}

func (he *httpCallError) Status() CallStatus {
	return he.status
}

func (he *httpCallError) RequestId() string {
	return he.requestId
}

func (he *httpCallError) StatusCode() int {
	return he.statusCode
}

func (he *httpCallError) Error() string {
	return fmt.Sprintf("error from moov - status: %s http.request_id: %s http.status_code: %d", he.status.Name, he.requestId, he.statusCode)
}
