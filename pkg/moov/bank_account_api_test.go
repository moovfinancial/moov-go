package moov_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

func newBankAccountTestClient(t *testing.T, srv *httptest.Server) *moov.Client {
	t.Helper()

	client, err := moov.NewClient(
		moov.WithCredentials(moov.Credentials{PublicKey: "pk", SecretKey: "sk"}),
		moov.WithMoovURLScheme("http"),
	)
	require.NoError(t, err)
	client.Credentials.Host = strings.TrimPrefix(srv.URL, "http://")

	return client
}

func TestCreateBankAccount_RequestRiskVerification(t *testing.T) {
	account := moov.BankAccountRequest{
		HolderName:    "Ada Lovelace",
		HolderType:    moov.HolderType_Individual,
		AccountType:   moov.BankAccountType_Checking,
		AccountNumber: "123456789",
		RoutingNumber: "273976369",
	}

	tests := []struct {
		name    string
		opts    []moov.CreateBankAccountType
		want    map[string]any
		waitFor string
	}{
		{
			name: "sets the flag on an existing bank account body",
			opts: []moov.CreateBankAccountType{
				moov.WithBankAccount(account),
				moov.WithBankAccountRequestRiskVerification(),
			},
			want: map[string]any{
				"account": map[string]any{
					"holderName":      "Ada Lovelace",
					"holderType":      "individual",
					"bankAccountType": "checking",
					"accountNumber":   "123456789",
					"routingNumber":   "273976369",
				},
				"requestRiskVerification": true,
			},
		},
		{
			name: "keeps the flag when the bank account option is applied later",
			opts: []moov.CreateBankAccountType{
				moov.WithBankAccountRequestRiskVerification(),
				moov.WithBankAccount(account),
			},
			want: map[string]any{
				"account": map[string]any{
					"holderName":      "Ada Lovelace",
					"holderType":      "individual",
					"bankAccountType": "checking",
					"accountNumber":   "123456789",
					"routingNumber":   "273976369",
				},
				"requestRiskVerification": true,
			},
		},
		{
			name: "skips a nil option",
			opts: []moov.CreateBankAccountType{
				nil,
				moov.WithBankAccount(account),
				nil,
			},
			want: map[string]any{
				"account": map[string]any{
					"holderName":      "Ada Lovelace",
					"holderType":      "individual",
					"bankAccountType": "checking",
					"accountNumber":   "123456789",
					"routingNumber":   "273976369",
				},
			},
		},
		{
			name: "omits the flag when the option is not passed",
			opts: []moov.CreateBankAccountType{
				moov.WithBankAccount(account),
			},
			want: map[string]any{
				"account": map[string]any{
					"holderName":      "Ada Lovelace",
					"holderType":      "individual",
					"bankAccountType": "checking",
					"accountNumber":   "123456789",
					"routingNumber":   "273976369",
				},
			},
		},
		{
			name: "sets the flag on an existing Plaid body",
			opts: []moov.CreateBankAccountType{
				moov.WithPlaid(moov.PlaidRequest{Token: "plaid-token"}),
				moov.WithBankAccountRequestRiskVerification(),
			},
			want: map[string]any{
				"plaid":                   map[string]any{"token": "plaid-token"},
				"requestRiskVerification": true,
			},
		},
		{
			name: "keeps WaitForPaymentMethod as a header",
			opts: []moov.CreateBankAccountType{
				moov.WithBankAccount(account),
				moov.WithBankAccountRequestRiskVerification(),
				moov.WaitForPaymentMethod(),
			},
			want: map[string]any{
				"account": map[string]any{
					"holderName":      "Ada Lovelace",
					"holderType":      "individual",
					"bankAccountType": "checking",
					"accountNumber":   "123456789",
					"routingNumber":   "273976369",
				},
				"requestRiskVerification": true,
			},
			waitFor: "payment-method",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				method    string
				path      string
				waitFor   string
				body      map[string]any
				decodeErr error
			)

			// The handler runs on its own goroutine, so it records what it saw rather than
			// asserting: require.* calls t.FailNow, which is only valid on the test goroutine.
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				method = r.Method
				path = r.URL.Path
				waitFor = r.Header.Get("X-Wait-For")
				decodeErr = json.NewDecoder(r.Body).Decode(&body)

				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{
					"bankAccountID": "bank-123",
					"riskVerificationOutcome": "success"
				}`))
			}))
			t.Cleanup(srv.Close)

			created, err := newBankAccountTestClient(t, srv).CreateBankAccount(BgCtx(), "account-123", tt.opts...)
			require.NoError(t, err)
			require.NoError(t, decodeErr)

			require.Equal(t, http.MethodPost, method)
			require.Equal(t, "/accounts/account-123/bank-accounts", path)
			require.Equal(t, tt.waitFor, waitFor)
			require.Equal(t, tt.want, body)

			require.NotNil(t, created)
			require.Equal(t, "bank-123", created.BankAccountID)
			if _, ok := tt.want["requestRiskVerification"]; ok {
				require.Equal(t, moov.RiskVerificationOutcomeSuccess, created.RiskVerificationOutcome)
			}
		})
	}
}

func TestBankAccount_RiskVerificationOutcomeUnmarshal(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want moov.BankAccountRiskVerificationOutcome
	}{
		{name: "notAttempted", raw: "notAttempted", want: moov.RiskVerificationOutcomeNotAttempted},
		{name: "success", raw: "success", want: moov.RiskVerificationOutcomeSuccess},
		{name: "inconclusive", raw: "inconclusive", want: moov.RiskVerificationOutcomeInconclusive},
		{name: "decline", raw: "decline", want: moov.RiskVerificationOutcomeDecline},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var account moov.BankAccount
			err := json.Unmarshal([]byte(`{"riskVerificationOutcome":"`+tt.raw+`"}`), &account)
			require.NoError(t, err)
			require.Equal(t, tt.want, account.RiskVerificationOutcome)
		})
	}
}
