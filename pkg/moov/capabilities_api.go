package moov

import (
	"context"
	"net/http"
)

// RequestCapabilities adds a new capability for the given account
func (c Client) RequestCapabilities(ctx context.Context, accountID string, capabilities []CapabilityName) ([]Capability, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathCapabilities, accountID),
		AcceptJson(),
		JsonBody(requestCapabilities{
			Capabilities: capabilities,
		}))
	if err != nil {
		return nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		return CompletedListOrError[Capability](resp)
	default:
		return nil, resp
	}
}

// ListCapabilities returns the capabilities for the given account
func (c Client) ListCapabilities(ctx context.Context, accountID string) ([]Capability, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathCapabilities, accountID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[Capability](resp)
}

// GetCapability returns a given capability for a given account
func (c Client) GetCapability(ctx context.Context, accountID string, capability CapabilityName) (*Capability, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathCapability, accountID, capability),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Capability](resp)
}

// DisableCapability disables a specific capability
func (c Client) DisableCapability(ctx context.Context, accountID string, capability CapabilityName) error {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodDelete, pathCapability, accountID, capability))
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}
