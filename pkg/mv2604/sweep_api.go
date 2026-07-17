package mv2604

import (
	"context"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

type SweepClient struct {
	*moov.Client
}

func NewSweepClient(client *moov.Client) SweepClient {
	return SweepClient{Client: client}
}

func (s SweepClient) ListSweepConfigs(ctx context.Context, accountID string) ([]moov.SweepConfig, error) {
	return moov.ListSweepConfigsGeneric(ctx, s.Client, moov.Version2026_04, accountID)
}

func (s SweepClient) GetSweepConfig(ctx context.Context, accountID string, sweepConfigID string) (*moov.SweepConfig, error) {
	return moov.GetSweepConfigGeneric(ctx, s.Client, moov.Version2026_04, accountID, sweepConfigID)
}

func (s SweepClient) CreateSweepConfig(ctx context.Context, create moov.CreateSweepConfig) (*moov.SweepConfig, error) {
	return moov.CreateSweepConfigGeneric(ctx, s.Client, moov.Version2026_04, create)
}

func (s SweepClient) UpdateSweepConfig(ctx context.Context, accountID, sweepConfigID string, update UpdateSweepConfig) (*moov.SweepConfig, error) {
	return moov.UpdateSweepConfigGeneric(ctx, s.Client, moov.Version2026_04, accountID, sweepConfigID, update)
}

type UpdateSweepConfig struct {
	Status              *moov.SweepConfigStatus `json:"status,omitempty"`
	PushPaymentMethodID *string                 `json:"pushPaymentMethodID,omitempty"`
	PullPaymentMethodID *string                 `json:"pullPaymentMethodID,omitempty"`
	// An optional field to specify an amount to maintain in the wallet. This is a decimal-formatted numerical string that represents up to 2 decimal place precision. In USD for example, 12.34 is $12.34 and 0.99 is $0.99. If not supplied, the default is 0.00.
	MinimumBalance *string `json:"minimumBalance,omitempty"`

	// An optional override of the default NACHA company entry description for sweep transfers.
	StatementDescriptor *moov.Nullable[string] `json:"statementDescriptor,omitempty"`
}
