package moov

import (
	"context"
	"net/http"
)

// CreateProduct creates a new product for the specified account.
// https://docs.moov.io/api/tools/products/post/
func (c Client) CreateProduct(ctx context.Context, accountID string, product ProductRequest) (*Product, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathProducts, accountID),
		AcceptJson(),
		JsonBody(product))
	if err != nil {
		return nil, err
	}

	return StartedObjectOrError[Product](resp)
}

// ListProducts lists products for an account.
// https://docs.moov.io/api/tools/products/list/
func (c Client) ListProducts(ctx context.Context, accountID string, filters ...ProductListFilter) ([]Product, error) {
	args := prependArgs(filters, AcceptJson())
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathProducts, accountID), args...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[Product](resp)
}

// GetProduct retrieves a product by ID.
// https://docs.moov.io/api/tools/products/get/
func (c Client) GetProduct(ctx context.Context, accountID string, productID string) (*Product, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathProduct, accountID, productID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Product](resp)
}

// UpdateProduct updates a product and its options.
// https://docs.moov.io/api/tools/products/put/
func (c Client) UpdateProduct(ctx context.Context, accountID string, productID string, product ProductRequest) (*Product, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPut, pathProduct, accountID, productID),
		AcceptJson(),
		JsonBody(product))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Product](resp)
}

// DisableProduct disables a product by ID.
// The product will no longer be available, but will remain in the system for historical and reporting purposes.
// https://docs.moov.io/api/tools/products/delete/
func (c Client) DisableProduct(ctx context.Context, accountID string, productID string) error {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodDelete, pathProduct, accountID, productID),
		AcceptJson())
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}
