package moov

type UpdateUnderwriting struct {
	AverageTransactionSize          int64 `json:"averageTransactionSize"`
	MaxTransactionSize              int64 `json:"maxTransactionSize"`
	AverageMonthlyTransactionVolume int64 `json:"averageMonthlyTransactionVolume"`
}

type Underwriting struct {
	AverageTransactionSize          int64              `json:"averageTransactionSize"`
	MaxTransactionSize              int64              `json:"maxTransactionSize"`
	AverageMonthlyTransactionVolume int64              `json:"averageMonthlyTransactionVolume"`
	Status                          UnderwritingStatus `json:"status"`
}

type UnderwritingStatus string

const (
	UnderwritingStatusApproved      UnderwritingStatus = "approved"
	UnderwritingStatusRejected      UnderwritingStatus = "rejected"
	UnderwritingStatusPendingReview UnderwritingStatus = "pendingReview"
	UnderwritingStatusPending       UnderwritingStatus = "pending"
	UnderwritingStatusNotRequested  UnderwritingStatus = "notRequested"
)
