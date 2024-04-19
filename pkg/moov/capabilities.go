package moov

import (
	"context"
	"errors"
	"net/http"
	"time"
)

var (
	ErrAccount = errors.New("no account with the specified accountID was found")
)

const (
	CAPABILITY_TRANSFERS     = "transfers"
	CAPABILITY_SEND_FUNDS    = "send-funds"
	CAPABILITY_COLLECT_FUNDS = "collect-funds"
	CAPABILITY_WALLET        = "wallet"
	CAPABILITY_CARD_ISSUING  = "card-issuing"
	CAPABILITY_ENBABLED      = "enabled"
	CAPABILITY_DISABLED      = "disabled"
	CAPABILITY_PENDING       = "pending"
)

// Capabilities a list of CAPABILITY_*
var Capabilities []string

type Capability struct {
	Capability   string `json:"capability"`
	AccountID    string `json:"accountID"`
	Status       string `json:"status,omitempty"`
	Requirements struct {
		CurrentlyDue []string `json:"currentlyDue,omitempty"`
		Errors       []struct {
			Requirement string `json:"requirement,omitempty"`
			ErrorCode   string `json:"errorCode,omitempty"`
		}
	}
	DisabledReason string    `json:"disabledReason,omitempty"`
	CreatedOn      time.Time `json:"createdOn,omitempty"`
	UpdatedOn      time.Time `json:"updatedOn,omitempty"`
	DisabledOn     time.Time `json:"disabledOn,omitempty"`
}

// RequestCapabilities adds a new capability for the given account
func (c Client) RequestCapabilities(ctx context.Context, accountID string, capabilities []Capability) ([]Capability, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, "/accounts/%s/capabilities", accountID),
		AcceptJson(),
		JsonBody(capabilities))
	if err != nil {
		return nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		return CompletedListOrError[Capability](resp)
	case StatusStateConflict:
		return nil, ErrAccount
	default:
		return nil, resp
	}
}

// ListCapabilities returns the capabilities for the given account
func (c Client) ListCapabilities(ctx context.Context, accountID string) ([]Capability, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, "/accounts/%s/capabilities", accountID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[Capability](resp)
}

// GetCapability returns a given capability for a given account
func (c Client) GetCapability(ctx context.Context, accountID string, capability string) (*Capability, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, "/accounts/%s/capabilities/%s", accountID, capability),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Capability](resp)
}

// DisableCapability disables a specific capability
func (c Client) DisableCapability(ctx context.Context, accountID string, capability string) error {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodDelete, "/accounts/%s/capabilities/%s", accountID, capability))
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}
