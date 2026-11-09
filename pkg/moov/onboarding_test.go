package moov_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

func Test_OnboardingInvite_CreateGetListRevoke(t *testing.T) {
	var (
		mc  = NewTestClient(t)
		ctx = t.Context()
	)

	create := moov.OnboardingInviteRequest{
		ReturnURL: "https://moov.io",
		Scopes: []moov.ApplicationScope{
			"accounts.read",
		},
		Capabilities: []moov.CapabilityName{
			moov.CapabilityName_Transfers,
			moov.CapabilityName_Wallet,
		},
		FeePlanCodes: []string{"merchant-direct"},
		Prefill: func() *moov.CreateAccount {
			account := createTestBusinessAccount()
			return &account
		}(),
	}

	// create the invite
	created, err := mc.CreateOnboardingInvite(ctx, create)
	require.NoError(t, err)
	require.NotEmpty(t, created.Code)
	require.NotEmpty(t, created.Link)
	require.Equal(t, create.ReturnURL, created.ReturnURL)

	code := string(created.Code)
	t.Cleanup(func() {
		_ = mc.RevokeOnboardingInvite(ctx, code)
	})

	// fetch by code
	fetched, err := mc.GetOnboardingInvite(ctx, code)
	require.NoError(t, err)
	require.Equal(t, created.Code, fetched.Code)

	// list invites and confirm the created one is present
	listed, err := mc.ListOnboardingInvites(ctx)
	require.NoError(t, err)
	require.Contains(t, listed, *created)

	// revoke the invite
	err = mc.RevokeOnboardingInvite(ctx, code)
	require.NoError(t, err)

	// revoked invite can still be fetched and has RevokedOn set
	fetched, err = mc.GetOnboardingInvite(ctx, code)
	require.NoError(t, err)
	require.NotNil(t, fetched.RevokedOn)
}
