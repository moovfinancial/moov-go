package mv2607_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2607"
)

// newUnderwritingTestClient points a client at srv so tests can assert on the raw request.
func newUnderwritingTestClient(t *testing.T, srv *httptest.Server) mv2607.UnderwritingClient {
	t.Helper()

	client, err := moov.NewClient(
		moov.WithCredentials(moov.Credentials{PublicKey: "pk", SecretKey: "sk"}),
		moov.WithMoovURLScheme("http"),
	)
	require.NoError(t, err)
	client.Credentials.Host = strings.TrimPrefix(srv.URL, "http://")

	return mv2607.NewUnderwritingClient(client)
}

const underwritingWithWireJSON = `{
	"geographicReach": "us-only",
	"sendFunds": {
		"instantBank": {
			"estimatedActivity": {
				"monthlyVolumeRange": "10k-50k"
			}
		},
		"wire": {
			"estimatedActivity": {
				"averageTransactionAmount": 250000,
				"maximumTransactionAmount": 1000000,
				"monthlyVolumeRange": "1m-5m"
			}
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
		_, _ = w.Write([]byte(underwritingWithWireJSON))
	}))
	t.Cleanup(srv.Close)

	underwriting := newUnderwritingTestClient(t, srv)

	actual, err := underwriting.UpsertUnderwriting(context.Background(), "account-123", mv2607.UpsertUnderwriting{
		GeographicReach: moov.PtrOf(mv2607.GeographicReachUsOnly),
		SendFunds: &mv2607.SendFunds{
			InstantBank: &mv2607.SendFundsInstantBank{
				EstimatedActivity: &mv2607.EstimatedActivity{
					MonthlyVolumeRange: moov.PtrOf(mv2607.MonthlyVolumeRange10K50K),
				},
			},
			Wire: &mv2607.SendFundsWire{
				EstimatedActivity: &mv2607.EstimatedActivity{
					AverageTransactionAmount: moov.PtrOf(int64(250_000)),
					MaximumTransactionAmount: moov.PtrOf(int64(1_000_000)),
					MonthlyVolumeRange:       moov.PtrOf(mv2607.MonthlyVolumeRange1M5M),
				},
			},
		},
		SubmissionIntent: moov.PtrOf(mv2607.SubmissionIntentSubmit),
	})
	require.NoError(t, err)
	require.NoError(t, decodeErr)

	require.Equal(t, http.MethodPost, method)
	require.Equal(t, "/accounts/account-123/underwriting", path)
	require.Equal(t, moov.Version2026_07.String(), version)

	sendFunds, ok := body["sendFunds"].(map[string]any)
	require.True(t, ok)
	wire, ok := sendFunds["wire"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, map[string]any{
		"averageTransactionAmount": float64(250_000),
		"maximumTransactionAmount": float64(1_000_000),
		"monthlyVolumeRange":       "1m-5m",
	}, wire["estimatedActivity"])

	require.NotNil(t, actual)
	require.NotNil(t, actual.SendFunds)
	require.NotNil(t, actual.SendFunds.Wire)
	require.Equal(t, moov.PtrOf(mv2607.MonthlyVolumeRange1M5M), actual.SendFunds.Wire.EstimatedActivity.MonthlyVolumeRange)
	require.Equal(t, moov.PtrOf(int64(250_000)), actual.SendFunds.Wire.EstimatedActivity.AverageTransactionAmount)
	require.Equal(t, moov.PtrOf(int64(1_000_000)), actual.SendFunds.Wire.EstimatedActivity.MaximumTransactionAmount)
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
		_, _ = w.Write([]byte(underwritingWithWireJSON))
	}))
	t.Cleanup(srv.Close)

	underwriting := newUnderwritingTestClient(t, srv)

	actual, err := underwriting.GetUnderwriting(context.Background(), "account-123")
	require.NoError(t, err)

	require.Equal(t, http.MethodGet, method)
	require.Equal(t, "/accounts/account-123/underwriting", path)
	require.Equal(t, moov.Version2026_07.String(), version)

	require.NotNil(t, actual)
	require.Equal(t, moov.PtrOf(mv2607.GeographicReachUsOnly), actual.GeographicReach)
	require.NotNil(t, actual.SendFunds)
	require.NotNil(t, actual.SendFunds.InstantBank)
	require.NotNil(t, actual.SendFunds.Wire)
	require.Equal(t, moov.PtrOf(mv2607.MonthlyVolumeRange1M5M), actual.SendFunds.Wire.EstimatedActivity.MonthlyVolumeRange)
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

	_, err = underwriting.UpsertUnderwriting(context.Background(), "", mv2607.UpsertUnderwriting{})
	require.EqualError(t, err, "accountID is required")

	require.False(t, called, "no request should reach the server")
}

// SendFunds must not emit a wire key when the rail isn't set, so pre-v2026.07 shapes
// round-trip unchanged.
func TestSendFundsOmitsWireWhenUnset(t *testing.T) {
	b, err := json.Marshal(mv2607.SendFunds{
		InstantBank: &mv2607.SendFundsInstantBank{},
	})
	require.NoError(t, err)
	require.NotContains(t, string(b), `"wire"`)
}
