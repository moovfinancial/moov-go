package mv2604_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2604"
)

func TestLinkGooglePayToken(t *testing.T) {
	client := mv2604.NewGooglePayClient(newTestClient(t, &expectedRequest{
		method:   http.MethodPost,
		path:     "/accounts/a/google-pay/tokens",
		jsonBody: `{"merchantAccountID":"m","paymentMethodData":{"info":{},"tokenizationData":{"token":"fixture"}}}`,
	}, http.StatusOK, `[{"paymentMethodID":"p","paymentMethodType":"google-pay"}]`))

	got, err := client.LinkGooglePayToken(t.Context(), "a", mv2604.LinkGooglePay{
		MerchantAccountID: "m",
		PaymentMethodData: mv2604.GooglePayPaymentMethodData{
			TokenizationData: mv2604.GooglePayTokenizationData{Token: "fixture"},
		},
	})
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, "p", got[0].PaymentMethodID)
	require.Equal(t, "google-pay", got[0].PaymentMethodType)
}

func TestLinkGooglePayTokenRequiresAccountID(t *testing.T) {
	client := mv2604.NewGooglePayClient(newTestClient(t, nil, http.StatusBadRequest, ""))

	got, err := client.LinkGooglePayToken(t.Context(), "", mv2604.LinkGooglePay{})
	require.ErrorIs(t, err, moov.ErrAccountIDRequired)
	require.Nil(t, got)
}
