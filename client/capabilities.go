package moov

import "time"

const (
	CAPABILITY_TRANSFERS     = "transfers"
	CAPABILITY_SEND_FUNDS    = "send-funds"
	CAPABILITY_COLLECT_FUNDS = "collect-funds"
	CAPABILITY_WALLET        = "wallet"
	CAPABILITY_CARD_ISSUING  = "card-issuing"
	CAPABILITY_ENBABLED      = "enabled"
	CAPABILITY_DISABLED      = "disabled"
	CAPABILITY_PENDING       = "pending"
)

// Capabilities a list of CAPABILITY_*
var Capabilities []string

type Captability struct {
	Capability   string `json:"capability"`
	AccountID    string `json:"accountID"`
	Status       string `json:"status,omitempty"`
	Requirements struct {
		CurrentlyDue []string `json:"currentlyDue,omitempty"`
		Errors       []struct {
			Requirement string `json:"requirement,omitempty"`
			ErrorCode   string `json:"errorCode,omitempty"`
		}
	}
	DisabledReason string    `json:"disabledReason,omitempty"`
	CreatedOn      time.Time `json:"createdOn,omitempty"`
	UpdatedOn      time.Time `json:"updatedOn,omitempty"`
	DisabledOn     time.Time `json:"disabledOn,omitempty"`
}
