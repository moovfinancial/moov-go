package moov_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

func Test_GetTransferRiskOutcomes_NotFound(t *testing.T) {
	mc := NewTestClient(t)

	// The endpoint is partner-gated and collapses every denial path (not
	// allowlisted, transfer not owned, no persisted outcome) to the same 404,
	// so an unknown transfer ID is a deterministic not-found regardless of the
	// calling account's allowlist status.
	transferID := uuid.NewString()
	outcomes, err := mc.GetTransferRiskOutcomes(context.Background(), transferID)
	require.Nil(t, outcomes)

	var httpErr moov.HttpCallResponse
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, moov.StatusNotFound, httpErr.Status())
}
