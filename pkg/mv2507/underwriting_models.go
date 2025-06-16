package mv2507

import "github.com/moovfinancial/moov-go/pkg/moov"

var Underwriting moov.UnderwritingClient[UpsertUnderwriting, UnderwritingResp] = moov.UnderwritingClient[UpsertUnderwriting, UnderwritingResp]{Version: moov.Q3_2025}

type UpsertUnderwriting struct {
	GeographicReach           *GeographicReach           `json:"geographicReach,omitempty"`
	BusinessPresence          *BusinessPresence          `json:"businessPresence,omitempty"`
	PendingLitigation         *PendingLitigation         `json:"pendingLitigation,omitempty"`
	VolumeShareByCustomerType *VolumeShareByCustomerType `json:"volumeShareByCustomerType,omitempty"`
	CollectFunds              *CollectFunds              `json:"collectFunds,omitempty"`
	MoneyTransfer             *MoneyTransfer             `json:"moneyTransfer,omitempty"`
	SendFunds                 *SendFunds                 `json:"sendFunds,omitempty"`
}

// UnderwritingV2507 will have both legacy and V2507 fields for backward compatibility
type UnderwritingResp struct {
	// Legacy fields
	// TEMP: TODO: Status will be removed once we release to prod
	Status                          UnderwritingStatus     `json:"status"`
	AverageTransactionSize          int64                  `json:"averageTransactionSize"`
	MaxTransactionSize              int64                  `json:"maxTransactionSize"`
	AverageMonthlyTransactionVolume int64                  `json:"averageMonthlyTransactionVolume"`
	VolumeByCustomerType            VolumeByCustomerType   `json:"volumeByCustomerType"`
	CardVolumeDistribution          CardVolumeDistribution `json:"cardVolumeDistribution"`
	Fulfillment                     Fulfillment            `json:"fulfillment"`

	//V2507
	GeographicReach           *GeographicReach           `json:"geographicReach,omitempty"`
	BusinessPresence          *BusinessPresence          `json:"businessPresence,omitempty"`
	PendingLitigation         *PendingLitigation         `json:"pendingLitigation,omitempty"`
	VolumeShareByCustomerType *VolumeShareByCustomerType `json:"volumeShareByCustomerType,omitempty"`
	CollectFunds              *CollectFunds              `json:"collectFunds,omitempty"`
	MoneyTransfer             *MoneyTransfer             `json:"moneyTransfer,omitempty"`
	SendFunds                 *SendFunds                 `json:"sendFunds,omitempty"`
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

const FulfillmentMethodBillOrDebtPayment FulfillmentMethod = "bill-or-debt-payment"
const FulfillmentMethodDigitalContent FulfillmentMethod = "digital-content"
const FulfillmentMethodDonation FulfillmentMethod = "donation"
const FulfillmentMethodInPersonService FulfillmentMethod = "in-person-service"
const FulfillmentMethodLocalPickupOrDelivery FulfillmentMethod = "local-pickup-or-delivery"
const FulfillmentMethodOther FulfillmentMethod = "other"
const FulfillmentMethodRemoteService FulfillmentMethod = "remote-service"
const FulfillmentMethodShippedPhysicalGoods FulfillmentMethod = "shipped-physical-goods"
const FulfillmentMethodSubscriptionOrMembership FulfillmentMethod = "subscription-or-membership"

type FulfillmentTimeframe string

const FulfillmentTimeframeImmediate FulfillmentTimeframe = "immediate"
const FulfillmentTimeframeOther FulfillmentTimeframe = "other"
const FulfillmentTimeframeOver30Days FulfillmentTimeframe = "over-30-days"
const FulfillmentTimeframePreOrder FulfillmentTimeframe = "pre-order"
const FulfillmentTimeframeRecurringSchedule FulfillmentTimeframe = "recurring-schedule"
const FulfillmentTimeframeScheduledEvent FulfillmentTimeframe = "scheduled-event"
const FulfillmentTimeframeWithin30Days FulfillmentTimeframe = "within-30-days"
const FulfillmentTimeframeWithin7Days FulfillmentTimeframe = "within-7-days"

type RefundPolicy string

const RefundPolicyConditionalRefund RefundPolicy = "conditional-refund"
const RefundPolicyCustomPolicy RefundPolicy = "custom-policy"
const RefundPolicyEventBasedPolicy RefundPolicy = "event-based-policy"
const RefundPolicyFullRefundExtendedWindow RefundPolicy = "full-refund-extended-window"
const RefundPolicyFullRefundWithin30Days RefundPolicy = "full-refund-within-30-days"
const RefundPolicyNoRefunds RefundPolicy = "no-refunds"
const RefundPolicyPartialRefund RefundPolicy = "partial-refund"
const RefundPolicyProratedRefund RefundPolicy = "prorated-refund"
const RefundPolicyStoreCreditOnly RefundPolicy = "store-credit-only"

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
	Ach        *SendFundsAch        `json:"ach,omitempty"`
	PushToCard *SendFundsPushToCard `json:"pushToCard,omitempty"`
	Rtp        *SendFundsRtp        `json:"rtp,omitempty"`
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

// v1
type UnderwritingStatus string

const (
	UnderwritingStatusApproved      UnderwritingStatus = "approved"
	UnderwritingStatusRejected      UnderwritingStatus = "rejected"
	UnderwritingStatusPendingReview UnderwritingStatus = "pendingReview"
	UnderwritingStatusPending       UnderwritingStatus = "pending"
	UnderwritingStatusNotRequested  UnderwritingStatus = "notRequested"
)

type VolumeByCustomerType struct {
	BusinessToBusinessPercentage int32 `json:"businessToBusinessPercentage"`
	ConsumerToBusinessPercentage int32 `json:"consumerToBusinessPercentage"`
}

type CardVolumeDistribution struct {
	EcommercePercentage     int32 `json:"ecommercePercentage"`
	CardPresentPercentage   int32 `json:"cardPresentPercentage"`
	MailOrPhonePercentage   int32 `json:"mailOrPhonePercentage"`
	DebtRepaymentPercentage int32 `json:"debtRepaymentPercentage"`
}

type ReturnPolicy string

const (
	NONE               ReturnPolicy = "none"
	EXCHANGE_ONLY      ReturnPolicy = "exchangeOnly"
	WITHIN_THIRTY_DAYS ReturnPolicy = "withinThirtyDays"
	OTHER              ReturnPolicy = "other"
)

type Fulfillment struct {
	HasPhysicalGoods     bool         `json:"hasPhysicalGoods"`
	IsShippingProduct    bool         `json:"isShippingProduct"`
	ShipmentDurationDays int32        `json:"shipmentDurationDays"`
	ReturnPolicy         ReturnPolicy `json:"returnPolicy"`
}
