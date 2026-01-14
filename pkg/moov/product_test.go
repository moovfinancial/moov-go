package moov_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductMarshal(t *testing.T) {
	input := []byte(`{
		"productID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43",
		"title": "Test Product",
		"description": "A test product description",
		"basePrice": {
			"currency": "USD",
			"valueDecimal": "19.99"
		},
		"optionGroups": [
			{
				"name": "Size",
				"description": "Choose a size",
				"minSelect": 1,
				"maxSelect": 1,
				"options": [
					{
						"name": "Small",
						"description": "Small size",
						"priceModifier": {
							"currency": "USD",
							"valueDecimal": "0.00"
						}
					},
					{
						"name": "Large",
						"priceModifier": {
							"currency": "USD",
							"valueDecimal": "5.00"
						}
					}
				]
			}
		],
		"images": [
			{
				"imageID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43",
				"altText": "Product image",
				"link": "https://api.moov.io/images/qJRAaAwwF5hmfeAFdHjIb",
				"publicID": "qJRAaAwwF5hmfeAFdHjIb"
			}
		],
		"createdOn": "2024-01-15T10:30:00Z",
		"updatedOn": "2024-01-15T10:30:00Z",
		"disabledOn": null
	}`)

	createdOn, _ := time.Parse(time.RFC3339, "2024-01-15T10:30:00Z")
	updatedOn, _ := time.Parse(time.RFC3339, "2024-01-15T10:30:00Z")

	want := moov.Product{
		ProductID:   "ec7e1848-dc80-4ab0-8827-dd7fc0737b43",
		Title:       "Test Product",
		Description: moov.PtrOf("A test product description"),
		BasePrice: moov.AmountDecimal{
			Currency:     "USD",
			ValueDecimal: "19.99",
		},
		OptionGroups: []moov.ProductOptionGroup{
			{
				Name:        "Size",
				Description: moov.PtrOf("Choose a size"),
				MinSelect:   1,
				MaxSelect:   1,
				Options: []moov.ProductOption{
					{
						Name:        "Small",
						Description: moov.PtrOf("Small size"),
						PriceModifier: &moov.AmountDecimal{
							Currency:     "USD",
							ValueDecimal: "0.00",
						},
					},
					{
						Name: "Large",
						PriceModifier: &moov.AmountDecimal{
							Currency:     "USD",
							ValueDecimal: "5.00",
						},
					},
				},
			},
		},
		Images: []moov.ProductImageMetadata{
			{
				ImageID:  "ec7e1848-dc80-4ab0-8827-dd7fc0737b43",
				AltText:  moov.PtrOf("Product image"),
				Link:     "https://api.moov.io/images/qJRAaAwwF5hmfeAFdHjIb",
				PublicID: "qJRAaAwwF5hmfeAFdHjIb",
			},
		},
		CreatedOn:  createdOn,
		UpdatedOn:  updatedOn,
		DisabledOn: nil,
	}

	got := new(moov.Product)

	dec := json.NewDecoder(bytes.NewReader(input))
	dec.DisallowUnknownFields()

	err := dec.Decode(&got)
	require.NoError(t, err)
	require.Equal(t, want, *got)
}

func Test_Products(t *testing.T) {
	mc := NewTestClient(t)
	ctx := context.Background()
	accountID := MERCHANT_ID

	// upload image for use with product test
	_, imgReader := randomImage(t, 100, 100, encodePNG)
	metadata := &moov.ImageMetadataRequest{
		AltText: moov.PtrOf("Test product image"),
	}
	image, err := mc.UploadImage(ctx, accountID, imgReader, metadata)
	require.NoError(t, err)
	require.NotNil(t, image)
	t.Cleanup(func() {
		_ = mc.DeleteImage(ctx, accountID, image.ImageID)
	})

	var createdProduct moov.Product

	t.Run("create product", func(t *testing.T) {
		product := moov.ProductRequest{
			Title:       "SDK Test Product",
			Description: moov.PtrOf("A product with size options"),
			BasePrice: moov.AmountDecimal{
				Currency:     "USD",
				ValueDecimal: "19.99",
			},
			Images: []moov.AssignProductImage{
				{ImageID: image.ImageID},
			},
			OptionGroups: []moov.CreateProductOptionGroup{
				{
					Name:        "Size",
					Description: moov.PtrOf("Select a size"),
					MinSelect:   1,
					MaxSelect:   1,
					Options: []moov.CreateProductOption{
						{
							Name:        "Small",
							Description: moov.PtrOf("Small size"),
							PriceModifier: &moov.AmountDecimal{
								Currency:     "USD",
								ValueDecimal: "4.00",
							},
							Images: []moov.AssignProductImage{{ImageID: image.ImageID}},
						},
						{
							Name: "Large",
							PriceModifier: &moov.AmountDecimal{
								Currency:     "USD",
								ValueDecimal: "5.00",
							},
						},
					},
				},
			},
		}

		// Products may take a moment to be aware of an uploaded image
		var created *moov.Product
		success := assert.Eventually(t, func() bool {
			var err error
			created, err = mc.CreateProduct(ctx, accountID, product)
			return err == nil
		}, 5*time.Second, 250*time.Millisecond)
		require.True(t, success, "failed to create product")
		require.NotNil(t, created)
		createdProduct = *created
	})

	t.Run("list products", func(t *testing.T) {
		products, err := mc.ListProducts(ctx, accountID)
		require.NoError(t, err)
		require.NotNil(t, products)
		var found bool
		for _, p := range products {
			if p.ProductID == createdProduct.ProductID {
				found = true
				break
			}
		}
		require.True(t, found, "created product not found in list")
	})

	t.Run("get product", func(t *testing.T) {
		got, err := mc.GetProduct(ctx, accountID, createdProduct.ProductID)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, createdProduct, *got)
	})

	t.Run("update product", func(t *testing.T) {
		product := moov.ProductRequest{
			Title:       "SDK Test Product - Updated",
			Description: moov.PtrOf("Updated description"),
			BasePrice: moov.AmountDecimal{
				Currency:     "USD",
				ValueDecimal: "39.99",
			},
		}

		updated, err := mc.UpdateProduct(ctx, accountID, createdProduct.ProductID, product)
		require.NoError(t, err)
		require.NotNil(t, updated)
		require.Equal(t, product.Title, updated.Title)
		require.Equal(t, product.Description, updated.Description)
		require.Equal(t, product.BasePrice.ValueDecimal, updated.BasePrice.ValueDecimal)
		require.Empty(t, updated.Images)
		require.Empty(t, updated.OptionGroups)
	})

	t.Run("disable product", func(t *testing.T) {
		err := mc.DisableProduct(ctx, accountID, createdProduct.ProductID)
		require.NoError(t, err)
	})
}
