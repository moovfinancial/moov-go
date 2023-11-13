package moov

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
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
	transfer := new(Transfer)

	dec := json.NewDecoder(bytes.NewReader(input))
	dec.DisallowUnknownFields()

	err := dec.Decode(&transfer)
	if err != nil {
		require.NoError(t, err)
	}

	assert.Equal(t, "Gym Dues", transfer.Source.AchDetails.CompanyEntryDescription)
	t.Logf("%#v", transfer)
}
