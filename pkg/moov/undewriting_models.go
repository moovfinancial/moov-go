package moov

type UpdateUnderwriting struct {
	AverageTransactionSize          int64 `json:"averageTransactionSize"`
	MaxTransactionSize              int64 `json:"maxTransactionSize"`
	AverageMonthlyTransactionVolume int64 `json:"averageMonthlyTransactionVolume"`
}

type Underwriting struct {
	AverageTransactionSize          int64  `json:"averageTransactionSize"`
	MaxTransactionSize              int64  `json:"maxTransactionSize"`
	AverageMonthlyTransactionVolume int64  `json:"averageMonthlyTransactionVolume"`
	Status                          Status `json:"status"`
}

type Status string

const (
	UnderwritingStatusApproved      Status = "approved"
	UnderwritingStatusRejected      Status = "rejected"
	UnderwritingStatusPendingReview Status = "pendingReview"
	UnderwritingStatusPending       Status = "pending"
	UnderwritingStatusNotRequested  Status = "notRequested"
)
