package mv2610

type CaptureStatus string

// List of CaptureStatus
const (
	CaptureStatus_Created   CaptureStatus = "created"
	CaptureStatus_Pending   CaptureStatus = "pending"
	CaptureStatus_Completed CaptureStatus = "completed"
	CaptureStatus_Failed    CaptureStatus = "failed"
	CaptureStatus_Canceled  CaptureStatus = "canceled"
)
