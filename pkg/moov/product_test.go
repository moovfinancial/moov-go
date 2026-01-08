package moov_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
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

	product := new(moov.Product)

	dec := json.NewDecoder(bytes.NewReader(input))
	dec.DisallowUnknownFields()

	err := dec.Decode(&product)
	require.NoError(t, err)
	require.Equal(t, "ec7e1848-dc80-4ab0-8827-dd7fc0737b43", product.ProductID)
	require.Equal(t, "Test Product", product.Title)
	require.Equal(t, "A test product description", *product.Description)
	require.Equal(t, "USD", product.BasePrice.Currency)
	require.Equal(t, "19.99", product.BasePrice.ValueDecimal)
	require.Len(t, product.OptionGroups, 1)
	require.Equal(t, "Size", product.OptionGroups[0].Name)
	require.Equal(t, "Choose a size", *product.OptionGroups[0].Description)
	require.Equal(t, int32(1), product.OptionGroups[0].MinSelect)
	require.Equal(t, int32(1), product.OptionGroups[0].MaxSelect)
	require.Len(t, product.OptionGroups[0].Options, 2)
	require.Equal(t, "Small", product.OptionGroups[0].Options[0].Name)
	require.Equal(t, "Small size", *product.OptionGroups[0].Options[0].Description)
	require.Equal(t, "Large", product.OptionGroups[0].Options[1].Name)
	require.Nil(t, product.OptionGroups[0].Options[1].Description)
	require.Len(t, product.Images, 1)
	require.Equal(t, "ec7e1848-dc80-4ab0-8827-dd7fc0737b43", product.Images[0].ImageID)
	require.Equal(t, "Product image", *product.Images[0].AltText)
	require.NotNil(t, product.CreatedOn)
	require.NotNil(t, product.UpdatedOn)
	require.Nil(t, product.DisabledOn)
}

func Test_Products(t *testing.T) {
	mc := NewTestClient(t)
	ctx := context.Background()
	accountID := MERCHANT_ID

	// upload image for use with product test
	_, imgReader := randomImage(t, 100, 100, encodePNG)
	metadata := &moov.ImageMetadataRequest{
		AltText: "Test product image",
	}
	image, err := mc.UploadImage(ctx, accountID, imgReader, metadata)
	require.NoError(t, err)
	require.NotNil(t, image)

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

		created, err := mc.CreateProduct(ctx, accountID, product)
		require.NoError(t, err)
		require.NotNil(t, created)

		createdProduct = *created
	})

	t.Run("list products", func(t *testing.T) {
		products, err := mc.ListProducts(ctx, accountID)
		require.NoError(t, err)
		require.NotNil(t, products)
		require.Greater(t, len(products), 0)
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
	})

	t.Run("disable product", func(t *testing.T) {
		err := mc.DisableProduct(ctx, accountID, createdProduct.ProductID)
		require.NoError(t, err)
	})
}
