package mv2604_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2604"
)

func TestLinkGooglePayToken(t *testing.T) {
	var (
		method  string
		path    string
		version string
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		method = r.Method
		path = r.URL.Path
		version = r.Header.Get(moov.VersionHeader)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"paymentMethodID":"pm-1","paymentMethodType":"google-pay","googlePay":{"tokenID":"tok-1"}}]`))
	}))
	t.Cleanup(srv.Close)

	client, err := moov.NewClient(
		moov.WithCredentials(moov.Credentials{PublicKey: "pk", SecretKey: "sk"}),
		moov.WithMoovURLScheme("http"),
	)
	require.NoError(t, err)
	client.Credentials.Host = strings.TrimPrefix(srv.URL, "http://")

	googlePay := mv2604.NewGooglePayClient(client)

	methods, err := googlePay.LinkGooglePayToken(context.Background(), "account-123", mv2604.LinkGooglePay{
		MerchantAccountID: "merchant-1",
		PaymentMethodData: mv2604.GooglePayPaymentMethodData{
			TokenizationData: mv2604.GooglePayTokenizationData{Token: "gateway-token"},
		},
	})
	require.NoError(t, err)

	require.Equal(t, http.MethodPost, method)
	require.Equal(t, "/accounts/account-123/google-pay/tokens", path)
	require.Equal(t, moov.Version2026_04.String(), version)

	require.Len(t, methods, 1)
	require.Equal(t, "pm-1", methods[0].PaymentMethodID)
	require.NotNil(t, methods[0].GooglePay)
	require.Equal(t, "tok-1", methods[0].GooglePay.TokenID)
}
