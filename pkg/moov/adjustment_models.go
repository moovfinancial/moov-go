package moov

import (
	"time"
)

type Adjustment struct {
	AdjustmentID string        `json:"adjustmentID,omitempty"`
	WalletID     string        `json:"walletID,omitempty"`
	CreatedOn    time.Time     `json:"createdOn,omitempty"`
	Amount       AmountDecimal `json:"amount,omitempty"`
}

type AdjustmentListFilter callArg

func WithAdjustmentWalletID(walletID string) AdjustmentListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["walletID"] = walletID
		return nil
	})
}
