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

func TestNewEvent(t *testing.T) {
	accountCreated := AccountCreated{
		AccountID: uuid.NewString(),
	}

	transferCreated := TransferCreated{
		AccountID:  accountCreated.AccountID,
		TransferID: uuid.NewString(),
		Status:     moov.TransferStatus_Created,
	}

	// Initialize the HTTP handler func for the target webhook URL
	webhookHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
		event, err := NewEvent(r.Body)
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

	rec := httptest.NewRecorder()

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
				CreatedOn: time.Now().UTC(),
			}

			var body bytes.Buffer
			err = json.NewEncoder(&body).Encode(event)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/my-awesome-webhook-url", &body)
			webhookHandlerFunc(rec, req)
		})
	}
}
