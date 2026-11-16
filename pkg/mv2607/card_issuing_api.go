package mv2607

import (
	"context"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

type CardIssuingClient struct {
	*moov.Client
}

func NewCardIssuingClient(client *moov.Client) CardIssuingClient {
	return CardIssuingClient{Client: client}
}

// ListMerchantCategories returns the issuing merchant categories and the merchant
// category codes (MCCs) each one covers.
// https://docs.moov.io/api/money-movement/issuing/list-merchant-categories/
func (i CardIssuingClient) ListMerchantCategories(ctx context.Context) (*MerchantCategories, error) {
	return moov.ListIssuingMerchantCategoriesGeneric[MerchantCategories](ctx, i.Client, moov.Version2026_07)
}
