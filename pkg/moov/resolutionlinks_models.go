package moov

import "time"

type ResolutionLinkRecord struct {
	ResolutionLinkCode string                `json:"resolution_link_code"`
	PartnerAccountID   string                `json:"partner_account_id"`
	Mode               string                `json:"account_mode"`
	AccountID          string                `json:"account_id"`
	CreatedOn          time.Time             `json:"created_on"`
	UpdatedOn          time.Time             `json:"updated_on,omitempty"`
	DisabledOn         *time.Time            `json:"disabled_on,omitempty"`
	ExpiresOn          time.Time             `json:"expires_on"`
	Recipient          Recipient             `json:"recipient"`
	Options            ResolutionLinkOptions `json:"options"`
}

type CreateResolutionLink struct {
	PartnerAccountID string                `json:"partner_account_id,omitempty"`
	Recipient        Recipient             `json:"recipient"`
	Options          ResolutionLinkOptions `json:"options,omitempty"`
	AccountID        string
}

type Recipient struct {
	Email string `json:"email"`
	Phone *Phone `json:"phone,omitempty"`
}

type ResolutionLinkOptions struct {
	MerchantName string `json:"merchant_name,omitempty"`
	AccountName  string `json:"account_name,omitempty"`
}

type FileUploadRequest struct {
	FileContents []byte
	FileName     string `json:"file_name,omitempty"`
	Purpose      string `json:"purpose,omitempty"`
	Size         int    `json:"file_size,omitempty"`
	Metadata     string `json:"metadata,omitempty"`
}
