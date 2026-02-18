package moov

import (
	"context"
	"net/http"
)

func (c Client) ListResolutionLinks(ctx context.Context, accountID string) ([]ResolutionLinkRecord, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathResolutionLinks, accountID))
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[ResolutionLinkRecord](resp)
}

func (c Client) GetResolutionLink(ctx context.Context, accountID string, resolutionLinkCode string) (*ResolutionLinkRecord, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathResolutionLink, accountID, resolutionLinkCode))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[ResolutionLinkRecord](resp)
}

func (c Client) CreateResolutionLink(ctx context.Context, accountID string) (*ResolutionLinkRecord, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathResolutionLinks, accountID))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[ResolutionLinkRecord](resp)
}

func (c Client) DeleteResolutionLink(ctx context.Context, accountID string, resolutionLinkCode string) error {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodDelete, pathResolutionLink, accountID, resolutionLinkCode))
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}
