package moov

import (
	"strings"
	"time"
)

type IssuedCard struct {
	IssuedCardID            string               `json:"issuedCardID"`
	Brand                   IssuedCardBrand      `json:"brand"`
	LastFourCardNumber      string               `json:"lastFourCardNumber"`
	Expiration              IssuedCardExpiration `json:"expiration"`
	AuthorizedUserAccountID *string              `json:"authorizedUserAccountID,omitempty"`
	Nickname                *string              `json:"nickname,omitempty"`
	FundingWalletID         string               `json:"fundingWalletID"`
	State                   IssuedCardState      `json:"state"`
	FormFactor              IssuedCardFormFactor `json:"formFactor"`
	BillingAddress          *Address             `json:"billingAddress,omitempty"`
	Controls                *IssuedControls      `json:"controls,omitempty"`
	Metadata                map[string]string    `json:"metadata,omitempty"`
	CreatedOn               time.Time            `json:"createdOn"`
	UpdatedOn               time.Time            `json:"updatedOn"`
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

// IssuedCardState represents the operational status of an IssuedCard
type IssuedCardState string

const (
	// operational and can approve incoming authorizations
	IssuedCardState_Active IssuedCardState = "active"

	// permanently deactivated, either by request or because it expired, and cannot approve incoming authorizations
	IssuedCardState_Closed IssuedCardState = "closed"
)

// IssuedCardFormFactor specifies the type of IssuedCard
type IssuedCardFormFactor string

const (
	// provides a digital number without a physical card
	IssuedCardFormFactor_Virtual IssuedCardFormFactor = "virtual"
)

// IssuingControls specifies any controls that should apply to the IssuedCard on create
type IssuingControls struct {
	// if true, the card closes after the first authorization
	SingleUse *bool `json:"singleUse,omitempty"`

	VelocityLimits []IssuingVelocityLimit `json:"velocityLimits,omitempty"`

	// restricts card usage by merchant category; when not set, all categories are allowed
	MerchantCategoryRestrictions *MerchantCategoryRestrictions `json:"merchantCategoryRestrictions,omitempty"`

	// restricts card usage to specific merchants, or blocks specific merchants
	MerchantRestrictions *MerchantRestrictions `json:"merchantRestrictions,omitempty"`

	// limits card usage to specific days and times
	AllowedSchedule *AllowedSchedule `json:"allowedSchedule,omitempty"`

	// a spend cutoff; all authorizations after this datetime are declined regardless of other controls
	ExpiresOn *time.Time `json:"expiresOn,omitempty"`
}

// IssuedControls specifies the controls applied to an IssuedCard on read, including velocity runtime state
type IssuedControls struct {
	// if true, the card closes after the first authorization
	SingleUse *bool `json:"singleUse,omitempty"`

	VelocityLimits []IssuedVelocityLimit `json:"velocityLimits,omitempty"`

	// restricts card usage by merchant category; when not set, all categories are allowed
	MerchantCategoryRestrictions *MerchantCategoryRestrictions `json:"merchantCategoryRestrictions,omitempty"`

	// restricts card usage to specific merchants, or blocks specific merchants
	MerchantRestrictions *MerchantRestrictions `json:"merchantRestrictions,omitempty"`

	// limits card usage to specific days and times
	AllowedSchedule *AllowedSchedule `json:"allowedSchedule,omitempty"`

	// a spend cutoff; all authorizations after this datetime are declined regardless of other controls
	ExpiresOn *time.Time `json:"expiresOn,omitempty"`
}

// UpdateIssuingControls specifies the mutable controls on a PATCH. Each field replaces the entire
// corresponding value.
type UpdateIssuingControls struct {
	// replaces the entire set of velocity limits; nil leaves them unchanged, a non-nil empty slice clears them
	VelocityLimits *[]IssuingVelocityLimit `json:"velocityLimits,omitempty"`

	// replaces the merchant category restrictions; use SetNull to remove
	MerchantCategoryRestrictions *Nullable[MerchantCategoryRestrictions] `json:"merchantCategoryRestrictions,omitempty"`

	// replaces the merchant restrictions; use SetNull to remove
	MerchantRestrictions *Nullable[MerchantRestrictions] `json:"merchantRestrictions,omitempty"`

	// replaces the allowed schedule; use SetNull to remove all schedule restrictions
	AllowedSchedule *Nullable[AllowedSchedule] `json:"allowedSchedule,omitempty"`

	// a spend cutoff; use SetNull to remove the cutoff
	ExpiresOn *Nullable[time.Time] `json:"expiresOn,omitempty"`
}

// IssuingVelocityLimit specifies any spending limits per time Interval that should apply to the IssuedCard
type IssuingVelocityLimit struct {
	// the maximum amount in cents that can be spent in a given interval
	Amount *int32 `json:"amount,omitempty"`

	// the maximum number of transactions allowed in the given interval; at least one of Amount or Count must be set
	Count *int32 `json:"count,omitempty"`

	Interval *IssuingIntervalLimit `json:"interval,omitempty"`
}

// IssuedVelocityLimit is a velocity limit with its current runtime state, returned on read
type IssuedVelocityLimit struct {
	// the maximum amount in cents that can be spent in a given interval
	Amount *int32 `json:"amount,omitempty"`

	// the maximum number of transactions allowed in the given interval
	Count *int32 `json:"count,omitempty"`

	Interval *IssuingIntervalLimit `json:"interval,omitempty"`

	// the amount in cents already spent in the current interval
	AmountUsed *int32 `json:"amountUsed,omitempty"`

	// the amount in cents remaining in the current interval
	AmountRemaining *int32 `json:"amountRemaining,omitempty"`

	// the number of transactions already made in the current interval
	CountUsed *int32 `json:"countUsed,omitempty"`

	// the number of transactions remaining in the current interval
	CountRemaining *int32 `json:"countRemaining,omitempty"`

	// when the current interval resets; absent for per-transaction limits
	ResetsOn *time.Time `json:"resetsOn,omitempty"`
}

// IssuingIntervalLimit specifies the time frame for the IssuingVelocityLimit
type IssuingIntervalLimit string

const (
	IssuingIntervalLimit_PerTransaction IssuingIntervalLimit = "per-transaction"
	IssuingIntervalLimit_Daily          IssuingIntervalLimit = "daily"
	IssuingIntervalLimit_Weekly         IssuingIntervalLimit = "weekly"
	IssuingIntervalLimit_Monthly        IssuingIntervalLimit = "monthly"
)

// IssuingControlsRestrictionMode indicates whether the listed items are the only ones allowed, or the ones to block
type IssuingControlsRestrictionMode string

const (
	IssuingControlsRestrictionMode_Allow IssuingControlsRestrictionMode = "allow"
	IssuingControlsRestrictionMode_Block IssuingControlsRestrictionMode = "block"
)

// MerchantEntry identifies a merchant by ID, descriptor pattern, or both. At least one of Mid or
// DescriptorPattern must be set.
type MerchantEntry struct {
	// the merchant's unique identifier (ISO 8583 DE42), matched exactly
	Mid *string `json:"mid,omitempty"`

	// a case-insensitive RE2 regular expression matched against the merchant descriptor (ISO 8583 DE43)
	DescriptorPattern *string `json:"descriptorPattern,omitempty"`

	// an optional label for this entry
	Name *string `json:"name,omitempty"`
}

// MerchantCategoryRestrictions restricts card usage by merchant category
type MerchantCategoryRestrictions struct {
	// whether the listed categories are the only ones allowed, or the ones to block
	Mode IssuingControlsRestrictionMode `json:"mode"`

	// predefined category groups to allow or block
	Categories []IssuingMerchantCategory `json:"categories,omitempty"`

	// individual merchant category codes (MCCs) to allow or block, for codes not covered by a predefined category
	CustomMCCs []string `json:"customMCCs,omitempty"`

	// merchants that are exempt from category restrictions regardless of their category
	ExemptMerchants []MerchantEntry `json:"exemptMerchants,omitempty"`
}

// MerchantRestrictions restricts card usage to specific merchants, independent of merchant category
type MerchantRestrictions struct {
	// whether the listed merchants are the only ones allowed, or the ones to block
	Mode IssuingControlsRestrictionMode `json:"mode"`

	// the merchants to allow or block
	Merchants []MerchantEntry `json:"merchants"`
}

// IssuingScheduleDay is a day of the week used by an AllowedSchedule window
type IssuingScheduleDay string

const (
	IssuingScheduleDay_Monday    IssuingScheduleDay = "monday"
	IssuingScheduleDay_Tuesday   IssuingScheduleDay = "tuesday"
	IssuingScheduleDay_Wednesday IssuingScheduleDay = "wednesday"
	IssuingScheduleDay_Thursday  IssuingScheduleDay = "thursday"
	IssuingScheduleDay_Friday    IssuingScheduleDay = "friday"
	IssuingScheduleDay_Saturday  IssuingScheduleDay = "saturday"
	IssuingScheduleDay_Sunday    IssuingScheduleDay = "sunday"
)

// ScheduleWindow is a window of time during which the card may authorize
type ScheduleWindow struct {
	// the days of the week this window applies to
	Days []IssuingScheduleDay `json:"days"`

	// inclusive window start time in 24-hour HH:MM format
	StartTime string `json:"startTime"`

	// exclusive window end time in 24-hour HH:MM format; if earlier than StartTime, the window wraps past midnight
	EndTime string `json:"endTime"`
}

// AllowedSchedule limits card usage to specific days and times
type AllowedSchedule struct {
	// IANA timezone string used to evaluate window boundaries against the authorization time
	Timezone string `json:"timezone"`

	// time windows during which the card may authorize; any matching window allows the transaction
	Windows []ScheduleWindow `json:"windows"`
}

// IssuingMerchantCategory is a predefined merchant category group
type IssuingMerchantCategory string

const (
	IssuingMerchantCategory_Advertising          IssuingMerchantCategory = "advertising"
	IssuingMerchantCategory_Airlines             IssuingMerchantCategory = "airlines"
	IssuingMerchantCategory_AlcoholAndBars       IssuingMerchantCategory = "alcohol-and-bars"
	IssuingMerchantCategory_CarRental            IssuingMerchantCategory = "car-rental"
	IssuingMerchantCategory_Education            IssuingMerchantCategory = "education"
	IssuingMerchantCategory_Electronics          IssuingMerchantCategory = "electronics"
	IssuingMerchantCategory_FuelAndGas           IssuingMerchantCategory = "fuel-and-gas"
	IssuingMerchantCategory_Gambling             IssuingMerchantCategory = "gambling"
	IssuingMerchantCategory_Groceries            IssuingMerchantCategory = "groceries"
	IssuingMerchantCategory_GroundTransportation IssuingMerchantCategory = "ground-transportation"
	IssuingMerchantCategory_HardwareAndHome      IssuingMerchantCategory = "hardware-and-home"
	IssuingMerchantCategory_Healthcare           IssuingMerchantCategory = "healthcare"
	IssuingMerchantCategory_LiveEntertainment    IssuingMerchantCategory = "live-entertainment"
	IssuingMerchantCategory_Lodging              IssuingMerchantCategory = "lodging"
	IssuingMerchantCategory_Movies               IssuingMerchantCategory = "movies"
	IssuingMerchantCategory_OfficeSupplies       IssuingMerchantCategory = "office-supplies"
	IssuingMerchantCategory_Parking              IssuingMerchantCategory = "parking"
	IssuingMerchantCategory_PersonalCare         IssuingMerchantCategory = "personal-care"
	IssuingMerchantCategory_ProfessionalServices IssuingMerchantCategory = "professional-services"
	IssuingMerchantCategory_RestaurantsAndDining IssuingMerchantCategory = "restaurants-and-dining"
	IssuingMerchantCategory_RetailGeneral        IssuingMerchantCategory = "retail-general"
	IssuingMerchantCategory_RideshareAndTaxis    IssuingMerchantCategory = "rideshare-and-taxis"
	IssuingMerchantCategory_SoftwareAndSaas      IssuingMerchantCategory = "software-and-saas"
	IssuingMerchantCategory_SportsAndRecreation  IssuingMerchantCategory = "sports-and-recreation"
	IssuingMerchantCategory_Subscriptions        IssuingMerchantCategory = "subscriptions"
	IssuingMerchantCategory_TravelAgencies       IssuingMerchantCategory = "travel-agencies"
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
	AuthorizedUserAccountID *string               `json:"authorizedUserAccountID,omitempty"`
	Nickname                *string               `json:"nickname,omitempty"`
	Metadata                map[string]string     `json:"metadata,omitempty"`
	BillingAddress          *Address              `json:"billingAddress,omitempty"`
	Expiration              *IssuedCardExpiration `json:"expiration,omitempty"`
	Controls                *IssuingControls      `json:"controls,omitempty"`
}

type UpdateIssuedCard struct {
	State          *UpdateIssuedCardState `json:"state,omitempty"`
	Nickname       *string                `json:"nickname,omitempty"`
	Metadata       map[string]string      `json:"metadata,omitempty"`
	BillingAddress *AddressPatch          `json:"billingAddress,omitempty"`
	Controls       *UpdateIssuingControls `json:"controls,omitempty"`
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
