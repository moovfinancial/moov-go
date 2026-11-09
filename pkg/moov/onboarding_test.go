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
		Prefill:      moov.PtrOf(createTestBusinessAccount()),
	}

	created, err := mc.CreateOnboardingInvite(ctx, create)
	require.NoError(t, err)
	require.NotEmpty(t, created.Code)
	require.NotEmpty(t, created.Link)
	require.Equal(t, create.ReturnURL, created.ReturnURL)

	code := string(created.Code)
	t.Cleanup(func() {
		_ = mc.RevokeOnboardingInvite(ctx, code)
	})

	fetched, err := mc.GetOnboardingInvite(ctx, code)
	require.NoError(t, err)
	require.Equal(t, created.Code, fetched.Code)

	listed, err := mc.ListOnboardingInvites(ctx)
	require.NoError(t, err)
	require.Contains(t, listed, *created)

	err = mc.RevokeOnboardingInvite(ctx, code)
	require.NoError(t, err)

	fetched, err = mc.GetOnboardingInvite(ctx, code)
	require.NoError(t, err)
	require.NotNil(t, fetched.RevokedOn)
}
