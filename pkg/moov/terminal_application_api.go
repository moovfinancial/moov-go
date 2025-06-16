package moov

import (
	"context"
	"net/http"
)

// CreateTerminalApplication creates a new terminal application.
func (c Client) CreateTerminalApplication(ctx context.Context, terminalApplication TerminalApplicationRequest) (*TerminalApplication, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathTerminalApplications),
		AcceptJson(),
		JsonBody(terminalApplication))
	if err != nil {
		return nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		a, err := UnmarshalObjectResponse[TerminalApplication](resp)
		return a, err
	default:
		return nil, resp
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

func (c Client) CreateTerminalApplicationVersion(ctx context.Context, terminalApplicationID, version string) (*TerminalApplicationVersion, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathTerminalApplicationVersions, terminalApplicationID),
		AcceptJson(),
		JsonBody(TerminalApplicationVersion{
			Version: version,
		}))
	if err != nil {
		return nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		a, err := UnmarshalObjectResponse[TerminalApplicationVersion](resp)
		return a, err
	default:
		return nil, resp
	}
}
