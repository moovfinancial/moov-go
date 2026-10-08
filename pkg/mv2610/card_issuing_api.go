package mv2610

import (
	"context"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

// CardIssuingClient calls the card issuing API pinned to moov.Version2026_10.
type CardIssuingClient struct{ *moov.Client }

// NewCardIssuingClient wraps client in a CardIssuingClient.
func NewCardIssuingClient(client *moov.Client) CardIssuingClient {
	return CardIssuingClient{Client: client}
}

// ListIssuedCardActivity lists issued card activity for the given account.
func (c CardIssuingClient) ListIssuedCardActivity(ctx context.Context, accountID string, filters ...moov.ListIssuedCardActivityFilter) ([]IssuedCardActivity, error) {
	return moov.ListIssuedCardActivityGeneric[IssuedCardActivity](ctx, c.Client, moov.Version2026_10, accountID, filters...)
}

// SimulateAuthorization creates a simulated authorization for an issued card in test mode.
// Exactly one of the returned authorization, started response, or error is non-nil.
func (c CardIssuingClient) SimulateAuthorization(ctx context.Context, accountID string, simulation CreateAuthorizationSimulation) (*moov.IssuedCardAuthorization, *AuthorizationSimulationAsyncResponse, error) {
	return moov.SimulateIssuedCardAuthorizationGeneric[CreateAuthorizationSimulation, AuthorizationSimulationAsyncResponse](ctx, c.Client, moov.Version2026_10, accountID, simulation)
}

// SimulateClearing creates a simulated clearing for an issued card authorization in test mode.
// Exactly one of the returned authorization, started response, or error is non-nil.
func (c CardIssuingClient) SimulateClearing(ctx context.Context, accountID string, authorizationID string, simulation CreateClearingSimulation) (*moov.IssuedCardAuthorization, *AuthorizationSimulationAsyncResponse, error) {
	return moov.SimulateIssuedCardClearingGeneric[CreateClearingSimulation, AuthorizationSimulationAsyncResponse](ctx, c.Client, moov.Version2026_10, accountID, authorizationID, simulation)
}

// SimulateReversal creates a simulated reversal for an issued card authorization in test mode.
// Exactly one of the returned authorization, started response, or error is non-nil.
func (c CardIssuingClient) SimulateReversal(ctx context.Context, accountID string, authorizationID string) (*moov.IssuedCardAuthorization, *AuthorizationSimulationAsyncResponse, error) {
	return moov.SimulateIssuedCardReversalGeneric[AuthorizationSimulationAsyncResponse](ctx, c.Client, moov.Version2026_10, accountID, authorizationID)
}
