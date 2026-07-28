package mv2607

import (
	"context"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

type ProductClient struct {
	*moov.Client
}

func NewProductClient(client *moov.Client) ProductClient {
	return ProductClient{Client: client}
}

// CreateProduct creates a new product for the specified account.
// https://docs.moov.io/api/tools/products/post/
func (p ProductClient) CreateProduct(ctx context.Context, accountID string, product ProductRequest) (*Product, error) {
	return moov.CreateProductGeneric[ProductRequest, Product](ctx, p.Client, moov.Version2026_07, accountID, product)
}

// ListProducts lists products for an account.
// https://docs.moov.io/api/tools/products/list/
func (p ProductClient) ListProducts(ctx context.Context, accountID string, filters ...moov.ProductListFilter) ([]Product, error) {
	return moov.ListProductsGeneric[Product](ctx, p.Client, moov.Version2026_07, accountID, filters...)
}

// GetProduct retrieves a product by ID.
// https://docs.moov.io/api/tools/products/get/
func (p ProductClient) GetProduct(ctx context.Context, accountID string, productID string) (*Product, error) {
	return moov.GetProductGeneric[Product](ctx, p.Client, moov.Version2026_07, accountID, productID)
}

// UpdateProduct updates a product and its options.
// https://docs.moov.io/api/tools/products/put/
func (p ProductClient) UpdateProduct(ctx context.Context, accountID string, productID string, product ProductRequest) (*Product, error) {
	return moov.UpdateProductGeneric[ProductRequest, Product](ctx, p.Client, moov.Version2026_07, accountID, productID, product)
}

// DisableProduct disables a product by ID.
// The product will no longer be available, but will remain in the system for historical and reporting purposes.
// https://docs.moov.io/api/tools/products/delete/
func (p ProductClient) DisableProduct(ctx context.Context, accountID string, productID string) error {
	return moov.DisableProductGeneric(ctx, p.Client, moov.Version2026_07, accountID, productID)
}

// ListProductCategories returns the full, read-only list of product categories from the
// product taxonomy.
func (p ProductClient) ListProductCategories(ctx context.Context) ([]ProductCategory, error) {
	categories, err := moov.ListProductCategoriesGeneric[ProductCategories](ctx, p.Client, moov.Version2026_07)
	if err != nil {
		return nil, err
	}

	return categories.Categories, nil
}
