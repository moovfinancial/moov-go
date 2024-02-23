package examples

import (
	"context"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"

	"github.com/stretchr/testify/require"
)

func TestCheckout(t *testing.T) {
	// Create a new moov client
	mc, err := moov.NewClient(moov.WithCredentials(moov.CredentialsFromEnv()))
	require.NoError(t, err)

	// Create a new context or use an existing one
	ctx := context.Background()

	// Ping the server to check credentials
	err = mc.Ping(ctx)
	require.NoError(t, err)

	// Create access token to instantiate a moov.js drop for creating an account
	token, err := mc.AccountCreationToken(ctx)
	require.NoError(t, err)

	// load this into a moov.js drop to create an account
	t.Logf("AccountCreationToken=%s", token.AccessToken)

	// Create a new moov account server side
	account := moov.Account{
		AccountType: moov.INDIVIDUAL,
		Profile: moov.Profile{
			Individual: moov.Individual{
				Name: moov.Name{
					FirstName: "Wade",
					LastName:  "Arnold",
				},
				Email: "wade@arnold.com",
				Phone: moov.Phone{
					Number:      "555-555-5555",
					CountryCode: "1",
				},
			},
		},
	}
	completedAccount, _, err := mc.CreateAccount(ctx, account)
	require.NoError(t, err)

	t.Logf("AccountID=%s", completedAccount.AccountID)

	// You can search for all existing accounts with a variety of query parameter's. We will use this account as the destination to receive money.
	accounts, err := mc.ListAccounts(ctx, moov.WithAccountName("Lincoln National Corporation"), moov.WithAccountCount(10))
	require.NoError(t, err)

	var lincolnAccount moov.Account
	for _, account := range accounts {
		t.Logf("AccountID=%s \t Display Name=%s", account.AccountID, account.DisplayName)
		if account.DisplayName == "Lincoln National Corporation" {
			lincolnAccount = account
		}
	}
	// Get the paymentMethodID for the Lincoln National wallet
	lincolnPaymentMethods, err := mc.ListPaymentMethods(ctx, lincolnAccount.AccountID, moov.WithPaymentMethodType(("moov-wallet")))
	require.NoError(t, err)

	// As of the writing of this example, every account has one and only one wallet
	if len(lincolnPaymentMethods) == 0 {
		t.Fatal("No wallet found for Lincoln National Corporation")
	}
	lincolnWallet := lincolnPaymentMethods[0]

	t.Logf("Lincoln AccountID=%s \t PaymentMethodID=%s", lincolnAccount.AccountID, lincolnWallet.PaymentMethodID)

	// To link a card, apple pay, or bank account to an account you can use a Moov drop to stay out of PCI scope.
	// Create a scoped access token for the account we just created. This token can link every payment type from Moov.js
	token, err = mc.AccessToken(ctx, moov.Scopes.CardsWrite(completedAccount.AccountID),
		moov.Scopes.BankAccountsWrite(completedAccount.AccountID),
		moov.Scopes.ApplePayWrite(completedAccount.AccountID))
	require.NoError(t, err)

	// Use this token in the card or bank account linking Drop.
	// Token's allow the browser to link cards and bank accounts to the account directly with Moov.
	//fmt.Println(completedAccount.AccountID) // cardInput.AccountID
	t.Logf("access token: %v", token.AccessToken) // cardInput.AccessToken

	// The drop's onSuccess event for the card linking drop returns a Card object the cardID and paymentMethodID
	// The paymentMethodID is used to transfer money from the card to the moov wallet and needs to be sent to the server

	// If you are PCI compliant you can send a card directly to Moov
	card := moov.CreateCard{
		CardNumber: "4111111111111111", // Moov test card for sandbox
		CardCvv:    "123",
		Expiration: moov.Expiration{
			Month: "12",
			Year:  "28",
		},
		HolderName: "Wade Arnold",
		BillingAddress: moov.Address{
			PostalCode: "80401",
		},
	}
	linkCard, err := mc.CreateCard(ctx, completedAccount.AccountID, card)
	require.NoError(t, err)

	// The same card can be used for card-payment, push-to-card or pull-from-card and each of their own paymentMethodID
	var cardPaymentMethodID string
	for _, pm := range linkCard.PaymentMethods {
		if pm.PaymentMethodType == "card-payment" {
			cardPaymentMethodID = pm.PaymentMethodID
		}
	}

	t.Logf("PaymentMethodID=%s \t lastFourCardNumber=%s \t brand=%s", cardPaymentMethodID, linkCard.LastFourCardNumber, linkCard.Brand)

	// An existing linked card can be used to transfer money to the moov wallet by finding it's paymentMethodID
	paymentMethods, err := mc.ListPaymentMethods(ctx, completedAccount.AccountID, moov.WithPaymentMethodType("card-payment"))
	require.NoError(t, err)

	for _, paymentMethod := range paymentMethods {
		t.Logf("PaymentMethodID=%s \t lastFourCardNumber=%s \t PaymentMethodType=%s", paymentMethod.PaymentMethodID, paymentMethod.Card.LastFourCardNumber, paymentMethod.PaymentMethodType)
	}

	// We showed three ways to get a PaymentMethodID to transfer money from a card to the moov wallet.
	// The following is a card to wallet transfer but any payment method can be used to transfer money to the moov wallet

	// Amount to be transferred from the card to the wallet.
	amount := moov.Amount{
		Currency: "USD",
		Value:    9900, // $99.00
	}

	source := moov.Source{
		PaymentMethodID: cardPaymentMethodID,
		CardDetails: moov.CardDetails{
			DynamicDescriptor: "WhlBdy *Yoga 11-12",
			TransactionSource: "first-recurring",
		},
	}
	destination := moov.Destination{
		PaymentMethodID: lincolnPaymentMethods[0].PaymentMethodID,
	}
	facilitatorFee := moov.FacilitatorFee{
		Total: 8, // $0.08
	}
	description := "Pay Instructor for May 15 Class"

	// Create a transfer from the card to the wallet
	completedTransfer, _, err := mc.CreateTransfer(context.Background(), moov.CreateTransfer{
		Source:         source,
		Destination:    destination,
		Amount:         amount,
		FacilitatorFee: facilitatorFee,
		Description:    description,
	}, moov.WithTransferWaitForRailResponse())
	require.NoError(t, err)

	t.Logf("TransferID=%s \t Status=%s \n", completedTransfer.TransferID, completedTransfer.Status)

	// We can check the status of the transfer. More important for batched ACH transactions
	transferStatus, err := mc.GetTransfer(ctx, completedTransfer.TransferID)
	require.NoError(t, err)

	t.Logf("TransferID=%s \t Status=%s \n", transferStatus.TransferID, transferStatus.Status)
}
