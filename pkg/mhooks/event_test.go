package mhooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

func TestParseEvent(t *testing.T) {
	var (
		timestamp = "2024-04-26T21:20:55Z"
		secret    = "my-webhook-signing-secret"
		signature = "6231d03752de6963087e6aea1c78a27a0617b6df1c071195f30ed85defe34e02fd0bf3995949fe12dafd747c42de9cfae03b8aafcf69cceba5495f4c7b719d82"

		accountCreated = AccountCreated{
			AccountID: uuid.NewString(),
		}

		transferCreated = TransferCreated{
			AccountID:  accountCreated.AccountID,
			TransferID: uuid.NewString(),
			Status:     moov.TransferStatus_Created,
		}
	)
	createdOn, err := time.Parse(time.RFC3339, timestamp)
	require.NoError(t, err)

	// Initialize the HTTP handler func for the target webhook URL
	webhookHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
		event, err := ParseEvent(r, secret)
		require.NoError(t, err)

		//nolint:exhaustive
		switch event.EventType {
		case EventTypeAccountCreated:
			got, err := event.AccountCreated()
			require.NoError(t, err)

			t.Logf("Got AccountCreated webhook with accountID=%v", got.AccountID)
			require.Equal(t, accountCreated, *got)
		case EventTypeTransferCreated:
			got, err := event.TransferCreated()
			require.NoError(t, err)

			t.Logf("Got TransferCreated webhook with transferID=%v\n", got.TransferID)
			require.Equal(t, transferCreated, *got)
		}

		w.WriteHeader(200)
	}

	for i, tt := range []struct {
		eventType EventType
		data      any
	}{
		{
			eventType: EventTypeAccountCreated,
			data:      accountCreated,
		},
		{
			eventType: EventTypeTransferCreated,
			data:      transferCreated,
		},
	} {
		t.Run(fmt.Sprintf("%d %v", i, tt.eventType), func(t *testing.T) {
			dataBytes, err := json.Marshal(tt.data)
			require.NoError(t, err)

			event := Event{
				EventID:   uuid.NewString(),
				EventType: tt.eventType,
				Data:      dataBytes,
				CreatedOn: createdOn,
			}

			var body bytes.Buffer
			err = json.NewEncoder(&body).Encode(event)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/my-awesome-webhook-url", &body)
			req.Header.Set("x-timestamp", timestamp)
			req.Header.Set("x-nonce", "LwxF1Uk7QOeDq2nzB3theslHbtAo7y3uuncB1PoijwCZZaRZsd8DOtffBT7p")
			req.Header.Set("x-webhook-id", "dff0a709-f982-4475-81e8-214b435c74ab")
			req.Header.Set("x-signature", signature)

			rec := httptest.NewRecorder()
			webhookHandlerFunc(rec, req)
		})
	}
}
