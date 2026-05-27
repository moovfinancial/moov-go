package moov_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGooglePayPaymentMethodMarshal(t *testing.T) {
	input := []byte(`{
		"paymentMethodID": "c3f8c4d0-7b1a-4e2a-9c5b-1234567890ab",
		"paymentMethodType": "google-pay",
		"googlePay": {
			"tokenID": "b0a1c2d3-e4f5-6789-abcd-ef0123456789",
			"brand": "Visa",
			"cardType": "debit",
			"cardDisplayName": "Visa 1234",
			"fingerprint": "9948962d92a1ce40c9f918cd9ece3a22bde62fb325a2f1fe2e833969de672ba3",
			"expiration": {
				"month": "01",
				"year": "29"
			},
			"dynamicLastFour": "1234",
			"issuerCountry": "US",
			"authMethod": "CRYPTOGRAM_3DS"
		}
	}`)

	pm := new(moov.PaymentMethod)

	dec := json.NewDecoder(bytes.NewReader(input))
	dec.DisallowUnknownFields()

	err := dec.Decode(pm)

	require.NoError(t, err)
	require.NotNil(t, pm.GooglePay)
	assert.Equal(t, "b0a1c2d3-e4f5-6789-abcd-ef0123456789", pm.GooglePay.TokenID)
	assert.Equal(t, moov.CardBrand("Visa"), pm.GooglePay.Brand)
	assert.Equal(t, moov.CardType("debit"), pm.GooglePay.CardType)
	assert.Equal(t, "Visa 1234", pm.GooglePay.CardDisplayName)
	assert.Equal(t, "9948962d92a1ce40c9f918cd9ece3a22bde62fb325a2f1fe2e833969de672ba3", pm.GooglePay.Fingerprint)
	assert.Equal(t, "01", pm.GooglePay.Expiration.Month)
	assert.Equal(t, "29", pm.GooglePay.Expiration.Year)
	assert.Equal(t, "1234", pm.GooglePay.DynamicLastFour)
	assert.Equal(t, "US", pm.GooglePay.IssuerCountry)
	assert.Equal(t, moov.GooglePayAuthMethod_Cryptogram3DS, pm.GooglePay.AuthMethod)
}

func TestGooglePayTransferMarshal(t *testing.T) {
	googlePayJSON := `{
		"tokenID": "b0a1c2d3-e4f5-6789-abcd-ef0123456789",
		"brand": "Visa",
		"cardType": "debit",
		"cardDisplayName": "Visa 1234",
		"fingerprint": "9948962d92a1ce40c9f918cd9ece3a22bde62fb325a2f1fe2e833969de672ba3",
		"expiration": {"month": "01", "year": "29"},
		"dynamicLastFour": "1234",
		"issuerCountry": "US",
		"authMethod": "CRYPTOGRAM_3DS"
	}`

	transferInput := []byte(`{
		"transferID": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
		"createdOn": "2026-04-01T00:00:00Z",
		"status": "completed",
		"amount": {"currency": "USD", "value": 1000},
		"source": {
			"paymentMethodID": "c3f8c4d0-7b1a-4e2a-9c5b-1234567890ab",
			"paymentMethodType": "google-pay",
			"account": {"accountID": "acc-111", "email": "payer@example.com", "displayName": "Payer"},
			"googlePay": ` + googlePayJSON + `
		},
		"destination": {
			"paymentMethodID": "d4e5f6a7-b8c9-0123-def0-123456789abc",
			"paymentMethodType": "moov-wallet",
			"account": {"accountID": "acc-222", "email": "payee@example.com", "displayName": "Payee"},
			"wallet": {"walletID": "wallet-123"}
		}
	}`)

	xfr := new(moov.Transfer)

	dec := json.NewDecoder(bytes.NewReader(transferInput))
	dec.DisallowUnknownFields()

	err := dec.Decode(xfr)

	require.NoError(t, err)
	assert.Equal(t, "a1b2c3d4-e5f6-7890-abcd-ef1234567890", xfr.TransferID)
	assert.Equal(t, time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC), xfr.CreatedOn)

	require.NotNil(t, xfr.Source.GooglePay)
	assert.Equal(t, "b0a1c2d3-e4f5-6789-abcd-ef0123456789", xfr.Source.GooglePay.TokenID)
	assert.Equal(t, moov.CardBrand("Visa"), xfr.Source.GooglePay.Brand)
	assert.Equal(t, "9948962d92a1ce40c9f918cd9ece3a22bde62fb325a2f1fe2e833969de672ba3", xfr.Source.GooglePay.Fingerprint)
	assert.Equal(t, "US", xfr.Source.GooglePay.IssuerCountry)

	assert.Nil(t, xfr.Destination.GooglePay)
	require.NotNil(t, xfr.Destination.Wallet)
	assert.Equal(t, "wallet-123", xfr.Destination.Wallet.WalletID)
}
