package moov

import (
	"context"
	"fmt"
	"net/http"
)

// CreateOnboardingInvite creates a new onboarding invite. The returned invite
// contains a unique link the recipient can follow to onboard a Moov account.
// https://docs.moov.io/api/money-movement/onboarding/create/
func (c Client) CreateOnboardingInvite(ctx context.Context, invite OnboardingInviteRequest) (*OnboardingInvite, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathOnboardingInvites),
		AcceptJson(),
		JsonBody(invite),
	)
	if err != nil {
		return nil, fmt.Errorf("calling http: %w", err)
	}

	return CompletedObjectOrError[OnboardingInvite](resp)
}

// GetOnboardingInvite retrieves an onboarding invite by its code.
// https://docs.moov.io/api/money-movement/onboarding/get/
func (c Client) GetOnboardingInvite(ctx context.Context, code string) (*OnboardingInvite, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathOnboardingInvite, code),
		AcceptJson(),
	)
	if err != nil {
		return nil, fmt.Errorf("calling http: %w", err)
	}

	return CompletedObjectOrError[OnboardingInvite](resp)
}

// ListOnboardingInvites lists all onboarding invites created by the caller.
// https://docs.moov.io/api/money-movement/onboarding/list/
func (c Client) ListOnboardingInvites(ctx context.Context) ([]OnboardingInvite, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathOnboardingInvites),
		AcceptJson(),
	)
	if err != nil {
		return nil, fmt.Errorf("calling http: %w", err)
	}

	return CompletedListOrError[OnboardingInvite](resp)
}

// RevokeOnboardingInvite revokes an onboarding invite by its code. Once revoked,
// an invite can no longer be redeemed.
// https://docs.moov.io/api/money-movement/onboarding/revoke/
func (c Client) RevokeOnboardingInvite(ctx context.Context, code string) error {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodDelete, pathOnboardingInvite, code),
	)
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}
