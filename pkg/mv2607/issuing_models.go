package mv2607

import "github.com/moovfinancial/moov-go/pkg/moov"

type IssuingMerchantCategoryMCCs struct {
	Category moov.IssuingMerchantCategory `json:"category"`
	// MCCs are the merchant category codes covered by the category, which may include ranges.
	MCCs []string `json:"mccs"`
}

type IssuingMerchantCategoryMCCsList struct {
	MerchantCategories []IssuingMerchantCategoryMCCs `json:"merchantCategories"`
}
