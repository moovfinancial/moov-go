package moov

type UpdateUnderwriting struct {
	AverageTransactionSize          int64                  `json:"averageTransactionSize"`
	MaxTransactionSize              int64                  `json:"maxTransactionSize"`
	AverageMonthlyTransactionVolume int64                  `json:"averageMonthlyTransactionVolume"`
	VolumeByCustomerType            VolumeByCustomerType   `json:"volumeByCustomerType"`
	CardVolumeDistribution          CardVolumeDistribution `json:"cardVolumeDistribution"`
	Fulfillment                     Fulfillment            `json:"fulfillment"`
}

type Underwriting struct {
	AverageTransactionSize          int64                  `json:"averageTransactionSize"`
	MaxTransactionSize              int64                  `json:"maxTransactionSize"`
	AverageMonthlyTransactionVolume int64                  `json:"averageMonthlyTransactionVolume"`
	Status                          UnderwritingStatus     `json:"status"`
	VolumeByCustomerType            VolumeByCustomerType   `json:"volumeByCustomerType"`
	CardVolumeDistribution          CardVolumeDistribution `json:"cardVolumeDistribution"`
	Fulfillment                     Fulfillment            `json:"fulfillment"`
}

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
