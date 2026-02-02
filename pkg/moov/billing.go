package moov

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type FeePlan struct {
	BillableFees       []BillableFee `json:"billableFees,omitempty"`
	CardAcquiringModel string        `json:"cardAcquiringModel,omitempty"`
	CreatedAt          time.Time     `json:"createdAt,omitempty"`
	Description        string        `json:"description,omitempty"`
	MinimumCommitment  AmountDecimal `json:"minimumCommitment,omitempty"`
	MonthlyPlatformFee AmountDecimal `json:"monthlyPlatformFee,omitempty"`
	Name               string        `json:"name,omitempty"`
	PlanID             string        `json:"planID,omitempty"`
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

// Residual represents a partner residual payment calculation
type Residual struct {
	ResidualID       string        `json:"residualID"`
	PartnerAccountID string        `json:"partnerAccountID"`
	PeriodStart      time.Time     `json:"periodStart"`
	PeriodEnd        time.Time     `json:"periodEnd"`
	MerchantFees     AmountDecimal `json:"merchantFees"`
	PartnerCost      AmountDecimal `json:"partnerCost"`
	NetIncome        AmountDecimal `json:"netIncome"`
	RevenueShare     string        `json:"revenueShare"`
	ResidualAmount   AmountDecimal `json:"residualAmount"`
	CreatedOn        time.Time     `json:"createdOn"`
	UpdatedOn        time.Time     `json:"updatedOn"`
}

// PartnerPricing represents a partner pricing plan
type PartnerPricing struct {
	PlanID             string             `json:"planID"`
	Name               string             `json:"name"`
	Description        *string            `json:"description,omitempty"`
	RevenueShare       string             `json:"revenueShare"`
	CardAcquiringModel CardAcquiringModel `json:"cardAcquiringModel"`
	BillableFees       []BillableFee      `json:"billableFees"`
	MinimumCommitment  AmountDecimal      `json:"minimumCommitment"`
	MonthlyPlatformFee AmountDecimal      `json:"monthlyPlatformFee"`
	CreatedAt          time.Time          `json:"createdAt"`
}

// PartnerPricingAgreement represents a partner pricing agreement
type PartnerPricingAgreement struct {
	AgreementID        string                 `json:"agreementID"`
	PlanID             string                 `json:"planID"`
	AccountID          string                 `json:"accountID,omitempty"`
	Name               string                 `json:"name"`
	Description        *string                `json:"description,omitempty"`
	AcceptedOn         time.Time              `json:"acceptedOn"`
	Status             FeePlanAgreementStatus `json:"status"`
	RevenueShare       string                 `json:"revenueShare"`
	CardAcquiringModel CardAcquiringModel     `json:"cardAcquiringModel"`
	BillableFees       []BillableFee          `json:"billableFees"`
	MinimumCommitment  AmountDecimal          `json:"minimumCommitment"`
	MonthlyPlatformFee AmountDecimal          `json:"monthlyPlatformFee"`
}

type FeePlanAgreementListFilter callArg

type FeePlanListFilter callArg

type ResidualListFilter callArg

type ResidualFeeListFilter callArg

type PartnerPricingAgreementListFilter callArg

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

func WithResidualCount(c int) ResidualListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["count"] = fmt.Sprintf("%d", c)
		return nil
	})
}

func WithResidualSkip(c int) ResidualListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["skip"] = fmt.Sprintf("%d", c)
		return nil
	})
}

func WithResidualStartDateTime(start string) ResidualListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["startDateTime"] = start
		return nil
	})
}

func WithResidualEndDateTime(end string) ResidualListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["endDateTime"] = end
		return nil
	})
}

func WithResidualFeeCount(c int) ResidualFeeListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["count"] = fmt.Sprintf("%d", c)
		return nil
	})
}

func WithResidualFeeSkip(c int) ResidualFeeListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["skip"] = fmt.Sprintf("%d", c)
		return nil
	})
}

func WithResidualFeeStartDateTime(start string) ResidualFeeListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["startDateTime"] = start
		return nil
	})
}

func WithResidualFeeEndDateTime(end string) ResidualFeeListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["endDateTime"] = end
		return nil
	})
}

func WithPartnerPricingAgreementCount(c int) PartnerPricingAgreementListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["count"] = fmt.Sprintf("%d", c)
		return nil
	})
}

func WithPartnerPricingAgreementSkip(c int) PartnerPricingAgreementListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["skip"] = fmt.Sprintf("%d", c)
		return nil
	})
}

func WithPartnerPricingAgreementStatuses(statuses []FeePlanAgreementStatus) PartnerPricingAgreementListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		statusStrings := make([]string, len(statuses))
		for i, status := range statuses {
			statusStrings[i] = string(status)
		}
		call.params["status"] = strings.Join(statusStrings, ",")
		return nil
	})
}

func WithPartnerPricingAgreementIds(agreementIds []string) PartnerPricingAgreementListFilter {
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

// ListResiduals lists all residuals associated with an account
// TODO: Add docs.moov.io link when published
func (c Client) ListResiduals(ctx context.Context, accountID string, filters ...ResidualListFilter) ([]Residual, error) {
	args := prependArgs(filters, AcceptJson())
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathResiduals, accountID), args...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[Residual](resp)
}

// GetResidual retrieves a specific residual by ID
// TODO: Add docs.moov.io link when published
func (c Client) GetResidual(ctx context.Context, accountID, residualID string) (*Residual, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathResidual, accountID, residualID),
		AcceptJson(),
	)
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Residual](resp)
}

// ListResidualFees lists all fees associated with a residual
// TODO: Add docs.moov.io link when published
func (c Client) ListResidualFees(ctx context.Context, accountID, residualID string, filters ...ResidualFeeListFilter) ([]IncurredFee, error) {
	args := prependArgs(filters, AcceptJson())
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathResidualFees, accountID, residualID), args...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[IncurredFee](resp)
}

// ListPartnerPricingAgreements lists all partner pricing agreements associated with an account
// TODO: Add docs.moov.io link when published
func (c Client) ListPartnerPricingAgreements(ctx context.Context, accountID string, filters ...PartnerPricingAgreementListFilter) ([]PartnerPricingAgreement, error) {
	args := prependArgs(filters, AcceptJson())
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathPartnerPricingAgreements, accountID), args...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[PartnerPricingAgreement](resp)
}
