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
		VolumeByCustomerType: moov.VolumeByCustomerType{
			BusinessToBusinessPercentage: 50,
			ConsumerToBusinessPercentage: 50,
		},
		CardVolumeDistribution: moov.CardVolumeDistribution{
			EcommercePercentage:     50,
			CardPresentPercentage:   50,
			MailOrPhonePercentage:   0,
			DebtRepaymentPercentage: 0,
		},
		Fulfillment: moov.Fulfillment{
			HasPhysicalGoods:     true,
			IsShippingProduct:    true,
			ShipmentDurationDays: 3,
			ReturnPolicy:         moov.WITHIN_THIRTY_DAYS,
		},
	}

	t.Run("insert", func(t *testing.T) {
		actual, err := mc.UpsertUnderwriting(context.Background(), account.AccountID, create)

		NoResponseError(t, err)
		require.NotNil(t, actual)
		require.Equal(t, create.AverageTransactionSize, actual.AverageTransactionSize)
		require.Equal(t, create.MaxTransactionSize, actual.MaxTransactionSize)
		require.Equal(t, create.AverageMonthlyTransactionVolume, actual.AverageMonthlyTransactionVolume)
		require.Equal(t, create.Fulfillment, actual.Fulfillment)
		require.Equal(t, create.VolumeByCustomerType, actual.VolumeByCustomerType)
		require.Equal(t, create.CardVolumeDistribution, actual.CardVolumeDistribution)
		require.Equal(t, moov.UnderwritingStatusNotRequested, actual.Status)
	})

	t.Run("get", func(t *testing.T) {
		actual, err := mc.GetUnderwriting(context.Background(), account.AccountID)

		NoResponseError(t, err)
		require.NotNil(t, actual)
		require.Equal(t, create.AverageTransactionSize, actual.AverageTransactionSize)
		require.Equal(t, create.MaxTransactionSize, actual.MaxTransactionSize)
		require.Equal(t, create.AverageMonthlyTransactionVolume, actual.AverageMonthlyTransactionVolume)
		require.Equal(t, create.Fulfillment, actual.Fulfillment)
		require.Equal(t, create.VolumeByCustomerType, actual.VolumeByCustomerType)
		require.Equal(t, create.CardVolumeDistribution, actual.CardVolumeDistribution)
		require.Equal(t, moov.UnderwritingStatusNotRequested, actual.Status)
	})

	update := moov.UpdateUnderwriting{
		AverageTransactionSize:          1_500,
		MaxTransactionSize:              1_200,
		AverageMonthlyTransactionVolume: 11_000,
		VolumeByCustomerType: moov.VolumeByCustomerType{
			BusinessToBusinessPercentage: 60,
			ConsumerToBusinessPercentage: 40,
		},
		CardVolumeDistribution: moov.CardVolumeDistribution{
			EcommercePercentage:     60,
			CardPresentPercentage:   40,
			MailOrPhonePercentage:   0,
			DebtRepaymentPercentage: 0,
		},
		Fulfillment: moov.Fulfillment{
			HasPhysicalGoods:     true,
			IsShippingProduct:    true,
			ShipmentDurationDays: 3,
			ReturnPolicy:         moov.EXCHANGE_ONLY,
		},
	}

	t.Run("update", func(t *testing.T) {
		actual, err := mc.UpsertUnderwriting(context.Background(), account.AccountID, update)

		NoResponseError(t, err)
		require.NotNil(t, actual)
		require.Equal(t, update.AverageTransactionSize, actual.AverageTransactionSize)
		require.Equal(t, update.MaxTransactionSize, actual.MaxTransactionSize)
		require.Equal(t, update.AverageMonthlyTransactionVolume, actual.AverageMonthlyTransactionVolume)
		require.Equal(t, update.Fulfillment, actual.Fulfillment)
		require.Equal(t, update.VolumeByCustomerType, actual.VolumeByCustomerType)
		require.Equal(t, update.CardVolumeDistribution, actual.CardVolumeDistribution)
		require.Equal(t, moov.UnderwritingStatusNotRequested, actual.Status)
	})

	t.Run("get", func(t *testing.T) {
		actual, err := mc.GetUnderwriting(context.Background(), account.AccountID)

		NoResponseError(t, err)
		require.NotNil(t, actual)
		require.Equal(t, update.AverageTransactionSize, actual.AverageTransactionSize)
		require.Equal(t, update.MaxTransactionSize, actual.MaxTransactionSize)
		require.Equal(t, update.AverageMonthlyTransactionVolume, actual.AverageMonthlyTransactionVolume)
		require.Equal(t, update.Fulfillment, actual.Fulfillment)
		require.Equal(t, update.VolumeByCustomerType, actual.VolumeByCustomerType)
		require.Equal(t, update.CardVolumeDistribution, actual.CardVolumeDistribution)
		require.Equal(t, moov.UnderwritingStatusNotRequested, actual.Status)
	})
}

func TestUpsertUnderwritingV2(t *testing.T) {
	mc := NewTestClient(t)

	account := CreateTemporaryTestAccount(t, mc, createTestBusinessAccount())

	create := moov.UpsertUnderwriting{
		GeographicReach: func() *moov.GeographicReach {
			value := moov.GeographicReachUsAndInternational
			return &value
		}(),
		CollectFunds: &moov.CollectFunds{
			CardPayments: &moov.CollectFundsCardPayments{
				EstimatedActivity: &moov.EstimatedActivity{
					MonthlyVolumeRange: func() *moov.MonthlyVolumeRange {
						value := moov.MonthlyVolumeRangeUnder10K
						return &value
					}(),
				},
			},
		},
	}

	t.Run("insert", func(t *testing.T) {
		actual, err := mc.UpsertUnderwritingV2(context.Background(), "v2025.07.00", account.AccountID, create)

		NoResponseError(t, err)
		require.NotNil(t, actual)
		require.Equal(t, create.GeographicReach, actual.GeographicReach)
		require.Equal(t, moov.UnderwritingStatusPending, actual.Status)
	})
}
