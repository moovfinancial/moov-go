package mv2610

// TransferAuthorizationStatus is the status of a card authorization on a transfer.
type TransferAuthorizationStatus string

// List of TransferAuthorizationStatus
const (
	TransferAuthorizationStatus_Approved TransferAuthorizationStatus = "approved"
	TransferAuthorizationStatus_Declined TransferAuthorizationStatus = "declined"
	TransferAuthorizationStatus_Reversed TransferAuthorizationStatus = "reversed"
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
