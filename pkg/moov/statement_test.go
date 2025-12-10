package moov_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

func Test_Statements(t *testing.T) {
	c := NewTestClient(t)

	statements, err := c.ListStatements(t.Context(), FACILITATOR_ID)
	require.NoError(t, err)
	require.NotNil(t, statements)
	if len(statements) > 0 {
		// as JSON
		statement, err := c.GetStatement(t.Context(), FACILITATOR_ID, statements[0].StatementID)
		require.NoError(t, err)
		require.NotNil(t, statement)
		require.Equal(t, statement.StatementID, statements[0].StatementID)

		// as PDF
		pdfData, err := c.GetStatementPDF(t.Context(), FACILITATOR_ID, statements[0].StatementID)
		require.NoError(t, err)
		require.NotEmpty(t, pdfData)
	}
}

func Test_ListStatements_WithCount(t *testing.T) {
	mc := NewTestClient(t)

	statements, err := mc.ListStatements(
		t.Context(),
		FACILITATOR_ID,
		moov.WithStatementCount(5),
	)
	require.NoError(t, err)
	require.NotNil(t, statements)
	if len(statements) > 0 {
		require.LessOrEqual(t, len(statements), 5)
	}
}

func Test_ListStatements_WithSkip(t *testing.T) {
	mc := NewTestClient(t)

	statements, err := mc.ListStatements(
		t.Context(),
		FACILITATOR_ID,
		moov.WithStatementSkip(2),
	)
	require.NoError(t, err)
	require.NotNil(t, statements)
}

func Test_ListStatements_WithBillingPeriodStartDateTime(t *testing.T) {
	mc := NewTestClient(t)

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	statements, err := mc.ListStatements(
		t.Context(),
		FACILITATOR_ID,
		moov.WithBillingPeriodStartDateTime(startDate),
	)
	require.NoError(t, err)
	require.NotNil(t, statements)

	// Verify all returned statements are after or on the start date
	for _, stmt := range statements {
		if !stmt.BillingPeriodStartDateTime.IsZero() {
			require.True(t, stmt.BillingPeriodStartDateTime.Equal(startDate) || stmt.BillingPeriodStartDateTime.After(startDate),
				"Statement billing period start should be on or after filter start date")
		}
	}
}

func Test_ListStatements_WithBillingPeriodEndDateTime(t *testing.T) {
	mc := NewTestClient(t)

	endDate := time.Now().UTC()

	statements, err := mc.ListStatements(
		t.Context(),
		FACILITATOR_ID,
		moov.WithBillingPeriodEndDateTime(endDate),
	)
	require.NoError(t, err)
	require.NotNil(t, statements)

	// Verify all returned statements are before or on the end date
	for _, stmt := range statements {
		if !stmt.BillingPeriodEndDateTime.IsZero() {
			require.True(t, stmt.BillingPeriodEndDateTime.Equal(endDate) || stmt.BillingPeriodEndDateTime.Before(endDate),
				"Statement billing period end should be on or before filter end date")
		}
	}
}

func Test_ListStatements_WithDateRange(t *testing.T) {
	mc := NewTestClient(t)

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Now().UTC()

	statements, err := mc.ListStatements(
		t.Context(),
		FACILITATOR_ID,
		moov.WithBillingPeriodStartDateTime(startDate),
		moov.WithBillingPeriodEndDateTime(endDate),
	)
	require.NoError(t, err)
	require.NotNil(t, statements)

	// Verify all returned statements fall within the date range
	for _, stmt := range statements {
		if !stmt.BillingPeriodStartDateTime.IsZero() {
			require.True(t, stmt.BillingPeriodStartDateTime.Equal(startDate) || stmt.BillingPeriodStartDateTime.After(startDate),
				"Statement billing period start should be on or after filter start date")
		}
		if !stmt.BillingPeriodEndDateTime.IsZero() {
			require.True(t, stmt.BillingPeriodEndDateTime.Equal(endDate) || stmt.BillingPeriodEndDateTime.Before(endDate),
				"Statement billing period end should be on or before filter end date")
		}
	}
}
