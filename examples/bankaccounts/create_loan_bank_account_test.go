package bankaccounts

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/moovfinancial/moov-go/pkg/moov"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
)

// Sets up an account with linked loan bank account
func TestBankAccount_CreateLoan(t *testing.T) {
	// Step 1: create Moov client and set some variables

	// The following code shows how you can configure the moov client with
	// your credentials, if you don't want to use environment variables.
	// However, it is recommended to load the credentials from the
	// configuration file.

	mc, err := moov.NewClient(
		moov.WithCredentials(moov.CredentialsFromEnv()), // optional, default is to read from environment
	)
	require.NoError(t, err)

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

	// Step 3: add (link) user's bank account

	// You can manually supply bank account information or pass tokens from
	// IAV providers Plaid or MX
	bankAccountPayload := moov.BankAccountRequest{
		HolderName:    fmt.Sprintf("%s %s", account.Profile.Individual.Name.FirstName, account.Profile.Individual.Name.LastName),
		HolderType:    moov.HolderType_Individual,
		AccountType:   moov.BankAccountType_Loan,
		RoutingNumber: "271071321",                          // this is a real routing number
		AccountNumber: fmt.Sprintf("%d", time.Now().Unix()), // fake account number
	}
	bankAccount, err := mc.CreateBankAccount(ctx, account.AccountID, moov.WithBankAccount(bankAccountPayload), moov.WaitForPaymentMethod())
	require.NoError(t, err)

	t.Logf("Created Bank Account: %v", bankAccount.BankAccountID)
	t.Logf("Bank Account Type: %v", bankAccount.BankAccountType)
	t.Logf("Bank Account Status: %v", bankAccount.Status)

	// Loans only have credit payment methods available.
	// Get payment methods for the created bank account.
	paymentMethods, err := mc.ListPaymentMethods(ctx, account.AccountID, moov.WithPaymentMethodSourceID(bankAccount.BankAccountID))
	require.NoError(t, err)

	t.Logf("Found %d payment methods", len(paymentMethods))
	for _, pm := range paymentMethods {
		t.Logf("  ID: %v  Type: %v", pm.PaymentMethodID, pm.PaymentMethodType)
	}
}
