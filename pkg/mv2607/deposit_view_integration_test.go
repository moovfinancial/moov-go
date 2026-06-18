package mv2607_test

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2607"
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

func TestCreateDepositView_Integration(t *testing.T) {
	client := newIntegrationClient(t)
	ctx := context.Background()

	// Create a temporary account to attach the deposit view to.
	created, started, err := client.CreateAccount(ctx, moov.CreateAccount{
		Type: moov.AccountType_Business,
		Profile: moov.CreateProfile{
			Business: &moov.CreateBusinessProfile{
				Name:        "moov-go deposit view SDK test",
				Type:        moov.BusinessType_Llc,
				Description: "moov-go SDK deposit view integration test",
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

	strPtr := func(s string) *string { return &s }
	numPtr := func(s string) *json.Number { n := json.Number(s); return &n }
	intPtr := func(i int) *int { return &i }

	document, err := json.Marshal(jhSilverlakeRecord{
		SrcKey:    strPtr("src-1"),
		FornKey:   strPtr("cust-1"),
		AcctId:    numPtr("987654321"),
		AcctType:  strPtr("DDA"),
		CurBal:    numPtr("1234.56"),
		AvlBal:    numPtr("1000.00"),
		NumCrMTD:  intPtr(3),
		NumDrMTD:  intPtr(7),
		BrandCode: strPtr("MOOV"),
	})
	require.NoError(t, err)

	depositView := mv2607.NewDepositViewClient(client)

	resp, err := depositView.CreateDepositView(ctx, account.AccountID, mv2607.SourceSystemJHSilverlake, document)
	// The deposit view endpoint is access-gated; partners without it (e.g. the
	// generic CI credentials) receive 403. Skip rather than fail in that case.
	if httpResp, ok := err.(moov.HttpCallResponse); ok && httpResp.StatusCode() == http.StatusForbidden {
		t.Skipf("credentials not authorized for the deposit view endpoint, skipping: %v", err)
	}
	require.NoError(t, err)

	require.NotNil(t, resp)
	require.Equal(t, account.AccountID, resp.MoovAccountID)
	require.Equal(t, mv2607.SourceSystemJHSilverlake, resp.SourceSystem)
	require.Equal(t, "987654321", resp.SourceAccountID)
	require.False(t, resp.IngestedAt.IsZero())
}
