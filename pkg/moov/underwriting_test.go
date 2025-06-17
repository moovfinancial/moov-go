package moov_test

import (
	"context"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/Q3_2025"
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

func TestUpsertUnderwriting_V2507(t *testing.T) {
	mc := NewTestClient(t)

	account := CreateTemporaryTestAccount(t, mc, createTestBusinessAccount())

	create := Q3_2025.UpsertUnderwriting{
		GeographicReach: func() *Q3_2025.GeographicReach {
			value := Q3_2025.GeographicReachUsAndInternational
			return &value
		}(),
		CollectFunds: &Q3_2025.CollectFunds{
			CardPayments: &Q3_2025.CollectFundsCardPayments{
				EstimatedActivity: &Q3_2025.EstimatedActivity{
					MonthlyVolumeRange: func() *Q3_2025.MonthlyVolumeRange {
						value := Q3_2025.MonthlyVolumeRangeUnder10K
						return &value
					}(),
				},
			},
		},
	}

	t.Run("insert", func(t *testing.T) {
		actual, err := Q3_2025.Underwriting.Upsert(context.Background(), *mc, account.AccountID, create)

		NoResponseError(t, err)
		require.NotNil(t, actual)
		require.Equal(t, create.GeographicReach, actual.GeographicReach)
		require.Equal(t, create.CollectFunds, actual.CollectFunds)
		require.Nil(t, actual.BusinessPresence)
		require.Nil(t, actual.PendingLitigation)
		require.Nil(t, actual.VolumeShareByCustomerType)
		require.Nil(t, actual.SendFunds)
		require.Nil(t, actual.MoneyTransfer)
	})

	t.Run("get", func(t *testing.T) {
		actual, err := Q3_2025.Underwriting.Get(context.Background(), *mc, account.AccountID)

		NoResponseError(t, err)
		require.NotNil(t, actual)
		require.Equal(t, create.GeographicReach, actual.GeographicReach)
		require.Equal(t, create.CollectFunds, actual.CollectFunds)
		require.Nil(t, actual.BusinessPresence)
		require.Nil(t, actual.PendingLitigation)
		require.Nil(t, actual.VolumeShareByCustomerType)
		require.Nil(t, actual.SendFunds)
		require.Nil(t, actual.MoneyTransfer)
	})

	update := Q3_2025.UpsertUnderwriting{
		GeographicReach:   func() *Q3_2025.GeographicReach { value := Q3_2025.GeographicReachUsOnly; return &value }(),
		BusinessPresence:  func() *Q3_2025.BusinessPresence { value := Q3_2025.BusinessPresenceHomeBased; return &value }(),
		PendingLitigation: func() *Q3_2025.PendingLitigation { value := Q3_2025.PendingLitigationNone; return &value }(),
		VolumeShareByCustomerType: &Q3_2025.VolumeShareByCustomerType{
			Business: func() *int { value := 70; return &value }(),
			Consumer: func() *int { value := 30; return &value }(),
			P2P:      func() *int { value := 0; return &value }(),
		},
		CollectFunds: &Q3_2025.CollectFunds{
			CardPayments: &Q3_2025.CollectFundsCardPayments{
				EstimatedActivity: &Q3_2025.EstimatedActivity{MonthlyVolumeRange: func() *Q3_2025.MonthlyVolumeRange { value := Q3_2025.MonthlyVolumeRange10K50K; return &value }()},
			},
			Ach: &Q3_2025.CollectFundsAch{
				EstimatedActivity: &Q3_2025.EstimatedActivity{MonthlyVolumeRange: func() *Q3_2025.MonthlyVolumeRange { value := Q3_2025.MonthlyVolumeRangeUnder10K; return &value }()},
			},
		},
		SendFunds: &Q3_2025.SendFunds{
			Ach: &Q3_2025.SendFundsAch{
				EstimatedActivity: &Q3_2025.EstimatedActivity{MonthlyVolumeRange: func() *Q3_2025.MonthlyVolumeRange { value := Q3_2025.MonthlyVolumeRangeUnder10K; return &value }()},
			},
		},
		MoneyTransfer: &Q3_2025.MoneyTransfer{
			PullFromCard: &Q3_2025.MoneyTransferPullFromCard{
				EstimatedActivity: &Q3_2025.EstimatedActivity{MonthlyVolumeRange: func() *Q3_2025.MonthlyVolumeRange { value := Q3_2025.MonthlyVolumeRange10K50K; return &value }()},
			},
		},
	}

	t.Run("update", func(t *testing.T) {
		actual, err := Q3_2025.Underwriting.Upsert(context.Background(), *mc, account.AccountID, update)

		NoResponseError(t, err)
		require.NotNil(t, actual)
		require.Equal(t, update.GeographicReach, actual.GeographicReach)
		require.Equal(t, update.BusinessPresence, actual.BusinessPresence)
		require.Equal(t, update.PendingLitigation, actual.PendingLitigation)
		require.Equal(t, update.VolumeShareByCustomerType, actual.VolumeShareByCustomerType)
		require.Equal(t, update.CollectFunds, actual.CollectFunds)
		require.Equal(t, update.SendFunds, actual.SendFunds)
		require.Equal(t, update.MoneyTransfer, actual.MoneyTransfer)
	})

	t.Run("get after update", func(t *testing.T) {
		actual, err := Q3_2025.Underwriting.Get(context.Background(), *mc, account.AccountID)

		NoResponseError(t, err)
		require.NotNil(t, actual)
		require.Equal(t, update.GeographicReach, actual.GeographicReach)
		require.Equal(t, update.BusinessPresence, actual.BusinessPresence)
		require.Equal(t, update.PendingLitigation, actual.PendingLitigation)
		require.Equal(t, update.VolumeShareByCustomerType, actual.VolumeShareByCustomerType)
		require.Equal(t, update.CollectFunds, actual.CollectFunds)
		require.Equal(t, update.SendFunds, actual.SendFunds)
		require.Equal(t, update.MoneyTransfer, actual.MoneyTransfer)
	})
}
