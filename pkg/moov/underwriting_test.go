package moov_test

import (
	"context"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2507"
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

	create := mv2507.UpsertUnderwriting{
		GeographicReach: moov.PtrOf(mv2507.GeographicReachUsAndInternational),
		CollectFunds: &mv2507.CollectFunds{
			CardPayments: &mv2507.CollectFundsCardPayments{
				EstimatedActivity: &mv2507.EstimatedActivity{
					MonthlyVolumeRange: moov.PtrOf(mv2507.MonthlyVolumeRangeUnder10K),
				},
			},
		},
	}

	t.Run("insert", func(t *testing.T) {
		actual, err := mv2507.Underwriting.Upsert(context.Background(), *mc, account.AccountID, create)

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
		actual, err := mv2507.Underwriting.Get(context.Background(), *mc, account.AccountID)

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

	update := mv2507.UpsertUnderwriting{
		GeographicReach:   moov.PtrOf(mv2507.GeographicReachUsOnly),
		BusinessPresence:  moov.PtrOf(mv2507.BusinessPresenceHomeBased),
		PendingLitigation: moov.PtrOf(mv2507.PendingLitigationNone),
		VolumeShareByCustomerType: &mv2507.VolumeShareByCustomerType{
			Business: moov.PtrOf(70),
			Consumer: moov.PtrOf(30),
			P2P:      moov.PtrOf(0),
		},
		CollectFunds: &mv2507.CollectFunds{
			CardPayments: &mv2507.CollectFundsCardPayments{
				EstimatedActivity: &mv2507.EstimatedActivity{MonthlyVolumeRange: moov.PtrOf(mv2507.MonthlyVolumeRange10K50K)},
			},
			Ach: &mv2507.CollectFundsAch{
				EstimatedActivity: &mv2507.EstimatedActivity{MonthlyVolumeRange: moov.PtrOf(mv2507.MonthlyVolumeRangeUnder10K)},
			},
		},
		SendFunds: &mv2507.SendFunds{
			Ach: &mv2507.SendFundsAch{
				EstimatedActivity: &mv2507.EstimatedActivity{MonthlyVolumeRange: moov.PtrOf(mv2507.MonthlyVolumeRangeUnder10K)},
			},
		},
		MoneyTransfer: &mv2507.MoneyTransfer{
			PullFromCard: &mv2507.MoneyTransferPullFromCard{
				EstimatedActivity: &mv2507.EstimatedActivity{MonthlyVolumeRange: moov.PtrOf(mv2507.MonthlyVolumeRange10K50K)},
			},
		},
	}

	t.Run("update", func(t *testing.T) {
		actual, err := mv2507.Underwriting.Upsert(context.Background(), *mc, account.AccountID, update)

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
		actual, err := mv2507.Underwriting.Get(context.Background(), *mc, account.AccountID)

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

func TestUpsertUnderwriting_V2507(t *testing.T) {
	mc := NewTestClient(t)

	account := CreateTemporaryTestAccount(t, mc, createTestBusinessAccount())

	create := mv2507.UpsertUnderwriting{
		GeographicReach: moov.PtrOf(mv2507.GeographicReachUsAndInternational),
		CollectFunds: &mv2507.CollectFunds{
			CardPayments: &mv2507.CollectFundsCardPayments{
				EstimatedActivity: &mv2507.EstimatedActivity{
					MonthlyVolumeRange: moov.PtrOf(mv2507.MonthlyVolumeRangeUnder10K),
				},
			},
		},
		SubmissionIntent: moov.PtrOf(mv2507.SubmissionIntentWait),
	}

	t.Run("insert", func(t *testing.T) {
		actual, err := mv2507.Underwriting.Upsert(context.Background(), *mc, account.AccountID, create)

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
		actual, err := mv2507.Underwriting.Get(context.Background(), *mc, account.AccountID)

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

	update := mv2507.UpsertUnderwriting{
		GeographicReach:   moov.PtrOf(mv2507.GeographicReachUsOnly),
		BusinessPresence:  moov.PtrOf(mv2507.BusinessPresenceHomeBased),
		PendingLitigation: moov.PtrOf(mv2507.PendingLitigationNone),
		VolumeShareByCustomerType: &mv2507.VolumeShareByCustomerType{
			Business: moov.PtrOf(70),
			Consumer: moov.PtrOf(30),
			P2P:      moov.PtrOf(0),
		},
		CollectFunds: &mv2507.CollectFunds{
			CardPayments: &mv2507.CollectFundsCardPayments{
				EstimatedActivity: &mv2507.EstimatedActivity{MonthlyVolumeRange: moov.PtrOf(mv2507.MonthlyVolumeRange10K50K)},
			},
			Ach: &mv2507.CollectFundsAch{
				EstimatedActivity: &mv2507.EstimatedActivity{MonthlyVolumeRange: moov.PtrOf(mv2507.MonthlyVolumeRangeUnder10K)},
			},
		},
		SendFunds: &mv2507.SendFunds{
			Ach: &mv2507.SendFundsAch{
				EstimatedActivity: &mv2507.EstimatedActivity{MonthlyVolumeRange: moov.PtrOf(mv2507.MonthlyVolumeRangeUnder10K)},
			},
		},
		MoneyTransfer: &mv2507.MoneyTransfer{
			PullFromCard: &mv2507.MoneyTransferPullFromCard{
				EstimatedActivity: &mv2507.EstimatedActivity{MonthlyVolumeRange: moov.PtrOf(mv2507.MonthlyVolumeRange10K50K)},
			},
		},
		SubmissionIntent: moov.PtrOf(mv2507.SubmissionIntentSubmit),
	}

	t.Run("update", func(t *testing.T) {
		actual, err := mv2507.Underwriting.Upsert(context.Background(), *mc, account.AccountID, update)

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
		actual, err := mv2507.Underwriting.Get(context.Background(), *mc, account.AccountID)

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
