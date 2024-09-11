package moov

import (
	"context"
	"fmt"
	"net/http"
)

// ListSweepConfigs lists sweep configs associated with the accountID
// https://docs.moov.io/api/money-movement/sweeps/list-configs/
func (c Client) ListSweepConfigs(ctx context.Context, accountID string) ([]SweepConfig, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathSweepConfigs, accountID), AcceptJson())
	if err != nil {
		return nil, fmt.Errorf("listing sweep configs: %v", err)
	}

	return CompletedListOrError[SweepConfig](resp)
}

// GetSweepConfig retrieves a sweep config for the given accountID and sweepConfigID
// https://docs.moov.io/api/money-movement/sweeps/get-config/
func (c Client) GetSweepConfig(ctx context.Context, accountID string, sweepConfigID string) (*SweepConfig, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathSweepConfig, accountID, sweepConfigID), AcceptJson())
	if err != nil {
		return nil, fmt.Errorf("getting sweep config: %v", err)
	}

	return CompletedObjectOrError[SweepConfig](resp)
}

type CreateSweepConfig struct {
	AccountID           string            `json:"-"`
	WalletID            string            `json:"walletID"`
	Status              SweepConfigStatus `json:"status"`
	PushPaymentMethodID string            `json:"pushPaymentMethodID"`
	PullPaymentMethodID string            `json:"pullPaymentMethodID"`
	// An optional override of the default NACHA company entry description for sweep transfers.
	StatementDescriptor *string `json:"statementDescriptor,omitempty"`
	// An optional field to specify an amount to maintain in the wallet. This is a decimal-formatted numerical string that represents up to 2 decimal place precision. In USD for example, 12.34 is $12.34 and 0.99 is $0.99. If not supplied, the default is 0.00.
	MinimumBalance *string `json:"minimumBalance,omitempty"`
}

// CreateSweepConfig creates a sweep config for a wallet.
// https://docs.moov.io/api/money-movement/sweeps/create-config/
func (c Client) CreateSweepConfig(ctx context.Context, create CreateSweepConfig) (*SweepConfig, error) {
	resp, err := c.CallHttp(
		ctx,
		Endpoint(http.MethodPost, pathSweepConfigs, create.AccountID),
		AcceptJson(),
		JsonBody(create),
	)
	if err != nil {
		return nil, fmt.Errorf("creating sweep config: %v", err)
	}

	return CompletedObjectOrError[SweepConfig](resp)
}

type UpdateSweepConfig struct {
	AccountID           string             `json:"-"`
	SweepConfigID       string             `json:"-"`
	Status              *SweepConfigStatus `json:"status,omitempty"`
	PushPaymentMethodID *string            `json:"pushPaymentMethodID,omitempty"`
	PullPaymentMethodID *string            `json:"pullPaymentMethodID,omitempty"`
	// An optional override of the default NACHA company entry description for sweep transfers.
	StatementDescriptor *string `json:"statementDescriptor,omitempty"`
	// An optional field to specify an amount to maintain in the wallet. This is a decimal-formatted numerical string that represents up to 2 decimal place precision. In USD for example, 12.34 is $12.34 and 0.99 is $0.99. If not supplied, the default is 0.00.
	MinimumBalance *string `json:"minimumBalance,omitempty"`
}

// UpdateSweepConfig updates a sweep config.
// https://docs.moov.io/api/money-movement/sweeps/patch-config/
func (c Client) UpdateSweepConfig(ctx context.Context, update UpdateSweepConfig) (*SweepConfig, error) {
	resp, err := c.CallHttp(
		ctx,
		Endpoint(http.MethodPatch, pathSweepConfig, update.AccountID, update.SweepConfigID),
		AcceptJson(),
		JsonBody(update),
	)
	if err != nil {
		return nil, fmt.Errorf("updating sweep config: %v", err)
	}

	return CompletedObjectOrError[SweepConfig](resp)
}
