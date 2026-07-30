package mv2607_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2607"
)

// newProductTestClient points a client at srv so tests can assert on the raw request.
func newProductTestClient(t *testing.T, srv *httptest.Server) mv2607.ProductClient {
	t.Helper()

	client, err := moov.NewClient(
		moov.WithCredentials(moov.Credentials{PublicKey: "pk", SecretKey: "sk"}),
		moov.WithMoovURLScheme("http"),
	)
	require.NoError(t, err)
	client.Credentials.Host = strings.TrimPrefix(srv.URL, "http://")

	return mv2607.NewProductClient(client)
}

const productWithCategoryJSON = `{
	"productID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43",
	"title": "Cold Brew",
	"basePrice": {
		"currency": "USD",
		"valueDecimal": "4.50"
	},
	"category": {
		"categoryID": "413",
		"name": "Beverages",
		"fullName": "Food, Beverages & Tobacco > Beverages",
		"level": 2,
		"parentID": "412"
	},
	"createdOn": "2026-06-18T12:00:00Z",
	"updatedOn": "2026-06-18T12:00:00Z"
}`

func TestCreateProduct(t *testing.T) {
	var (
		method    string
		path      string
		version   string
		body      map[string]any
		decodeErr error
	)

	// The handler runs on its own goroutine, so it records what it saw rather than
	// asserting: require.* calls t.FailNow, which is only valid on the test goroutine.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		path = r.URL.Path
		version = r.Header.Get(moov.VersionHeader)
		decodeErr = json.NewDecoder(r.Body).Decode(&body)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(productWithCategoryJSON))
	}))
	t.Cleanup(srv.Close)

	products := newProductTestClient(t, srv)

	created, err := products.CreateProduct(context.Background(), "account-123", mv2607.ProductRequest{
		Title: "Cold Brew",
		BasePrice: moov.AmountDecimal{
			Currency:     "USD",
			ValueDecimal: "4.50",
		},
		CategoryID: moov.PtrOf("413"),
	})
	require.NoError(t, err)
	require.NoError(t, decodeErr)

	require.Equal(t, http.MethodPost, method)
	require.Equal(t, "/accounts/account-123/products", path)
	require.Equal(t, moov.Version2026_07.String(), version)
	require.Equal(t, "413", body["categoryID"])

	require.NotNil(t, created)
	require.Equal(t, "Cold Brew", created.Title)
	require.NotNil(t, created.Category)
	require.Equal(t, "413", created.Category.CategoryID)
	require.Equal(t, "Beverages", created.Category.Name)
	require.Equal(t, "Food, Beverages & Tobacco > Beverages", created.Category.FullName)
	require.Equal(t, int32(2), created.Category.Level)
	require.Equal(t, moov.PtrOf("412"), created.Category.ParentID)
}

func TestProductRequestOmitsCategoryID(t *testing.T) {
	b, err := json.Marshal(mv2607.ProductRequest{
		Title:     "Cold Brew",
		BasePrice: moov.AmountDecimal{Currency: "USD", ValueDecimal: "4.50"},
	})
	require.NoError(t, err)
	require.NotContains(t, string(b), `"categoryID"`)
}

func TestUpdateProduct(t *testing.T) {
	var (
		method    string
		path      string
		version   string
		body      map[string]any
		decodeErr error
	)

	// The handler runs on its own goroutine, so it records what it saw rather than
	// asserting: require.* calls t.FailNow, which is only valid on the test goroutine.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		path = r.URL.Path
		version = r.Header.Get(moov.VersionHeader)
		decodeErr = json.NewDecoder(r.Body).Decode(&body)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(productWithCategoryJSON))
	}))
	t.Cleanup(srv.Close)

	products := newProductTestClient(t, srv)

	updated, err := products.UpdateProduct(context.Background(), "account-123", "ec7e1848-dc80-4ab0-8827-dd7fc0737b43", mv2607.ProductRequest{
		Title:      "Cold Brew",
		BasePrice:  moov.AmountDecimal{Currency: "USD", ValueDecimal: "4.50"},
		CategoryID: moov.PtrOf("413"),
	})
	require.NoError(t, err)
	require.NoError(t, decodeErr)

	require.Equal(t, http.MethodPut, method)
	require.Equal(t, "/accounts/account-123/products/ec7e1848-dc80-4ab0-8827-dd7fc0737b43", path)
	require.Equal(t, moov.Version2026_07.String(), version)
	require.Equal(t, "413", body["categoryID"])

	require.NotNil(t, updated)
	require.NotNil(t, updated.Category)
	require.Equal(t, "413", updated.Category.CategoryID)
}

func TestGetProduct(t *testing.T) {
	var version string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		version = r.Header.Get(moov.VersionHeader)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(productWithCategoryJSON))
	}))
	t.Cleanup(srv.Close)

	products := newProductTestClient(t, srv)

	got, err := products.GetProduct(context.Background(), "account-123", "ec7e1848-dc80-4ab0-8827-dd7fc0737b43")
	require.NoError(t, err)

	require.Equal(t, moov.Version2026_07.String(), version)
	require.NotNil(t, got.Category)
	require.Equal(t, "Beverages", got.Category.Name)
}

func TestListProductsFilteredByCategory(t *testing.T) {
	var (
		path    string
		query   string
		version string
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		query = r.URL.Query().Get("category")
		version = r.Header.Get(moov.VersionHeader)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("[" + productWithCategoryJSON + "]"))
	}))
	t.Cleanup(srv.Close)

	products := newProductTestClient(t, srv)

	list, err := products.ListProducts(context.Background(), "account-123", moov.WithProductCategory("412"))
	require.NoError(t, err)

	require.Equal(t, "/accounts/account-123/products", path)
	require.Equal(t, "412", query)
	require.Equal(t, moov.Version2026_07.String(), version)

	require.Len(t, list, 1)
	require.NotNil(t, list[0].Category)
	require.Equal(t, "413", list[0].Category.CategoryID)
}

func TestListProductCategories(t *testing.T) {
	var (
		method  string
		path    string
		version string
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		path = r.URL.Path
		version = r.Header.Get(moov.VersionHeader)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"categories": [
				{
					"categoryID": "412",
					"name": "Food, Beverages & Tobacco",
					"fullName": "Food, Beverages & Tobacco",
					"level": 1
				},
				{
					"categoryID": "413",
					"name": "Beverages",
					"fullName": "Food, Beverages & Tobacco > Beverages",
					"level": 2,
					"parentID": "412"
				}
			]
		}`))
	}))
	t.Cleanup(srv.Close)

	products := newProductTestClient(t, srv)

	categories, err := products.ListProductCategories(context.Background())
	require.NoError(t, err)

	require.Equal(t, http.MethodGet, method)
	require.Equal(t, "/product-categories", path)
	require.Equal(t, moov.Version2026_07.String(), version)

	require.Len(t, categories, 2)
	require.Equal(t, "412", categories[0].CategoryID)
	require.Equal(t, int32(1), categories[0].Level)
	require.Nil(t, categories[0].ParentID)
	require.Equal(t, "413", categories[1].CategoryID)
	require.Equal(t, moov.PtrOf("412"), categories[1].ParentID)
}

func TestDisableProduct(t *testing.T) {
	var (
		method  string
		path    string
		version string
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		path = r.URL.Path
		version = r.Header.Get(moov.VersionHeader)
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(srv.Close)

	products := newProductTestClient(t, srv)

	require.NoError(t, products.DisableProduct(context.Background(), "account-123", "ec7e1848-dc80-4ab0-8827-dd7fc0737b43"))

	require.Equal(t, http.MethodDelete, method)
	require.Equal(t, "/accounts/account-123/products/ec7e1848-dc80-4ab0-8827-dd7fc0737b43", path)
	require.Equal(t, moov.Version2026_07.String(), version)
}
