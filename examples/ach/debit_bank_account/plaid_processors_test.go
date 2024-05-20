package debit_bank_account

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
)

// Sets up an account with linked bank account to be debited via ACH
func TestPlaidProcessorExample(t *testing.T) {
	// Step 1: create Moov client and set some variables

	// The following code shows how you can configure the moov client with
	// your credentials, if you don't want to use environment variables.
	// However, it is recommended to load the credentials from the
	// configuration file.

	mc, err := moov.NewClient() // reads credentials from Environmental variables
	require.NoError(t, err)

	// The account we'll send funds to
	destinationAccountID := "ebbf46c6-122a-4367-bc45-7dd555e1d3b9" // example

	// Create a new context or use an existing one
	ctx := context.Background()

	// Ping the server to check credentials
	err = mc.Ping(ctx)
	require.NoError(t, err)

	// Step 2: create account for the user

	// Add new account
	account, _, err := mc.CreateAccount(ctx, moov.CreateAccount{
		Type: moov.AccountType_Individual,
		Profile: moov.CreateProfile{
			Individual: &moov.CreateIndividualProfile{
				Name: moov.Name{
					FirstName: faker.FirstName(),
					LastName:  faker.LastName(),
				},
				Email: faker.Email(),
			},
		},
	})
	require.NoError(t, err)

	t.Logf("Created Account: %v", account.AccountID)

	// Step 3: add (link) user's bank account with Plaid Processor
	//
	// See https://plaid.com/docs/auth/partnerships/moov/ for Plaid's documentation
	// See https://docs.moov.io/guides/sources/bank-accounts/plaid/ for Moov's documentation
	plaidProcessorToken := strings.TrimSpace(os.Getenv("PLAID_PROCESSOR_TOKEN"))
	if plaidProcessorToken == "" {
		plaidProcessorToken = createLinkToken(t)
	}
	require.NotEmpty(t, plaidProcessorToken)

	plaid := moov.PlaidRequest{
		Token: plaidProcessorToken,
	}
	bankAccount, err := mc.CreateBankAccount(ctx, account.AccountID, moov.WithPlaid(plaid), moov.WaitForPaymentMethod())
	require.NoError(t, err)

	t.Logf("Created Bank Account: %v", bankAccount.BankAccountID)

	// Step 4: find (pull) payment method for the linked bank account

	// When we have only one bank account linked, we can avoid checking that the
	// payment method is for user's bank account and just use the first one.
	paymentMethods, err := mc.ListPaymentMethods(ctx, account.AccountID, moov.WithPaymentMethodType("ach-debit-collect"))
	require.NoError(t, err)

	pullPaymentMethod := paymentMethods[0]

	// Step 5: configure destination payment method

	// We can pull money from the bank account and send to the
	// destination Moov wallet ("moov-wallet" payment method).
	paymentMethods, err = mc.ListPaymentMethods(ctx, destinationAccountID, moov.WithPaymentMethodType("moov-wallet"))
	require.NoError(t, err)
	require.Len(t, paymentMethods, 1)

	// This is the destination payment method (Moov wallet)
	destinationPaymentMethod := paymentMethods[0]

	// Step 6: create transfer
	completedTransfer, _, err := mc.CreateTransfer(
		ctx,
		moov.CreateTransfer{
			Source: moov.CreateTransfer_Source{
				PaymentMethodID: pullPaymentMethod.PaymentMethodID,
			},
			Destination: moov.CreateTransfer_Destination{
				PaymentMethodID: destinationPaymentMethod.PaymentMethodID,
			},
			Amount: moov.Amount{
				Currency: "USD",
				Value:    2717, // $27.17
			},
		}).
		// not required since ACH is processed in batches,
		// but useful in getting the full transfer model
		WaitForRailResponse()

	require.NoError(t, err)

	t.Logf("Transfer %s created", completedTransfer.TransferID)
	t.Logf("Amount: %#v", completedTransfer.Amount)
	t.Logf("Status: %v", completedTransfer.Status)
	t.Logf("CreatedOn: %v", completedTransfer.CreatedOn)
}

type plaidLinkTokenRequest struct {
	ClientID   string `json:"client_id"`
	Secret     string `json:"secret"`
	ClientName string `json:"client_name"`
	User       struct {
		ClientUserID string `json:"client_user_id"`
	} `json:"user"`
	Products     []string `json:"products"`
	CountryCodes []string `json:"country_codes"`
	Language     string   `json:"language"`
	Webhook      string   `json:"webhook"`
	RedirectURI  string   `json:"redirect_uri"`
}

type plaidLinkTokenResponse struct {
	LinkToken string `json:"link_token"`
}

func createLinkToken(t testing.TB) string {
	t.Helper()

	linkToken := plaidLinkTokenRequest{
		ClientID:     os.Getenv("PLAID_CLIENT_ID"),
		Secret:       os.Getenv("PLAID_SECRET"),
		ClientName:   "Moov Debit BankAccount Example",
		Products:     []string{"auth"},
		CountryCodes: []string{"US"},
		Language:     "en",
		Webhook:      "https://plaid-processor-test.moov.io",
		RedirectURI:  "https://local.moov.io",
	}
	linkToken.User.ClientUserID = "4c7ab47a-487a-4c91-92a7-7fda891b2bb0"

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(linkToken)

	req, err := http.NewRequest("POST", "https://sandbox.plaid.com/link/token/create", &buf)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	t.Cleanup(func() { resp.Body.Close() })

	bs, _ := io.ReadAll(resp.Body)

	var response plaidLinkTokenResponse
	json.NewDecoder(bytes.NewReader(bs)).Decode(&response)

	if response.LinkToken == "" {
		t.Log("Plaid Link Token Response")
		t.Logf("%v", string(bs))
	}
	return response.LinkToken
}
