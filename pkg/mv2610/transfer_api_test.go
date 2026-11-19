package mv2610_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2610"
)

// The response covers optional nested details: an ACH return and a wallet captureID.
func TestListTransferEvents(t *testing.T) {
	var path, version string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		version = r.Header.Get(moov.VersionHeader)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{
				"eventID": "evt-capture",
				"eventType": "capture",
				"amount": {"currency": "USD", "valueDecimal": "100.00"},
				"details": {"capture": {"captureID": "capt-1", "status": "completed"}},
				"occurredOn": "2026-08-25T16:01:01Z"
			},
			{
				"eventID": "evt-ach",
				"eventType": "ach-debit",
				"amount": {"currency": "USD", "valueDecimal": "100.00"},
				"details": {"achDebit": {"status": "returned", "return": {"code": "R01", "reason": "Insufficient funds"}}},
				"occurredOn": "2026-08-25T14:05:00Z"
			},
			{
				"eventID": "evt-wallet",
				"eventType": "wallet-credit",
				"amount": {"currency": "USD", "valueDecimal": "100.00"},
				"details": {"walletCredit": {"accountID": "acct-1", "walletID": "wallet-1", "transactionID": "txn-1", "transactionType": "payment", "captureID": "capt-1"}},
				"occurredOn": "2026-08-25T16:01:02Z"
			}
		]`))
	}))
	t.Cleanup(srv.Close)

	client, err := moov.NewClient(
		moov.WithCredentials(moov.Credentials{PublicKey: "pk", SecretKey: "sk"}),
		moov.WithMoovURLScheme("http"),
	)
	require.NoError(t, err)
	client.Credentials.Host = strings.TrimPrefix(srv.URL, "http://")

	transfers := mv2610.NewTransferClient(client)

	_, err = transfers.ListTransferEvents(context.Background(), "account-123", "")
	require.EqualError(t, err, "accountID and transferID are required")
	require.Empty(t, path, "an empty transferID must not send a request")

	events, err := transfers.ListTransferEvents(context.Background(), "account-123", "transfer-123")
	require.NoError(t, err)

	require.Equal(t, "/accounts/account-123/transfers/transfer-123/events", path)
	require.Equal(t, moov.Version2026_10.String(), version)

	amount := moov.AmountDecimal{Currency: "USD", ValueDecimal: "100.00"}
	captureID := "capt-1"
	require.Equal(t, []mv2610.TransferEvent{
		{
			EventID:   "evt-capture",
			EventType: mv2610.TransferEventType_Capture,
			Amount:    amount,
			Details: mv2610.TransferEventDetails{
				Capture: &mv2610.TransferEventCaptureDetails{CaptureID: "capt-1", Status: mv2610.CaptureStatus_Completed},
			},
			OccurredOn: time.Date(2026, 8, 25, 16, 1, 1, 0, time.UTC),
		},
		{
			EventID:   "evt-ach",
			EventType: mv2610.TransferEventType_AchDebit,
			Amount:    amount,
			Details: mv2610.TransferEventDetails{
				AchDebit: &mv2610.TransferEventAchDetails{
					Status: moov.AchStatus_Returned,
					Return: &moov.AchException{Code: "R01", Reason: "Insufficient funds"},
				},
			},
			OccurredOn: time.Date(2026, 8, 25, 14, 5, 0, 0, time.UTC),
		},
		{
			EventID:   "evt-wallet",
			EventType: mv2610.TransferEventType_WalletCredit,
			Amount:    amount,
			Details: mv2610.TransferEventDetails{
				WalletCredit: &mv2610.TransferEventWalletTransactionDetails{
					AccountID:       "acct-1",
					WalletID:        "wallet-1",
					TransactionID:   "txn-1",
					TransactionType: moov.WalletTransactionTypePayment,
					CaptureID:       &captureID,
				},
			},
			OccurredOn: time.Date(2026, 8, 25, 16, 1, 2, 0, time.UTC),
		},
	}, events)
}

// The transfer body is the moov-api v2026.10 getTransfer example, so the expected values come from the spec.
func TestTransferClient(t *testing.T) {
	type request struct {
		method         string
		path           string
		query          string
		version        string
		idempotencyKey string
		body           string
	}
	var requests []request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		requests = append(requests, request{
			method:         r.Method,
			path:           r.URL.Path,
			query:          r.URL.RawQuery,
			version:        r.Header.Get(moov.VersionHeader),
			idempotencyKey: r.Header.Get("X-Idempotency-Key"),
			body:           string(body),
		})
		w.Header().Set("Content-Type", "application/json")

		switch r.Method + " " + r.URL.Path {
		case "POST /accounts/acct-1/transfers", "GET /accounts/acct-1/transfers/xfr-1":
			_, _ = w.Write([]byte(awaitingCaptureTransferJSON))
		case "GET /accounts/acct-1/transfers":
			_, _ = w.Write([]byte("[" + awaitingCaptureTransferJSON + "]"))
		case "POST /accounts/acct-1/transfers/.fetch":
			_, _ = w.Write([]byte(`{"xfr-1":` + awaitingCaptureTransferJSON + `}`))
		case "POST /accounts/acct-1/transfers/xfr-1/captures":
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"captureID":"capt-1","amount":{"currency":"USD","valueDecimal":"25.00"},"isFinal":false,"status":"completed","createdOn":"2024-08-24T14:15:22Z"}`))
		case "GET /accounts/acct-1/transfers/xfr-1/captures":
			_, _ = w.Write([]byte(`[{"captureID":"capt-1","amount":{"currency":"USD","valueDecimal":"25.00"},"isFinal":false,"status":"completed","createdOn":"2024-08-24T14:15:22Z"}]`))
		case "GET /accounts/acct-1/transfers/xfr-1/captures/capt-1":
			_, _ = w.Write([]byte(`{"captureID":"capt-1","amount":{"currency":"USD","valueDecimal":"25.00"},"isFinal":false,"status":"completed","createdOn":"2024-08-24T14:15:22Z"}`))
		case "POST /accounts/acct-1/transfers/xfr-1/cancellations", "GET /accounts/acct-1/transfers/xfr-1/cancellations/canc-1":
			_, _ = w.Write([]byte(`{"cancellationID":"canc-1","status":"completed","createdOn":"2024-08-24T14:15:22Z"}`))
		case "GET /accounts/acct-1/transfers/xfr-1/cancellations":
			_, _ = w.Write([]byte(`[{"cancellationID":"canc-1","status":"completed","createdOn":"2024-08-24T14:15:22Z"}]`))
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
	transfers := mv2610.NewTransferClient(client)
	ctx := context.Background()

	manual := mv2610.CaptureMethod_Manual
	created, _, err := transfers.CreateTransfer(ctx, "acct-1", mv2610.CreateTransfer{
		Source:      mv2610.CreateTransferSource{PaymentMethodID: "pm-card"},
		Destination: mv2610.CreateTransferDestination{PaymentMethodID: "pm-wallet"},
		Amount:      moov.AmountDecimal{Currency: "USD", ValueDecimal: "329.45"},
		Options:     &mv2610.CreateTransferRailOptions{CardPayment: &mv2610.CreateCardPaymentOptions{CaptureMethod: &manual}},
	}).WaitForRailResponse()
	require.NoError(t, err)

	got, err := transfers.GetTransfer(ctx, "acct-1", "xfr-1")
	require.NoError(t, err)
	require.Equal(t, created, got)
	require.Equal(t, mv2610.TransferType_CardPayment, got.TransferType)
	require.Equal(t, moov.TransferStatus_AwaitingCapture, got.Status)
	require.Equal(t, moov.AmountDecimal{Currency: "USD", ValueDecimal: "329.45"}, got.Amount)
	require.Equal(t, "8.25", got.AmountDetails.Tax.ValueDecimal)
	require.Equal(t, "2000", got.Source.Card.LastFourCardNumber)
	require.Equal(t, mv2610.CaptureMethod_Manual, *got.Options.CardPayment.CaptureMethod)
	require.True(t, *got.Options.CardPayment.IsExtendedAuthorizationRequested)
	require.Equal(t, moov.CardTransactionStatus_Confirmed, *got.ProcessingDetails.CardPayment.Status)
	require.Equal(t, "615318762094", *got.ProcessingDetails.CardPayment.RetrievalReferenceNumber)
	require.Equal(t, "329.45", got.Authorization.CapturableAmount.ValueDecimal)
	require.Equal(t, "0.00", got.Authorization.CapturedAmount.ValueDecimal)
	require.False(t, *got.Authorization.IsExtendedAuthorizationGranted)

	listed, err := transfers.ListTransfers(ctx, "acct-1", moov.WithTransferTypes([]string{
		string(mv2610.TransferType_CardPayment), string(mv2610.TransferType_AchDebit),
	}))
	require.NoError(t, err)
	require.Equal(t, []mv2610.Transfer{*got}, listed)

	fetched, err := transfers.FetchTransfers(ctx, "acct-1", moov.FetchTransfersRequest{TransferIDs: []string{"xfr-1"}})
	require.NoError(t, err)
	require.Equal(t, map[string]mv2610.Transfer{"xfr-1": *got}, fetched)

	isFinal := true
	capture, err := transfers.CreateCapture(ctx, "acct-1", "xfr-1", mv2610.CreateCapture{
		Amount:  &moov.AmountDecimal{Currency: "USD", ValueDecimal: "25.00"},
		IsFinal: &isFinal,
	})
	require.NoError(t, err, "a 201 capture response is a success")
	require.Equal(t, "capt-1", capture.CaptureID)
	require.Equal(t, mv2610.CaptureStatus_Completed, capture.Status)

	cancellation, err := transfers.CancelTransfer(ctx, "acct-1", "xfr-1", mv2610.CreateCancellation{
		Amount: &moov.AmountDecimal{Currency: "USD", ValueDecimal: "304.45"},
	})
	require.NoError(t, err)
	wantCancellation := moov.Cancellation{
		CancellationID: "canc-1",
		Status:         moov.CancellationStatus_Completed,
		CreatedOn:      time.Date(2024, 8, 24, 14, 15, 22, 0, time.UTC),
	}
	require.Equal(t, &wantCancellation, cancellation)

	// Auth-capture transfers exist only in v2026.10, so the capture and cancellation reads must send its version header.
	cancellations, err := transfers.ListCancellations(ctx, "acct-1", "xfr-1")
	require.NoError(t, err)
	require.Equal(t, []moov.Cancellation{wantCancellation}, cancellations)

	gotCancellation, err := transfers.GetCancellation(ctx, "acct-1", "xfr-1", "canc-1")
	require.NoError(t, err)
	require.Equal(t, &wantCancellation, gotCancellation)

	captures, err := transfers.ListCaptures(ctx, "acct-1", "xfr-1")
	require.NoError(t, err)
	require.Equal(t, []mv2610.Capture{*capture}, captures)

	gotCapture, err := transfers.GetCapture(ctx, "acct-1", "xfr-1", "capt-1")
	require.NoError(t, err)
	require.Equal(t, capture, gotCapture)

	require.Len(t, requests, 10)
	for _, r := range requests {
		require.Equal(t, moov.Version2026_10.String(), r.version, r.method+" "+r.path)
	}
	require.JSONEq(t, `{
		"source":{"paymentMethodID":"pm-card"},
		"destination":{"paymentMethodID":"pm-wallet"},
		"amount":{"currency":"USD","valueDecimal":"329.45"},
		"options":{"cardPayment":{"captureMethod":"manual"}}
	}`, requests[0].body)
	require.Equal(t, "transferTypes=card-payment%2Cach-debit", requests[2].query)
	require.JSONEq(t, `{"amount":{"currency":"USD","valueDecimal":"25.00"},"isFinal":true}`, requests[4].body)
	require.JSONEq(t, `{"amount":{"currency":"USD","valueDecimal":"304.45"}}`, requests[5].body)
	for _, i := range []int{0, 4, 5} {
		_, err := uuid.Parse(requests[i].idempotencyKey)
		require.NoError(t, err, "%s %s needs an idempotency key", requests[i].method, requests[i].path)
	}
}

// The server answers 202 when the wait for the rail times out. The refund or reversal exists, so the SDK must
// return it, not an error. The bodies are the moov-api v2026.10 CardAcquiringRefund and Reversal shapes.
func TestTransferClientRefundsAndReversals(t *testing.T) {
	const refundJSON = `{"refundID":"ref-1","createdOn":"2024-08-24T14:15:22Z","updatedOn":"2024-08-24T14:15:22Z",` +
		`"status":"pending","amount":{"currency":"USD","valueDecimal":"10.00"},"processingDetails":{}}`
	type request struct {
		version        string
		idempotencyKey string
		waitFor        string
		body           string
	}
	requestByRoute := map[string]request{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		key := r.Method + " " + r.URL.Path
		requestByRoute[key] = request{
			version:        r.Header.Get(moov.VersionHeader),
			idempotencyKey: r.Header.Get("X-Idempotency-Key"),
			waitFor:        r.Header.Get("X-Wait-For"),
			body:           string(body),
		}
		w.Header().Set("Content-Type", "application/json")

		switch key {
		case "POST /accounts/acct-1/transfers/xfr-1/refunds":
			w.WriteHeader(http.StatusAccepted)
			_, _ = w.Write([]byte(refundJSON))
		case "GET /accounts/acct-1/transfers/xfr-1/refunds":
			_, _ = w.Write([]byte("[" + refundJSON + "]"))
		case "GET /accounts/acct-1/transfers/xfr-1/refunds/ref-1":
			_, _ = w.Write([]byte(refundJSON))
		case "POST /accounts/acct-1/transfers/xfr-1/reversals":
			w.WriteHeader(http.StatusAccepted)
			_, _ = w.Write([]byte(`{"refund":` + refundJSON + `}`))
		case "POST /accounts/acct-1/transfers/xfr-2/reversals":
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(`{"error":"transfer cannot be reversed"}`))
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
	transfers := mv2610.NewTransferClient(client)
	ctx := context.Background()

	wantRefund := mv2610.Refund{
		RefundID:  "ref-1",
		CreatedOn: time.Date(2024, 8, 24, 14, 15, 22, 0, time.UTC),
		UpdatedOn: time.Date(2024, 8, 24, 14, 15, 22, 0, time.UTC),
		Status:    moov.RefundStatus_Pending,
		Amount:    moov.AmountDecimal{Currency: "USD", ValueDecimal: "10.00"},
	}

	refund, err := transfers.RefundTransfer(ctx, "acct-1", "xfr-1", mv2610.CreateRefund{}, moov.WithRefundWaitForRailResponse())
	require.NoError(t, err, "a 202 refund is a success")
	require.Equal(t, &wantRefund, refund)

	refunds, err := transfers.ListRefunds(ctx, "acct-1", "xfr-1")
	require.NoError(t, err)
	require.Equal(t, []mv2610.Refund{wantRefund}, refunds)

	gotRefund, err := transfers.GetRefund(ctx, "acct-1", "xfr-1", "ref-1")
	require.NoError(t, err)
	require.Equal(t, &wantRefund, gotRefund)

	reversal, err := transfers.ReverseTransfer(ctx, "acct-1", "xfr-1", mv2610.CreateReversal{})
	require.NoError(t, err, "a 202 reversal is a success")
	require.Equal(t, &mv2610.CreatedReversal{Refund: &wantRefund}, reversal)

	_, err = transfers.ReverseTransfer(ctx, "acct-1", "xfr-2", mv2610.CreateReversal{})
	require.Error(t, err, "a 409 reversal is not a success")

	require.Len(t, requestByRoute, 5)
	for key, r := range requestByRoute {
		require.Equal(t, moov.Version2026_10.String(), r.version, key)
	}
	for _, key := range []string{"POST /accounts/acct-1/transfers/xfr-1/refunds", "POST /accounts/acct-1/transfers/xfr-1/reversals"} {
		_, err := uuid.Parse(requestByRoute[key].idempotencyKey)
		require.NoError(t, err, "%s needs an idempotency key", key)
		require.JSONEq(t, `{}`, requestByRoute[key].body, "%s without an amount refunds or reverses the full amount", key)
	}
	require.Equal(t, "rail-response", requestByRoute["POST /accounts/acct-1/transfers/xfr-1/refunds"].waitFor)
}

// awaitingCaptureTransferJSON is the moov-api v2026.10 getTransfer example "Awaiting capture card payment transfer".
const awaitingCaptureTransferJSON = `{
	"transferID": "d835gf30-4b19-4850-a9b2-c0624c41ecb3",
	"transferType": "card-payment",
	"source": {
		"paymentMethodID": "9506dbf6-4208-44c3-ad8a-e4431660e1f2",
		"paymentMethodType": "card-payment",
		"account": {"accountID": "7e4b26c2-b399-49ef-8390-50e1ea44d550", "email": "jules@julesjackson.com", "displayName": "Jules Jackson"},
		"card": {
			"cardID": "aefd5563-93c6-413c-875e-1bd0ebfc116d",
			"fingerprint": "2f5d782ceef1c3bd31ed5...",
			"brand": "Visa",
			"cardType": "credit",
			"lastFourCardNumber": "2000",
			"bin": "400020",
			"expiration": {"month": "01", "year": "28"},
			"holderName": "Jules Jackson",
			"billingAddress": {"postalCode": "80301"},
			"domesticPushToCard": "standard",
			"domesticPullFromCard": "supported",
			"cardVerification": {
				"cvv": "unavailable",
				"addressLine1": "unavailable",
				"postalCode": "unavailable",
				"accountName": {"firstName": "unavailable", "lastName": "unavailable", "middleName": "unavailable", "fullName": "unavailable"}
			},
			"issuer": "Moov Visa Sandbox",
			"issuerCountry": "US",
			"cardAccountUpdater": {}
		}
	},
	"destination": {
		"account": {"accountID": "34233b72-780c-4a0d-8b08-cbbe238878f8", "email": "john@wholebodyfitness.io", "displayName": "Whole Body Fitness"}
	},
	"status": "awaiting-capture",
	"amount": {"currency": "USD", "valueDecimal": "329.45"},
	"amountDetails": {
		"tip": {"currency": "USD", "valueDecimal": "3.50"},
		"tax": {"currency": "USD", "valueDecimal": "8.25"}
	},
	"description": "Transfer from card to wallet",
	"createdOn": "2025-01-21T21:32:16Z",
	"moovFees": [
		{
			"accountID": "7e4b26c2-b399-49ef-8390-50e1ea44d550",
			"transferParty": "source",
			"totalAmount": {"currency": "USD", "valueDecimal": "0.10"},
			"feeIDs": ["9d957d33-1a9a-47aa-9460-fe1a90f003dd"]
		}
	],
	"options": {
		"cardPayment": {"dynamicDescriptor": "WhlBdy *Yoga 11-12", "captureMethod": "manual", "isExtendedAuthorizationRequested": true}
	},
	"processingDetails": {
		"cardPayment": {"status": "confirmed", "authorizationCode": "A1B2C3", "networkTransactionID": "123456789012345", "retrievalReferenceNumber": "615318762094"}
	},
	"authorization": {
		"authorizationID": "twkx2my2tjd2vfda7ynjb4ad3u_auth",
		"requestedAmount": {"currency": "USD", "valueDecimal": "329.45"},
		"authorizedAmount": {"currency": "USD", "valueDecimal": "329.45"},
		"capturedAmount": {"currency": "USD", "valueDecimal": "0.00"},
		"capturableAmount": {"currency": "USD", "valueDecimal": "329.45"},
		"expiresOn": "2025-01-28T21:32:16Z",
		"isExtendedAuthorizationGranted": false
	}
}`
