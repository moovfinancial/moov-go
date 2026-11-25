package mhooks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

const (
	eventID        = "evt-01"
	eventCreatedOn = "2026-04-01T12:00:00Z"
)

func signedWebhookRequest(body string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(body))
	req.Header.Set("x-timestamp", testTimestamp)
	req.Header.Set("x-nonce", testNonce)
	req.Header.Set("x-webhook-id", testWebhookID)
	req.Header.Set("x-signature", testSignature)
	return req
}

func webhookBody(eventType EventType, data string) string {
	return fmt.Sprintf(`{"eventID":%q,"type":%q,"createdOn":%q,"data":%s}`, eventID, eventType, eventCreatedOn, data)
}

type eventCase struct {
	eventType    EventType
	data         string
	checkPayload func(t *testing.T, event, otherEvent *Event)
}

func newEventCase[T any](eventType EventType, getter func(Event) (*T, error), data string, want T) eventCase {
	return eventCase{
		eventType: eventType,
		data:      data,
		checkPayload: func(t *testing.T, event, otherEvent *Event) {
			t.Helper()
			got, err := getter(*event)
			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, want, *got)

			got, err = getter(*otherEvent)
			require.Nil(t, got)
			require.EqualError(t, err, fmt.Sprintf("invalid event type: expected %s but got %s", eventType, otherEvent.EventType))
		},
	}
}

func TestParseEvent(t *testing.T) {
	visaMID := "v1"
	networkIDUpdatedOn := time.Date(2026, time.April, 2, 0, 0, 0, 0, time.UTC)
	transferID := "t"

	tests := []eventCase{
		newEventCase("account.created", Event.AccountCreated,
			`{"accountID":"a1","foreignID":"f1"}`,
			AccountCreated{AccountID: "a1", ForeignID: "f1"},
		),
		newEventCase("account.disconnected", Event.AccountDisconnected,
			`{"accountID":"a2","foreignID":"f2"}`,
			AccountDisconnected{AccountID: "a2", ForeignID: "f2"},
		),
		newEventCase("account.updated", Event.AccountUpdated,
			`{"accountID":"a3","foreignID":"f3"}`,
			AccountUpdated{AccountID: "a3", ForeignID: "f3"},
		),
		newEventCase("authorization.expiring", Event.AuthorizationExpiring,
			`{"merchantAccountID":"m","partnerAccountID":"p","authorizationID":"u","expiresOn":"2026-04-09T15:00:00Z"}`,
			AuthorizationExpiring{
				MerchantAccountID: "m",
				PartnerAccountID:  "p",
				AuthorizationID:   "u",
				ExpiresOn:         time.Date(2026, time.April, 9, 15, 0, 0, 0, time.UTC),
			},
		),
		newEventCase("balance.updated", Event.BalanceUpdated,
			`{"accountID":"a","walletID":"w"}`,
			BalanceUpdated{AccountID: "a", WalletID: "w"},
		),
		newEventCase("bankAccount.created", Event.BankAccountCreated,
			`{"accountID":"a","bankAccountID":"b1","status":"new"}`,
			BankAccountCreated{AccountID: "a", BankAccountID: "b1", Status: "new"},
		),
		newEventCase("bankAccount.deleted", Event.BankAccountDeleted,
			`{"accountID":"a","bankAccountID":"b2"}`,
			BankAccountDeleted{AccountID: "a", BankAccountID: "b2"},
		),
		newEventCase("bankAccount.updated", Event.BankAccountUpdated,
			`{"bankAccountID":"b3","status":"errored","statusReason":"other","exceptionDetails":{"description":"fixture"}}`,
			BankAccountUpdated{
				BankAccountID:    "b3",
				Status:           "errored",
				StatusReason:     "other",
				ExceptionDetails: &moov.ExceptionDetails{Description: "fixture"},
			},
		),
		newEventCase("billingStatement.created", Event.BillingStatementCreated,
			`{"statementID":"s1"}`,
			BillingStatementCreated{StatementID: "s1"},
		),
		newEventCase("cancellation.created", Event.CancellationCreated,
			`{"cancellationID":"c1","transferID":"t1","status":"pending"}`,
			CancellationCreated{CancellationID: "c1", TransferID: "t1", Status: "pending"},
		),
		newEventCase("cancellation.updated", Event.CancellationUpdated,
			`{"cancellationID":"c2","transferID":"t2","status":"completed"}`,
			CancellationUpdated{CancellationID: "c2", TransferID: "t2", Status: "completed"},
		),
		newEventCase("card.autoUpdated", Event.CardAutoUpdated,
			`{"accountID":"a","cardID":"c","updateType":"expiration-update"}`,
			CardAutoUpdated{AccountID: "a", CardID: "c", UpdateType: "expiration-update"},
		),
		newEventCase("capability.requested", Event.CapabilityRequested,
			`{"accountID":"a","capabilityID":"transfers"}`,
			CapabilityRequested{AccountID: "a", Capability: "transfers"},
		),
		newEventCase("capability.updated", Event.CapabilityUpdated,
			`{"accountID":"a","capabilityID":"send-funds","status":"enabled"}`,
			CapabilityUpdated{AccountID: "a", Capability: "send-funds", Status: "enabled"},
		),
		newEventCase("dispute.created", Event.DisputeCreated,
			`{"accountID":"a","disputeID":"d1","transferID":"t","transactionID":"x"}`,
			DisputeCreated{AccountID: "a", DisputeID: "d1", TransferID: "t", TransactionID: "x"},
		),
		newEventCase("dispute.updated", Event.DisputeUpdated,
			`{"disputeID":"d2","status":"under-review","phase":"chargeback"}`,
			DisputeUpdated{DisputeID: "d2", Status: "under-review", Phase: "chargeback"},
		),
		newEventCase("invoice.created", Event.InvoiceCreated,
			`{"accountID":"a","invoiceID":"i1"}`,
			InvoiceCreated{AccountID: "a", InvoiceID: "i1"},
		),
		newEventCase("invoice.updated", Event.InvoiceUpdated,
			`{"accountID":"a","invoiceID":"i2","status":"paid"}`,
			InvoiceUpdated{AccountID: "a", InvoiceID: "i2", Status: "paid"},
		),
		newEventCase("networkID.updated", Event.NetworkIDUpdated,
			`{"accountID":"a","visaMid":"v1","updatedOn":"2026-04-02T00:00:00Z"}`,
			NetworkIDUpdated{AccountID: "a", VisaMid: &visaMID, UpdatedOn: &networkIDUpdatedOn},
		),
		newEventCase("paymentMethod.disabled", Event.PaymentMethodDisabled,
			`{"accountID":"a","paymentMethodID":"p1","sourceID":"s1"}`,
			PaymentMethodDisabled{AccountID: "a", PaymentMethodID: "p1", SourceID: "s1"},
		),
		newEventCase("paymentMethod.enabled", Event.PaymentMethodEnabled,
			`{"accountID":"a","paymentMethodID":"p2","sourceID":"s2"}`,
			PaymentMethodEnabled{AccountID: "a", PaymentMethodID: "p2", SourceID: "s2"},
		),
		newEventCase("refund.created", Event.RefundCreated,
			`{"accountID":"a","transferID":"t1","refundID":"r1"}`,
			RefundCreated{AccountID: "a", TransferID: "t1", RefundID: "r1"},
		),
		newEventCase("refund.updated", Event.RefundUpdated,
			`{"transferID":"t2","refundID":"r2","status":"failed"}`,
			RefundUpdated{TransferID: "t2", RefundID: "r2", Status: "failed"},
		),
		newEventCase("representative.created", Event.RepresentativeCreated,
			`{"accountID":"a","representativeID":"r1"}`,
			RepresentativeCreated{AccountID: "a", RepresentativeID: "r1"},
		),
		newEventCase("representative.deleted", Event.RepresentativeDeleted,
			`{"accountID":"b","representativeID":"r2"}`,
			RepresentativeDeleted{AccountID: "b", RepresentativeID: "r2"},
		),
		newEventCase("representative.updated", Event.RepresentativeUpdated,
			`{"accountID":"c","representativeID":"r3"}`,
			RepresentativeUpdated{AccountID: "c", RepresentativeID: "r3"},
		),
		newEventCase("sweep.created", Event.SweepCreated,
			`{"sweepID":"s1","walletID":"w"}`,
			SweepCreated{SweepID: "s1", WalletID: "w"},
		),
		newEventCase("sweep.updated", Event.SweepUpdated,
			`{"sweepID":"s2","walletID":"w","status":"paid","transferID":"t"}`,
			SweepUpdated{SweepID: "s2", WalletID: "w", Status: "paid", TransferID: &transferID},
		),
		newEventCase("event.test", Event.TestPing,
			`{"ping":true}`,
			TestPing{Ping: true},
		),
		newEventCase("ticket.created", Event.TicketCreated,
			`{"accountID":"a","ticketID":"k1","foreignID":"f1"}`,
			TicketCreated{AccountID: "a", TicketID: "k1", ForeignID: "f1"},
		),
		newEventCase("ticket.updated", Event.TicketUpdated,
			`{"ticketID":"k2","status":"in-progress"}`,
			TicketUpdated{TicketID: "k2", Status: "in-progress"},
		),
		newEventCase("ticket.messageAdded", Event.TicketMessageAdded,
			`{"accountID":"a","ticketID":"k3","foreignID":"f3"}`,
			TicketMessageAdded{AccountID: "a", TicketID: "k3", ForeignID: "f3"},
		),
		newEventCase("transfer.created", Event.TransferCreated,
			`{"accountID":"a","transferID":"t1","status":"created"}`,
			TransferCreated{AccountID: "a", TransferID: "t1", Status: "created"},
		),
		newEventCase("transfer.updated", Event.TransferUpdated,
			`{"transferID":"t","status":"source.settled","source":{"accountID":"a1","paymentMethodID":"p1"},"destination":{"accountID":"a2","paymentMethodID":"p2"}}`,
			TransferUpdated{
				TransferID:  "t",
				Status:      "source.settled",
				Source:      PaymentMethodPartial{AccountID: "a1", PaymentMethodID: "p1"},
				Destination: PaymentMethodPartial{AccountID: "a2", PaymentMethodID: "p2"},
			},
		),
		newEventCase("wallet.created", Event.WalletCreated,
			`{"accountID":"a","walletID":"w1"}`,
			WalletCreated{AccountID: "a", WalletID: "w1"},
		),
		newEventCase("wallet.updated", Event.WalletUpdated,
			`{"walletID":"w2","status":"closed"}`,
			WalletUpdated{WalletID: "w2", Status: "closed"},
		),
		newEventCase("walletTransaction.updated", Event.WalletTransactionUpdated,
			`{"walletID":"w","transactionID":"x","status":"completed","availableBalance":{"currency":"USD","value":1250,"valueDecimal":"12.50"}}`,
			WalletTransactionUpdated{
				WalletID: "w", TransactionID: "x", Status: "completed",
				AvailableBalance: &moov.AvailableBalance{Currency: "USD", Value: 1250, ValueDecimal: "12.50"},
			},
		),
	}

	createdOn := time.Date(2026, time.April, 1, 12, 0, 0, 0, time.UTC)
	ping, err := ParseEvent(signedWebhookRequest(webhookBody("event.test", `{"ping":true}`)), testSecret)
	require.NoError(t, err)
	require.NotNil(t, ping)
	account, err := ParseEvent(signedWebhookRequest(webhookBody("account.created", `{"accountID":"a1"}`)), testSecret)
	require.NoError(t, err)
	require.NotNil(t, account)

	for _, tc := range tests {
		t.Run(string(tc.eventType), func(t *testing.T) {
			event, err := ParseEvent(signedWebhookRequest(webhookBody(tc.eventType, tc.data)), testSecret)
			require.NoError(t, err)
			require.NotNil(t, event)
			require.Equal(t, eventID, event.EventID)
			require.Equal(t, tc.eventType, event.EventType)
			require.Equal(t, createdOn, event.CreatedOn.UTC())
			require.JSONEq(t, tc.data, string(event.Data))

			otherEvent := ping
			if tc.eventType == "event.test" {
				otherEvent = account
			}
			tc.checkPayload(t, event, otherEvent)
		})
	}
}

func TestParseEvent_SweepUpdatedWithoutTransferID(t *testing.T) {
	event, err := ParseEvent(signedWebhookRequest(webhookBody("sweep.updated", `{"sweepID":"s2"}`)), testSecret)
	require.NoError(t, err)
	require.NotNil(t, event)

	got, err := event.SweepUpdated()
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, "s2", got.SweepID)
	require.Nil(t, got.TransferID)
}

func TestParseEvent_Errors(t *testing.T) {
	tests := []struct {
		name      string
		body      string
		signature string
		want      func(*testing.T, error)
	}{
		{
			name:      "invalid signature takes precedence over invalid JSON",
			body:      "{",
			signature: "abcd",
			want: func(t *testing.T, err error) {
				require.ErrorIs(t, err, ErrInvalidSignature)
			},
		},
		{
			name: "empty body",
			want: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "decoding event:")
				require.ErrorIs(t, err, io.EOF)
			},
		},
		{
			name: "invalid JSON",
			body: `{"type":!}`,
			want: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "decoding event:")
				var syntaxErr *json.SyntaxError
				require.ErrorAs(t, err, &syntaxErr)
			},
		},
		{
			name: "unknown event type",
			body: webhookBody("custom.unknown", "{}"),
			want: func(t *testing.T, err error) {
				require.EqualError(t, err, "invalid event type: custom.unknown")
			},
		},
		{
			name: "missing event type",
			body: `{"eventID":"evt-01","createdOn":"2026-04-01T12:00:00Z","data":{}}`,
			want: func(t *testing.T, err error) {
				require.EqualError(t, err, "invalid event type: ")
			},
		},
		{
			name: "wrong payload field type",
			body: webhookBody("event.test", `{"ping":"yes"}`),
			want: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "unmarshalling event data:")
				var typeErr *json.UnmarshalTypeError
				require.ErrorAs(t, err, &typeErr)
				require.Equal(t, "ping", typeErr.Field)
			},
		},
		{
			name: "missing payload",
			body: `{"eventID":"evt-01","type":"event.test","createdOn":"2026-04-01T12:00:00Z"}`,
			want: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "unmarshalling event data:")
				var syntaxErr *json.SyntaxError
				require.ErrorAs(t, err, &syntaxErr)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := signedWebhookRequest(tc.body)
			if tc.signature != "" {
				req.Header.Set("x-signature", tc.signature)
			}
			event, err := ParseEvent(req, testSecret)
			require.Nil(t, event)
			tc.want(t, err)
		})
	}
}
