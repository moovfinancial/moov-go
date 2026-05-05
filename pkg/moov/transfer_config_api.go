package moov

import (
	"context"
	"net/http"
)

// CreateTransferConfig creates a transfer config for the specified account.
func (c Client) CreateTransferConfig(ctx context.Context, accountID string, config UpsertTransferConfig) (*TransferConfig, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathTransferConfig, accountID),
		MoovVersion(Version2026_04),
		AcceptJson(),
		JsonBody(config),
	)
	if err != nil {
		return nil, err
	}

	return StartedObjectOrError[TransferConfig](resp)
}

// GetTransferConfig retrieves the transfer config for the specified account.
func (c Client) GetTransferConfig(ctx context.Context, accountID string) (*TransferConfig, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathTransferConfig, accountID),
		MoovVersion(Version2026_04),
		AcceptJson(),
	)
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[TransferConfig](resp)
}

// UpdateTransferConfig replaces the transfer config for the specified account.
func (c Client) UpdateTransferConfig(ctx context.Context, accountID string, config UpsertTransferConfig) (*TransferConfig, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPut, pathTransferConfig, accountID),
		MoovVersion(Version2026_04),
		AcceptJson(),
		JsonBody(config),
	)
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[TransferConfig](resp)
}
