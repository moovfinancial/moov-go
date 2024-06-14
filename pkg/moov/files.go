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
	FileID         string      `json:"fileID,omitempty"`
	FileName       string      `json:"fileName,omitempty"`
	FilePurpose    FilePurpose `json:"filePurpose,omitempty"`
	FileStatus     FileStatus  `json:"fileStatus,omitempty"`
	DecisionReason *string     `json:"decisionReason,omitempty"`
	Size           int         `json:"fileSizeBytes,omitempty"`
	Metadata       string      `json:"metadata,omitempty"`
	AccountID      string      `json:"accountID,omitempty"`
	CreatedOn      time.Time   `json:"createdOn,omitempty"`
	UpdatedOn      time.Time   `json:"updatedOn,omitempty"`
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
			MultipartFile("file", upload.Filename, upload.File),
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
