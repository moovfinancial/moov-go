package moov

import (
	"context"
	"net/http"
)

// CreateTerminalApplication creates a new terminal application.
func (c Client) LinkAccountTerminalApplication(ctx context.Context, accountID, terminalApplicationID string) (*AccountTerminalApplication, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathAccountTerminalApplications, accountID),
		AcceptJson(),
		JsonBody(LinkAccountTerminalApplicationRequest{
			TerminalApplicationID: terminalApplicationID,
		}))
	if err != nil {
		return nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		a, err := UnmarshalObjectResponse[AccountTerminalApplication](resp)
		return a, err
	default:
		return nil, resp
	}
}

// GetTerminalApplication returns a terminal application based on terminalApplicationID.
func (c Client) GetAccountTerminalApplication(ctx context.Context, accountID, terminalApplicationID string) (*AccountTerminalApplication, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathAccountTerminalApplication, accountID, terminalApplicationID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[AccountTerminalApplication](resp)
}

// ListTerminalApplications returns a list of terminalApplications.
func (c Client) ListAccountTerminalApplications(ctx context.Context, accountID string) ([]AccountTerminalApplication, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathAccountTerminalApplications, accountID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[AccountTerminalApplication](resp)
}

// GetTerminalApplication returns a terminal application based on terminalApplicationID.
func (c Client) GetAccountTerminalApplicationConfiguration(ctx context.Context, accountID, terminalApplicationID string) (*AccountTerminalApplicationConfiguration, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathAccountTerminalApplicationConfiguration, accountID, terminalApplicationID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[AccountTerminalApplicationConfiguration](resp)
}
