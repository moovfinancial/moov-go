package moov

import (
	"context"
	"io"
	"net/http"
)

// UploadAvatar uploads a user avatar image for an account.
// The accountID is used as the uniqueID; only accountID values are accepted for
// writes. The image must be a JPEG, PNG, or WebP file and will be normalized to
// a 512x512 PNG. User-uploaded avatars take precedence over enriched avatars at
// read time.
// https://docs.moov.io/api/enrichment/form-shortening/avatars/put/
func (c Client) UploadAvatar(ctx context.Context, accountID string, file io.Reader) error {
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
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodDelete, pathAvatar, accountID),
		MoovVersion(Version2026_07),
		AcceptJson())
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}
