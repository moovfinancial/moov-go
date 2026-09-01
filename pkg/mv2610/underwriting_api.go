package mv2610

import (
	"context"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

type UnderwritingClient struct {
	*moov.Client
}

func NewUnderwritingClient(client *moov.Client) UnderwritingClient {
	return UnderwritingClient{Client: client}
}

// GetUnderwriting returns the underwriting information for the given account.
// https://docs.moov.io/guides/accounts/requirements/underwriting/
func (u UnderwritingClient) GetUnderwriting(ctx context.Context, accountID string) (*Underwriting, error) {
	return moov.GetUnderwritingGeneric[Underwriting](ctx, u.Client, moov.Version2026_10, accountID)
}

// UpsertUnderwriting adds or updates underwriting information for the given account.
// Returns the underwriting information for the account.
// https://docs.moov.io/guides/accounts/requirements/underwriting/
func (u UnderwritingClient) UpsertUnderwriting(ctx context.Context, accountID string, underwriting UpsertUnderwriting) (*Underwriting, error) {
	return moov.UpsertUnderwritingGeneric[UpsertUnderwriting, Underwriting](ctx, u.Client, moov.Version2026_10, accountID, underwriting)
}
