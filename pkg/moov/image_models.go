package moov

import (
	"fmt"
	"time"
)

// ImageMetadataRequest represents the request body for creating or updating image metadata.
type ImageMetadataRequest struct {
	AltText *string `json:"altText,omitempty"`
}

// ImageMetadata represents metadata about an uploaded image.
type ImageMetadata struct {
	ImageID    string     `json:"imageID"`
	PublicID   string     `json:"publicID"`
	AltText    *string    `json:"altText,omitempty"`
	Link       string     `json:"link"`
	CreatedOn  time.Time  `json:"createdOn"`
	UpdatedOn  time.Time  `json:"updatedOn"`
	DisabledOn *time.Time `json:"disabledOn,omitempty"`
}

// ImageListFilter is used to filter the list of images.
type ImageListFilter callArg

// WithImageSkip sets the skip parameter for pagination.
func WithImageSkip(skip int) ImageListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["skip"] = fmt.Sprintf("%d", skip)
		return nil
	})
}

// WithImageCount sets the count parameter for pagination (max 200, default 20).
func WithImageCount(count int) ImageListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["count"] = fmt.Sprintf("%d", count)
		return nil
	})
}
