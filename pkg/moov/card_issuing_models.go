package moov

import (
	"strings"
	"time"
)

type IssuedCard struct {
	IssuedCardID       string               `json:"issuedCardID"`
	Brand              IssuedCardBrand      `json:"brand"`
	LastFourCardNumber string               `json:"lastFourCardNumber"`
	Expiration         IssuedCardExpiration `json:"expiration"`
	AuthorizedUser     AuthorizedUser       `json:"authorizedUser"`
	Memo               *string              `json:"memo,omitempty"`
	FundingWalletID    string               `json:"fundingWalletID"`
	State              IssuedCardState      `json:"state"`
	FormFactor         IssuedCardFormFactor `json:"formFactor"`
	Controls           *IssuingControls     `json:"controls,omitempty"`
	CreatedOn          time.Time            `json:"createdOn"`
}

type IssuedCardBrand string

const (
	IssuedCardBrand_AmericanExpress IssuedCardBrand = "American Express"
	IssuedCardBrand_Discover        IssuedCardBrand = "Discover"
	IssuedCardBrand_Mastercard      IssuedCardBrand = "Mastercard"
	IssuedCardBrand_Visa            IssuedCardBrand = "Visa"
)

type IssuedCardExpiration struct {
	Month string `json:"month,omitempty"`
	Year  string `json:"year,omitempty"`
}

type AuthorizedUser struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// IssuedCardState represents the operational status of an IssuedCard
type IssuedCardState string

const (
	// operational and can approve incoming authorizations
	IssuedCardState_Active IssuedCardState = "active"

	// still in the activation process and cannot yet approve incoming authorizations
	IssuedCardState_Inactive IssuedCardState = "inactive"

	// awaiting additional AuthorizedUser verification before the card can become active
	IssuedCardState_PendingVerification IssuedCardState = "pending-verification"

	// permanently deactivated, either by request or because it expired, and cannot approve incoming authorizations
	IssuedCardState_Closed IssuedCardState = "closed"
)

// IssuedCardFormFactor specifies the type of IssuedCard
type IssuedCardFormFactor string

const (
	// provides a digital number without a physical card
	IssuedCardFormFactor_Virtual IssuedCardFormFactor = "virtual"
)

// IssuingControls specifies any controls that should apply to the IssuedCard
type IssuingControls struct {
	// if true, the card closes after the first authorization
	SingleUse *bool `json:"singleUse,omitempty"`

	VelocityLimits []IssuingVelocityLimit `json:"velocityLimits,omitempty"`
}

// IssuingVelocityLimit specifies any spending limits per time Interval that should apply to the IssuedCard
type IssuingVelocityLimit struct {
	// the maximum amount in cents that can be spent in a given interval
	Amount *int32 `json:"amount,omitempty"`

	Interval *IssuingIntervalLimit `json:"interval,omitempty"`
}

// IssuingIntervalLimit specifies the time frame for the IssuingVelocityLimit
type IssuingIntervalLimit string

const (
	IssuingIntervalLimit_PerTransaction IssuingIntervalLimit = "per-transaction"
)

type ListIssuedCardsFilter callArg

func WithIssuedCardStates(states []IssuedCardState) ListIssuedCardsFilter {
	return callBuilderFn(func(call *callBuilder) error {
		stateStrings := make([]string, len(states))
		for i, state := range states {
			stateStrings[i] = string(state)
		}
		call.params["states"] = strings.Join(stateStrings, ",")
		return nil
	})
}

func WithIssuedCardSkip(skip int) ListIssuedCardsFilter {
	return Skip(skip)
}

func WithIssuedCardCount(count int) ListIssuedCardsFilter {
	return Count(count)
}

type CreateIssuedCard struct {
	FundingWalletID string                `json:"fundingWalletID"`
	AuthorizedUser  CreateAuthorizedUser  `json:"authorizedUser"`
	FormFactor      IssuedCardFormFactor  `json:"formFactor"`
	Memo            *string               `json:"memo,omitempty"`
	Expiration      *IssuedCardExpiration `json:"expiration,omitempty"`
	Controls        *IssuingControls      `json:"controls,omitempty"`
}

type CreateAuthorizedUser struct {
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	BirthDate *BirthDate `json:"birthDate,omitempty"`
}

type BirthDate struct {
	Day   int32 `json:"day"`
	Month int32 `json:"month"`
	Year  int32 `json:"year"`
}

type UpdateIssuedCard struct {
	State          *UpdateIssuedCardState `json:"state,omitempty"`
	Memo           *string                `json:"memo,omitempty"`
	AuthorizedUser *CreateAuthorizedUser  `json:"authorizedUser,omitempty"`
}

type UpdateIssuedCardState string

const (
	UpdateIssuedCardState_Closed UpdateIssuedCardState = "closed"
)
