package debit_bank_account

import (
	"context"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

// Sets up an account with linked bank account to be debited via ACH
func TestMicroDepositExample(t *testing.T) {
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

	// Step 3: add (link) user's bank account

	// You can manually supply bank account information or pass tokens from
	// IAV providers Plaid or MX
	bankAccountPayload := moov.BankAccountRequest{
		HolderName:    "Jules Jackson",
		HolderType:    moov.HolderType_Individual,
		AccountType:   moov.BankAccountType_Checking,
		RoutingNumber: "273976369", // this is a real routing number
		AccountNumber: "123456789", // fake it great!
	}
	bankAccount, err := mc.CreateBankAccount(ctx, account.AccountID, moov.WithBankAccount(bankAccountPayload), moov.WaitForPaymentMethod())
	require.NoError(t, err)

	t.Logf("Created Bank Account: %v", bankAccount.BankAccountID)

	// Initiate micro-deposits
	baErr := mc.MicroDepositInitiate(ctx, account.AccountID, bankAccount.BankAccountID)
	require.NoError(t, baErr)

	// Verify micro-deposits (later)
	amounts := []int{0, 0} // Sandbox amounts are always [0, 0]
	verifyErr := mc.MicroDepositConfirm(ctx, account.AccountID, bankAccount.BankAccountID, amounts)
	require.NoError(t, verifyErr)

	// Step 4: find (pull) payment method for the linked bank account

	// When we have only one bank account linked, we can avoid checking that the
	// payment method is for user's bank account and just use the first one.
	paymentMethods, err := mc.ListPaymentMethods(ctx, account.AccountID, moov.WithPaymentMethodType("ach-debit-collect"))
	require.NoError(t, err)

	// We expect to have only one `ach-debit-collect` payment method as we added
	// only one bank account
	require.Len(t, paymentMethods, 1)

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
				Value:    4328, // $43.28
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
