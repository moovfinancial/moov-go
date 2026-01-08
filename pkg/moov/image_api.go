package moov

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// UploadImage uploads a new PNG, JPEG, or WebP image with optional metadata.
// Duplicate images, and requests larger than 16MB will be rejected.
// https://docs.moov.io/api/tools/images/post/
func (c Client) UploadImage(ctx context.Context, accountID string, file io.Reader, metadata *ImageMetadataRequest) (*ImageMetadata, error) {
	var multiParts []multipartFn
	multiParts = append(multiParts, MultipartBinary("image", file))

	if metadata != nil {
		mdJson, err := json.Marshal(metadata)
		if err != nil {
			return nil, err
		}
		multiParts = append(multiParts, MultipartField("metadata", string(mdJson)))
	}

	resp, err := c.CallHttp(ctx, Endpoint(http.MethodPost, pathImages, accountID), MultipartBody(multiParts...))
	if err != nil {
		return nil, err
	}

	return StartedObjectOrError[ImageMetadata](resp)
}

// ListImageMetadata lists metadata for all images in the specified account.
// https://docs.moov.io/api/tools/images/list/
func (c Client) ListImageMetadata(ctx context.Context, accountID string, filters ...ImageListFilter) ([]ImageMetadata, error) {
	args := prependArgs(filters, AcceptJson())
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathImages, accountID), args...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[ImageMetadata](resp)
}

// GetImageMetadata retrieves metadata for a specific image by its ID.
// https://docs.moov.io/api/tools/images/get/
func (c Client) GetImageMetadata(ctx context.Context, accountID string, imageID string) (*ImageMetadata, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathImage, accountID, imageID), AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[ImageMetadata](resp)
}

// UpdateImage replaces an existing image and, optionally, its metadata.
// This endpoint replaces the existing image with the new PNG, JPEG, or WebP. Omit
// the metadata parameter to keep existing metadata. Duplicate images, and requests
// larger than 16MB will be rejected.
// https://docs.moov.io/api/tools/images/put-image/
func (c Client) UpdateImage(ctx context.Context, accountID string, imageID string, file io.Reader, metadata *ImageMetadataRequest) (*ImageMetadata, error) {
	var multiParts []multipartFn
	multiParts = append(multiParts, MultipartBinary("image", file))

	if metadata != nil {
		mdJson, err := json.Marshal(metadata)
		if err != nil {
			return nil, err
		}
		multiParts = append(multiParts, MultipartField("metadata", string(mdJson)))
	}

	resp, err := c.CallHttp(ctx, Endpoint(http.MethodPut, pathImage, accountID, imageID), MultipartBody(multiParts...))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[ImageMetadata](resp)
}

// UpdateImageMetadata replaces the metadata for an existing image.
// https://docs.moov.io/api/tools/images/put-metadata/
func (c Client) UpdateImageMetadata(ctx context.Context, accountID string, imageID string, metadata ImageMetadataRequest) (*ImageMetadata, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodPut, pathImageMetadata, accountID, imageID), AcceptJson(), JsonBody(metadata))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[ImageMetadata](resp)
}

// DeleteImage disables an image by its ID.
// Disabled images are still accessible via their public URL, and cannot be assigned
// to products or line-items.
// https://docs.moov.io/api/tools/images/delete/
func (c Client) DeleteImage(ctx context.Context, accountID string, imageID string) error {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodDelete, pathImage, accountID, imageID), AcceptJson())
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}
