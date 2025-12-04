package moov

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type FeePlan struct {
	BillableFees       []BillableFee `json:"billableFees"`
	CardAcquiringModel string        `json:"cardAcquiringModel"`
	CreatedAt          time.Time     `json:"createdAt"`
	Description        string        `json:"description"`
	MinimumCommitment  AmountDecimal `json:"minimumCommitment"`
	MonthlyPlatformFee AmountDecimal `json:"monthlyPlatformFee"`
	Name               string        `json:"name"`
	PlanID             string        `json:"planID"`
}

// FeePlanAgreement represents a billing fee plan agreement for a Moov account
type FeePlanAgreement struct {
	AgreementID        string                 `json:"agreementID,omitempty"`
	PlanID             string                 `json:"planID,omitempty"`
	AccountID          string                 `json:"accountID,omitempty"`
	Name               string                 `json:"name,omitempty"`
	Description        *string                `json:"description,omitempty"`
	AcceptedOn         time.Time              `json:"acceptedOn,omitempty"`
	Status             FeePlanAgreementStatus `json:"status,omitempty"`
	CardAcquiringModel CardAcquiringModel     `json:"cardAcquiringModel,omitempty"`
	BillableFees       []BillableFee          `json:"billableFees,omitempty"`
	MinimumCommitment  AmountDecimal          `json:"minimumCommitment,omitempty"`
	MonthlyPlatformFee AmountDecimal          `json:"monthlyPlatformFee,omitempty"`
}

type FeePlanAgreementRequest struct {
	PlanID string `json:"planID"`
}

// BillableFee represents a billable fee within a fee plan agreement
type BillableFee struct {
	BillableFeeID string         `json:"billableFeeID,omitempty"`
	BillableEvent string         `json:"billableEvent,omitempty"`
	FeeName       string         `json:"feeName,omitempty"`
	FeeModel      FeeModel       `json:"feeModel,omitempty"`
	FeeCategory   FeeCategory    `json:"feeCategory,omitempty"`
	FeeProperties *FeeProperties `json:"feeProperties,omitempty"`
	FeeConditions *FeeConditions `json:"feeConditions,omitempty"`
}

// FeeProperties represents the properties of a fee, including amounts and rates
type FeeProperties struct {
	FixedAmount       *AmountDecimal `json:"fixedAmount,omitempty"`
	VariableRate      *string        `json:"variableRate,omitempty"`
	MinPerTransaction *AmountDecimal `json:"minPerTransaction,omitempty"`
	MaxPerTransaction *AmountDecimal `json:"maxPerTransaction,omitempty"`
	VolumeRanges      []VolumeRange  `json:"volumeRanges,omitempty"`
}

// VolumeRange represents a volume-based pricing range
type VolumeRange struct {
	FromValue     *int          `json:"fromValue,omitempty"`
	ToValue       *int          `json:"toValue,omitempty"`
	FlatAmount    AmountDecimal `json:"flatAmount,omitempty"`
	PerUnitAmount AmountDecimal `json:"perUnitAmount,omitempty"`
}

// FeeConditions represents conditions under which a fee applies.
// This is a dynamic object that can contain any properties as conditions.
// Common examples include transactionType, cardBrand, etc.
type FeeConditions map[string]interface{}

// FeeModel represents the model used for calculating fees
type FeeModel string

// List of FeeModel
const (
	FeeModel_Fixed    FeeModel = "fixed"
	FeeModel_Blended  FeeModel = "blended"
	FeeModel_Variable FeeModel = "variable"
)

// FeeCategory represents the category of a fee
type FeeCategory string

// List of FeeCategory
const (
	FeeCategory_Ach                FeeCategory = "ach"
	FeeCategory_CardAcquiring      FeeCategory = "card-acquiring"
	FeeCategory_CardOther          FeeCategory = "card-other"
	FeeCategory_CardPull           FeeCategory = "card-pull"
	FeeCategory_CardPush           FeeCategory = "card-push"
	FeeCategory_MonthlyPlatform    FeeCategory = "monthly-platform"
	FeeCategory_NetworkPassthrough FeeCategory = "network-passthrough"
	FeeCategory_Other              FeeCategory = "other"
	FeeCategory_Rtp                FeeCategory = "rtp"
)

// CardAcquiringModel represents the pricing model for card acquiring
type CardAcquiringModel string

// List of CardAcquiringModel
const (
	CardAcquiringModel_CostPlus CardAcquiringModel = "cost-plus"
	CardAcquiringModel_FlatRate CardAcquiringModel = "flat-rate"
)

// FeePlanAgreementStatus represents the status of a fee plan agreement
type FeePlanAgreementStatus string

// List of FeePlanAgreementStatus
const (
	FeePlanAgreementStatus_Active     FeePlanAgreementStatus = "active"
	FeePlanAgreementStatus_Terminated FeePlanAgreementStatus = "terminated"
)

type FeePlanAgreementListFilter callArg

type FeePlanListFilter callArg

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

func WithFeePlanAgreementStatuses(statuses []FeePlanAgreementStatus) FeePlanAgreementListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		statusStrings := make([]string, len(statuses))
		for i, status := range statuses {
			statusStrings[i] = string(status)
		}
		call.params["status"] = strings.Join(statusStrings, ",")
		return nil
	})
}

func WithFeePlanAgreementIds(agreementIds []string) FeePlanAgreementListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["agreementID"] = strings.Join(agreementIds, ",")
		return nil
	})
}

func WithFeePlanIds(planIds []string) FeePlanListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["planIds"] = strings.Join(planIds, ",")
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

// ListFeePlans lists available FeePlans for a Moov account
// https://docs.moov.io/api/moov-accounts/billing/list-plans/
func (c Client) ListFeePlans(ctx context.Context, accountID string, filters ...FeePlanListFilter) ([]FeePlan, error) {
	args := prependArgs(filters, AcceptJson())
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathFeePlans, accountID), args...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[FeePlan](resp)
}

// CreateFeePlanAgreement creates a FeePlanAgreement for a Moov account
// https://docs.moov.io/api/moov-accounts/billing/create-agreement/
func (c Client) CreateFeePlanAgreement(ctx context.Context, accountID string, request FeePlanAgreementRequest) (*FeePlanAgreement, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathFeePlanAgreements, accountID),
		AcceptJson(),
		JsonBody(request))
	if err != nil {
		return nil, err
	}

	switch resp.Status() {
	case StatusStarted, StatusCompleted:
		return UnmarshalObjectResponse[FeePlanAgreement](resp)
	default:
		return nil, resp
	}
}
