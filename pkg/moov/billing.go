package moov

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// FeePlanAgreement represents a billing fee plan agreement for a Moov account
type FeePlanAgreement struct {
	AgreementID        string         `json:"agreementID,omitempty"`
	PlanID             string         `json:"planID,omitempty"`
	AccountID          string         `json:"accountID,omitempty"`
	Name               string         `json:"name,omitempty"`
	Description        string         `json:"description,omitempty"`
	AcceptedOn         time.Time      `json:"acceptedOn,omitempty"`
	Status             string         `json:"status,omitempty"`
	CardAcquiringModel string         `json:"cardAcquiringModel,omitempty"`
	BillableFees       []BillableFee  `json:"billableFees,omitempty"`
	MinimumCommitment  *AmountDecimal `json:"minimumCommitment,omitempty"`
	MonthlyPlatformFee *AmountDecimal `json:"monthlyPlatformFee,omitempty"`
}

// BillableFee represents a billable fee within a fee plan agreement
type BillableFee struct {
	BillableFeeID string         `json:"billableFeeID,omitempty"`
	BillableEvent string         `json:"billableEvent,omitempty"`
	FeeName       string         `json:"feeName,omitempty"`
	FeeModel      string         `json:"feeModel,omitempty"`
	FeeCategory   string         `json:"feeCategory,omitempty"`
	FeeProperties *FeeProperties `json:"feeProperties,omitempty"`
	FeeConditions *FeeConditions `json:"feeConditions,omitempty"`
}

// FeeProperties represents the properties of a fee, including amounts and rates
type FeeProperties struct {
	FixedAmount       *AmountDecimal `json:"fixedAmount,omitempty"`
	VariableRate      string         `json:"variableRate,omitempty"`
	MinPerTransaction *AmountDecimal `json:"minPerTransaction,omitempty"`
	MaxPerTransaction *AmountDecimal `json:"maxPerTransaction,omitempty"`
	VolumeRanges      []VolumeRange  `json:"volumeRanges,omitempty"`
}

// VolumeRange represents a volume-based pricing range
type VolumeRange struct {
	FromValue     int            `json:"fromValue,omitempty"`
	ToValue       int            `json:"toValue,omitempty"`
	FlatAmount    *AmountDecimal `json:"flatAmount,omitempty"`
	PerUnitAmount *AmountDecimal `json:"perUnitAmount,omitempty"`
}

// FeeConditions represents conditions under which a fee applies
type FeeConditions struct {
	TransactionType []string `json:"transactionType,omitempty"`
}

type FeePlanAgreementListFilter callArg

func WithFeePlanAgreementCount(c int) FeePlanAgreementListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["count"] = fmt.Sprintf("%d", c)
		return nil
	})
}

func WithFeePlanAgreementSkip(c int) FeePlanAgreementListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["skip"] = fmt.Sprintf("%d", c)
		return nil
	})
}

func WithFeePlanAgreementStatuses(statuses []string) FeePlanAgreementListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["status"] = strings.Join(statuses, ",")
		return nil
	})
}

func WithFeePlanAgreementIds(agreementIds []string) FeePlanAgreementListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["agreementID"] = strings.Join(agreementIds, ",")
		return nil
	})
}

// ListFeePlanAgreements lists of FeePlanAgreements that are associated with a Moov account
// https://docs.moov.io/api/moov-accounts/billing/list-agreements/
func (c Client) ListFeePlanAgreements(ctx context.Context, accountID string, filters ...FeePlanAgreementListFilter) ([]FeePlanAgreement, error) {
	args := prependArgs(filters, AcceptJson())
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathFeePlanAgreements, accountID), args...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[FeePlanAgreement](resp)
}
