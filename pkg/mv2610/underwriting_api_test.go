package mv2610_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2610"
)

// newUnderwritingTestClient points a client at srv so tests can assert on the raw request.
func newUnderwritingTestClient(t *testing.T, srv *httptest.Server) mv2610.UnderwritingClient {
	t.Helper()

	client, err := moov.NewClient(
		moov.WithCredentials(moov.Credentials{PublicKey: "pk", SecretKey: "sk"}),
		moov.WithMoovURLScheme("http"),
	)
	require.NoError(t, err)
	client.Credentials.Host = strings.TrimPrefix(srv.URL, "http://")

	return mv2610.NewUnderwritingClient(client)
}

const underwritingWithCardIssuingJSON = `{
	"geographicReach": "us-only",
	"sendFunds": {
		"wire": {
			"estimatedActivity": {
				"monthlyVolumeRange": "1m-5m"
			}
		}
	},
	"cardIssuing": {
		"estimatedActivity": {
			"averageTransactionAmount": 5000,
			"maximumTransactionAmount": 50000,
			"monthlyVolumeRange": "50k-100k"
		}
	}
}`

func TestUpsertUnderwriting(t *testing.T) {
	var (
		method    string
		path      string
		version   string
		body      map[string]any
		decodeErr error
	)

	// The handler runs on its own goroutine, so it records what it saw rather than
	// asserting: require.* calls t.FailNow, which is only valid on the test goroutine.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		path = r.URL.Path
		version = r.Header.Get(moov.VersionHeader)
		decodeErr = json.NewDecoder(r.Body).Decode(&body)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(underwritingWithCardIssuingJSON))
	}))
	t.Cleanup(srv.Close)

	underwriting := newUnderwritingTestClient(t, srv)

	actual, err := underwriting.UpsertUnderwriting(context.Background(), "account-123", mv2610.UpsertUnderwriting{
		GeographicReach: moov.PtrOf(mv2610.GeographicReachUsOnly),
		CardIssuing: &mv2610.CardIssuing{
			EstimatedActivity: &mv2610.EstimatedActivity{
				AverageTransactionAmount: moov.PtrOf(int64(5_000)),
				MaximumTransactionAmount: moov.PtrOf(int64(50_000)),
				MonthlyVolumeRange:       moov.PtrOf(mv2610.MonthlyVolumeRange50K100K),
			},
		},
		SubmissionIntent: moov.PtrOf(mv2610.SubmissionIntentSubmit),
	})
	require.NoError(t, err)
	require.NoError(t, decodeErr)

	require.Equal(t, http.MethodPost, method)
	require.Equal(t, "/accounts/account-123/underwriting", path)
	require.Equal(t, moov.Version2026_10.String(), version)

	cardIssuing, ok := body["cardIssuing"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, map[string]any{
		"averageTransactionAmount": float64(5_000),
		"maximumTransactionAmount": float64(50_000),
		"monthlyVolumeRange":       "50k-100k",
	}, cardIssuing["estimatedActivity"])

	require.NotNil(t, actual)
	require.NotNil(t, actual.CardIssuing)
	require.NotNil(t, actual.CardIssuing.EstimatedActivity)
	require.Equal(t, moov.PtrOf(mv2610.MonthlyVolumeRange50K100K), actual.CardIssuing.EstimatedActivity.MonthlyVolumeRange)
	require.Equal(t, moov.PtrOf(int64(5_000)), actual.CardIssuing.EstimatedActivity.AverageTransactionAmount)
	require.Equal(t, moov.PtrOf(int64(50_000)), actual.CardIssuing.EstimatedActivity.MaximumTransactionAmount)
}

func TestGetUnderwriting(t *testing.T) {
	var (
		method  string
		path    string
		version string
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		path = r.URL.Path
		version = r.Header.Get(moov.VersionHeader)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(underwritingWithCardIssuingJSON))
	}))
	t.Cleanup(srv.Close)

	underwriting := newUnderwritingTestClient(t, srv)

	actual, err := underwriting.GetUnderwriting(context.Background(), "account-123")
	require.NoError(t, err)

	require.Equal(t, http.MethodGet, method)
	require.Equal(t, "/accounts/account-123/underwriting", path)
	require.Equal(t, moov.Version2026_10.String(), version)

	require.NotNil(t, actual)
	require.Equal(t, moov.PtrOf(mv2610.GeographicReachUsOnly), actual.GeographicReach)
	require.NotNil(t, actual.SendFunds)
	require.NotNil(t, actual.SendFunds.Wire)
	require.NotNil(t, actual.CardIssuing)
	require.Equal(t, moov.PtrOf(mv2610.MonthlyVolumeRange50K100K), actual.CardIssuing.EstimatedActivity.MonthlyVolumeRange)
}

// An empty accountID would produce /accounts//underwriting, so both calls must fail
// before reaching the network.
func TestUnderwritingRequiresAccountID(t *testing.T) {
	called := false

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	t.Cleanup(srv.Close)

	underwriting := newUnderwritingTestClient(t, srv)

	_, err := underwriting.GetUnderwriting(context.Background(), "")
	require.EqualError(t, err, "accountID is required")

	_, err = underwriting.UpsertUnderwriting(context.Background(), "", mv2610.UpsertUnderwriting{})
	require.EqualError(t, err, "accountID is required")

	require.False(t, called, "no request should reach the server")
}

// UpsertUnderwriting must not emit a cardIssuing key when it isn't set, so pre-v2026.10
// shapes round-trip unchanged.
func TestUpsertOmitsCardIssuingWhenUnset(t *testing.T) {
	b, err := json.Marshal(mv2610.UpsertUnderwriting{
		SendFunds: &mv2610.SendFunds{Wire: &mv2610.SendFundsWire{}},
	})
	require.NoError(t, err)
	require.NotContains(t, string(b), `"cardIssuing"`)
}
