package moov

import (
	"fmt"
	"time"
)

// Product represents a good or service offered by a merchant.
type Product struct {
	ProductID    string                 `json:"productID"`
	Title        string                 `json:"title"`
	Description  *string                `json:"description,omitempty"`
	BasePrice    AmountDecimal          `json:"basePrice"`
	OptionGroups []ProductOptionGroup   `json:"optionGroups,omitempty"`
	Images       []ProductImageMetadata `json:"images,omitempty"`
	CreatedOn    time.Time              `json:"createdOn"`
	UpdatedOn    time.Time              `json:"updatedOn"`
	DisabledOn   *time.Time             `json:"disabledOn,omitempty"`
}

// ProductOptionGroup represents a group of product configuration options.
type ProductOptionGroup struct {
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	MinSelect   int32           `json:"minSelect"`
	MaxSelect   int32           `json:"maxSelect"`
	Options     []ProductOption `json:"options"`
}

// ProductOption represents a single product option within a group.
type ProductOption struct {
	Name          string                 `json:"name"`
	Description   *string                `json:"description,omitempty"`
	PriceModifier *AmountDecimal         `json:"priceModifier,omitempty"`
	Images        []ProductImageMetadata `json:"images,omitempty"`
}

// ProductImageMetadata represents metadata about an image associated with a product or option.
type ProductImageMetadata struct {
	ImageID  string  `json:"imageID"`
	AltText  *string `json:"altText,omitempty"`
	Link     string  `json:"link"`
	PublicID string  `json:"publicID"`
}

// ProductRequest represents a request to create or update a product.
type ProductRequest struct {
	Title        string                     `json:"title"`
	Description  *string                    `json:"description,omitempty"`
	BasePrice    AmountDecimal              `json:"basePrice"`
	Images       []AssignProductImage       `json:"images,omitempty"`
	OptionGroups []CreateProductOptionGroup `json:"optionGroups,omitempty"`
}

// CreateProductOptionGroup represents a product option group to be created.
type CreateProductOptionGroup struct {
	Name        string                `json:"name"`
	Description *string               `json:"description,omitempty"`
	MinSelect   int32                 `json:"minSelect"`
	MaxSelect   int32                 `json:"maxSelect"`
	Options     []CreateProductOption `json:"options"`
}

// CreateProductOption represents a product option to be created.
type CreateProductOption struct {
	Name          string               `json:"name"`
	Description   *string              `json:"description,omitempty"`
	PriceModifier *AmountDecimal       `json:"priceModifier,omitempty"`
	Images        []AssignProductImage `json:"images,omitempty"`
}

// AssignProductImage assigns a previously uploaded image to a product or option.
type AssignProductImage struct {
	ImageID string `json:"imageID"`
}

// ProductListFilter is used to filter the list of products.
type ProductListFilter callArg

// WithProductSkip sets the skip parameter for pagination.
func WithProductSkip(skip int) ProductListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["skip"] = fmt.Sprintf("%d", skip)
		return nil
	})
}

// WithProductCount sets the count parameter for pagination (max 200, default 20).
func WithProductCount(count int) ProductListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["count"] = fmt.Sprintf("%d", count)
		return nil
	})
}
