package mv2607_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2607"
)

func TestUpsertUnderwriting_Integration(t *testing.T) {
	client := newIntegrationClient(t)
	ctx := context.Background()

	created, started, err := client.CreateAccount(ctx, moov.CreateAccount{
		Type: moov.AccountType_Business,
		Profile: moov.CreateProfile{
			Business: &moov.CreateBusinessProfile{
				Name:        "moov-go underwriting SDK test",
				Type:        moov.BusinessType_Llc,
				Description: "moov-go SDK underwriting integration test",
				IndustryCodes: &moov.IndustryCodes{
					Mcc:   "6012",
					Naics: "522110",
					Sic:   "6021",
				},
				Industry: "electronics-appliances",
			},
		},
	})
	require.NoError(t, err)

	account := created
	if account == nil {
		account = started
	}
	require.NotNil(t, account)

	t.Cleanup(func() {
		_ = client.DisconnectAccount(context.Background(), account.AccountID)
	})

	underwriting := mv2607.NewUnderwritingClient(client)

	upsert := mv2607.UpsertUnderwriting{
		GeographicReach:   moov.PtrOf(mv2607.GeographicReachUsOnly),
		BusinessPresence:  moov.PtrOf(mv2607.BusinessPresenceHomeBased),
		PendingLitigation: moov.PtrOf(mv2607.PendingLitigationNone),
		VolumeShareByCustomerType: &mv2607.VolumeShareByCustomerType{
			Business: moov.PtrOf(70),
			Consumer: moov.PtrOf(30),
			P2P:      moov.PtrOf(0),
		},
		SendFunds: &mv2607.SendFunds{
			Ach: &mv2607.SendFundsAch{
				EstimatedActivity: &mv2607.EstimatedActivity{MonthlyVolumeRange: moov.PtrOf(mv2607.MonthlyVolumeRangeUnder10K)},
			},
			InstantBank: &mv2607.SendFundsInstantBank{
				EstimatedActivity: &mv2607.EstimatedActivity{MonthlyVolumeRange: moov.PtrOf(mv2607.MonthlyVolumeRangeUnder10K)},
			},
			Wire: &mv2607.SendFundsWire{
				EstimatedActivity: &mv2607.EstimatedActivity{MonthlyVolumeRange: moov.PtrOf(mv2607.MonthlyVolumeRange1M5M)},
			},
		},
		SubmissionIntent: moov.PtrOf(mv2607.SubmissionIntentWait),
	}

	t.Run("upsert", func(t *testing.T) {
		actual, err := underwriting.UpsertUnderwriting(ctx, account.AccountID, upsert)
		require.NoError(t, err)

		require.NotNil(t, actual)
		require.Equal(t, upsert.GeographicReach, actual.GeographicReach)
		require.Equal(t, upsert.BusinessPresence, actual.BusinessPresence)
		require.Equal(t, upsert.PendingLitigation, actual.PendingLitigation)
		require.Equal(t, upsert.VolumeShareByCustomerType, actual.VolumeShareByCustomerType)
		require.Equal(t, upsert.SendFunds, actual.SendFunds)
	})

	t.Run("get", func(t *testing.T) {
		actual, err := underwriting.GetUnderwriting(ctx, account.AccountID)
		require.NoError(t, err)

		require.NotNil(t, actual)
		require.Equal(t, upsert.GeographicReach, actual.GeographicReach)
		require.Equal(t, upsert.SendFunds, actual.SendFunds)
	})
}
