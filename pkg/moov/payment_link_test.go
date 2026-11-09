package moov_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/internal/testtools"
	"github.com/moovfinancial/moov-go/pkg/moov"
)

// merchantWalletPaymentMethod creates a merchant wallet because payment links
// require a wallet owned by the merchant account, then returns its payment method ID.
func merchantWalletPaymentMethod(t *testing.T, mc *moov.Client, accountID string) string {
	t.Helper()

	wallet, err := mc.CreateWallet(context.Background(), accountID, moov.CreateWallet{
		Name: fmt.Sprintf("moov-go payment link %d", time.Now().UnixNano()),
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		closed := moov.WalletStatus_Closed
		_, _ = mc.UpdateWallet(context.Background(), accountID, wallet.WalletID, moov.UpdateWallet{Status: &closed})
	})

	pms, err := mc.ListPaymentMethods(context.Background(), accountID, moov.WithPaymentMethodType(string(moov.PaymentMethodType_MoovWallet)))
	require.NoError(t, err)
	for _, pm := range pms {
		if pm.Wallet != nil && pm.Wallet.WalletID == wallet.WalletID {
			return pm.PaymentMethodID
		}
	}
	require.FailNow(t, "no payment method found for created wallet")
	return ""
}

func Test_PaymentLink_CreateGetListUpdateDisable(t *testing.T) {
	var (
		mc        = NewTestClient(t)
		ctx       = t.Context()
		accountID = testtools.MERCHANT_ID
	)

	create := moov.CreatePaymentLink{
		PartnerAccountID:        testtools.PARTNER_ID,
		MerchantPaymentMethodID: merchantWalletPaymentMethod(t, mc, accountID),
		Amount: &moov.Amount{
			Currency: "USD",
			Value:    1050,
		},
		Display: moov.PaymentLinkDisplayOptions{
			Title:        "Testing moov-go",
			Description:  "Payment link created by the moov-go test suite",
			CallToAction: moov.CallToAction_Pay,
		},
		Payment: &moov.PaymentLinkPaymentDetails{
			AllowedMethods: []moov.CollectionPaymentMethodType{
				moov.CollectionPaymentMethodType_CardPayment,
			},
		},
	}

	// create the payment link
	created, err := mc.CreatePaymentLink(ctx, accountID, create)
	require.NoError(t, err)
	require.NotEmpty(t, created.Code)
	require.NotEmpty(t, created.Link)
	require.Equal(t, moov.PaymentLinkType_Payment, created.PaymentLinkType)
	require.Equal(t, moov.PaymentLinkStatus_Active, created.Status)
	require.Equal(t, testtools.PARTNER_ID, created.PartnerAccountID)
	require.Equal(t, accountID, created.MerchantAccountID)

	code := string(created.Code)
	t.Cleanup(func() {
		_ = mc.DisablePaymentLink(ctx, accountID, code)
	})

	// fetch by code
	fetched, err := mc.GetPaymentLink(ctx, accountID, code)
	require.NoError(t, err)
	require.Equal(t, created.Code, fetched.Code)

	// list filtered by type
	listed, err := mc.ListPaymentLinks(ctx, accountID,
		moov.WithPaymentLinkTypes(moov.PaymentLinkType_Payment),
		moov.WithPaymentLinkCount(100))
	require.NoError(t, err)
	require.Contains(t, listed, *created)

	// update the display title and set an expiration
	expiresOn := time.Now().UTC().Add(24 * time.Hour).Truncate(time.Second)
	updated, err := mc.UpdatePaymentLink(ctx, accountID, code, moov.UpdatePaymentLink{
		Display: &moov.PaymentLinkDisplayOptionsUpdate{
			Title: "Updated moov-go title",
		},
		ExpiresOn: moov.Set(expiresOn),
	})
	require.NoError(t, err)
	require.Equal(t, "Updated moov-go title", updated.Display.Title)
	require.NotNil(t, updated.ExpiresOn)

	// clear the expiration with an explicit null
	updated, err = mc.UpdatePaymentLink(ctx, accountID, code, moov.UpdatePaymentLink{
		ExpiresOn: moov.SetNull[time.Time](),
	})
	require.NoError(t, err)
	require.Nil(t, updated.ExpiresOn)

	// fetch the QR code image
	img, contentType, err := mc.GetPaymentLinkQRCode(ctx, accountID, code)
	require.NoError(t, err)
	require.NotEmpty(t, img)
	require.Contains(t, contentType, "image/png")

	// disable the payment link
	err = mc.DisablePaymentLink(ctx, accountID, code)
	require.NoError(t, err)

	// disabled link can still be fetched and reports its status
	fetched, err = mc.GetPaymentLink(ctx, accountID, code)
	require.NoError(t, err)
	require.Equal(t, moov.PaymentLinkStatus_Disabled, fetched.Status)
	require.NotNil(t, fetched.DisabledOn)
}

func Test_PaymentLink_CustomAmount(t *testing.T) {
	var (
		mc        = NewTestClient(t)
		ctx       = t.Context()
		accountID = testtools.MERCHANT_ID
	)

	create := moov.CreatePaymentLink{
		PartnerAccountID:        testtools.PARTNER_ID,
		MerchantPaymentMethodID: merchantWalletPaymentMethod(t, mc, accountID),
		Display: moov.PaymentLinkDisplayOptions{
			Title:        "Custom amount moov-go test",
			Description:  "Pay what you want",
			CallToAction: moov.CallToAction_Donate,
		},
		CustomAmountPayment: &moov.PaymentLinkCustomAmountPaymentDetails{
			AllowedMethods: []moov.CollectionPaymentMethodType{
				moov.CollectionPaymentMethodType_CardPayment,
			},
			AmountRange: &moov.AmountDecimalRange{
				Minimum: &moov.AmountDecimal{Currency: "USD", ValueDecimal: "1.00"},
				Maximum: &moov.AmountDecimal{Currency: "USD", ValueDecimal: "100.00"},
			},
		},
	}

	created, err := mc.CreatePaymentLink(ctx, accountID, create)
	require.NoError(t, err)
	require.Equal(t, moov.PaymentLinkType_CustomAmountPayment, created.PaymentLinkType)
	require.Nil(t, created.Amount)

	t.Cleanup(func() {
		_ = mc.DisablePaymentLink(ctx, accountID, string(created.Code))
	})
}
