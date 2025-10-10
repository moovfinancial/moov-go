package moov

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	moovgo "github.com/moovfinancial/moov-go"
)

func DefaultHttpClient() *http.Client {
	return &http.Client{
		Transport: http.DefaultTransport,
	}
}

func (c *Client) CallHttp(ctx context.Context, endpoint EndpointArg, args ...callArg) (CallResponse, error) {
	// Request a slot from the rate limiter
	c.waitForSlot(ctx)

	call, err := newCall(endpoint, args...)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s://%s%s", c.moovURLScheme, c.Credentials.Host, call.path)

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
	req.Header.Add("User-Agent", fmt.Sprintf("moov-go/%s", moovgo.Version()))

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

	decoder := standardDecoder
	if c.decoder != nil {
		decoder = c.decoder
	}

	return &httpCallResponse{
		resp: resp,
		body: body,

		decoder: decoder,
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

	decoder Decoder
}

func (r *httpCallResponse) Status() CallStatus {
	if r == nil || r.resp == nil {
		return StatusServerError
	}

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

func standardDecoder(r io.Reader, contentType string, item any) error {
	if strings.Contains(contentType, "application/json") {
		return json.NewDecoder(r).Decode(item)
	}
	return fmt.Errorf("unknown content-type %s", contentType)
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

	return r.decoder(bytes.NewReader(r.body), ct, item)
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

type errorResponse struct {
	Error string `json:"error"`
}

func (r *httpCallResponse) Error() string {
	generic := fmt.Sprintf("error from moov - status: %s http.request_id: %s http.status_code: %d", r.Status().Name, r.RequestId(), r.StatusCode())

	switch r.StatusCode() {
	case http.StatusBadRequest, http.StatusConflict, http.StatusUnprocessableEntity:
		// Check if the response contains an error in JSON
		if json.Valid(r.body) {
			var wrapper errorResponse
			json.Unmarshal(r.body, &wrapper)
			if wrapper.Error != "" {
				return fmt.Sprintf("%s\n  %s", generic, wrapper.Error)
			}

			// a JSON response like
			//  {"profile":{"business":{"taxID":{"ein":{"number":"must be a valid employer identification number"}}}}}
			// gets transformed into
			//  profile.business.taxID.ein.number: must be a valid employer identification number
			out := strings.TrimPrefix(string(r.body), `{"`)
			out = strings.ReplaceAll(out, `":{"`, ".")
			out = strings.ReplaceAll(out, `":"`, ": ")
			out = strings.ReplaceAll(out, `"}`, "\n")
			out = strings.ReplaceAll(out, "}", "")
			return fmt.Sprintf("%s - %v", generic, strings.TrimSpace(out))
		}
	}
	return generic
}
