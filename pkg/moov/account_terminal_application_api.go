package moov

import (
	"context"
	"net/http"
)

// LinkAccountTerminalApplication links an account to a terminal application.
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

	return CompletedObjectOrError[AccountTerminalApplication](resp)
}

// GetAccountTerminalApplication gets an account-terminal application link.
func (c Client) GetAccountTerminalApplication(ctx context.Context, accountID, terminalApplicationID string) (*AccountTerminalApplication, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathAccountTerminalApplication, accountID, terminalApplicationID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[AccountTerminalApplication](resp)
}

// ListAccountTerminalApplications lists account-terminal application links for an account.
func (c Client) ListAccountTerminalApplications(ctx context.Context, accountID string) ([]AccountTerminalApplication, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathAccountTerminalApplications, accountID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[AccountTerminalApplication](resp)
}

// GetAccountTerminalApplicationConfiguration gets configuration for a linked terminal application.
func (c Client) GetAccountTerminalApplicationConfiguration(ctx context.Context, accountID, terminalApplicationID string) (*AccountTerminalApplicationConfiguration, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathAccountTerminalApplicationConfiguration, accountID, terminalApplicationID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[AccountTerminalApplicationConfiguration](resp)
}
