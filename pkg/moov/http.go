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
var _ HttpCallResponse = &httpCallResponse{}

type HttpCallResponse interface {
	error

	Status() CallStatus
	Unmarshal(item any) error

	RequestId() string
	StatusCode() int
}

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

	if sb, ok := item.(*strings.Builder); ok {
		_, err := sb.Write(r.body)
		return err
	}

	if bb, ok := item.(*bytes.Buffer); ok {
		_, err := bb.Write(r.body)
		return err
	}

	if strings.Contains(ct, "json") {
		// content type checking here...
		return json.Unmarshal(r.body, item)
	}

	return fmt.Errorf("unknown content-type: %s", ct)
}

func (r *httpCallResponse) StatusCode() int {
	if r.resp != nil {
		return r.resp.StatusCode
	}
	return 0
}

func (r *httpCallResponse) RequestId() string {
	if r.resp != nil {
		return r.resp.Header.Get("X-Request-ID")
	}
	return ""
}

func (r *httpCallResponse) Error() string {
	return fmt.Sprintf("error from moov - status: %s http.request_id: %s http.status_code: %d", r.Status().Name, r.RequestId(), r.StatusCode())
}

// var _ HttpCallError = &httpCallError{}
// var _ CallResponseError = &httpCallError{}

// type httpCallError struct {
// 	resp *httpCallResponse

// 	requestId  string
// 	statusCode int
// }

// func (he *httpCallError) Status() CallStatus {
// 	return he.Status()
// }

// func (he *httpCallError) Unmarshal(item any) error {
// 	return he.Unmarshal(item)
// }

// func (he *httpCallError) RequestId() string {
// 	return he.requestId
// }

// func (he *httpCallError) StatusCode() int {
// 	return he.statusCode
// }

// func (he *httpCallError) Error() string {
// 	return fmt.Sprintf("error from moov - status: %s http.request_id: %s http.status_code: %d", he.Status().Name, he.requestId, he.statusCode)
// }
