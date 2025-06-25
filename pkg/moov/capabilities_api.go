package moov

import (
	"context"
	"net/http"
)

type CapabilityClient[T any, V any] struct {
	Version Version
}

// Request adds one or more new capability to the given account
func (cc CapabilityClient[T, V]) Request(ctx context.Context, client Client, accountID string, requested T) ([]V, error) {
	resp, err := client.CallHttp(ctx,
		Endpoint(http.MethodPost, pathCapabilities, accountID),
		MoovVersion(cc.Version),
		AcceptJson(),
		JsonBody(requested))
	if err != nil {
		return nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		return CompletedListOrError[V](resp)
	default:
		return nil, resp
	}
}

// List returns the capabilities for the given account and version
func (cc CapabilityClient[T, V]) List(ctx context.Context, client Client, accountID string) ([]V, error) {
	resp, err := client.CallHttp(ctx,
		Endpoint(http.MethodGet, pathCapabilities, accountID),
		MoovVersion(cc.Version),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[V](resp)
}

// Get returns a given capability for a given account
func (cc CapabilityClient[T, V]) Get(ctx context.Context, client Client, accountID string, capability CapabilityName) (*V, error) {
	resp, err := client.CallHttp(ctx,
		Endpoint(http.MethodGet, pathCapability, accountID, capability),
		MoovVersion(cc.Version),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[V](resp)
}

// Disable disables a specific capability for a given account
func (cc CapabilityClient[T, V]) Disable(ctx context.Context, client Client, accountID string, capability CapabilityName) error {
	resp, err := client.CallHttp(ctx,
		Endpoint(http.MethodDelete, pathCapability, accountID, capability),
		MoovVersion(cc.Version))
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}

// Legacy

// Use only for Preversioned API calls. Use mvxxxx.Capabilities.Request(...) instead.
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

// Use only for Preversioned API calls. Use mvxxxx.Capabilities.List(...) instead.
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

// Use only for Preversioned API calls. Use mvxxxx.Capabilities.Get(...) instead.
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

// Use only for Preversioned API calls. Use mvxxxx.Capabilities.Disable(...) instead.
// DisableCapability disables a specific capability
func (c Client) DisableCapability(ctx context.Context, accountID string, capability CapabilityName) error {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodDelete, pathCapability, accountID, capability))
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}
