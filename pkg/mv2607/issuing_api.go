package mv2607

import (
	"context"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

type IssuingClient struct {
	*moov.Client
}

func NewIssuingClient(client *moov.Client) IssuingClient {
	return IssuingClient{Client: client}
}

// ListMerchantCategories returns the issuing merchant categories and the merchant
// category codes (MCCs) each one covers.
// https://docs.moov.io/api/money-movement/issuing/list-merchant-categories/
func (i IssuingClient) ListMerchantCategories(ctx context.Context) ([]MerchantCategory, error) {
	categories, err := moov.ListIssuingMerchantCategoriesGeneric[MerchantCategories](ctx, i.Client, moov.Version2026_07)
	if err != nil {
		return nil, err
	}

	return categories.MerchantCategories, nil
}
