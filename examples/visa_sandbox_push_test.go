package examples

import (
	"context"
	"fmt"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
)

func TestVisaSandboxPush(t *testing.T) {
	// Step 1: create Moov client and set some variables

	// The following code shows how you can configure the moov client with
	// your credentials, if you don't want to use environment variables.
	// However, it is recommended to load the credentials from the
	// configuration file.
	mc, err := moov.NewClient(moov.WithCredentials(moov.CredentialsFromEnv()))
	require.NoError(t, err)

	// sourceAccountID is the account from which we will pull money. This
	// account should be properly configured (e.g. MCC should be set)
	sourceAccountID := "xxxxx"

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
				Phone: moov.Phone{
					Number:      faker.Phonenumber(),
					CountryCode: "1",
				},
			},
		},
	})
	require.NoError(t, err)

	// Step 3: add (link) user's card

	// You can make the direct call only if you are PCI compliant,
	// otherwise you need to use the Moov.js library
	_, err = mc.CreateCard(ctx, account.AccountID, moov.CreateCard{
		CardNumber: "4111111111111111", // Moov test card for sandbox
		CardCvv:    "123",
		Expiration: moov.Expiration{
			Month: "12",
			Year:  "28",
		},
		HolderName: "John Doe",
		BillingAddress: moov.Address{
			PostalCode: "80401",
		},
	})
	require.NoError(t, err)

	// Step 4: find (push) payment method for the linked card

	// When we have only one card linked, we can avoid checking that the
	// payment method is for user's card and just use the first one.
	paymentMethods, err := mc.ListPaymentMethods(ctx, account.AccountID, moov.WithPaymentMethodType("push-to-card"))
	require.NoError(t, err)

	// We expect to have only one `push-to-card` payment method as we added
	// only one card
	require.Len(t, paymentMethods, 1)

	pushPaymentMethod := paymentMethods[0]

	// Step 3: configure source payment method

	// We can pull money from the Moov wallet ("moov-wallet" payment
	// methoD) and push it to the card ("push-to-card" payment method).
	paymentMethods, err = mc.ListPaymentMethods(ctx, sourceAccountID, moov.WithPaymentMethodType("moov-wallet"))
	require.NoError(t, err)

	require.Len(t, paymentMethods, 1)

	// This is the source payment method (Moov wallet)
	sourcePaymentMethod := paymentMethods[0]

	// Step 4: create transfer
	completedTransfer, _, err := mc.CreateTransfer(
		ctx,
		moov.CreateTransfer{
			Source: moov.Source{
				PaymentMethodID: sourcePaymentMethod.PaymentMethodID,
			},
			Destination: moov.Destination{
				PaymentMethodID: pushPaymentMethod.PaymentMethodID,
				CardDetails: moov.CardDetails{
					DynamicDescriptor: "Test transfer",
				},
			},
			Amount: moov.Amount{
				Currency: "USD",
				Value:    99, // $0.99
			},
			FacilitatorFee: moov.FacilitatorFee{
				Total: 2, // $0.02
			},
			Description: "Push to card",
		},
		moov.WithTransferWaitForRailResponse(),
	)
	require.NoError(t, err)

	fmt.Printf("Transfer: %+v\n", completedTransfer.TransferID)
}
