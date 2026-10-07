package mv2610_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2610"
)

const issuedCardAuthorizationJSON = `{
	"authorizationID": "auth-1",
	"issuedCardID": "card-1",
	"fundingWalletID": "wallet-1",
	"network": "shazam",
	"authorizedAmount": "12.34",
	"status": "pending",
	"merchantData": {
		"networkID": "net-1",
		"name": "Coffee Shop",
		"country": "US",
		"mcc": "5814"
	},
	"createdOn": "2026-10-01T12:00:00Z"
}`

type simulationCall func(c mv2610.CardIssuingClient) (*moov.IssuedCardAuthorization, *mv2610.AuthorizationSimulationAsyncResponse, error)

var simulationCases = []struct {
	name     string
	path     string
	wantBody string
	call     simulationCall
}{
	{
		name:     "authorization",
		path:     "/issuing/simulations/account-123/authorizations",
		wantBody: `{"issuedCardID":"card-1","amount":"12.34","merchantData":{"networkID":"net-1","name":"Coffee Shop","country":"US","mcc":"5814"}}`,
		call: func(c mv2610.CardIssuingClient) (*moov.IssuedCardAuthorization, *mv2610.AuthorizationSimulationAsyncResponse, error) {
			return c.SimulateAuthorization(context.Background(), "account-123", mv2610.CreateAuthorizationSimulation{
				IssuedCardID: "card-1",
				Amount:       "12.34",
				MerchantData: &moov.IssuedCardTransactionMerchant{
					NetworkID: "net-1",
					Name:      moov.PtrOf("Coffee Shop"),
					Country:   "US",
					Mcc:       "5814",
				},
			})
		},
	},
	{
		name:     "clearing",
		path:     "/issuing/simulations/account-123/authorizations/auth-1/clearings",
		wantBody: `{"amount":"5.00"}`,
		call: func(c mv2610.CardIssuingClient) (*moov.IssuedCardAuthorization, *mv2610.AuthorizationSimulationAsyncResponse, error) {
			return c.SimulateClearing(context.Background(), "account-123", "auth-1", mv2610.CreateClearingSimulation{
				Amount: moov.PtrOf("5.00"),
			})
		},
	},
	{
		name:     "clearing without amount",
		path:     "/issuing/simulations/account-123/authorizations/auth-1/clearings",
		wantBody: `{}`,
		call: func(c mv2610.CardIssuingClient) (*moov.IssuedCardAuthorization, *mv2610.AuthorizationSimulationAsyncResponse, error) {
			return c.SimulateClearing(context.Background(), "account-123", "auth-1", mv2610.CreateClearingSimulation{})
		},
	},
	{
		name:     "reversal",
		path:     "/issuing/simulations/account-123/authorizations/auth-1/reversals",
		wantBody: "",
		call: func(c mv2610.CardIssuingClient) (*moov.IssuedCardAuthorization, *mv2610.AuthorizationSimulationAsyncResponse, error) {
			return c.SimulateReversal(context.Background(), "account-123", "auth-1")
		},
	},
}

func TestSimulations_Completed(t *testing.T) {
	for _, tc := range simulationCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				method  string
				path    string
				version string
				body    string
			)

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				b, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				method = r.Method
				path = r.URL.Path
				version = r.Header.Get(moov.VersionHeader)
				body = string(b)

				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(issuedCardAuthorizationJSON))
			}))
			t.Cleanup(srv.Close)

			authorization, started, err := tc.call(newCardIssuingTestClient(t, srv))
			require.NoError(t, err)
			require.Nil(t, started)

			require.Equal(t, http.MethodPost, method)
			require.Equal(t, tc.path, path)
			require.Equal(t, moov.Version2026_10.String(), version)
			if tc.wantBody == "" {
				require.Empty(t, body)
			} else {
				require.JSONEq(t, tc.wantBody, body)
			}

			require.Equal(t, &moov.IssuedCardAuthorization{
				AuthorizationID:  "auth-1",
				IssuedCardID:     "card-1",
				FundingWalletID:  "wallet-1",
				Network:          moov.IssuedCardTransactionNetwork_Shazam,
				AuthorizedAmount: "12.34",
				Status:           moov.IssuedCardAuthorizationStatus_Pending,
				MerchantData: moov.IssuedCardTransactionMerchant{
					NetworkID: "net-1",
					Name:      moov.PtrOf("Coffee Shop"),
					Country:   "US",
					Mcc:       "5814",
				},
				CreatedOn: time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC),
			}, authorization)
		})
	}
}

func TestSimulations_Started(t *testing.T) {
	for _, tc := range simulationCases {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusAccepted)
				_, _ = w.Write([]byte(`{"authorizationID":"auth-1"}`))
			}))
			t.Cleanup(srv.Close)

			authorization, started, err := tc.call(newCardIssuingTestClient(t, srv))
			require.NoError(t, err)
			require.Nil(t, authorization)
			require.Equal(t, &mv2610.AuthorizationSimulationAsyncResponse{AuthorizationID: "auth-1"}, started)
		})
	}
}

func TestSimulations_Error(t *testing.T) {
	for _, tc := range simulationCases {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnprocessableEntity)
				_, _ = w.Write([]byte(`{"amount":"must be greater than zero"}`))
			}))
			t.Cleanup(srv.Close)

			authorization, started, err := tc.call(newCardIssuingTestClient(t, srv))
			require.Nil(t, authorization)
			require.Nil(t, started)

			var httpErr moov.HttpCallResponse
			require.ErrorAs(t, err, &httpErr)
			require.Equal(t, moov.StatusFailedValidation, httpErr.Status())
		})
	}
}

func TestSimulations_NilClient(t *testing.T) {
	authorization, started, err := mv2610.NewCardIssuingClient(nil).SimulateReversal(context.Background(), "account-123", "auth-1")
	require.Nil(t, authorization)
	require.Nil(t, started)
	require.EqualError(t, err, "client is nil")
}

func TestSimulations_RequiresIDs(t *testing.T) {
	cases := []struct {
		name    string
		wantErr string
		call    simulationCall
	}{
		{
			name:    "authorization without accountID",
			wantErr: "accountID is required",
			call: func(c mv2610.CardIssuingClient) (*moov.IssuedCardAuthorization, *mv2610.AuthorizationSimulationAsyncResponse, error) {
				return c.SimulateAuthorization(context.Background(), "", mv2610.CreateAuthorizationSimulation{})
			},
		},
		{
			name:    "clearing without accountID",
			wantErr: "accountID and authorizationID are required",
			call: func(c mv2610.CardIssuingClient) (*moov.IssuedCardAuthorization, *mv2610.AuthorizationSimulationAsyncResponse, error) {
				return c.SimulateClearing(context.Background(), "", "auth-1", mv2610.CreateClearingSimulation{})
			},
		},
		{
			name:    "clearing without authorizationID",
			wantErr: "accountID and authorizationID are required",
			call: func(c mv2610.CardIssuingClient) (*moov.IssuedCardAuthorization, *mv2610.AuthorizationSimulationAsyncResponse, error) {
				return c.SimulateClearing(context.Background(), "account-123", "", mv2610.CreateClearingSimulation{})
			},
		},
		{
			name:    "reversal without accountID",
			wantErr: "accountID and authorizationID are required",
			call: func(c mv2610.CardIssuingClient) (*moov.IssuedCardAuthorization, *mv2610.AuthorizationSimulationAsyncResponse, error) {
				return c.SimulateReversal(context.Background(), "", "auth-1")
			},
		},
		{
			name:    "reversal without authorizationID",
			wantErr: "accountID and authorizationID are required",
			call: func(c mv2610.CardIssuingClient) (*moov.IssuedCardAuthorization, *mv2610.AuthorizationSimulationAsyncResponse, error) {
				return c.SimulateReversal(context.Background(), "account-123", "")
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			called := false

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
			}))
			t.Cleanup(srv.Close)

			authorization, started, err := tc.call(newCardIssuingTestClient(t, srv))
			require.Nil(t, authorization)
			require.Nil(t, started)
			require.EqualError(t, err, tc.wantErr)
			require.False(t, called, "no request should reach the server")
		})
	}
}
