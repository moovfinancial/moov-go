package moov

import (
	"context"
	"net/http"
)

// CreateTerminalApplication creates a new terminal application.
func (c Client) CreateTerminalApplication(ctx context.Context, terminalApplication TerminalApplicationRequest) (*TerminalApplication, *TerminalApplication, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathTerminalApplications),
		AcceptJson(),
		WaitFor("connection"),
		JsonBody(terminalApplication))
	if err != nil {
		return nil, nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		a, err := UnmarshalObjectResponse[TerminalApplication](resp)
		return a, nil, err
	case StatusStarted:
		a, err := UnmarshalObjectResponse[TerminalApplication](resp)
		return nil, a, err
	default:
		return nil, nil, resp
	}
}

// GetTerminalApplication returns a terminal application based on terminalApplicationID.
func (c Client) GetTerminalApplication(ctx context.Context, terminalApplicationID string) (*TerminalApplication, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathTerminalApplication, terminalApplicationID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[TerminalApplication](resp)
}

// ListTerminalApplications returns a list of terminalApplications.
func (c Client) ListTerminalApplications(ctx context.Context) ([]TerminalApplication, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathTerminalApplications),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[TerminalApplication](resp)
}

// DeleteTerminalApplication deletes a terminal application.
func (c Client) DeleteTerminalApplication(ctx context.Context, terminalApplicationID string) error {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodDelete, pathTerminalApplication, terminalApplicationID))
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}
