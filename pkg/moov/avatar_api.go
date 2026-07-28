package moov

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
)

// GetAvatar retrieves an avatar image for an account and returns the raw image
// bytes along with the image's content type (e.g. "image/png"). Unlike the
// write endpoints, uniqueID accepts any unique ID associated with an account
// such as accountID, representativeID, routing number, or userID. The response
// is a user-uploaded avatar if one exists, otherwise the enriched avatar or an
// account-type-aware fallback icon.
// https://docs.moov.io/api/enrichment/form-shortening/avatars/get/
func (c Client) GetAvatar(ctx context.Context, uniqueID string) ([]byte, string, error) {
	if uniqueID == "" {
		return nil, "", errors.New("uniqueID cannot be empty")
	}

	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathAvatar, uniqueID),
		MoovVersion(Version2026_07),
		AcceptContentType("image/*"))
	if err != nil {
		return nil, "", err
	}

	buf, err := CompletedObjectOrError[bytes.Buffer](resp)
	if err != nil {
		return nil, "", err
	}

	contentType := ""
	if hcr, ok := resp.(*httpCallResponse); ok {
		contentType = hcr.ContentType()
	}
	return buf.Bytes(), contentType, nil
}

// UploadAvatar uploads a user avatar image for an account.
// The accountID is used as the uniqueID; only accountID values are accepted for
// writes. The image must be a JPEG, PNG, or WebP file and will be normalized to
// a 512x512 PNG. User-uploaded avatars take precedence over enriched avatars at
// read time.
// https://docs.moov.io/api/enrichment/form-shortening/avatars/put/
func (c Client) UploadAvatar(ctx context.Context, accountID string, file io.Reader) error {
	if accountID == "" {
		return errors.New("accountID cannot be empty")
	}
	if file == nil {
		return errors.New("file cannot be nil")
	}

	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPut, pathAvatar, accountID),
		MoovVersion(Version2026_07),
		MultipartBody(MultipartFile("file", "avatar", file, "application/octet-stream")))
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}

// DeleteAvatar deletes a user-uploaded avatar for an account.
// The accountID is used as the uniqueID; only accountID values are accepted for
// writes. After deletion, the avatar endpoint falls back to the enriched avatar
// or an account-type-aware fallback icon.
// https://docs.moov.io/api/enrichment/form-shortening/avatars/delete/
func (c Client) DeleteAvatar(ctx context.Context, accountID string) error {
	if accountID == "" {
		return errors.New("accountID cannot be empty")
	}

	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodDelete, pathAvatar, accountID),
		MoovVersion(Version2026_07),
		AcceptJson())
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}
