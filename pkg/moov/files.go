package moov

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime"
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

// FileContents is the raw content of a file along with what the response
// reports about it.
type FileContents struct {
	// FileName is the name the file was uploaded under, read from the
	// Content-Disposition header. Empty if that header is missing or malformed.
	FileName string

	// ContentType is detected from the file's contents at upload time rather
	// than taken from the uploader, so a csv arrives as
	// "text/plain; charset=utf-8" instead of "text/csv".
	ContentType string

	// Data is the raw file content.
	Data []byte
}

// DownloadFile retrieves the contents of a file linked to a Moov account.
// Files reserved for internal Moov use are not returned by this endpoint.
//
// Requires the /accounts/{accountID}/files.download scope, which is granted per
// partner rather than included in the default connection scopes.
// https://docs.moov.io/api/moov-accounts/files/download/
func (c Client) DownloadFile(ctx context.Context, accountID string, fileID string) (*FileContents, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathFileContents, accountID, fileID),
		MoovVersion(Version2026_10),
		AcceptContentType("*/*"))
	if err != nil {
		return nil, err
	}

	buf, err := CompletedObjectOrError[bytes.Buffer](resp)
	if err != nil {
		return nil, err
	}

	contents := &FileContents{Data: buf.Bytes()}

	// CallResponse exposes no header access, so reach the two headers through
	// the concrete type, as GetAvatar does. CallHttp only ever returns
	// *httpCallResponse -- the sole CallResponse implementation -- so the
	// assertion holds and ContentType and FileName are always populated here.
	hcr, ok := resp.(*httpCallResponse)
	if !ok {
		return contents, nil
	}
	contents.ContentType = hcr.ContentType()

	// File names are only length checked on upload, so the API quotes and
	// escapes them. Parse the header rather than splitting on "filename=".
	if _, params, err := mime.ParseMediaType(hcr.ContentDisposition()); err == nil {
		contents.FileName = params["filename"]
	}

	return contents, nil
}
