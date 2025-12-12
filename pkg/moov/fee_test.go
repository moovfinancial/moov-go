package moov_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

func Test_ListFees(t *testing.T) {
	mc := NewTestClient(t)

	fees, err := mc.GetFees(t.Context(), FACILITATOR_ID)
	require.NoError(t, err)
	require.NotNil(t, fees)
}

func Test_ListFees_WithCount(t *testing.T) {
	mc := NewTestClient(t)

	fees, err := mc.GetFees(
		t.Context(),
		FACILITATOR_ID,
		moov.WithFeeCount(5),
	)
	require.NoError(t, err)
	require.NotNil(t, fees)
	if len(fees) > 0 {
		assert.LessOrEqual(t, len(fees), 5)
	}
}

func Test_ListFees_WithSkip(t *testing.T) {
	mc := NewTestClient(t)

	fees, err := mc.GetFees(
		t.Context(),
		FACILITATOR_ID,
		moov.WithFeeSkip(2),
	)
	require.NoError(t, err)
	require.NotNil(t, fees)
}

func Test_ListFees_WithTransferID(t *testing.T) {
	mc := NewTestClient(t)

	// First get some fees to find a transferID
	allFees, err := mc.GetFees(t.Context(), FACILITATOR_ID, moov.WithFeeCount(100))
	require.NoError(t, err)

	// Find a fee with a transferID
	var transferID string
	for _, fee := range allFees {
		if fee.GeneratedBy != nil && fee.GeneratedBy.TransferID != nil {
			transferID = *fee.GeneratedBy.TransferID
			break
		}
	}

	if transferID != "" {
		fees, err := mc.GetFees(
			t.Context(),
			FACILITATOR_ID,
			moov.WithFeeTransferID(transferID),
		)
		require.NoError(t, err)
		require.NotNil(t, fees)

		// Verify all returned fees are for the specified transfer
		for _, fee := range fees {
			require.NotNil(t, fee.GeneratedBy)
			require.NotNil(t, fee.GeneratedBy.TransferID)
			assert.Equal(t, transferID, *fee.GeneratedBy.TransferID)
		}
	}
}

func Test_ListFees_WithDisputeID(t *testing.T) {
	mc := NewTestClient(t)

	// First get some fees to find a disputeID
	allFees, err := mc.GetFees(t.Context(), FACILITATOR_ID, moov.WithFeeCount(100))
	require.NoError(t, err)

	// Find a fee with a disputeID
	var disputeID string
	for _, fee := range allFees {
		if fee.GeneratedBy != nil && fee.GeneratedBy.DisputeID != nil {
			disputeID = *fee.GeneratedBy.DisputeID
			break
		}
	}

	if disputeID != "" {
		fees, err := mc.GetFees(
			t.Context(),
			FACILITATOR_ID,
			moov.WithFeeDisputeID(disputeID),
		)
		require.NoError(t, err)
		require.NotNil(t, fees)

		// Verify all returned fees are for the specified dispute
		for _, fee := range fees {
			require.NotNil(t, fee.GeneratedBy)
			require.NotNil(t, fee.GeneratedBy.DisputeID)
			assert.Equal(t, disputeID, *fee.GeneratedBy.DisputeID)
		}
	}
}

func Test_ListFees_WithDateTimeRange(t *testing.T) {
	mc := NewTestClient(t)

	fees, err := mc.GetFees(
		t.Context(),
		FACILITATOR_ID,
		moov.WithFeeStartDateTime("2024-01-01T00:00:00Z"),
		moov.WithFeeEndDateTime("2025-12-31T23:59:59Z"),
	)
	require.NoError(t, err)
	require.NotNil(t, fees)
}

func Test_FetchFees(t *testing.T) {
	mc := NewTestClient(t)

	// First list some fees to get feeIDs
	allFees, err := mc.GetFees(t.Context(), FACILITATOR_ID, moov.WithFeeCount(5))
	require.NoError(t, err)

	if len(allFees) > 0 {
		// Collect feeIDs
		feeIDs := make([]string, 0, len(allFees))
		for _, fee := range allFees {
			if fee.FeeID != "" {
				feeIDs = append(feeIDs, fee.FeeID)
			}
		}

		if len(feeIDs) > 0 {
			// Fetch fees by IDs
			request := moov.FeeListRequest{
				FeeIDs: feeIDs,
			}

			fetchedFees, err := mc.ListFees(t.Context(), FACILITATOR_ID, request)
			require.NoError(t, err)
			require.NotNil(t, fetchedFees)
			assert.Equal(t, len(feeIDs), len(fetchedFees))
		}
	}
}
