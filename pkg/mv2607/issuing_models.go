package mv2607

import "github.com/moovfinancial/moov-go/pkg/moov"

// MerchantCategory is an issuing merchant category with the merchant category codes (MCCs) it covers.
type MerchantCategory struct {
	// Category is the predefined merchant category group.
	Category moov.IssuingMerchantCategory `json:"category"`
	// MCCs are the merchant category codes covered by the category, which may include ranges.
	MCCs []string `json:"mccs"`
}

// MerchantCategories is a list of issuing merchant categories.
type MerchantCategories struct {
	MerchantCategories []MerchantCategory `json:"merchantCategories"`
}
