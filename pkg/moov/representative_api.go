package moov

import (
	"context"
	"errors"
	"net/http"
)

// CreateRepresentative creates a new representative for a given account
func (c Client) CreateRepresentative(ctx context.Context, accountID string, representative CreateRepresentative) (*Representative, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathRepresentatives, accountID),
		AcceptJson(),
		WaitFor("connection"),
		JsonBody(representative))
	if err != nil {
		return nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		return CompletedObjectOrError[Representative](resp)
	case StatusStateConflict:
		return nil, errors.Join(ErrAlreadyExists, resp)
	default:
		return nil, resp
	}
}

// ListRepresentatives returns the representatives for the given account
func (c Client) ListRepresentatives(ctx context.Context, accountID string) ([]Representative, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathRepresentatives, accountID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[Representative](resp)
}

// GetRepresentatives returns a representative for a given account
func (c Client) GetRepresentative(ctx context.Context, accountID string, representativeID string) (*Representative, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathRepresentative, accountID, representativeID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Representative](resp)
}

// UpdateRepresentative updates a representative for a given account
func (c Client) UpdateRepresentative(ctx context.Context, representativeAccountID string, representativeID string, representative UpdateRepresentative) (*Representative, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPatch, pathRepresentative, representativeAccountID, representativeID),
		AcceptJson(),
		JsonBody(representative))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Representative](resp)
}

// DeleteRepresentative deletes a representative for the given account
func (c Client) DeleteRepresentative(ctx context.Context, accountID string, representativeAccountID string) error {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodDelete, pathRepresentative, accountID, representativeAccountID))
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}
