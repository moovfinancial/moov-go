package mv2610_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2610"
)

func newCardIssuingTestClient(t *testing.T, srv *httptest.Server) mv2610.CardIssuingClient {
	t.Helper()

	client, err := moov.NewClient(
		moov.WithCredentials(moov.Credentials{PublicKey: "pk", SecretKey: "sk"}),
		moov.WithMoovURLScheme("http"),
	)
	require.NoError(t, err)
	client.Credentials.Host = strings.TrimPrefix(srv.URL, "http://")

	return mv2610.NewCardIssuingClient(client)
}

const issuedCardActivityJSON = `[
	{
		"authorizationID": "auth-1",
		"status": "declined",
		"issuedCardID": "card-1",
		"lastFourCardNumber": "1234",
		"authorizedUserAccountID": "user-1",
		"authorizedAmount": "12.34",
		"declineReason": "insufficient-funds",
		"createdOn": "2026-09-01T12:00:00Z",
		"merchantData": {
			"networkID": "net-1",
			"name": "Coffee Shop",
			"country": "US",
			"mcc": "5814"
		}
	},
	{
		"cardTransactionID": "txn-1",
		"status": "cleared",
		"issuedCardID": "card-1",
		"clearedAmount": "5.00",
		"createdOn": "2026-09-02T12:00:00Z",
		"merchantData": {
			"networkID": "net-2",
			"country": "US",
			"mcc": "5411"
		}
	}
]`

func TestListIssuedCardActivity(t *testing.T) {
	var (
		method  string
		path    string
		version string
		accept  string
		query   url.Values
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		path = r.URL.Path
		version = r.Header.Get(moov.VersionHeader)
		accept = r.Header.Get("Accept")
		query = r.URL.Query()

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(issuedCardActivityJSON))
	}))
	t.Cleanup(srv.Close)

	start := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 9, 30, 0, 0, 0, 750*int(time.Millisecond), time.UTC)

	actual, err := newCardIssuingTestClient(t, srv).ListIssuedCardActivity(context.Background(), "account-123",
		moov.WithIssuedCardActivitySkip(10),
		moov.WithIssuedCardActivityCount(20),
		moov.WithIssuedCardActivityCardID("card-1"),
		moov.WithIssuedCardActivityAuthorizedUserAccountID("user-1"),
		moov.WithIssuedCardActivityStatuses([]moov.IssuedCardAuthorizationStatus{
			moov.IssuedCardAuthorizationStatus_Declined,
			moov.IssuedCardAuthorizationStatus_Cleared,
		}),
		moov.WithIssuedCardActivityStartDate(start),
		moov.WithIssuedCardActivityEndDate(end),
	)
	require.NoError(t, err)

	require.Equal(t, http.MethodGet, method)
	require.Equal(t, "/issuing/account-123/activity", path)
	require.Equal(t, moov.Version2026_10.String(), version)
	require.Equal(t, "application/json", accept)
	require.Equal(t, url.Values{
		"skip":                    {"10"},
		"count":                   {"20"},
		"issuedCardID":            {"card-1"},
		"authorizedUserAccountID": {"user-1"},
		"statuses":                {"declined,cleared"},
		"startDateTime":           {"2026-09-01T00:00:00Z"},
		"endDateTime":             {"2026-09-30T00:00:00.75Z"},
	}, query)

	require.Len(t, actual, 2)

	declined := actual[0]
	require.Equal(t, moov.PtrOf("auth-1"), declined.AuthorizationID)
	require.Nil(t, declined.CardTransactionID)
	require.Equal(t, moov.IssuedCardAuthorizationStatus_Declined, declined.Status)
	require.Equal(t, "card-1", declined.IssuedCardID)
	require.Equal(t, moov.PtrOf("1234"), declined.LastFourCardNumber)
	require.Equal(t, moov.PtrOf("user-1"), declined.AuthorizedUserAccountID)
	require.Equal(t, moov.PtrOf("12.34"), declined.AuthorizedAmount)
	require.Nil(t, declined.ClearedAmount)
	require.Equal(t, moov.PtrOf(moov.IssuingDeclineReason_InsufficientFunds), declined.DeclineReason)
	require.Equal(t, start.Add(12*time.Hour), declined.CreatedOn)
	require.Equal(t, moov.IssuedCardTransactionMerchant{
		NetworkID: "net-1",
		Name:      moov.PtrOf("Coffee Shop"),
		Country:   "US",
		Mcc:       "5814",
	}, declined.MerchantData)

	cleared := actual[1]
	require.Nil(t, cleared.AuthorizationID)
	require.Equal(t, moov.PtrOf("txn-1"), cleared.CardTransactionID)
	require.Equal(t, moov.IssuedCardAuthorizationStatus_Cleared, cleared.Status)
	require.Equal(t, "card-1", cleared.IssuedCardID)
	require.Nil(t, cleared.LastFourCardNumber)
	require.Nil(t, cleared.AuthorizedUserAccountID)
	require.Nil(t, cleared.AuthorizedAmount)
	require.Equal(t, moov.PtrOf("5.00"), cleared.ClearedAmount)
	require.Nil(t, cleared.DeclineReason)
	require.Equal(t, start.Add(36*time.Hour), cleared.CreatedOn)
	require.Equal(t, moov.IssuedCardTransactionMerchant{
		NetworkID: "net-2",
		Country:   "US",
		Mcc:       "5411",
	}, cleared.MerchantData)
}

func TestListIssuedCardActivity_NoFilters(t *testing.T) {
	var rawQuery string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawQuery = r.URL.RawQuery

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	t.Cleanup(srv.Close)

	actual, err := newCardIssuingTestClient(t, srv).ListIssuedCardActivity(context.Background(), "account-123")
	require.NoError(t, err)
	require.Empty(t, actual)
	require.Empty(t, rawQuery)
}

func TestListIssuedCardActivity_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(srv.Close)

	actual, err := newCardIssuingTestClient(t, srv).ListIssuedCardActivity(context.Background(), "account-123")
	require.Nil(t, actual)

	var httpErr moov.HttpCallResponse
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, moov.StatusNotFound, httpErr.Status())
}

func TestListIssuedCardActivity_NilClient(t *testing.T) {
	actual, err := mv2610.NewCardIssuingClient(nil).ListIssuedCardActivity(context.Background(), "account-123")
	require.Nil(t, actual)
	require.EqualError(t, err, "client is nil")
}

func TestListIssuedCardActivity_RequiresAccountID(t *testing.T) {
	called := false

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	t.Cleanup(srv.Close)

	actual, err := newCardIssuingTestClient(t, srv).ListIssuedCardActivity(context.Background(), "")
	require.Nil(t, actual)
	require.EqualError(t, err, "accountID is required")
	require.False(t, called, "no request should reach the server")
}
