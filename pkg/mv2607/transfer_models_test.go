package mv2607_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2607"
)

func surchargeAmount() moov.AmountDecimal {
	return moov.AmountDecimal{
		Currency:     "USD",
		ValueDecimal: "0.03",
	}
}

func TestCreateRefundAmountDetailsJSON(t *testing.T) {
	t.Run("with surcharge", func(t *testing.T) {
		surcharge := surchargeAmount()
		in := mv2607.CreateRefund{
			Amount: 97,
			AmountDetails: &mv2607.RefundAmountDetails{
				Surcharge: &surcharge,
			},
		}

		b, err := json.Marshal(in)
		require.NoError(t, err)
		require.Contains(t, string(b), `"amountDetails"`)
		require.Contains(t, string(b), `"surcharge"`)

		var out mv2607.CreateRefund
		require.NoError(t, json.Unmarshal(b, &out))
		require.NotNil(t, out.AmountDetails)
		require.NotNil(t, out.AmountDetails.Surcharge)
		require.Equal(t, surcharge, *out.AmountDetails.Surcharge)
	})

	t.Run("without surcharge", func(t *testing.T) {
		in := mv2607.CreateRefund{Amount: 97}

		b, err := json.Marshal(in)
		require.NoError(t, err)
		require.NotContains(t, string(b), `"amountDetails"`)

		var out mv2607.CreateRefund
		require.NoError(t, json.Unmarshal(b, &out))
		require.Nil(t, out.AmountDetails)
	})
}

func TestCreateReversalAmountDetailsJSON(t *testing.T) {
	t.Run("with surcharge", func(t *testing.T) {
		surcharge := surchargeAmount()
		in := mv2607.CreateReversal{
			Amount: 97,
			AmountDetails: &mv2607.RefundAmountDetails{
				Surcharge: &surcharge,
			},
		}

		b, err := json.Marshal(in)
		require.NoError(t, err)
		require.Contains(t, string(b), `"amountDetails"`)
		require.Contains(t, string(b), `"surcharge"`)

		var out mv2607.CreateReversal
		require.NoError(t, json.Unmarshal(b, &out))
		require.NotNil(t, out.AmountDetails)
		require.NotNil(t, out.AmountDetails.Surcharge)
		require.Equal(t, surcharge, *out.AmountDetails.Surcharge)
	})

	t.Run("without surcharge", func(t *testing.T) {
		in := mv2607.CreateReversal{Amount: 97}

		b, err := json.Marshal(in)
		require.NoError(t, err)
		require.NotContains(t, string(b), `"amountDetails"`)

		var out mv2607.CreateReversal
		require.NoError(t, json.Unmarshal(b, &out))
		require.Nil(t, out.AmountDetails)
	})
}

func TestRefundAmountDetailsJSON(t *testing.T) {
	var refund mv2607.Refund
	require.NoError(t, json.Unmarshal([]byte(`{
		"refundID": "refund-123",
		"status": "completed",
		"amount": {
			"currency": "USD",
			"value": 97
		},
		"amountDetails": {
			"surcharge": {
				"currency": "USD",
				"valueDecimal": "0.03"
			}
		}
	}`), &refund))

	require.NotNil(t, refund.AmountDetails)
	require.NotNil(t, refund.AmountDetails.Surcharge)
	require.Equal(t, surchargeAmount(), *refund.AmountDetails.Surcharge)
}

func TestTransferAmountDetailsSurchargeJSON(t *testing.T) {
	surcharge := surchargeAmount()
	create := mv2607.CreateTransfer{
		Source: moov.CreateTransfer_Source{
			PaymentMethodID: "source-payment-method",
		},
		Destination: moov.CreateTransfer_Destination{
			PaymentMethodID: "destination-payment-method",
		},
		Amount: moov.Amount{
			Currency: "USD",
			Value:    100,
		},
		AmountDetails: &mv2607.CreateTransferAmountDetails{
			Surcharge: &surcharge,
		},
	}

	b, err := json.Marshal(create)
	require.NoError(t, err)

	var decodedCreate mv2607.CreateTransfer
	require.NoError(t, json.Unmarshal(b, &decodedCreate))
	require.NotNil(t, decodedCreate.AmountDetails)
	require.NotNil(t, decodedCreate.AmountDetails.Surcharge)
	require.Equal(t, surcharge, *decodedCreate.AmountDetails.Surcharge)

	var transfer mv2607.Transfer
	require.NoError(t, json.Unmarshal([]byte(`{
		"transferID": "transfer-123",
		"amount": {
			"currency": "USD",
			"value": 100
		},
		"amountDetails": {
			"surcharge": {
				"currency": "USD",
				"valueDecimal": "0.03"
			}
		},
		"refunds": [{
			"refundID": "refund-123",
			"amountDetails": {
				"surcharge": {
					"currency": "USD",
					"valueDecimal": "0.03"
				}
			}
		}]
	}`), &transfer))
	require.NotNil(t, transfer.AmountDetails)
	require.NotNil(t, transfer.AmountDetails.Surcharge)
	require.Equal(t, surcharge, *transfer.AmountDetails.Surcharge)
	require.Len(t, transfer.Refunds, 1)
	require.NotNil(t, transfer.Refunds[0].AmountDetails)
	require.Equal(t, surcharge, *transfer.Refunds[0].AmountDetails.Surcharge)
}

func TestCreateTransferFeePaidByJSON(t *testing.T) {
	payout := mv2607.FeePaidBy_Destination
	create := mv2607.CreateTransfer{
		Source: moov.CreateTransfer_Source{
			PaymentMethodID: "source-payment-method",
		},
		Destination: moov.CreateTransfer_Destination{
			PaymentMethodID: "destination-payment-method",
		},
		Amount: moov.Amount{
			Currency: "USD",
			Value:    100,
		},
		FeePaidBy: &mv2607.TransferFeePaidBy{
			Payout: &payout,
		},
	}

	b, err := json.Marshal(create)
	require.NoError(t, err)
	require.JSONEq(t, `{
		"source": {"paymentMethodID": "source-payment-method"},
		"destination": {"paymentMethodID": "destination-payment-method"},
		"amount": {"currency": "USD", "value": 100},
		"facilitatorFee": {},
		"feePaidBy": {"payout": "destination"}
	}`, string(b))

	var decoded mv2607.CreateTransfer
	require.NoError(t, json.Unmarshal(b, &decoded))
	require.NotNil(t, decoded.FeePaidBy)
	require.NotNil(t, decoded.FeePaidBy.Payout)
	require.Equal(t, mv2607.FeePaidBy_Destination, *decoded.FeePaidBy.Payout)

	// feePaidBy is omitted when unset.
	b, err = json.Marshal(mv2607.CreateTransfer{
		Source:      moov.CreateTransfer_Source{PaymentMethodID: "s"},
		Destination: moov.CreateTransfer_Destination{PaymentMethodID: "d"},
		Amount:      moov.Amount{Currency: "USD", Value: 100},
	})
	require.NoError(t, err)
	require.NotContains(t, string(b), "feePaidBy")
}
