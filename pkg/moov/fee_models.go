package moov

import (
	"fmt"
	"time"
)

type IncurredFee struct {
	FeeID       string        `json:"feeID,omitempty"`
	AccountID   string        `json:"accountID,omitempty"`
	WalletID    string        `json:"walletID,omitempty"`
	CreatedOn   time.Time     `json:"createdOn,omitempty"`
	FeeName     string        `json:"feeName,omitempty"`
	Amount      AmountDecimal `json:"amount,omitempty"`
	GeneratedBy *GeneratedBy  `json:"generatedBy,omitempty"`
	FeeGroup    string        `json:"feeGroup,omitempty"`
	ResidualID  *string       `json:"residualID,omitempty"`
}

type GeneratedBy struct {
	TransferID    *string `json:"transferID,omitempty"`
	CardID        *string `json:"cardID,omitempty"`
	DisputeID     *string `json:"disputeID,omitempty"`
	AccountID     *string `json:"accountID,omitempty"`
	BankAccountID *string `json:"bankAccountID,omitempty"`
}

func (g *GeneratedBy) GetTransferID() *string {
	if g == nil {
		return nil
	}
	return g.TransferID
}

func (g *GeneratedBy) GetCardID() *string {
	if g == nil {
		return nil
	}
	return g.CardID
}

func (g *GeneratedBy) GetDisputeID() *string {
	if g == nil {
		return nil
	}
	return g.DisputeID
}

func (g *GeneratedBy) GetAccountID() *string {
	if g == nil {
		return nil
	}
	return g.AccountID
}

func (g *GeneratedBy) GetBankAccountID() *string {
	if g == nil {
		return nil
	}
	return g.BankAccountID
}

type FeeGetFilter callArg

func WithFeeCount(c int) FeeGetFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["count"] = fmt.Sprintf("%d", c)
		return nil
	})
}

func WithFeeSkip(c int) FeeGetFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["skip"] = fmt.Sprintf("%d", c)
		return nil
	})
}

func WithFeeTransferID(transferID string) FeeGetFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["transferID"] = transferID
		return nil
	})
}

func WithFeeDisputeID(disputeID string) FeeGetFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["disputeID"] = disputeID
		return nil
	})
}

func WithFeeResidualID(residualID string) FeeGetFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["residualID"] = residualID
		return nil
	})
}

func WithFeeStartDateTime(start string) FeeGetFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["startDateTime"] = start
		return nil
	})
}

func WithFeeEndDateTime(end string) FeeGetFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["endDateTime"] = end
		return nil
	})
}

type FeeListRequest struct {
	FeeIDs []string `json:"feeIDs"`
}
