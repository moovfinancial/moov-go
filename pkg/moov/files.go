package moov

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type FilePurpose string

const (
	FilePurpose_IdentityVerification       FilePurpose = "identity_verification"
	FilePurpose_BusinessVerification       FilePurpose = "business_verification"
	FilePurpose_RepresentativeVerification FilePurpose = "representative_verification"
	FilePurpose_IndividualVerification     FilePurpose = "individual_verification"
	FilePurpose_MerchantUnderwriting       FilePurpose = "merchant_underwriting"
	FilePurpose_AccountRequirement         FilePurpose = "account_requirement"
)

type UploadFile struct {
	FilePurpose FilePurpose
	Metadata    map[string]string

	Filename string
	File     io.Reader
}

type FileStatus string

const (
	FileStatus_Pending  FileStatus = "pending"
	FileStatus_Approved FileStatus = "approved"
	FileStatus_Rejected FileStatus = "rejected"
)

type File struct {
	FileID         string      `json:"fileID"`
	FileName       string      `json:"fileName"`
	FilePurpose    FilePurpose `json:"filePurpose"`
	FileStatus     FileStatus  `json:"fileStatus"`
	DecisionReason *string     `json:"decisionReason"`
	Size           int         `json:"fileSizeBytes"`
	Metadata       string      `json:"metadata"`
	AccountID      string      `json:"accountID"`
	CreatedOn      time.Time   `json:"createdOn"`
	UpdatedOn      time.Time   `json:"updatedOn"`
}

func (c Client) UploadFile(ctx context.Context, accountID string, upload UploadFile) (*File, error) {
	mdJson, err := json.Marshal(upload.Metadata)
	if err != nil {
		return nil, err
	}

	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathFiles, accountID),
		AcceptJson(),
		MultipartBody(
			MultipartField("filePurpose", string(upload.FilePurpose)),
			MultipartField("metadata", string(mdJson)),
			MultipartFile("file", upload.Filename, upload.File, "application/octet-stream"),
		))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[File](resp)
}

func (c Client) ListFiles(ctx context.Context, accountID string) ([]File, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathFiles, accountID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[File](resp)
}

func (c Client) GetFile(ctx context.Context, accountID string, fileID string) (*File, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathFile, accountID, fileID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[File](resp)
}
