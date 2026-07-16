package moov_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

func TestLatestAPIFields(t *testing.T) {
	t.Run("transfer card fields", func(t *testing.T) {
		var transfer moov.Transfer
		err := json.Unmarshal([]byte(`{
			"source": {"cardDetails": {"deferredOn": "2026-07-16T12:00:00Z"}},
			"destination": {"cardDetails": {"scheduledDeliveryOn": "2026-07-17T12:00:00Z"}}
		}`), &transfer)
		require.NoError(t, err)
		require.Equal(t, "2026-07-16T12:00:00Z", transfer.Source.CardDetails.DeferredOn.Format(time.RFC3339))
		require.Equal(t, "2026-07-17T12:00:00Z", transfer.Destination.CardDetails.ScheduledDeliveryOn.Format(time.RFC3339))
	})

	t.Run("statement card cancellations", func(t *testing.T) {
		var statement moov.Statement
		err := json.Unmarshal([]byte(`{"cardAcquiringFees":{"visa":{"cardCancellations":{"count":2}}}}`), &statement)
		require.NoError(t, err)
		require.EqualValues(t, 2, statement.CardAcquiringFees.Visa.CardCancellations.Count)
	})

	t.Run("fee attribution", func(t *testing.T) {
		var fee moov.IncurredFee
		err := json.Unmarshal([]byte(`{"feePaidBy":"destination"}`), &fee)
		require.NoError(t, err)
		require.Equal(t, moov.FeePaidByDestination, *fee.FeePaidBy)
	})
}
