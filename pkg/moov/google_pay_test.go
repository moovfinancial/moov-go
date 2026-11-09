package moov_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestLinkGooglePayTokenGeneric(t *testing.T) {
	var (
		method  string
		path    string
		version string
		body    []byte
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		method = r.Method
		path = r.URL.Path
		version = r.Header.Get(moov.VersionHeader)
		body = b

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"paymentMethodID":"pm-1","paymentMethodType":"google-pay","googlePay":{"tokenID":"tok-1","brand":"Visa"}}]`))
	}))
	t.Cleanup(srv.Close)

	client, err := moov.NewClient(
		moov.WithCredentials(moov.Credentials{PublicKey: "pk", SecretKey: "sk"}),
		moov.WithMoovURLScheme("http"),
	)
	require.NoError(t, err)
	client.Credentials.Host = strings.TrimPrefix(srv.URL, "http://")

	link := moov.LinkGooglePay{
		MerchantAccountID: "merchant-1",
		PaymentMethodData: moov.GooglePayPaymentMethodData{
			Type:             "CARD",
			Info:             moov.GooglePayCardInfo{CardNetwork: "VISA", CardDetails: "1234"},
			TokenizationData: moov.GooglePayTokenizationData{Token: "gateway-token"},
		},
	}

	methods, err := moov.LinkGooglePayTokenGeneric(context.Background(), client, moov.Version2026_04, "account-123", link)
	require.NoError(t, err)

	require.Equal(t, http.MethodPost, method)
	require.Equal(t, "/accounts/account-123/google-pay/tokens", path)
	require.Equal(t, moov.Version2026_04.String(), version)
	require.JSONEq(t, `{"merchantAccountID":"merchant-1","paymentMethodData":{"type":"CARD","info":{"cardNetwork":"VISA","cardDetails":"1234"},"tokenizationData":{"token":"gateway-token"}}}`, string(body))

	require.Len(t, methods, 1)
	require.Equal(t, "pm-1", methods[0].PaymentMethodID)
	require.Equal(t, moov.PaymentMethodType_GooglePay, methods[0].PaymentMethodType)
	require.NotNil(t, methods[0].GooglePay)
	require.Equal(t, "tok-1", methods[0].GooglePay.TokenID)

	t.Run("rejected token preserves the API error", func(t *testing.T) {
		rejected := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"token rejected"}`))
		}))
		t.Cleanup(rejected.Close)

		client.Credentials.Host = strings.TrimPrefix(rejected.URL, "http://")

		_, err := moov.LinkGooglePayTokenGeneric(context.Background(), client, moov.Version2026_04, "account-123", link)
		require.Error(t, err)

		httpErr := moov.ErrorAsHttpCallResponse(err)
		require.NotNil(t, httpErr)
		require.Equal(t, moov.StatusBadRequest, httpErr.Status())
	})
}
