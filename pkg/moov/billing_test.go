package moov_test

import (
	"context"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ListFeePlanAgreements(t *testing.T) {
	mc := NewTestClient(t)

	agreements, err := mc.ListFeePlanAgreements(context.Background(), FACILITATOR_ID)
	require.NoError(t, err)
	require.NotNil(t, agreements)
}

func Test_ListFeePlanAgreements_WithFilters(t *testing.T) {
	mc := NewTestClient(t)

	agreements, err := mc.ListFeePlanAgreements(
		context.Background(),
		FACILITATOR_ID,
		moov.WithFeePlanAgreementCount(10),
		moov.WithFeePlanAgreementSkip(0),
		moov.WithFeePlanAgreementStatuses([]moov.FeePlanAgreementStatus{moov.FeePlanAgreementStatus_Active}),
	)
	require.NoError(t, err)
	require.NotNil(t, agreements)
}

func Test_ListFeePlanAgreements_WithCount(t *testing.T) {
	mc := NewTestClient(t)

	agreements, err := mc.ListFeePlanAgreements(
		context.Background(),
		FACILITATOR_ID,
		moov.WithFeePlanAgreementCount(5),
	)
	require.NoError(t, err)
	require.NotNil(t, agreements)
}

func Test_ListFeePlanAgreements_WithSkip(t *testing.T) {
	mc := NewTestClient(t)

	agreements, err := mc.ListFeePlanAgreements(
		context.Background(),
		FACILITATOR_ID,
		moov.WithFeePlanAgreementSkip(10),
	)
	require.NoError(t, err)
	require.NotNil(t, agreements)
}

func Test_ListFeePlanAgreements_WithStatuses(t *testing.T) {
	mc := NewTestClient(t)

	agreements, err := mc.ListFeePlanAgreements(
		context.Background(),
		FACILITATOR_ID,
		moov.WithFeePlanAgreementStatuses([]moov.FeePlanAgreementStatus{moov.FeePlanAgreementStatus_Active}),
	)
	require.NoError(t, err)
	require.NotNil(t, agreements)
}

func Test_ListFeePlanAgreements_WithMultipleStatuses(t *testing.T) {
	mc := NewTestClient(t)

	agreements, err := mc.ListFeePlanAgreements(
		context.Background(),
		FACILITATOR_ID,
		moov.WithFeePlanAgreementStatuses([]moov.FeePlanAgreementStatus{
			moov.FeePlanAgreementStatus_Active,
			moov.FeePlanAgreementStatus_Terminated,
		}),
	)
	require.NoError(t, err)
	require.NotNil(t, agreements)
}

func Test_ListFeePlanAgreements_WithAgreementIds(t *testing.T) {
	mc := NewTestClient(t)

	// First, get all agreements to find valid IDs
	allAgreements, err := mc.ListFeePlanAgreements(context.Background(), FACILITATOR_ID)
	require.NoError(t, err)

	if len(allAgreements) > 0 {
		// Use the first agreement ID for testing
		agreementIDs := []string{allAgreements[0].AgreementID}

		agreements, err := mc.ListFeePlanAgreements(
			context.Background(),
			FACILITATOR_ID,
			moov.WithFeePlanAgreementIds(agreementIDs),
		)
		require.NoError(t, err)
		require.NotNil(t, agreements)

		if len(agreements) > 0 {
			assert.Equal(t, agreementIDs[0], agreements[0].AgreementID)
		}
	}
}
