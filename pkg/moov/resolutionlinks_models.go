package moov

import "time"

type ResolutionLinkRecord struct {
	ResolutionLinkCode string     `json:"resolution_link_code"`
	PartnerAccountID   string     `json:"partner_account_id"`
	Mode               string     `json:"account_mode"`
	AccountID          string     `json:"account_id"`
	CreatedOn          time.Time  `json:"created_on"`
	UpdatedOn          time.Time  `json:"updated_on,omitempty"`
	DisabledOn         *time.Time `json:"disabled_on,omitempty"`
	ExpiresOn          time.Time  `json:"expires_on"`
	Recipient          string
	Options            ResolutionLinkOptions `json:"options"`
}

type CreateResolutionLink struct {
	PartnerAccountID string
	AccountID        string
	Recipient        Recipient
	Options          ResolutionLinkOptions
}

type Recipient struct {
	Email string
	Phone *Phone `json:"phone,omitempty"`
}

type ResolutionLinkOptions struct {
	MerchantName string `json:"merchantName,omitempty" otel:"merchant_name"`
	AccountName  string `json:"accountName,omitempty" otel:"account_name,omitempty"`
}

type FileUploadRequest struct {
	FileContents []byte
	FileName     string `otel:"file_name, omitempty"`
	Purpose      string `otel:"purpose, omitempty"`
	Size         int    `otel:"file_size, omitempty"`
	Metadata     string `otel:"metadata, omitempty"`
}
