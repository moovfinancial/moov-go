package debit_card_pull

import (
	"context"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

// Pull transfer is a transfer from the user's card to the Moov wallet.
func TestDebitPullWithRefund(t *testing.T) {
	// Step 1: create Moov client and set some variables

	// The following code shows how you can configure the moov client with
	// your credentials, if you don't want to use environment variables.
	// However, it is recommended to load the credentials from the
	// configuration file.

	mc, err := moov.NewClient(moov.WithCredentials(moov.CredentialsFromEnv()))
	require.NoError(t, err)

	// Create a new context or use an existing one
	ctx := context.Background()

	// Ping the server to check credentials
	err = mc.Ping(ctx)
	require.NoError(t, err)

	// Step 2: create account for the user
	// For now just using a known existing account
	accountID := "ebbf46c6-122a-4367-bc45-7dd555e1d3b9"

	// Step 3: add (link) user's card

	// You can make the direct call only if you are PCI compliant,
	// otherwise you need to use the Moov.js library
	card, err := mc.CreateCard(ctx, accountID, moov.CreateCard{
		CardNumber: "4111100010002000", // Moov test card for sandbox
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

	// Step 4: find (pull) payment method for the linked card

	// When we have only one card linked, we can avoid checking that the
	// payment method is for user's card and just use the first one.
	paymentMethods, err := mc.ListPaymentMethods(ctx, accountID, moov.WithPaymentMethodType("pull-from-card"))
	require.NoError(t, err)

	// We expect to have only one `pull-to-card` payment method as we added
	// only one card
	require.Len(t, paymentMethods, 1)

	pullPaymentMethod := paymentMethods[0]

	// Step 5: configure destination payment method

	// We can pull money from the card ("pull-from-card" payment method),
	// and to the Moov wallet ("moov-wallet" payment method).
	paymentMethods, err = mc.ListPaymentMethods(ctx, accountID, moov.WithPaymentMethodType("moov-wallet"))
	require.NoError(t, err)

	require.Len(t, paymentMethods, 1)

	// This is the source payment method (Moov wallet)
	destinationPaymentMethod := paymentMethods[0]

	// Step 6: create transfer
	completedTransfer, _, err := mc.CreateTransfer(
		ctx,
		moov.CreateTransfer{
			Source: moov.CreateTransfer_Source{
				PaymentMethodID: pullPaymentMethod.PaymentMethodID,
				CardDetails: &moov.CreateTransfer_CardDetailsSource{
					DynamicDescriptor: "Test pull transfer",
				},
			},
			Destination: moov.CreateTransfer_Destination{
				PaymentMethodID: destinationPaymentMethod.PaymentMethodID,
			},
			Amount: moov.Amount{
				Currency: "USD",
				Value:    97, // $0.97
			},
			FacilitatorFee: moov.CreateTransfer_FacilitatorFee{
				Total: moov.PtrOf(int64(2)), // $0.02
			},
			Description: "Pull from card",
		}).WaitForRailResponse()
	require.NoError(t, err)

	t.Logf("Transfer %s created", completedTransfer.TransferID)
	t.Logf("Amount: %#v", completedTransfer.Amount)
	t.Logf("Status: %v", completedTransfer.Status)
	t.Logf("CreatedOn: %v", completedTransfer.CreatedOn)
	t.Log("")

	// Step 7: refund transfer
	refund, _, err := mc.RefundTransfer(
		ctx,
		completedTransfer.TransferID,
		moov.CreateRefund{
			Amount: 97,
		},
		moov.WithRefundWaitForRailResponse(),
	)
	require.NoError(t, err)

	t.Logf("Transfer %s created", refund.RefundID)
	t.Logf("Refund %s", refund.RefundID)
	t.Logf("Amount: %#v", refund.Amount)
	t.Logf("Status: %v", refund.Status)
	t.Logf("CreatedOn: %v", refund.CreatedOn)

	// Step 8: disable the card
	err = mc.DisableCard(ctx, accountID, card.CardID)
	require.NoError(t, err)
}
