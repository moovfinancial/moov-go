package moov_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
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

func Test_ListResiduals(t *testing.T) {
	mc := NewTestClient(t)

	residuals, err := mc.ListResiduals(context.Background(), FACILITATOR_ID)
	require.NoError(t, err)
	require.NotNil(t, residuals)
}

func Test_ListResiduals_WithCount(t *testing.T) {
	mc := NewTestClient(t)

	residuals, err := mc.ListResiduals(
		context.Background(),
		FACILITATOR_ID,
		moov.WithResidualCount(5),
	)
	require.NoError(t, err)
	require.NotNil(t, residuals)
	if len(residuals) > 0 {
		assert.LessOrEqual(t, len(residuals), 5)
	}
}

func Test_ListResiduals_WithSkip(t *testing.T) {
	mc := NewTestClient(t)

	residuals, err := mc.ListResiduals(
		context.Background(),
		FACILITATOR_ID,
		moov.WithResidualSkip(2),
	)
	require.NoError(t, err)
	require.NotNil(t, residuals)
}

func Test_ListResiduals_WithDateTimeRange(t *testing.T) {
	mc := NewTestClient(t)

	residuals, err := mc.ListResiduals(
		context.Background(),
		FACILITATOR_ID,
		moov.WithResidualStartDateTime("2024-01-01T00:00:00Z"),
		moov.WithResidualEndDateTime("2099-12-31T23:59:59Z"),
	)
	require.NoError(t, err)
	require.NotNil(t, residuals)
}

func Test_GetResidual(t *testing.T) {
	mc := NewTestClient(t)

	// First, list residuals to find a valid ID
	residuals, err := mc.ListResiduals(context.Background(), FACILITATOR_ID)
	require.NoError(t, err)

	if len(residuals) > 0 {
		residual, err := mc.GetResidual(context.Background(), FACILITATOR_ID, residuals[0].ResidualID)
		require.NoError(t, err)
		require.NotNil(t, residual)
		assert.Equal(t, residuals[0].ResidualID, residual.ResidualID)
	}
}

func Test_ListResidualFees(t *testing.T) {
	mc := NewTestClient(t)

	// First, list residuals to find a valid ID
	residuals, err := mc.ListResiduals(context.Background(), FACILITATOR_ID)
	require.NoError(t, err)

	if len(residuals) > 0 {
		fees, err := mc.ListResidualFees(context.Background(), FACILITATOR_ID, residuals[0].ResidualID)
		require.NoError(t, err)
		require.NotNil(t, fees)
	}
}

func Test_ListResidualFees_WithFilters(t *testing.T) {
	mc := NewTestClient(t)

	// First, list residuals to find a valid ID
	residuals, err := mc.ListResiduals(context.Background(), FACILITATOR_ID)
	require.NoError(t, err)

	if len(residuals) > 0 {
		fees, err := mc.ListResidualFees(
			context.Background(),
			FACILITATOR_ID,
			residuals[0].ResidualID,
			moov.WithResidualFeeCount(10),
			moov.WithResidualFeeSkip(0),
		)
		require.NoError(t, err)
		require.NotNil(t, fees)
	}
}

func Test_ListPartnerPricingAgreements(t *testing.T) {
	mc := NewTestClient(t)

	agreements, err := mc.ListPartnerPricingAgreements(context.Background(), FACILITATOR_ID)
	require.NoError(t, err)
	require.NotNil(t, agreements)
}

func Test_ListPartnerPricingAgreements_WithCount(t *testing.T) {
	mc := NewTestClient(t)

	agreements, err := mc.ListPartnerPricingAgreements(
		context.Background(),
		FACILITATOR_ID,
		moov.WithPartnerPricingAgreementCount(5),
	)
	require.NoError(t, err)
	require.NotNil(t, agreements)
	if len(agreements) > 0 {
		assert.LessOrEqual(t, len(agreements), 5)
	}
}

func Test_ListPartnerPricingAgreements_WithSkip(t *testing.T) {
	mc := NewTestClient(t)

	agreements, err := mc.ListPartnerPricingAgreements(
		context.Background(),
		FACILITATOR_ID,
		moov.WithPartnerPricingAgreementSkip(2),
	)
	require.NoError(t, err)
	require.NotNil(t, agreements)
}

func Test_ListPartnerPricingAgreements_WithStatuses(t *testing.T) {
	mc := NewTestClient(t)

	agreements, err := mc.ListPartnerPricingAgreements(
		context.Background(),
		FACILITATOR_ID,
		moov.WithPartnerPricingAgreementStatuses([]moov.FeePlanAgreementStatus{moov.FeePlanAgreementStatus_Active}),
	)
	require.NoError(t, err)
	require.NotNil(t, agreements)
}

func Test_ListPartnerPricingAgreements_WithAgreementIds(t *testing.T) {
	mc := NewTestClient(t)

	// First, get all agreements to find valid IDs
	allAgreements, err := mc.ListPartnerPricingAgreements(context.Background(), FACILITATOR_ID)
	require.NoError(t, err)

	if len(allAgreements) > 0 {
		// Use the first agreement ID for testing
		agreementIDs := []string{allAgreements[0].AgreementID}

		agreements, err := mc.ListPartnerPricingAgreements(
			context.Background(),
			FACILITATOR_ID,
			moov.WithPartnerPricingAgreementIds(agreementIDs),
		)
		require.NoError(t, err)
		require.NotNil(t, agreements)

		if len(agreements) > 0 {
			assert.Equal(t, agreementIDs[0], agreements[0].AgreementID)
		}
	}
}
