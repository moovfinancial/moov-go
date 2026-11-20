package mv2610

// AuthorizationStatus is the status of a card authorization on a transfer.
type AuthorizationStatus string

// List of AuthorizationStatus
const (
	AuthorizationStatus_Approved AuthorizationStatus = "approved"
	AuthorizationStatus_Declined AuthorizationStatus = "declined"
	AuthorizationStatus_Reversed AuthorizationStatus = "reversed"
)

type CaptureStatus string

// List of CaptureStatus
const (
	CaptureStatus_Created   CaptureStatus = "created"
	CaptureStatus_Pending   CaptureStatus = "pending"
	CaptureStatus_Completed CaptureStatus = "completed"
	CaptureStatus_Failed    CaptureStatus = "failed"
	CaptureStatus_Canceled  CaptureStatus = "canceled"
)
