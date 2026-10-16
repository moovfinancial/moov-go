package mv2701_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2701"
)

func newCardIssuingTestClient(t *testing.T, srv *httptest.Server) mv2701.CardIssuingClient {
	t.Helper()

	client, err := moov.NewClient(
		moov.WithCredentials(moov.Credentials{PublicKey: "pk", SecretKey: "sk"}),
		moov.WithMoovURLScheme("http"),
	)
	require.NoError(t, err)
	client.Credentials.Host = strings.TrimPrefix(srv.URL, "http://")

	return mv2701.NewCardIssuingClient(client)
}

const issuedCardJSON = `{
	"issuedCardID": "card-1",
	"brand": "Visa",
	"lastFourCardNumber": "1234",
	"expiration": {
		"month": "01",
		"year": "30"
	},
	"nickname": "Travel card",
	"fundingWalletID": "wallet-1",
	"state": "active",
	"formFactor": "virtual",
	"controls": {
		"expiresOn": "2027-06-01T00:00:00Z"
	},
	"createdOn": "2027-01-15T12:00:00Z",
	"updatedOn": "2027-01-15T12:00:00Z"
}`

func TestCreateIssuedCard(t *testing.T) {
	var (
		method  string
		path    string
		version string
		accept  string
		body    map[string]any
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		path = r.URL.Path
		version = r.Header.Get(moov.VersionHeader)
		accept = r.Header.Get("Accept")
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(issuedCardJSON))
	}))
	t.Cleanup(srv.Close)

	expiresOn := time.Date(2027, 6, 1, 0, 0, 0, 0, time.UTC)

	actual, err := newCardIssuingTestClient(t, srv).CreateIssuedCard(context.Background(), "account-123", mv2701.CreateIssuedCard{
		Nickname: moov.PtrOf("Travel card"),
		Controls: &moov.IssuingControls{
			ExpiresOn: moov.PtrOf(expiresOn),
		},
	})
	require.NoError(t, err)

	require.Equal(t, http.MethodPost, method)
	require.Equal(t, "/issuing/account-123/cards", path)
	require.Equal(t, moov.Version2027_01.String(), version)
	require.Equal(t, "application/json", accept)

	require.NotContains(t, body, "expiration")
	require.Equal(t, "Travel card", body["nickname"])
	require.Equal(t, map[string]any{"expiresOn": "2027-06-01T00:00:00Z"}, body["controls"])

	require.NotNil(t, actual)
	require.Equal(t, "card-1", actual.IssuedCardID)
	require.Equal(t, moov.IssuedCardExpiration{Month: "01", Year: "30"}, actual.Expiration)
	require.NotNil(t, actual.Controls)
	require.Equal(t, moov.PtrOf(expiresOn), actual.Controls.ExpiresOn)
}

func TestCreateIssuedCard_NilClient(t *testing.T) {
	actual, err := mv2701.NewCardIssuingClient(nil).CreateIssuedCard(context.Background(), "account-123", mv2701.CreateIssuedCard{})
	require.Nil(t, actual)
	require.EqualError(t, err, "client is nil")
}

func TestCreateIssuedCard_RequiresAccountID(t *testing.T) {
	called := false

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	t.Cleanup(srv.Close)

	actual, err := newCardIssuingTestClient(t, srv).CreateIssuedCard(context.Background(), "", mv2701.CreateIssuedCard{})
	require.Nil(t, actual)
	require.EqualError(t, err, "accountID is required")
	require.False(t, called, "no request should reach the server")
}
