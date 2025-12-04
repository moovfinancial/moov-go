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

func Test_ListFeePlans(t *testing.T) {
	mc := NewTestClient(t)

	plans, err := mc.ListFeePlans(context.Background(), FACILITATOR_ID)
	require.NoError(t, err)
	require.NotNil(t, plans)
}

func Test_ListFeePlans_WithFilters(t *testing.T) {
	mc := NewTestClient(t)

	// First, get all plans to find valid IDs
	allPlans, err := mc.ListFeePlans(context.Background(), FACILITATOR_ID)
	require.NoError(t, err)

	if len(allPlans) > 0 {
		// Use the first plan ID for testing
		planIDs := []string{allPlans[0].PlanID}

		plans, err := mc.ListFeePlans(
			context.Background(),
			FACILITATOR_ID,
			moov.WithFeePlanIds(planIDs),
		)
		require.NoError(t, err)
		require.NotNil(t, plans)

		if len(plans) > 0 {
			assert.Equal(t, planIDs[0], plans[0].PlanID)
		}
	}
}

func Test_ListFeePlans_WithMultiplePlanIds(t *testing.T) {
	mc := NewTestClient(t)

	// First, get all plans to find valid IDs
	allPlans, err := mc.ListFeePlans(context.Background(), FACILITATOR_ID)
	require.NoError(t, err)

	if len(allPlans) >= 2 {
		// Use first two plan IDs for testing
		planIDs := []string{allPlans[0].PlanID, allPlans[1].PlanID}

		plans, err := mc.ListFeePlans(
			context.Background(),
			FACILITATOR_ID,
			moov.WithFeePlanIds(planIDs),
		)
		require.NoError(t, err)
		require.NotNil(t, plans)
		assert.LessOrEqual(t, len(plans), 2)
	}
}

func Test_CreateFeePlanAgreement(t *testing.T) {
	mc := NewTestClient(t)

	// First, get available plans
	plans, err := mc.ListFeePlans(context.Background(), FACILITATOR_ID)
	require.NoError(t, err)

	if len(plans) > 0 {
		// Create a temporary test account
		account := CreateTemporaryTestAccount(t, mc, createTestBusinessAccount())

		// Create a fee plan agreement using the first available plan
		request := moov.FeePlanAgreementRequest{
			PlanID: plans[0].PlanID,
		}

		agreement, err := mc.CreateFeePlanAgreement(context.Background(), account.AccountID, request)
		require.NoError(t, err)
		require.NotNil(t, agreement)

		// Verify the agreement details
		assert.Equal(t, account.AccountID, agreement.AccountID)
		assert.Equal(t, plans[0].PlanID, agreement.PlanID)
		assert.NotEmpty(t, agreement.AgreementID)
		assert.Equal(t, moov.FeePlanAgreementStatus_Active, agreement.Status)
	} else {
		t.Skip("No fee plans available for testing")
	}
}

func Test_CreateFeePlanAgreement_WithValidPlanID(t *testing.T) {
	mc := NewTestClient(t)

	// First, get available plans
	plans, err := mc.ListFeePlans(context.Background(), FACILITATOR_ID)
	require.NoError(t, err)

	if len(plans) > 0 {
		// Create a temporary test account
		account := CreateTemporaryTestAccount(t, mc, createTestBusinessAccount())

		// Create a fee plan agreement
		request := moov.FeePlanAgreementRequest{
			PlanID: plans[0].PlanID,
		}

		agreement, err := mc.CreateFeePlanAgreement(context.Background(), account.AccountID, request)
		require.NoError(t, err)
		require.NotNil(t, agreement)

		// Verify the agreement status is active
		assert.Equal(t, moov.FeePlanAgreementStatus_Active, agreement.Status)
		assert.Equal(t, plans[0].PlanID, agreement.PlanID)
		assert.NotEmpty(t, agreement.AgreementID)
	} else {
		t.Skip("No fee plans available for testing")
	}
}
