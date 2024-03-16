package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

// Sets up an account with linked bank account to be debited via ACH
func TestACHTransferSetup(t *testing.T) {
	// Step 1: create Moov client and set some variables

	// The following code shows how you can configure the moov client with
	// your credentials, if you don't want to use environment variables.
	// However, it is recommended to load the credentials from the
	// configuration file.

	mc, err := moov.NewClient(moov.WithCredentials(moov.Credentials{
		PublicKey: os.Getenv("MOOV_PUBLIC_KEY"),
		SecretKey: os.Getenv("MOOV_SECRET_KEY"),
		Host:      os.Getenv("MOOV_HOST"), // api.moov.io
	}))
	require.NoError(t, err)

	// The account we'll send funds to
	destinationAccountID := "xxxxx"

	// Create a new context or use an existing one
	ctx := context.Background()

	// Ping the server to check credentials
	err = mc.Ping(ctx)
	require.NoError(t, err)

	// Step 2: create account for the user

	// Add new account
	account, _, err := mc.CreateAccount(ctx, moov.Account{
		AccountType: moov.INDIVIDUAL,
		Profile: moov.Profile{
			Individual: moov.Individual{
				Name: moov.Name{
					FirstName: faker.FirstName(),
					LastName:  faker.LastName(),
				},
				Email: faker.Email(),
			},
		},
	})
	require.NoError(t, err)

	// Step 3: add (link) user's bank account

	// You can manually supply bank account information or pass tokens from
	// IAV providers Plaid or MX
	_, err = mc.CreateBankAccount(ctx, account.AccountID, moov.BankAccount{
		HolderName:      "Wade Arnold",
		HolderType:      "individual",
		BankAccountType: "checking",
		RoutingNumber:   "273976369", // this is a real routing number
		AccountNumber:   "123456789", // fake it great!
	})
	require.NoError(t, err)

	// Step 4: find (pull) payment method for the linked bank account

	// When we have only one bank account linked, we can avoid checking that the
	// payment method is for user's bank account and just use the first one.
	paymentMethods, err := mc.ListPaymentMethods(ctx, account.AccountID, moov.WithPaymentMethodType("ach-debit-fund"))
	require.NoError(t, err)

	// We expect to have only one `ach-debit-fund` payment method as we added
	// only one bank account
	require.Len(t, paymentMethods, 1)

	pullPaymentMethod := paymentMethods[0]

	// Step 3: configure destination payment method

	// We can pull money from the bank account and send to the
	// destination Moov wallet ("moov-wallet" payment method).
	paymentMethods, err = mc.ListPaymentMethods(ctx, destinationAccountID, moov.WithPaymentMethodType("moov-wallet"))
	require.NoError(t, err)

	require.Len(t, paymentMethods, 1)

	// This is the destination payment method (Moov wallet)
	destinationPaymentMethod := paymentMethods[0]

	// Step 4: create transfer
	completedTransfer, _, err := mc.CreateTransfer(
		ctx,
		moov.CreateTransfer{
			Source: moov.Source{
				PaymentMethodID: pullPaymentMethod.PaymentMethodID,
			},
			Destination: moov.Destination{
				PaymentMethodID: destinationPaymentMethod.PaymentMethodID,
			},
			Amount: moov.Amount{
				Currency: "USD",
				Value:    5000, // $50.00
			},
		},
		// not required since ACH is processed in batches,
		// but useful in getting the full transfer model
		moov.WithTransferWaitForRailResponse(),
	)
	require.NoError(t, err)

	fmt.Printf("Transfer: %+v\n", completedTransfer.TransferID)

}