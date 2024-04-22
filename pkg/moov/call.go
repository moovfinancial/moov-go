package moov

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
)

type CallStatus struct {
	Name      string
	Retryable bool
}

func callStatus(name string, retryable bool) CallStatus {
	return CallStatus{
		Name:      name,
		Retryable: retryable,
	}
}

var (
	StatusCompleted = callStatus("completed", false) // Completely fully
	StatusStarted   = callStatus("started", true)    // Returned as async. This can be due to timing out, or started as async

	StatusBadRequest       = callStatus("bad_request", false)       // bad request, body, headers, etc...
	StatusStateConflict    = callStatus("state_conflict", false)    // violates some stateful constraint.
	StatusFailedValidation = callStatus("failed_validation", false) // request structure is valid but failed validation.
	StatusNotFound         = callStatus("not_found", false)         // did not find the specific item.

	StatusUnauthenticated = callStatus("unauthenticated", false) // Token or credentials don't exist or are expired.
	StatusUnauthorized    = callStatus("unauthorized", false)    // Not allowed to see or invalid method or url

	StatusRateLimited = callStatus("rate_limited", true) // Rate limited and able to be retried after a little bit

	StatusServerError = callStatus("server_error", true) // Unknown server issue...
)

type callBuilder struct {
	method string
	path   string
	params map[string]string

	headers map[string]string
	token   *string

	body io.Reader
}

func newCall(endpoint EndpointArg, args ...callArg) (*callBuilder, error) {
	call := &callBuilder{
		params:  make(map[string]string),
		headers: make(map[string]string),
	}

	args = prependArgs(args, endpoint)

	for _, args := range args {
		if err := args.apply(call); err != nil {
			return nil, err
		}
	}

	return call, nil
}

type callArg interface {
	apply(call *callBuilder) error
}

type callArgInstance struct {
	fn func(call *callBuilder) error
}

func (ca *callArgInstance) apply(call *callBuilder) error {
	return ca.fn(call)
}

func callBuilderFn(fn func(call *callBuilder) error) callArg {
	return &callArgInstance{fn: fn}
}

// Prepends the following arguments onto the options
func prependArgs[A callArg](opts []A, args ...callArg) []callArg {
	for _, o := range opts {
		args = append(args, o)
	}

	return args
}

// Request

type EndpointArg callArg

func Endpoint(method string, pathFmt string, args ...any) EndpointArg {
	return callBuilderFn(func(call *callBuilder) error {
		call.method = method
		call.path = fmt.Sprintf(pathFmt, args...)

		return nil
	})
}

func JsonBody(body any) callArg {
	return callBuilderFn(func(call *callBuilder) error {
		payload, err := json.Marshal(body)
		if err != nil {
			return err
		}

		call.headers["Content-Type"] = "application/json"
		call.body = bytes.NewBuffer(payload)

		return nil
	})
}

type multipartFn func(w *multipart.Writer) error

func MultipartField(key, value string) multipartFn {
	return func(w *multipart.Writer) error {
		return w.WriteField(key, value)
	}
}

func MultipartFile(key, filename string, file io.Reader) multipartFn {
	return func(w *multipart.Writer) error {
		part, err := w.CreateFormFile(key, filename)
		if err != nil {
			return err
		}

		_, err = io.Copy(part, file)
		if err != nil {
			return err
		}

		return nil
	}
}

func MultipartBody(parts ...multipartFn) callArg {
	return callBuilderFn(func(call *callBuilder) error {
		buf := &bytes.Buffer{}
		w := multipart.NewWriter(buf)

		for _, part := range parts {
			if err := part(w); err != nil {
				return err
			}
		}

		if err := w.Close(); err != nil {
			return err
		}

		call.headers["Content-Type"] = w.FormDataContentType()
		call.body = buf

		return nil
	})
}

func NoopArg() callArg {
	return callBuilderFn(func(call *callBuilder) error {
		return nil
	})
}

func AcceptJson() callArg {
	return callBuilderFn(func(call *callBuilder) error {
		call.headers["Accept"] = "application/json"
		return nil
	})
}

func WaitFor(state string) callArg {
	return callBuilderFn(func(call *callBuilder) error {
		call.headers["X-Wait-For"] = state
		return nil
	})
}

func IdempotencyKey(uuid string) callArg {
	return callBuilderFn(func(call *callBuilder) error {
		call.headers["X-Idempotency-Key"] = uuid
		return nil
	})
}

func Skip(skip int) ListTransferFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["skip"] = fmt.Sprintf("%d", skip)
		return nil
	})
}

func Count(count int) ListTransferFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["count"] = fmt.Sprintf("%d", count)
		return nil
	})
}

func Limit(limit int) ListTransferFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["limit"] = fmt.Sprintf("%d", limit)
		return nil
	})
}

// Response

type CallResponse interface {
	Status() CallStatus

	// Deserializes the body of the response into the item.
	// This is here so the response can handle any content type.
	Unmarshal(item any) error

	// Convert response into an golang error
	Error() string
}

func UnmarshalObjectResponse[A interface{}](resp CallResponse) (*A, error) {
	item := new(A)
	if err := resp.Unmarshal(item); err != nil {
		return nil, err
	}
	return item, nil
}

func UnmarshalListResponse[A interface{}](resp CallResponse) ([]A, error) {
	item := []A{}
	if err := resp.Unmarshal(&item); err != nil {
		return nil, err
	}

	return item, nil
}

// Helper function for a common pattern of API calls that return no body so its either an error or not.
func CompletedNilOrError(resp CallResponse) error {
	switch resp.Status() {
	case StatusCompleted:
		return nil
	default:
		return resp
	}
}

// Helper for a common pattern of successful API calls returning an object body or an error
func CompletedObjectOrError[A interface{}](resp CallResponse) (*A, error) {
	switch resp.Status() {
	case StatusCompleted:
		return UnmarshalObjectResponse[A](resp)
	default:
		return nil, resp
	}
}

// Helper for a common pattern of successful API calls returning a body with a slice of objects or an error
func CompletedListOrError[A interface{}](resp CallResponse) ([]A, error) {
	switch resp.Status() {
	case StatusCompleted:
		return UnmarshalListResponse[A](resp)
	default:
		return nil, resp
	}
}
