package moov_test

import (
	"context"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

func TestUpsertUnderwriting(t *testing.T) {
	mc := NewTestClient(t)

	account := CreateTemporaryTestAccount(t, mc, createTestBusinessAccount())

	create := moov.UpdateUnderwriting{
		AverageTransactionSize:          1_000,
		MaxTransactionSize:              500,
		AverageMonthlyTransactionVolume: 10_000,
	}

	t.Run("insert", func(t *testing.T) {
		actual, err := mc.UpsertUnderwriting(context.Background(), account.AccountID, create)

		NoResponseError(t, err)
		require.Equal(t, create.AverageTransactionSize, actual.AverageTransactionSize)
		require.Equal(t, create.MaxTransactionSize, actual.MaxTransactionSize)
		require.Equal(t, create.AverageMonthlyTransactionVolume, actual.AverageMonthlyTransactionVolume)
		require.Equal(t, moov.UnderwritingStatusNotRequested, actual.Status)
	})

	t.Run("get", func(t *testing.T) {
		actual, err := mc.GetUnderwriting(context.Background(), account.AccountID)

		NoResponseError(t, err)
		require.Equal(t, create.AverageTransactionSize, actual.AverageTransactionSize)
		require.Equal(t, create.MaxTransactionSize, actual.MaxTransactionSize)
		require.Equal(t, create.AverageMonthlyTransactionVolume, actual.AverageMonthlyTransactionVolume)
		require.Equal(t, moov.UnderwritingStatusNotRequested, actual.Status)
	})

	update := moov.UpdateUnderwriting{
		AverageTransactionSize:          1_500,
		MaxTransactionSize:              1_200,
		AverageMonthlyTransactionVolume: 11_000,
	}

	t.Run("update", func(t *testing.T) {
		actual, err := mc.UpsertUnderwriting(context.Background(), account.AccountID, update)

		NoResponseError(t, err)
		require.Equal(t, update.AverageTransactionSize, actual.AverageTransactionSize)
		require.Equal(t, update.MaxTransactionSize, actual.MaxTransactionSize)
		require.Equal(t, update.AverageMonthlyTransactionVolume, actual.AverageMonthlyTransactionVolume)
		require.Equal(t, moov.UnderwritingStatusNotRequested, actual.Status)
	})

	t.Run("get", func(t *testing.T) {
		actual, err := mc.GetUnderwriting(context.Background(), account.AccountID)

		NoResponseError(t, err)
		require.Equal(t, update.AverageTransactionSize, actual.AverageTransactionSize)
		require.Equal(t, update.MaxTransactionSize, actual.MaxTransactionSize)
		require.Equal(t, update.AverageMonthlyTransactionVolume, actual.AverageMonthlyTransactionVolume)
		require.Equal(t, moov.UnderwritingStatusNotRequested, actual.Status)
	})
}
