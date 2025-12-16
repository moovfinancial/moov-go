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

type IssuedCardAuthorization struct {
	AuthorizationID  string                        `json:"authorizationID"`
	IssuedCardID     string                        `json:"issuedCardID"`
	FundingWalletID  string                        `json:"fundingWalletID"`
	CreatedOn        time.Time                     `json:"createdOn"`
	Network          IssuedCardTransactionNetwork  `json:"network"`
	AuthorizedAmount string                        `json:"authorizedAmount"`
	Status           IssuedCardAuthorizationStatus `json:"status"`
	MerchantData     IssuedCardTransactionMerchant `json:"merchantData"`
	CardTransactions []string                      `json:"cardTransactions,omitempty"`
}

// IssuedCardTransactionNetwork represents name of the network a card transaction is routed through
type IssuedCardTransactionNetwork string

const (
	IssuedCardTransactionNetwork_Discover IssuedCardTransactionNetwork = "discover"
	IssuedCardTransactionNetwork_Shazam   IssuedCardTransactionNetwork = "shazam"
	IssuedCardTransactionNetwork_Visa     IssuedCardTransactionNetwork = "visa"
)

// IssuedCardAuthorizationStatus represents the status of the authorization
type IssuedCardAuthorizationStatus string

const (
	IssuedCardAuthorizationStatus_Pending  IssuedCardAuthorizationStatus = "pending"
	IssuedCardAuthorizationStatus_Declined IssuedCardAuthorizationStatus = "declined"
	IssuedCardAuthorizationStatus_Canceled IssuedCardAuthorizationStatus = "canceled"
	IssuedCardAuthorizationStatus_Cleared  IssuedCardAuthorizationStatus = "cleared"
	IssuedCardAuthorizationStatus_Expired  IssuedCardAuthorizationStatus = "expired"
)

type IssuedCardTransactionMerchant struct {
	// External identifier used to identify the merchant with the card brand
	NetworkID string  `json:"networkID"`
	Name      *string `json:"name,omitempty"`
	City      *string `json:"city,omitempty"`
	// Two-letter code of the merchant country
	Country string `json:"country"`
	// Five digit postal code of merchant
	PostalCode *string `json:"postalCode,omitempty"`
	// Two-letter code of merchant state
	State *string `json:"state,omitempty"`
	Mcc   string  `json:"mcc"`
}

type ListIssuedCardAuthorizationsFilter callArg

func WithIssuedCardAuthorizationStatuses(statuses []IssuedCardAuthorizationStatus) ListIssuedCardAuthorizationsFilter {
	return callBuilderFn(func(call *callBuilder) error {
		statusStrings := make([]string, len(statuses))
		for i, status := range statuses {
			statusStrings[i] = string(status)
		}
		call.params["statuses"] = strings.Join(statusStrings, ",")
		return nil
	})
}

func WithIssuedCardAuthorizationCardID(cardID string) ListIssuedCardAuthorizationsFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["issuedCardID"] = cardID
		return nil
	})
}

func WithIssuedCardAuthorizationSkip(skip int) ListIssuedCardAuthorizationsFilter {
	return Skip(skip)
}

func WithIssuedCardAuthorizationCount(count int) ListIssuedCardAuthorizationsFilter {
	return Count(count)
}

func WithIssuedCardAuthorizationStartDate(t time.Time) ListIssuedCardAuthorizationsFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["startDateTime"] = t.Format(time.RFC3339)
		return nil
	})
}

func WithIssuedCardAuthorizationEndDate(t time.Time) ListIssuedCardAuthorizationsFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["endDateTime"] = t.Format(time.RFC3339)
		return nil
	})
}

type IssuedCardAuthorizationEvent struct {
	EventID   string                             `json:"eventID"`
	EventType IssuedCardEventType                `json:"eventType"`
	CreatedOn time.Time                          `json:"createdOn"`
	Amount    string                             `json:"amount"`
	Result    IssuedCardAuthorizationEventResult `json:"result"`
}

type IssuedCardEventType string

const (
	IssuedCardEventType_Authorization            IssuedCardEventType = "authorization"
	IssuedCardEventType_Reversal                 IssuedCardEventType = "reversal"
	IssuedCardEventType_AuthorizationAdvice      IssuedCardEventType = "authorization-advice"
	IssuedCardEventType_AuthorizationExpiration  IssuedCardEventType = "authorization-expiration"
	IssuedCardEventType_AuthorizationIncremental IssuedCardEventType = "authorization-incremental"
	IssuedCardEventType_Clearing                 IssuedCardEventType = "clearing"
)

type IssuedCardAuthorizationEventResult string

const (
	IssuedCardAuthorizationEventResult_Approved  IssuedCardAuthorizationEventResult = "approved"
	IssuedCardAuthorizationEventResult_Declined  IssuedCardAuthorizationEventResult = "declined"
	IssuedCardAuthorizationEventResult_Processed IssuedCardAuthorizationEventResult = "processed"
)

type ListIssuedCardAuthorizationEventsFilter callArg

func WithIssuedCardAuthorizationEventSkip(skip int) ListIssuedCardAuthorizationEventsFilter {
	return Skip(skip)
}

func WithIssuedCardAuthorizationEventCount(count int) ListIssuedCardAuthorizationEventsFilter {
	return Count(count)
}

type IssuedCardTransaction struct {
	CardTransactionID string                        `json:"cardTransactionID"`
	IssuedCardID      string                        `json:"issuedCardID"`
	FundingWalletID   string                        `json:"fundingWalletID"`
	Amount            string                        `json:"amount"`
	AuthorizationID   *string                       `json:"authorizationID,omitempty"`
	CreatedOn         time.Time                     `json:"createdOn"`
	AuthorizedOn      time.Time                     `json:"authorizedOn"`
	MerchantData      IssuedCardTransactionMerchant `json:"merchantData"`
}

type ListIssuedCardTransactionsFilter callArg

func WithIssuedCardTransactionCardID(cardID string) ListIssuedCardTransactionsFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["issuedCardID"] = cardID
		return nil
	})
}

func WithIssuedCardTransactionSkip(skip int) ListIssuedCardTransactionsFilter {
	return Skip(skip)
}

func WithIssuedCardTransactionCount(count int) ListIssuedCardTransactionsFilter {
	return Count(count)
}

func WithIssuedCardTransactionStartDate(t time.Time) ListIssuedCardTransactionsFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["startDateTime"] = t.Format(time.RFC3339)
		return nil
	})
}

func WithIssuedCardTransactionEndDate(t time.Time) ListIssuedCardTransactionsFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["endDateTime"] = t.Format(time.RFC3339)
		return nil
	})
}
