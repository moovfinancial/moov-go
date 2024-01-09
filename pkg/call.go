package moov

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

// Request

type EndpointArg callArg

func Endpoint(method string, pathFmt string, args ...any) EndpointArg {
	return wrapCallFn(func(call *callBuilder) error {
		call.method = method
		call.path = fmt.Sprintf(pathFmt, args...)

		return nil
	})
}

// type callArg func(call *Call) error
type callArg interface {
	apply(call *callBuilder) error
}

type callArgInstance struct {
	fn func(call *callBuilder) error
}

func (ca *callArgInstance) apply(call *callBuilder) error {
	return ca.fn(call)
}

func wrapCallFn(fn func(call *callBuilder) error) callArg {
	return &callArgInstance{fn: fn}
}

// Prepends the following arguments onto the options
func prependArgs[A callArg](opts []A, args ...callArg) []callArg {
	for _, o := range opts {
		args = append(args, o)
	}

	return args
}

func JsonBody(body any) callArg {
	return wrapCallFn(func(call *callBuilder) error {
		payload, err := json.Marshal(body)
		if err != nil {
			return err
		}

		call.headers["Context-Type"] = "application/json"
		call.body = bytes.NewBuffer(payload)

		return nil
	})
}

func AcceptJson() callArg {
	return wrapCallFn(func(call *callBuilder) error {
		call.headers["Accept"] = "application/json"
		return nil
	})
}

func WaitFor(state string) callArg {
	return wrapCallFn(func(call *callBuilder) error {
		call.headers["X-Wait-For"] = state
		return nil
	})
}

// func Count(count int) CallArg {
// 	return wrapCallFn(func(call *callBuilder) error {
// 		call.params["count"] = strconv.Itoa(count)
// 		return nil
// 	})
// }

// func Skip(skip int) CallArg {
// 	return wrapCallFn(func(call *callBuilder) error {
// 		call.params["skip"] = strconv.Itoa(skip)
// 		return nil
// 	})
// }

// Response

type CallResponse interface {
	Status() CallStatus

	// Deserializes the body of the response into the item.
	// This is here so the response can handle any content type.
	Marshal(item any) error

	Error() error
}

func UnmarshalObjectResponse[A interface{}](resp CallResponse) (*A, error) {
	item := new(A)
	if err := resp.Marshal(item); err != nil {
		return nil, err
	}

	return item, nil
}

func UnmarshalListResponse[A interface{}](resp CallResponse) ([]A, error) {
	item := []A{}
	if err := resp.Marshal(&item); err != nil {
		return nil, err
	}

	return item, nil
}
