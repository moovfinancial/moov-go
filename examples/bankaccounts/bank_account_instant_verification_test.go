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

// Sets up an account with linked bank account and initiates instant verification
func TestBankAccount_InstantVerificationExample(t *testing.T) {
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
		AccountType:   moov.BankAccountType_Checking,
		RoutingNumber: "271071321",                          // this is a real routing number
		AccountNumber: fmt.Sprintf("%d", time.Now().Unix()), // fake account number
	}
	bankAccount, err := mc.CreateBankAccount(ctx, account.AccountID, moov.WithBankAccount(bankAccountPayload), moov.WaitForPaymentMethod())
	require.NoError(t, err)

	t.Logf("Created Bank Account: %v", bankAccount.BankAccountID)

	// Initiate instant verification
	baErr := mc.InstantVerificationInitiate(ctx, account.AccountID, bankAccount.BankAccountID)
	require.NoError(t, baErr)

	// Fetch the Bank Account Verification's status
	require.Eventually(t, func() bool {
		bav, err := mc.GetInstantBankAccountVerfication(ctx, account.AccountID, bankAccount.BankAccountID)
		require.NoError(t, err)

		require.NotNil(t, bav)
		require.Nil(t, bav.ExceptionDetails)
		require.Equal(t, moov.BankAccountVerificationMethodInstant, bav.VerificationMethod)

		if bav.Status == moov.BankAccountVerificationStatusSentCredit {
			return true
		}
		return false
	}, 20*time.Second, time.Second)

	// Complete instant verification
	code := "MV0001" // Sandbox code is always MV0001
	verifyErr := mc.InstantVerificationComplete(ctx, account.AccountID, bankAccount.BankAccountID, code)
	require.NoError(t, verifyErr)

	// Step 4: find (push) payment methods for the linked bank account

	sourceAccountID := "ebbf46c6-122a-4367-bc45-7dd555e1d3b9"

	// When we have only one bank account linked, we can avoid checking that the
	// payment method is for user's bank account and just use the first one.
	var paymentMethods []moov.PaymentMethod
	require.Eventually(t, func() bool {
		paymentMethods, err = mc.ListPaymentMethods(ctx, sourceAccountID, moov.WithPaymentMethodType("moov-wallet"))
		require.NoError(t, err)

		return len(paymentMethods) > 0
	}, 10*time.Second, time.Second)

	// We expect to have only one `moov-wallet` payment method on the connected account
	require.Len(t, paymentMethods, 1)

	sourcePaymentMethod := paymentMethods[0]

	// Step 5: configure destination payment method

	// We can pull money from the bank account and send to the
	// destination bank account over RTP ("rtp-credit" payment method).
	paymentMethods, err = mc.ListPaymentMethods(ctx, account.AccountID, moov.WithPaymentMethodType("rtp-credit"))
	require.NoError(t, err)

	require.Len(t, paymentMethods, 1)

	// This is the destination payment method (Bank Account)
	destinationPaymentMethod := paymentMethods[0]

	// Step 6: create transfer

	// The account facilitating the transfer
	partnerAccountID := "5352b013-ae58-4a63-8a3f-97f316a917cf" // example

	completedTransfer, _, err := mc.CreateTransfer(
		ctx,
		partnerAccountID,
		moov.CreateTransfer{
			Source: moov.CreateTransfer_Source{
				PaymentMethodID: sourcePaymentMethod.PaymentMethodID,
			},
			Destination: moov.CreateTransfer_Destination{
				PaymentMethodID: destinationPaymentMethod.PaymentMethodID,
			},
			Amount: moov.Amount{
				Currency: "USD",
				Value:    1245, // $12.45
			},
		}).
		// not required, but useful in getting the full transfer model
		WaitForRailResponse()

	require.NoError(t, err)

	t.Logf("Transfer %s created", completedTransfer.TransferID)
	t.Logf("Amount: %#v", completedTransfer.Amount)
	t.Logf("Status: %v", completedTransfer.Status)
	t.Logf("CreatedOn: %v", completedTransfer.CreatedOn)
}
