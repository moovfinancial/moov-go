package bankaccounts

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/moovfinancial/moov-go/pkg/moov"

	"github.com/stretchr/testify/require"
)

// Sets up an account with linked general-ledger bank account
func TestBankAccount_CreateGeneralLedger(t *testing.T) {
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
	// Only Businesses can create General Ledger bank accounts
	account, _, err := mc.CreateAccount(ctx, moov.CreateAccount{
		Type: moov.AccountType_Business,
		Profile: moov.CreateProfile{
			Business: &moov.CreateBusinessProfile{
				Address: &moov.Address{
					AddressLine1:    "123 Main Street",
					AddressLine2:    "Apt 302",
					City:            "Boulder",
					StateOrProvince: "CO",
					PostalCode:      "80301",
					Country:         "US",
				},
				Type:        "llc",
				Description: "Local fitness center paying out instructors",
				DBA:         "Whole Body Fitness",
				Email:       "amanda@classbooker.dev",
				IndustryCodes: &moov.IndustryCodes{
					Naics: "713940",
					Sic:   "7991",
					Mcc:   "7997",
				},
				Name: "Whole Body Fitness LLC",
				Phone: &moov.Phone{
					Number:      "8185551212",
					CountryCode: "1",
				},
				TaxID: &moov.TaxID{
					EIN: moov.EIN{
						Number: "123456789",
					},
				},
				Website: "www.wholebodyfitnessgym.com",
			},
		},
		RequestedCapabilities: []moov.CapabilityName{
			moov.CapabilityName_Transfers,
			moov.CapabilityName_SendFunds,
			moov.CapabilityName_CollectFunds,
			moov.CapabilityName_Wallet,
		},
	})
	require.NoError(t, err)

	t.Logf("Created Account: %v", account.AccountID)

	// Step 3: add (link) user's bank account

	// You can manually supply bank account information or pass tokens from
	// IAV providers Plaid or MX
	bankAccountPayload := moov.BankAccountRequest{
		HolderName:    account.Profile.Business.LegalBusinessName,
		HolderType:    moov.HolderType_Business,
		AccountType:   moov.BankAccountType_GeneralLedger,
		RoutingNumber: "271071321",                          // this is a real routing number
		AccountNumber: fmt.Sprintf("%d", time.Now().Unix()), // fake account number
	}
	bankAccount, err := mc.CreateBankAccount(ctx, account.AccountID, moov.WithBankAccount(bankAccountPayload), moov.WaitForPaymentMethod())
	require.NoError(t, err)

	t.Logf("Created Bank Account: %v", bankAccount.BankAccountID)
	t.Logf("Bank Account Type: %v", bankAccount.BankAccountType)

	// Initiate instant verification
	baErr := mc.InstantVerificationInitiate(ctx, account.AccountID, bankAccount.BankAccountID)
	require.NoError(t, baErr)
	time.Sleep(2 * time.Second)

	// Complete instant verification
	code := "MV0001" // Sandbox code is always MV0001
	verifyErr := mc.InstantVerificationComplete(ctx, account.AccountID, bankAccount.BankAccountID, code)
	require.NoError(t, verifyErr)

	// Get payment methods for the created bank account.
	paymentMethods, err := mc.ListPaymentMethods(ctx, account.AccountID, moov.WithPaymentMethodSourceID(bankAccount.BankAccountID))
	require.NoError(t, err)

	t.Logf("Found %d payment methods", len(paymentMethods))
	for _, pm := range paymentMethods {
		t.Logf("  ID: %v  Type: %v", pm.PaymentMethodID, pm.PaymentMethodType)
	}
}
