package moov

import (
	"context"
	"net/http"
)

// CreateOnboardingInvite creates a new onboarding invite.
func (c Client) CreateOnboardingInvite(ctx context.Context, invite OnboardingInviteRequest) (*OnboardingInvite, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathOnboardingInvites),
		AcceptJson(),
		JsonBody(invite))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[OnboardingInvite](resp)
}

// ListOnboardingInvites lists all onboarding invites.
func (c Client) ListOnboardingInvites(ctx context.Context) ([]OnboardingInvite, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathOnboardingInvites),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[OnboardingInvite](resp)
}

// GetOnboardingInvite retrieves an onboarding invite by its code.
func (c Client) GetOnboardingInvite(ctx context.Context, code string) (*OnboardingInvite, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathOnboardingInvite, code),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[OnboardingInvite](resp)
}

// RevokeOnboardingInvite revokes an onboarding invite by its code.
func (c Client) RevokeOnboardingInvite(ctx context.Context, code string) error {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodDelete, pathOnboardingInvite, code))
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}
