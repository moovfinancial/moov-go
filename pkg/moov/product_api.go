package moov

import (
	"context"
	"errors"
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

// CreateProductGeneric creates a product against a specific API version. Version-specific
// packages (e.g. mv2607) supply their own request and product types.
func CreateProductGeneric[TRequest any, TProduct any](ctx context.Context, client *Client, version Version, accountID string, product TRequest) (*TProduct, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}

	resp, err := client.CallHttp(ctx,
		Endpoint(http.MethodPost, pathProducts, accountID),
		MoovVersion(version),
		AcceptJson(),
		JsonBody(product))
	if err != nil {
		return nil, err
	}

	return StartedObjectOrError[TProduct](resp)
}

// ListProductsGeneric lists products against a specific API version.
func ListProductsGeneric[TProduct any](ctx context.Context, client *Client, version Version, accountID string, filters ...ProductListFilter) ([]TProduct, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}

	resp, err := client.CallHttp(ctx,
		Endpoint(http.MethodGet, pathProducts, accountID),
		prependArgs(filters, MoovVersion(version), AcceptJson())...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[TProduct](resp)
}

// GetProductGeneric retrieves a product against a specific API version.
func GetProductGeneric[TProduct any](ctx context.Context, client *Client, version Version, accountID string, productID string) (*TProduct, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}

	resp, err := client.CallHttp(ctx,
		Endpoint(http.MethodGet, pathProduct, accountID, productID),
		MoovVersion(version),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[TProduct](resp)
}

// UpdateProductGeneric updates a product against a specific API version.
func UpdateProductGeneric[TRequest any, TProduct any](ctx context.Context, client *Client, version Version, accountID string, productID string, product TRequest) (*TProduct, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}

	resp, err := client.CallHttp(ctx,
		Endpoint(http.MethodPut, pathProduct, accountID, productID),
		MoovVersion(version),
		AcceptJson(),
		JsonBody(product))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[TProduct](resp)
}

// DisableProductGeneric disables a product against a specific API version.
func DisableProductGeneric(ctx context.Context, client *Client, version Version, accountID string, productID string) error {
	if client == nil {
		return errors.New("client is nil")
	}

	resp, err := client.CallHttp(ctx,
		Endpoint(http.MethodDelete, pathProduct, accountID, productID),
		MoovVersion(version),
		AcceptJson())
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}

// ListProductCategoriesGeneric returns the full, read-only product taxonomy against a
// specific API version. TCategories is the version's wrapper type holding the category list.
// https://docs.moov.io/api/tools/products/list-categories/
func ListProductCategoriesGeneric[TCategories any](ctx context.Context, client *Client, version Version) (*TCategories, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}

	resp, err := client.CallHttp(ctx,
		Endpoint(http.MethodGet, pathProductCategories),
		MoovVersion(version),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[TCategories](resp)
}
