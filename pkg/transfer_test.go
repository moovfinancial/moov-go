package moov_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"

	moov "github.com/moovfinancial/moov-go/pkg"
	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/require"
)

func TestSynchronousTransferMarshalResponse(t *testing.T) {

	input := []byte(`{
		"transferID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43",
		"createdOn": "2019-08-24T14:15:22Z",
		"completedOn": "2019-08-24T14:15:22Z",
		"status": "pending",
		"failureReason": "wallet-insufficient-funds",
		"amount": { "currency": "USD","value": 1204 },
		"description": "Pay Instructor for May 15 Class",
		"metadata": {"property1": "string","property2": "string"},
		"facilitatorFee": { "total": 0, "totalDecimal": "12.987654321", "markup": 0, "markupDecimal": "0.987654321"},
		"moovFee": 0,
		"moovFeeDecimal": "0.987654321",
		"moovFeeDetails": {"cardScheme": "string", "interchange": "string", "moovProcessing": "string" },
		"groupID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43",
		"refundedAmount": {"currency": "USD","value": 1204},
		"refunds": [
			{"refundID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43", "createdOn": "2019-08-24T14:15:22Z", "updatedOn": "2019-08-24T14:15:22Z", 
			"status": "failed",
			"failureCode": "call-issuer",
			"amount": {"currency": "USD", "value": 1204 },
			"cardDetails": { 
				"status": "initiated","failureCode": "call-issuer",
				"statusUpdates": {"initiated": "2019-08-24T14:15:22Z","confirmed": "2019-08-24T14:15:22Z","settled": "2019-08-24T14:15:22Z","failed": "2019-08-24T14:15:22Z","completed": "2019-08-24T14:15:22Z"}}}],
		"disputedAmount": {"currency": "USD","value": 1204},
		"disputes": [{
			"disputeID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43","createdOn": "2019-08-24T14:15:22Z",
			"amount": {"currency": "USD","value": 1204}}],
		"source": {
			"paymentMethodID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43",
			"paymentMethodType": "moov-wallet",
			"account": {
				"accountID": "3dfff852-927d-47e8-822c-2fffc57ff6b9",
				"email": "amanda@classbooker.dev",
				"displayName": "Whole Body Fitness"},
			"bankAccount": {
				"bankAccountID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43",
				"fingerprint": "9948962d92a1ce40c9f918cd9ece3a22bde62fb325a2f1fe2e833969de672ba3",
				"status": "new",
				"holderName": "Jules Jackson",
				"holderType": "individual",
				"bankName": "Chase Bank",
				"bankAccountType": "checking",
				"routingNumber": "string",
				"lastFourAccountNumber": "7000"},
			"wallet": {
				"walletID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43"},
			"card": {
				"cardID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43",
				"fingerprint": "9948962d92a1ce40c9f918cd9ece3a22bde62fb325a2f1fe2e833969de672ba3",
				"brand": "Discover",
				"cardType": "debit",
				"lastFourCardNumber": "1234",
				"bin": "123456",
				"expiration": {"month": "01","year": "21"},
				"holderName": "Jules Jackson",
				"billingAddress": {
					"addressLine1": "123 Main Street",
					"addressLine2": "Apt 302",
					"city": "Boulder",
					"stateOrProvince": "CO",
					"postalCode": "80301",
					"country": "US"},
				"cardVerification": {
					"cvv": "match",
					"addressLine1": "match",
					"postalCode": "match"},
				"issuer": "GRINGOTTS BANK",
				"issuerCountry": "US",
				"cardOnFile": true,
				"merchantAccountID": "50469144-f859-46dc-bdbd-9587c2fa7b42",
				"cardAccountUpdater": {
					"updatedOn": "2019-08-24T14:15:22Z",
					"updateType": "number-update"
				},
				"domesticPushToCard": "fast-funds"},
			"applePay": {
				"brand": "Discover",
				"cardType": "debit",
				"cardDisplayName": "Visa 1234",
				"fingerprint": "9948962d92a1ce40c9f918cd9ece3a22bde62fb325a2f1fe2e833969de672ba3",
				"expiration": {"month": "01","year": "21"},
				"dynamicLastFour": "1234"},
			"achDetails": {
				"status": "initiated",
				"traceNumber": "124782618117",
				"return": {"code": "string","reason": "string","description": "string"},
				"correction": {"code": "string","reason": "string","description": "string"},
				"companyEntryDescription": "Gym Dues",
				"originatingCompanyName": "Whole Body Fit",
				"statusUpdates": {"initiated": "2019-08-24T14:15:22Z","originated": "2019-08-24T14:15:22Z","corrected": "2019-08-24T14:15:22Z","returned": "2019-08-24T14:15:22Z","completed": "2019-08-24T14:15:22Z"
				},
				"debitHoldPeriod": "2-days"},
			"cardDetails": {
				"status": "initiated",
				"failureCode": "call-issuer",
				"dynamicDescriptor": "WhlBdy *Yoga 11-12",
				"transactionSource": "first-recurring",
				"interchangeQualification": "Visa Signature and Visa Infinite (Spend not-qualified) Product 1",
				"statusUpdates": {"initiated": "2019-08-24T14:15:22Z","confirmed": "2019-08-24T14:15:22Z","settled": "2019-08-24T14:15:22Z","failed": "2019-08-24T14:15:22Z","canceled": "2019-08-24T14:15:22Z","completed": "2019-08-24T14:15:22Z"}},
			"transferID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43"},
		"destination": {
			"paymentMethodID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43",
			"paymentMethodType": "moov-wallet",
			"account": {
				"accountID": "3dfff852-927d-47e8-822c-2fffc57ff6b9",
				"email": "amanda@classbooker.dev",
				"displayName": "Whole Body Fitness"
			},
			"bankAccount": {
				"bankAccountID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43",
				"fingerprint": "9948962d92a1ce40c9f918cd9ece3a22bde62fb325a2f1fe2e833969de672ba3",
				"status": "new",
				"holderName": "Jules Jackson",
				"holderType": "individual",
				"bankName": "Chase Bank",
				"bankAccountType": "checking",
				"routingNumber": "string",
				"lastFourAccountNumber": "7000"},
			"wallet": {
				"walletID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43"
			},
			"card": {
				"cardID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43",
				"fingerprint": "9948962d92a1ce40c9f918cd9ece3a22bde62fb325a2f1fe2e833969de672ba3",
				"brand": "Discover",
				"cardType": "debit",
				"lastFourCardNumber": "1234",
				"bin": "123456",
				"expiration": {"month": "01","year": "21"},
				"holderName": "Jules Jackson",
				"billingAddress": {
					"addressLine1": "123 Main Street",
					"addressLine2": "Apt 302",
					"city": "Boulder",
					"stateOrProvince": "CO",
					"postalCode": "80301",
					"country": "US"},
				"cardVerification": {
					"cvv": "match",
					"addressLine1": "match",
					"postalCode": "match"},
				"issuer": "GRINGOTTS BANK",
				"issuerCountry": "US",
				"cardOnFile": true,
				"merchantAccountID": "50469144-f859-46dc-bdbd-9587c2fa7b42",
				"cardAccountUpdater": {
					"updatedOn": "2019-08-24T14:15:22Z",
					"updateType": "number-update"},
				"domesticPushToCard": "fast-funds"},
			"applePay": {
				"brand": "Discover",
				"cardType": "debit",
				"cardDisplayName": "Visa 1234",
				"fingerprint": "9948962d92a1ce40c9f918cd9ece3a22bde62fb325a2f1fe2e833969de672ba3",
				"expiration": {"month": "01","year": "21"},
				"dynamicLastFour": "1234"},
			"achDetails": {
				"status": "initiated",
				"traceNumber": "124782618117",
				"return": {
				"code": "string",
				"reason": "string",
				"description": "string"},
				"correction": {
					"code": "string",
					"reason": "string",
					"description": "string"},
				"companyEntryDescription": "Gym Dues",
				"originatingCompanyName": "Whole Body Fit",
				"statusUpdates": {
				"initiated": "2019-08-24T14:15:22Z","originated": "2019-08-24T14:15:22Z","corrected": "2019-08-24T14:15:22Z","returned": "2019-08-24T14:15:22Z","completed": "2019-08-24T14:15:22Z"}
			},
			"cardDetails": {
				"status": "initiated",
				"failureCode": "call-issuer",
				"dynamicDescriptor": "WhlBdy *Yoga 11-12",
				"transactionSource": "first-recurring",
				"interchangeQualification": "Visa Signature and Visa Infinite (Spend not-qualified) Product 1",
				"statusUpdates": {
					"initiated": "2019-08-24T14:15:22Z","confirmed": "2019-08-24T14:15:22Z","settled": "2019-08-24T14:15:22Z","failed": "2019-08-24T14:15:22Z","canceled": "2019-08-24T14:15:22Z","completed": "2019-08-24T14:15:22Z"}}}
		}`)
	transfer := new(moov.SynchronousTransfer)

	dec := json.NewDecoder(bytes.NewReader(input))
	dec.DisallowUnknownFields()

	err := dec.Decode(&transfer)
	require.NoError(t, err)

	require.Equal(t, "Gym Dues", transfer.Source.AchDetails.CompanyEntryDescription)
	t.Logf("%#v", transfer)
}

type TransferTestSuite struct {
	suite.Suite
	// values for testing will be set in init()
	accountID           string
	card                moov.Card
	deleteCardID        string
	paymentMethodSource moov.PaymentMethod
	paymentMethodDest   moov.PaymentMethod
	transfer            moov.SynchronousTransfer
}

// listen for 'go test' command --> run test methods
func TestTransferSuite(t *testing.T) {
	suite.Run(t, new(TransferTestSuite))
}

func (s *TransferTestSuite) SetupSuite() {

	// Sandbox accounts have a "Lincoln National Corporation" moov account added by default. Get it's AccountID so we can test against it
	mc := NewTestClient(s.T())

	accounts, err := mc.ListAccounts(context.Background(), moov.WithAccountName("Lincoln National Corporation"))
	s.NoError(err)

	for _, account := range accounts {
		if account.DisplayName == "Lincoln National Corporation" {
			// set the accountID for testing
			s.accountID = account.AccountID
		}
	}

	// get card to get paymentID
	cards, _ := mc.ListCards(context.Background(), s.accountID)
	if len(cards) == 0 {
		card := moov.CreateCard{
			CardNumber: "371111111111114",
			CardCvv:    "1234",
			Expiration: moov.Expiration{
				Month: "10",
				Year:  "28",
			},
			HolderName: "Wade Arnold",
			BillingAddress: moov.Address{
				AddressLine1:    "123 Main Street",
				City:            "Golden",
				StateOrProvince: "CO",
				PostalCode:      "80401",
				Country:         "US",
			},
		}

		respCard, err := mc.CreateCard(context.Background(), s.accountID, card)
		s.NoError(err, "Error creating card")

		s.card = *respCard
		s.deleteCardID = respCard.CardID
	} else {
		s.card = cards[0]
		s.deleteCardID = ""
	}

	// get payment method from card
	respPaymentMethods, err := mc.ListPaymentMethods(context.Background(), s.accountID, moov.WithPaymentMethodSourceID(s.card.CardID))
	s.NoError(err)
	s.Require().NotEmpty(respPaymentMethods)

	s.paymentMethodSource = respPaymentMethods[0]

	// get payment method of wallet
	respWallets, err := mc.ListWallets(s.accountID)
	s.NoError(err)

	respPaymentMethods1, err := mc.ListPaymentMethods(context.Background(), s.accountID, moov.WithPaymentMethodSourceID(respWallets[0].WalletID))
	s.NoError(err)
	s.Require().NotEmpty(respPaymentMethods1)
	s.paymentMethodDest = respPaymentMethods1[0]

	//	get sample transfer
	payload := moov.SearchQueryPayload{}
	respTransfers, err := mc.ListTransfers(payload)
	s.NoError(err)
	s.Require().NotEmpty(respTransfers)

	if len(respTransfers) > 0 {
		s.transfer = respTransfers[0]
	}
}

func (s *TransferTestSuite) TearDownSuite() {
	mc := NewTestClient(s.T())

	//delete the card we created
	if s.deleteCardID != "" {
		err := mc.DisableCard(context.Background(), s.accountID, s.deleteCardID)
		s.NoError(err)
	}
}

func (s *TransferTestSuite) TestCreateTransfer() {
	source := moov.Source{
		PaymentMethodID: s.paymentMethodSource.PaymentMethodID,
		Card:            s.card,
		CardDetails: moov.CardDetails{
			DynamicDescriptor: "WhlBdy *Yoga 11-12",
			TransactionSource: "first-recurring",
		},
	}
	destination := moov.Destination{
		PaymentMethodID: s.paymentMethodDest.PaymentMethodID,
		Wallet:          s.paymentMethodDest.Wallet,
		AchDetails: moov.AchDetails{
			CompanyEntryDescription: "Gym Dues",
			OriginatingCompanyName:  "Whole Body Fit",
		},
	}
	amount := moov.Amount{
		Currency: "USD",
		Value:    1204,
	}
	facilitatorFee := moov.FacilitatorFee{
		Total: 8, // $0.08
	}
	description := "Pay Instructor for May 15 Class"
	metadata := map[string]string{
		"property1": "string1",
		"property2": "string2",
	}

	mc := NewTestClient(s.T())

	completedTransfer, startedTransfer, err := mc.CreateTransfer(context.Background(), moov.CreateTransfer{
		Source:         source,
		Destination:    destination,
		Amount:         amount,
		FacilitatorFee: facilitatorFee,
		Description:    description,
		Metadata:       metadata,
	}, true)

	s.Require().NoError(err, "Error creating transfer")
	s.Require().Nil(startedTransfer) // We asked it to be synchronous so hopefully is nil
	s.Require().NotNil(completedTransfer)
	s.Require().NotEmpty(completedTransfer.TransferID)

	s.transfer = *completedTransfer
}

func (s *TransferTestSuite) TestListTransfers() {
	mc := NewTestClient(s.T())

	payload := moov.SearchQueryPayload{}
	transfers, err := mc.ListTransfers(payload)
	s.NoError(err)

	s.NotEmpty(transfers)
}

func (s *TransferTestSuite) TestGetTransfer() {
	mc := NewTestClient(s.T())

	transferID := "1ddcedb6-83a3-4ce9-89a4-c3d03106dda4"
	if s.transfer.TransferID != "" {
		transferID = s.transfer.TransferID
	}

	transfer, err := mc.GetTransfer(transferID, "")
	s.NoError(err)

	s.Equal(transferID, transfer.TransferID)
}

func (s *TransferTestSuite) TestUpdateTransferMetaData() {
	metadata := map[string]string{"property1": "property  1",
		"property2": "property  2"}

	mc := NewTestClient(s.T())

	transferID := "1ddcedb6-83a3-4ce9-89a4-c3d03106dda4"
	if s.transfer.TransferID != "" {
		transferID = s.transfer.TransferID
	}

	transfer, err := mc.UpdateTransferMetaData(transferID, "", metadata)
	s.NoError(err, "Error updating transfer metadata")

	s.Equal(transfer.Metadata, metadata)
}

func (s *TransferTestSuite) TestTransferOptions() {
	mc := NewTestClient(s.T())

	payload := moov.TransferOptionsPayload{
		Source: moov.TransferOptionsSourcePayload{
			PaymentMethodID: s.paymentMethodSource.PaymentMethodID,
			AccountID:       s.accountID,
		},
		Destination: moov.TransferOptionsDestinationPayload{
			PaymentMethodID: s.paymentMethodSource.PaymentMethodID,
			AccountID:       s.accountID,
		},
		Amount: moov.Amount{
			Currency: "USD",
			Value:    1204,
		},
	}

	options, err := mc.TransferOptions(payload)
	s.NoError(err)

	// @todo check if dest or origin are not empty?
	fmt.Println(options)
	// assert.NotNil(s.T(), options)
}

func (s *TransferTestSuite) TestRefundTransfer() {
	mc := NewTestClient(s.T())

	transferID := "1ddcedb6-83a3-4ce9-89a4-c3d03106dda4"
	if s.transfer.TransferID != "" {
		transferID = s.transfer.TransferID
	}

	refund, err := mc.RefundTransfer(transferID, true, 1000)
	s.NoError(err)

	s.NotEmpty(refund.RefundID)
}

func (s *TransferTestSuite) TestListRefunds() {
	mc := NewTestClient(s.T())

	transferID := "1ddcedb6-83a3-4ce9-89a4-c3d03106dda4"
	if s.transfer.TransferID != "" {
		transferID = s.transfer.TransferID
	}

	refunds, err := mc.ListRefunds(transferID)
	s.NoError(err)

	fmt.Println(len(refunds))
	s.NotNil(refunds)
}

func (s *TransferTestSuite) TestGetRefund() {
	mc := NewTestClient(s.T())

	transferID := "1ddcedb6-83a3-4ce9-89a4-c3d03106dda4"
	if s.transfer.TransferID != "" {
		transferID = s.transfer.TransferID
	}

	refundID := "8b491eb3-a262-4eba-a0ca-35983bef3262"
	refund, err := mc.GetRefund(transferID, refundID)
	s.NoError(err)

	s.Equal(refundID, refund.RefundID)
}

func (s *TransferTestSuite) TestReverseTransfer() {
	mc := NewTestClient(s.T())

	transferID := "1ddcedb6-83a3-4ce9-89a4-c3d03106dda4"
	if s.transfer.TransferID != "" {
		transferID = s.transfer.TransferID
	}

	reverse, err := mc.ReverseTransfer(transferID, 50)
	s.NoError(err)

	s.NotEmpty(reverse.Refund.RefundID)
}
