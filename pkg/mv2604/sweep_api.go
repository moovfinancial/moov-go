package mv2604

import (
	"context"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

type SweepService struct {
	client *moov.Client
}

func NewSweepService(client *moov.Client) SweepService {
	return SweepService{client: client}
}

func (s SweepService) UpdateSweepConfig(ctx context.Context, accountID, sweepConfigID string, update UpdateSweepConfig) (*moov.SweepConfig, error) {
	return moov.UpdateSweepConfigGeneric(ctx, s.client, moov.Version2026_04, accountID, sweepConfigID, update)
}

type UpdateSweepConfig struct {
	moov.UpdateSweepConfig

	// An optional override of the default NACHA company entry description for sweep transfers.
	StatementDescriptor *moov.Nullable[string] `json:"statementDescriptor,omitempty"`
}
