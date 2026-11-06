package moov_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

const paymentLinkJSON = `{
	"code": "plink-123",
	"paymentLinkType": "payment",
	"mode": "sandbox",
	"status": "active",
	"partnerAccountID": "partner-123",
	"merchantAccountID": "account-123",
	"ownerAccountID": "account-123",
	"merchantPaymentMethodID": "pm-123",
	"link": "https://checkout.moov.io/plink-123",
	"amount": {"currency": "USD", "value": 1050},
	"uses": 0,
	"display": {"title": "Test link", "description": "A test link", "callToAction": "pay"},
	"customer": {},
	"createdOn": "2026-07-01T00:00:00Z",
	"updatedOn": "2026-07-01T00:00:00Z"
}`

const customAmountPaymentLinkJSON = `{
	"code": "plink-custom",
	"paymentLinkType": "custom-amount-payment",
	"mode": "sandbox",
	"status": "active",
	"partnerAccountID": "partner-123",
	"merchantAccountID": "account-123",
	"ownerAccountID": "account-123",
	"merchantPaymentMethodID": "pm-123",
	"link": "https://checkout.moov.io/plink-custom",
	"uses": 0,
	"display": {"title": "Custom link", "description": "Pay what you want", "callToAction": "donate"},
	"customer": {},
	"createdOn": "2026-07-01T00:00:00Z",
	"updatedOn": "2026-07-01T00:00:00Z"
}`

func Test_PaymentLink_CreateGetListUpdateDisable(t *testing.T) {
	var (
		got struct {
			method  string
			path    string
			version string
			accept  string
			query   url.Values
			body    map[string]any
		}
		bodyErr error
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got.method = r.Method
		got.path = r.URL.Path
		got.version = r.Header.Get(moov.VersionHeader)
		got.accept = r.Header.Get("Accept")
		got.query = r.URL.Query()
		got.body = nil
		bodyErr = nil
		if r.Method == http.MethodPost || r.Method == http.MethodPatch {
			bodyErr = json.NewDecoder(r.Body).Decode(&got.body)
		}

		switch {
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/qrcode"):
			w.Header().Set("Content-Type", "image/png")
			_, _ = w.Write([]byte("PNGDATA"))
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/missing"):
			w.WriteHeader(http.StatusNotFound)
		case r.Method == http.MethodGet && r.URL.Path == "/accounts/account-123/payment-links":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte("[" + paymentLinkJSON + "]"))
		case r.Method == http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(paymentLinkJSON))
		case r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(paymentLinkJSON))
		}
	}))
	t.Cleanup(srv.Close)

	mc := newFilesTestClient(t, srv)
	ctx := BgCtx()

	created, err := mc.CreatePaymentLink(ctx, "account-123", moov.CreatePaymentLink{
		PartnerAccountID:        "partner-123",
		MerchantPaymentMethodID: "pm-123",
		Amount:                  &moov.Amount{Currency: "USD", Value: 1050},
		Display: moov.PaymentLinkDisplayOptions{
			Title:        "Test link",
			Description:  "A test link",
			CallToAction: moov.CallToAction_Pay,
		},
		Payment: &moov.PaymentLinkPaymentDetails{
			AllowedMethods: []moov.CollectionPaymentMethodType{moov.CollectionPaymentMethodType_CardPayment},
		},
	})
	require.NoError(t, err)
	require.NoError(t, bodyErr)
	require.Equal(t, http.MethodPost, got.method)
	require.Equal(t, "/accounts/account-123/payment-links", got.path)
	require.Equal(t, moov.Version2026_07.String(), got.version)
	require.Equal(t, float64(1050), got.body["amount"].(map[string]any)["value"])
	require.NotContains(t, got.body, "customAmountPayment")
	require.Equal(t, moov.PaymentLinkStatus_Active, created.Status)
	require.Equal(t, moov.PaymentLinkCode("plink-123"), created.Code)

	fetched, err := mc.GetPaymentLink(ctx, "account-123", "plink-123")
	require.NoError(t, err)
	require.Equal(t, http.MethodGet, got.method)
	require.Equal(t, "/accounts/account-123/payment-links/plink-123", got.path)
	require.Equal(t, moov.Version2026_07.String(), got.version)
	require.Equal(t, created.Code, fetched.Code)

	listed, err := mc.ListPaymentLinks(ctx, "account-123",
		moov.WithPaymentLinkCount(100),
		moov.WithPaymentLinkSkip(5),
		moov.WithPaymentLinkTypes(moov.PaymentLinkType_Payment, moov.PaymentLinkType_Payout),
	)
	require.NoError(t, err)
	require.Equal(t, "/accounts/account-123/payment-links", got.path)
	require.Equal(t, "100", got.query.Get("count"))
	require.Equal(t, "5", got.query.Get("skip"))
	require.Equal(t, "payment,payout", got.query.Get("types"))
	require.Len(t, listed, 1)
	require.Equal(t, created.Code, listed[0].Code)

	// A payout link with delivery options and line items serializes each section.
	_, err = mc.CreatePaymentLink(ctx, "account-123", moov.CreatePaymentLink{
		PartnerAccountID:        "partner-123",
		MerchantPaymentMethodID: "pm-123",
		Amount:                  &moov.Amount{Currency: "USD", Value: 5000},
		Display: moov.PaymentLinkDisplayOptions{
			Title:        "Payout link",
			Description:  "Send a payout",
			CallToAction: moov.CallToAction_Confirm,
		},
		Payout: &moov.PaymentLinkPayoutDetails{
			AllowedMethods: []moov.DisbursementPaymentMethodType{moov.DisbursementPaymentMethodType_PushToCard},
			Recipient:      moov.PayoutRecipient{Email: moov.PtrOf(moov.Email("recipient@example.com"))},
			PushOptions:    &moov.PushOptions{AllowedSpeeds: []moov.PushDeliverySpeed{moov.PushDeliverySpeed_Instant}},
			FeePaidBy:      map[string]string{"instant-push-to-card": "source"},
		},
		LineItems: &moov.CreatePaymentLinkLineItems{
			Items: []moov.CreatePaymentLinkLineItem{{
				Name:      "Widget",
				BasePrice: moov.AmountDecimal{Currency: "USD", ValueDecimal: "50.00"},
				Quantity:  1,
			}},
		},
	})
	require.NoError(t, err)
	require.NoError(t, bodyErr)
	payout := got.body["payout"].(map[string]any)
	require.Equal(t, "recipient@example.com", payout["recipient"].(map[string]any)["email"])
	require.Equal(t, "source", payout["feePaidBy"].(map[string]any)["instant-push-to-card"])
	require.Len(t, got.body["lineItems"].(map[string]any)["items"], 1)

	// An update sets a value, clears with an explicit null, or omits the field.
	expiresOn := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	updated, err := mc.UpdatePaymentLink(ctx, "account-123", "plink-123", moov.UpdatePaymentLink{
		Display:   &moov.PaymentLinkDisplayOptionsUpdate{Title: "Updated title"},
		ExpiresOn: moov.Set(expiresOn),
	})
	require.NoError(t, err)
	require.NoError(t, bodyErr)
	require.Equal(t, http.MethodPatch, got.method)
	require.Equal(t, "/accounts/account-123/payment-links/plink-123", got.path)
	require.Equal(t, moov.Version2026_07.String(), got.version)
	require.Equal(t, "Updated title", got.body["display"].(map[string]any)["title"])
	require.Equal(t, "2026-08-01T00:00:00Z", got.body["expiresOn"])
	require.Equal(t, created.Code, updated.Code)

	_, err = mc.UpdatePaymentLink(ctx, "account-123", "plink-123", moov.UpdatePaymentLink{
		ExpiresOn: moov.SetNull[time.Time](),
	})
	require.NoError(t, err)
	require.NoError(t, bodyErr)
	value, ok := got.body["expiresOn"]
	require.True(t, ok)
	require.Nil(t, value)
	require.NotContains(t, got.body, "amount")

	img, contentType, err := mc.GetPaymentLinkQRCode(ctx, "account-123", "plink-123")
	require.NoError(t, err)
	require.Equal(t, http.MethodGet, got.method)
	require.Equal(t, "/accounts/account-123/payment-links/plink-123/qrcode", got.path)
	require.Equal(t, "image/*", got.accept)
	require.Equal(t, moov.Version2026_07.String(), got.version)
	require.Equal(t, []byte("PNGDATA"), img)
	require.Equal(t, "image/png", contentType)

	err = mc.DisablePaymentLink(ctx, "account-123", "plink-123")
	require.NoError(t, err)
	require.Equal(t, http.MethodDelete, got.method)
	require.Equal(t, "/accounts/account-123/payment-links/plink-123", got.path)
	require.Equal(t, moov.Version2026_07.String(), got.version)

	_, err = mc.GetPaymentLink(ctx, "account-123", "missing")
	require.Error(t, err)
	var httpErr moov.HttpCallResponse
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, moov.StatusNotFound, httpErr.Status())
}

func Test_PaymentLink_CustomAmount(t *testing.T) {
	var (
		got struct {
			method string
			path   string
			body   map[string]any
		}
		bodyErr error
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got.method = r.Method
		got.path = r.URL.Path
		got.body = nil
		bodyErr = json.NewDecoder(r.Body).Decode(&got.body)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(customAmountPaymentLinkJSON))
	}))
	t.Cleanup(srv.Close)

	mc := newFilesTestClient(t, srv)

	created, err := mc.CreatePaymentLink(BgCtx(), "account-123", moov.CreatePaymentLink{
		PartnerAccountID:        "partner-123",
		MerchantPaymentMethodID: "pm-123",
		Display: moov.PaymentLinkDisplayOptions{
			Title:        "Custom link",
			Description:  "Pay what you want",
			CallToAction: moov.CallToAction_Donate,
		},
		CustomAmountPayment: &moov.PaymentLinkCustomAmountPaymentDetails{
			AllowedMethods: []moov.CollectionPaymentMethodType{moov.CollectionPaymentMethodType_CardPayment},
			AmountRange: &moov.AmountDecimalRange{
				Minimum: &moov.AmountDecimal{Currency: "USD", ValueDecimal: "1.00"},
				Maximum: &moov.AmountDecimal{Currency: "USD", ValueDecimal: "100.00"},
			},
			SuggestedAmounts: []moov.AmountDecimal{
				{Currency: "USD", ValueDecimal: "5.00"},
				{Currency: "USD", ValueDecimal: "10.00"},
			},
		},
	})
	require.NoError(t, err)
	require.NoError(t, bodyErr)
	require.Equal(t, http.MethodPost, got.method)
	require.Equal(t, "/accounts/account-123/payment-links", got.path)

	// A custom amount link must not send a fixed amount.
	require.NotContains(t, got.body, "amount")
	custom := got.body["customAmountPayment"].(map[string]any)
	require.Equal(t, "100.00", custom["amountRange"].(map[string]any)["maximum"].(map[string]any)["valueDecimal"])
	require.Len(t, custom["suggestedAmounts"], 2)

	require.Equal(t, moov.PaymentLinkType_CustomAmountPayment, created.PaymentLinkType)
	require.Nil(t, created.Amount)
}
