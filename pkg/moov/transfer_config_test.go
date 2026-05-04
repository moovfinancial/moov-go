package moov_test

import (
	"encoding/json"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

func Test_TransferConfig(t *testing.T) {
	t.Run("create get update", func(t *testing.T) {
		mc := NewTestClient(t)
		account := CreateTemporaryTestAccount(t, mc, createTestIndividualAccount())

		create := moov.CreateTransferConfig{
			TipPresets: &moov.CreateTipPresets{
				CalculationBasis:  moov.PtrOf(moov.TipCalculationBasis_PreTax),
				PercentageOptions: []int{10, 15, 20},
			},
		}

		created, err := mc.CreateTransferConfig(BgCtx(), account.AccountID, create)
		NoResponseError(t, err)
		require.NotNil(t, created)
		require.NotNil(t, created.TipPresets)
		require.Equal(t, moov.PtrOf(moov.TipCalculationBasis_PreTax), created.TipPresets.CalculationBasis)
		require.Equal(t, []int{10, 15, 20}, created.TipPresets.PercentageOptions)

		got, err := mc.GetTransferConfig(BgCtx(), account.AccountID)
		NoResponseError(t, err)
		require.NotNil(t, got)
		require.Equal(t, created, got)

		updateFixed := moov.PutTransferConfig{
			TipPresets: moov.PutTipPresets{
				FixedAmountOptions: []moov.AmountDecimal{
					{Currency: "USD", ValueDecimal: "1.00"},
					{Currency: "USD", ValueDecimal: "2.00"},
					{Currency: "USD", ValueDecimal: "5.00"},
				},
			},
		}

		updatedFixed, err := mc.UpdateTransferConfig(BgCtx(), account.AccountID, updateFixed)
		NoResponseError(t, err)
		require.NotNil(t, updatedFixed)
		require.NotNil(t, updatedFixed.TipPresets)
		require.Equal(t, updateFixed.TipPresets.FixedAmountOptions, updatedFixed.TipPresets.FixedAmountOptions)

		updatePercentages := moov.PutTransferConfig{
			TipPresets: moov.PutTipPresets{
				CalculationBasis:  moov.PtrOf(moov.TipCalculationBasis_PostTax),
				PercentageOptions: []int{12, 18, 25},
			},
		}

		updatedPercentages, err := mc.UpdateTransferConfig(BgCtx(), account.AccountID, updatePercentages)
		NoResponseError(t, err)
		require.NotNil(t, updatedPercentages)
		require.NotNil(t, updatedPercentages.TipPresets)
		require.Equal(t, moov.PtrOf(moov.TipCalculationBasis_PostTax), updatedPercentages.TipPresets.CalculationBasis)
		require.Equal(t, []int{12, 18, 25}, updatedPercentages.TipPresets.PercentageOptions)
	})

	t.Run("validation error json", func(t *testing.T) {
		payload := []byte(`{
			"TipPresets.CalculationBasis": "must be a valid value",
			"TipPresets.PercentageOptions": {
				"0": "must be between 0 and 100"
			},
			"TipPresets.FixedAmountOptions": {
				"0": {
					"currency": "must be present",
					"valueDecimal": "must be positive"
				}
			}
		}`)

		var actual moov.TransferConfigValidationError
		err := json.Unmarshal(payload, &actual)
		require.NoError(t, err)
		require.NotNil(t, actual.TipPresetsCalculationBasis)
		require.Equal(t, "must be a valid value", *actual.TipPresetsCalculationBasis)
		require.Equal(t, map[string]string{"0": "must be between 0 and 100"}, actual.TipPresetsPercentageOptions)
		require.Contains(t, actual.TipPresetsFixedAmountOptions, "0")
		require.Equal(t, "must be present", *actual.TipPresetsFixedAmountOptions["0"].Currency)
		require.Equal(t, "must be positive", *actual.TipPresetsFixedAmountOptions["0"].ValueDecimal)
	})
}
