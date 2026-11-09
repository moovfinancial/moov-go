package mv2607

import "github.com/moovfinancial/moov-go/pkg/moov"

type MerchantCategory struct {
	Category moov.IssuingMerchantCategory `json:"category"`
	// MCCs are the merchant category codes covered by the category, which may include ranges.
	MCCs []string `json:"mccs"`
}

type MerchantCategories struct {
	MerchantCategories []MerchantCategory `json:"merchantCategories"`
}
