package moov

import "time"

type ResolutionLinkResponse struct {
	ResolutionLinkCode string     `json:"code"`
	PartnerAccountID   string     `json:"partnerAccountID"`
	AccountID          string     `json:"accountID"`
	CreatedOn          time.Time  `json:"createdOn"`
	UpdatedOn          time.Time  `json:"updatedOn"`
	ExpiresOn          time.Time  `json:"expiresOn"`
	DisabledOn         *time.Time `json:"disabledOn,omitempty"`
	Recipient          string     `json:"recipient"`
	URL                string     `json:"url"`
}

type CreateResolutionLinkRequest struct {
	Recipient Recipient `json:"recipient"`
}

type Recipient struct {
	Email string `json:"email"`
	Phone *Phone `json:"phone,omitempty"`
}
