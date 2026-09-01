package mv2610

import "github.com/moovfinancial/moov-go/pkg/moov"

// UpsertUnderwriting is the v2026.10 underwriting request body.
type UpsertUnderwriting struct {
	GeographicReach           *GeographicReach           `json:"geographicReach,omitempty"`
	BusinessPresence          *BusinessPresence          `json:"businessPresence,omitempty"`
	PendingLitigation         *PendingLitigation         `json:"pendingLitigation,omitempty"`
	VolumeShareByCustomerType *VolumeShareByCustomerType `json:"volumeShareByCustomerType,omitempty"`
	CollectFunds              *CollectFunds              `json:"collectFunds,omitempty"`
	MoneyTransfer             *MoneyTransfer             `json:"moneyTransfer,omitempty"`
	SendFunds                 *SendFunds                 `json:"sendFunds,omitempty"`
	// CardIssuing was added in v2026.10.
	CardIssuing      *CardIssuing      `json:"cardIssuing,omitempty"`
	SubmissionIntent *SubmissionIntent `json:"submissionIntent,omitempty"`
}

// Underwriting is the v2026.10 underwriting response. The legacy fields are only
// returned for accounts that submitted underwriting through the pre-versioned API, so
// they're optional.
type Underwriting struct {
	AverageTransactionSize          *int64                       `json:"averageTransactionSize,omitempty"`
	MaxTransactionSize              *int64                       `json:"maxTransactionSize,omitempty"`
	AverageMonthlyTransactionVolume *int64                       `json:"averageMonthlyTransactionVolume,omitempty"`
	VolumeByCustomerType            *moov.VolumeByCustomerType   `json:"volumeByCustomerType,omitempty"`
	CardVolumeDistribution          *moov.CardVolumeDistribution `json:"cardVolumeDistribution,omitempty"`
	Fulfillment                     *moov.Fulfillment            `json:"fulfillment,omitempty"`

	GeographicReach           *GeographicReach           `json:"geographicReach,omitempty"`
	BusinessPresence          *BusinessPresence          `json:"businessPresence,omitempty"`
	PendingLitigation         *PendingLitigation         `json:"pendingLitigation,omitempty"`
	VolumeShareByCustomerType *VolumeShareByCustomerType `json:"volumeShareByCustomerType,omitempty"`
	CollectFunds              *CollectFunds              `json:"collectFunds,omitempty"`
	MoneyTransfer             *MoneyTransfer             `json:"moneyTransfer,omitempty"`
	SendFunds                 *SendFunds                 `json:"sendFunds,omitempty"`
	// CardIssuing was added in v2026.10.
	CardIssuing *CardIssuing `json:"cardIssuing,omitempty"`
}

type GeographicReach string

const (
	GeographicReachInternationalOnly  GeographicReach = "international-only"
	GeographicReachUsAndInternational GeographicReach = "us-and-international"
	GeographicReachUsOnly             GeographicReach = "us-only"
)

type PendingLitigation string

const (
	PendingLitigationBankruptcyOrInsolvency               PendingLitigation = "bankruptcy-or-insolvency"
	PendingLitigationConsumerProtectionOrClassAction      PendingLitigation = "consumer-protection-or-class-action"
	PendingLitigationDataBreachOrPrivacy                  PendingLitigation = "data-breach-or-privacy"
	PendingLitigationEmploymentOrWorkplaceDisputes        PendingLitigation = "employment-or-workplace-disputes"
	PendingLitigationFraudOrFinancialCrime                PendingLitigation = "fraud-or-financial-crime"
	PendingLitigationGovernmentEnforcementOrInvestigation PendingLitigation = "government-enforcement-or-investigation"
	PendingLitigationIntellectualProperty                 PendingLitigation = "intellectual-property"
	PendingLitigationNone                                 PendingLitigation = "none"
	PendingLitigationOther                                PendingLitigation = "other"
	PendingLitigationPersonalInjuryOrMedical              PendingLitigation = "personal-injury-or-medical"
)

type BusinessPresence string

const (
	BusinessPresenceCommercialOffice BusinessPresence = "commercial-office"
	BusinessPresenceHomeBased        BusinessPresence = "home-based"
	BusinessPresenceMixedPresence    BusinessPresence = "mixed-presence"
	BusinessPresenceMobileBusiness   BusinessPresence = "mobile-business"
	BusinessPresenceOnlineOnly       BusinessPresence = "online-only"
	BusinessPresenceRetailStorefront BusinessPresence = "retail-storefront"
)

type VolumeShareByCustomerType struct {
	Business *int `json:"business,omitempty"`
	Consumer *int `json:"consumer,omitempty"`
	P2P      *int `json:"p2p,omitempty"`
}

type EstimatedActivity struct {
	AverageTransactionAmount *int64              `json:"averageTransactionAmount,omitempty"`
	MaximumTransactionAmount *int64              `json:"maximumTransactionAmount,omitempty"`
	MonthlyVolumeRange       *MonthlyVolumeRange `json:"monthlyVolumeRange,omitempty"`
}

// MonthlyVolumeRange includes the low value in each range and excludes the high value.
type MonthlyVolumeRange string

const (
	MonthlyVolumeRangeUnder10K MonthlyVolumeRange = "under-10k"
	MonthlyVolumeRange10K50K   MonthlyVolumeRange = "10k-50k"
	MonthlyVolumeRange50K100K  MonthlyVolumeRange = "50k-100k"
	MonthlyVolumeRange100K250K MonthlyVolumeRange = "100k-250k"
	MonthlyVolumeRange250K500K MonthlyVolumeRange = "250k-500k"
	MonthlyVolumeRange500K1M   MonthlyVolumeRange = "500k-1m"
	MonthlyVolumeRange1M5M     MonthlyVolumeRange = "1m-5m"
	MonthlyVolumeRangeOver5M   MonthlyVolumeRange = "over-5m"
)

type CollectFunds struct {
	Ach          *CollectFundsAch          `json:"ach,omitempty"`
	CardPayments *CollectFundsCardPayments `json:"cardPayments,omitempty"`
}

type CollectFundsAch struct {
	EstimatedActivity *EstimatedActivity `json:"estimatedActivity,omitempty"`
}

type CollectFundsCardPayments struct {
	CardAcceptanceMethods *CardAcceptanceMethods  `json:"cardAcceptanceMethods,omitempty"`
	CurrentlyAcceptsCards *bool                   `json:"currentlyAcceptsCards,omitempty"`
	EstimatedActivity     *EstimatedActivity      `json:"estimatedActivity,omitempty"`
	Fulfillment           *CardPaymentFulfillment `json:"fulfillment,omitempty"`
	RefundPolicy          *RefundPolicy           `json:"refundPolicy,omitempty"`
}

type CardAcceptanceMethods struct {
	InPersonPercentage    *int `json:"inPersonPercentage,omitempty"`
	MailOrPhonePercentage *int `json:"mailOrPhonePercentage,omitempty"`
	OnlinePercentage      *int `json:"onlinePercentage,omitempty"`
}

type CardPaymentFulfillment struct {
	Method    FulfillmentMethod    `json:"method"`
	Timeframe FulfillmentTimeframe `json:"timeframe"`
}

type FulfillmentMethod string

const (
	FulfillmentMethodBillOrDebtPayment        FulfillmentMethod = "bill-or-debt-payment"
	FulfillmentMethodDigitalContent           FulfillmentMethod = "digital-content"
	FulfillmentMethodDonation                 FulfillmentMethod = "donation"
	FulfillmentMethodInPersonService          FulfillmentMethod = "in-person-service"
	FulfillmentMethodLocalPickupOrDelivery    FulfillmentMethod = "local-pickup-or-delivery"
	FulfillmentMethodOther                    FulfillmentMethod = "other"
	FulfillmentMethodRemoteService            FulfillmentMethod = "remote-service"
	FulfillmentMethodShippedPhysicalGoods     FulfillmentMethod = "shipped-physical-goods"
	FulfillmentMethodSubscriptionOrMembership FulfillmentMethod = "subscription-or-membership"
)

type FulfillmentTimeframe string

const (
	FulfillmentTimeframeImmediate         FulfillmentTimeframe = "immediate"
	FulfillmentTimeframeOther             FulfillmentTimeframe = "other"
	FulfillmentTimeframeOver30Days        FulfillmentTimeframe = "over-30-days"
	FulfillmentTimeframePreOrder          FulfillmentTimeframe = "pre-order"
	FulfillmentTimeframeRecurringSchedule FulfillmentTimeframe = "recurring-schedule"
	FulfillmentTimeframeScheduledEvent    FulfillmentTimeframe = "scheduled-event"
	FulfillmentTimeframeWithin30Days      FulfillmentTimeframe = "within-30-days"
	FulfillmentTimeframeWithin7Days       FulfillmentTimeframe = "within-7-days"
)

type RefundPolicy string

const (
	RefundPolicyConditionalRefund        RefundPolicy = "conditional-refund"
	RefundPolicyCustomPolicy             RefundPolicy = "custom-policy"
	RefundPolicyEventBasedPolicy         RefundPolicy = "event-based-policy"
	RefundPolicyFullRefundExtendedWindow RefundPolicy = "full-refund-extended-window"
	RefundPolicyFullRefundWithin30Days   RefundPolicy = "full-refund-within-30-days"
	RefundPolicyNoRefunds                RefundPolicy = "no-refunds"
	RefundPolicyPartialRefund            RefundPolicy = "partial-refund"
	RefundPolicyProratedRefund           RefundPolicy = "prorated-refund"
	RefundPolicyStoreCreditOnly          RefundPolicy = "store-credit-only"
)

type MoneyTransfer struct {
	PullFromCard *MoneyTransferPullFromCard `json:"pullFromCard,omitempty"`
	PushToCard   *MoneyTransferPushToCard   `json:"pushToCard,omitempty"`
}

type MoneyTransferPullFromCard struct {
	EstimatedActivity *EstimatedActivity `json:"estimatedActivity,omitempty"`
}

type MoneyTransferPushToCard struct {
	EstimatedActivity *EstimatedActivity `json:"estimatedActivity,omitempty"`
}

type SendFunds struct {
	Ach         *SendFundsAch         `json:"ach,omitempty"`
	PushToCard  *SendFundsPushToCard  `json:"pushToCard,omitempty"`
	Rtp         *SendFundsRtp         `json:"rtp,omitempty"`
	InstantBank *SendFundsInstantBank `json:"instantBank,omitempty"`
	Wire        *SendFundsWire        `json:"wire,omitempty"`
}

type SendFundsAch struct {
	EstimatedActivity *EstimatedActivity `json:"estimatedActivity,omitempty"`
}

type SendFundsPushToCard struct {
	EstimatedActivity *EstimatedActivity `json:"estimatedActivity,omitempty"`
}

type SendFundsRtp struct {
	EstimatedActivity *EstimatedActivity `json:"estimatedActivity,omitempty"`
}

type SendFundsInstantBank struct {
	EstimatedActivity *EstimatedActivity `json:"estimatedActivity,omitempty"`
}

type SendFundsWire struct {
	EstimatedActivity *EstimatedActivity `json:"estimatedActivity,omitempty"`
}

// CardIssuing is the underwriting data for the card-issuing capability, added in
// v2026.10. Card issuing has no granular rails, so it hangs directly off the
// underwriting root rather than acting as a container like SendFunds. Issued cards are
// funded from the account's wallet before they can be spent, so estimated activity is
// the only data collected.
type CardIssuing struct {
	EstimatedActivity *EstimatedActivity `json:"estimatedActivity,omitempty"`
}

type SubmissionIntent string

const (
	SubmissionIntentWait   SubmissionIntent = "wait"
	SubmissionIntentSubmit SubmissionIntent = "submit"
)
