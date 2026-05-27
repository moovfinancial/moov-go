package moov_test

import (
	"bytes"
	"encoding/json"
	"testing"

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
	assert.Equal(t, "CRYPTOGRAM_3DS", pm.GooglePay.AuthMethod)
}
