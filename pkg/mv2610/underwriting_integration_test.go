package mv2610_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2610"
)

// newIntegrationClient builds a client against the real Moov backend using
// credentials from the environment (optionally loaded from secrets.env at the repo
// root). The test is skipped when credentials are unavailable so it stays CI-safe.
func newIntegrationClient(t *testing.T) *moov.Client {
	t.Helper()

	secretsPath := filepath.Join("..", "..", "secrets.env")
	if secrets, err := godotenv.Read(secretsPath); err == nil {
		for k, v := range secrets {
			t.Setenv(k, v)
		}
	}

	if os.Getenv("MOOV_PUBLIC_KEY") == "" || os.Getenv("MOOV_SECRET_KEY") == "" {
		t.Skip("MOOV_PUBLIC_KEY/MOOV_SECRET_KEY not set; skipping integration test")
	}

	client, err := moov.NewClient()
	require.NoError(t, err)

	if err := client.Ping(context.Background()); err != nil {
		t.Skipf("unable to reach Moov backend, skipping integration test: %v", err)
	}

	return client
}

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

	underwriting := mv2610.NewUnderwritingClient(client)

	upsert := mv2610.UpsertUnderwriting{
		GeographicReach:   moov.PtrOf(mv2610.GeographicReachUsOnly),
		BusinessPresence:  moov.PtrOf(mv2610.BusinessPresenceHomeBased),
		PendingLitigation: moov.PtrOf(mv2610.PendingLitigationNone),
		VolumeShareByCustomerType: &mv2610.VolumeShareByCustomerType{
			Business: moov.PtrOf(70),
			Consumer: moov.PtrOf(30),
			P2P:      moov.PtrOf(0),
		},
		SendFunds: &mv2610.SendFunds{
			Ach: &mv2610.SendFundsAch{
				EstimatedActivity: &mv2610.EstimatedActivity{MonthlyVolumeRange: moov.PtrOf(mv2610.MonthlyVolumeRangeUnder10K)},
			},
			Wire: &mv2610.SendFundsWire{
				EstimatedActivity: &mv2610.EstimatedActivity{MonthlyVolumeRange: moov.PtrOf(mv2610.MonthlyVolumeRange1M5M)},
			},
		},
		CardIssuing: &mv2610.CardIssuing{
			EstimatedActivity: &mv2610.EstimatedActivity{MonthlyVolumeRange: moov.PtrOf(mv2610.MonthlyVolumeRange50K100K)},
		},
		SubmissionIntent: moov.PtrOf(mv2610.SubmissionIntentWait),
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
		require.Equal(t, upsert.CardIssuing, actual.CardIssuing)
	})

	t.Run("get", func(t *testing.T) {
		actual, err := underwriting.GetUnderwriting(ctx, account.AccountID)
		require.NoError(t, err)

		require.NotNil(t, actual)
		require.Equal(t, upsert.GeographicReach, actual.GeographicReach)
		require.Equal(t, upsert.SendFunds, actual.SendFunds)
		require.Equal(t, upsert.CardIssuing, actual.CardIssuing)
	})
}
