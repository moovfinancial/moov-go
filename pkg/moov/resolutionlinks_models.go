package moov

import "time"

type ResolutionLinkStatus string

const (
	ResolutionLinkStatus_Active    ResolutionLinkStatus = "active"
	ResolutionLinkStatus_Submitted ResolutionLinkStatus = "submitted"
	ResolutionLinkStatus_Completed ResolutionLinkStatus = "completed"
	ResolutionLinkStatus_Disabled  ResolutionLinkStatus = "disabled"
	ResolutionLinkStatus_Expired   ResolutionLinkStatus = "expired"
)

type ResolutionLinkResponse struct {
	ResolutionLinkCode string               `json:"code"`
	PartnerAccountID   string               `json:"partnerAccountID"`
	AccountID          string               `json:"accountID"`
	CreatedOn          time.Time            `json:"createdOn"`
	UpdatedOn          time.Time            `json:"updatedOn"`
	ExpiresOn          time.Time            `json:"expiresOn"`
	DisabledOn         *time.Time           `json:"disabledOn,omitempty"`
	Recipient          string               `json:"recipient"`
	URL                string               `json:"url"`
	Status             ResolutionLinkStatus `json:"status"`
}

type CreateResolutionLinkRequest struct {
	Recipient Recipient `json:"recipient"`
}

type Recipient struct {
	Email string `json:"email"`
	Phone *Phone `json:"phone,omitempty"`
}
