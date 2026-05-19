package mv2607_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2607"
)

func TestTransferClientSetsVersionHeader(t *testing.T) {
	type capturedRequest struct {
		method  string
		path    string
		version string
		body    string
	}

	var captured []capturedRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		captured = append(captured, capturedRequest{
			method:  r.Method,
			path:    r.URL.Path,
			version: r.Header.Get(moov.VersionHeader),
			body:    string(body),
		})

		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/accounts/account-123/transfers":
			_, _ = w.Write([]byte(`{
				"transferID": "transfer-123",
				"amountDetails": {
					"surcharge": {
						"currency": "USD",
						"valueDecimal": "0.03"
					}
				}
			}`))
		case "/accounts/account-123/transfers/transfer-123/refunds":
			_, _ = w.Write([]byte(`{
				"refundID": "refund-123",
				"amountDetails": {
					"surcharge": {
						"currency": "USD",
						"valueDecimal": "0.03"
					}
				}
			}`))
		case "/accounts/account-123/transfers/transfer-123/reversals":
			_, _ = w.Write([]byte(`{
				"refund": {
					"refundID": "refund-123",
					"amountDetails": {
						"surcharge": {
							"currency": "USD",
							"valueDecimal": "0.03"
						}
					}
				}
			}`))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)

	client, err := moov.NewClient(
		moov.WithCredentials(moov.Credentials{PublicKey: "pk", SecretKey: "sk"}),
		moov.WithMoovURLScheme("http"),
	)
	require.NoError(t, err)
	client.Credentials.Host = strings.TrimPrefix(srv.URL, "http://")
	transfers := mv2607.NewTransferClient(client)
	surcharge := surchargeAmount()

	transfer, _, err := transfers.CreateTransfer(context.Background(), "account-123", mv2607.CreateTransfer{
		Source: moov.CreateTransfer_Source{
			PaymentMethodID: "source-payment-method",
		},
		Destination: moov.CreateTransfer_Destination{
			PaymentMethodID: "destination-payment-method",
		},
		Amount: moov.Amount{
			Currency: "USD",
			Value:    100,
		},
		AmountDetails: &mv2607.CreateTransferAmountDetails{
			Surcharge: &surcharge,
		},
	}).WaitForRailResponse()
	require.NoError(t, err)
	require.NotNil(t, transfer.AmountDetails)
	require.Equal(t, surcharge, *transfer.AmountDetails.Surcharge)

	refund, _, err := transfers.RefundTransfer(context.Background(), "account-123", "transfer-123", mv2607.CreateRefund{
		Amount: 97,
		AmountDetails: &mv2607.RefundAmountDetails{
			Surcharge: &surcharge,
		},
	}, moov.WithRefundWaitForRailResponse())
	require.NoError(t, err)
	require.NotNil(t, refund.AmountDetails)
	require.Equal(t, surcharge, *refund.AmountDetails.Surcharge)

	reversal, err := transfers.ReverseTransfer(context.Background(), "account-123", "transfer-123", mv2607.CreateReversal{
		Amount: 97,
		AmountDetails: &mv2607.RefundAmountDetails{
			Surcharge: &surcharge,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, reversal.Refund.AmountDetails)
	require.Equal(t, surcharge, *reversal.Refund.AmountDetails.Surcharge)

	require.Len(t, captured, 3)
	for _, req := range captured {
		require.Equal(t, moov.Version2026_07.String(), req.version)
		require.Contains(t, req.body, `"amountDetails"`)
		require.Contains(t, req.body, `"surcharge"`)
	}
	require.Equal(t, http.MethodPost, captured[0].method)
	require.Equal(t, "/accounts/account-123/transfers", captured[0].path)
	require.Equal(t, http.MethodPost, captured[1].method)
	require.Equal(t, "/accounts/account-123/transfers/transfer-123/refunds", captured[1].path)
	require.Equal(t, http.MethodPost, captured[2].method)
	require.Equal(t, "/accounts/account-123/transfers/transfer-123/reversals", captured[2].path)
}
