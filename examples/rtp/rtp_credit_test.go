package rtp

import (
	"context"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
)

func TestRTPCreditExample(t *testing.T) {
	// Step 1: create Moov client

	// The following code shows how you can configure the moov client with
	// your credentials, if you don't want to use environment variables.
	// However, it is recommended to load the credentials from the
	// configuration file.

	mc, err := moov.NewClient() // reads credentials from Environmental variables
	require.NoError(t, err)

	// Create a new context or use an existing one
	ctx := context.Background()

	// Ping the server to check credentials
	err = mc.Ping(ctx)
	require.NoError(t, err)

	// Account IDs used
	sourceAccountID := "ebbf46c6-122a-4367-bc45-7dd555e1d3b9"

	// Step 2: Get the source wallet for our RTP transfer
	sourcePaymentMethods, err := mc.ListPaymentMethods(ctx, sourceAccountID, moov.WithPaymentMethodType("moov-wallet"))
	require.NoError(t, err)
	require.Len(t, sourcePaymentMethods, 1)

	sourcePaymentMethod := sourcePaymentMethods[0]

	// Step 3: Create account for the destination
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

	// Step 4: Add (link) a bank account enabled for RTP
	bankAccountPayload := moov.BankAccountRequest{
		HolderName:    "Jules Jackson",
		HolderType:    moov.HolderType_Individual,
		AccountType:   moov.BankAccountType_Checking,
		RoutingNumber: "011401533", // this is a real routing number
		AccountNumber: "123456789", // fake it great!
	}
	bankAccount, err := mc.CreateBankAccount(ctx, account.AccountID, moov.WithBankAccount(bankAccountPayload), moov.WaitForPaymentMethod())
	require.NoError(t, err)

	t.Logf("BankAccountID: %v", bankAccount.BankAccountID)

	// Step 5: Find the rtp-credit payment method
	destinationPaymentMethods, err := mc.ListPaymentMethods(ctx, account.AccountID, moov.WithPaymentMethodType("rtp-credit"))
	require.NoError(t, err)
	require.Greater(t, len(destinationPaymentMethods), 0)

	destinationPaymentMethod := destinationPaymentMethods[0]

	// Step 6: create transfer
	_, completedAsyncTransfer, err := mc.createTransfer(
		ctx,
		moov.CreateTransfer{
			Source: moov.CreateTransfer_Source{
				PaymentMethodID: sourcePaymentMethod.PaymentMethodID,
			},
			Destination: moov.CreateTransfer_Destination{
				PaymentMethodID: destinationPaymentMethod.PaymentMethodID,
			},
			Amount: moov.Amount{
				Currency: "USD",
				Value:    132, // $1.32
			},
		},
	)
	require.NoError(t, err)

	t.Logf("Transfer %s created", completedAsyncTransfer.TransferID)
	t.Logf("CreatedOn: %v", completedAsyncTransfer.CreatedOn)
}
