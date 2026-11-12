package mv2607

import "github.com/moovfinancial/moov-go/pkg/moov"

// MerchantCategory is a predefined merchant category group and the merchant category codes (MCCs) it covers.
type MerchantCategory struct {
	Category moov.IssuingMerchantCategory `json:"category"`
	// Mccs are the merchant category codes covered by the category, which may include ranges.
	Mccs []string `json:"mccs"`
}

// MerchantCategories is the list of predefined merchant category groups available for issued card spend controls.
type MerchantCategories struct {
	MerchantCategories []MerchantCategory `json:"merchantCategories"`
}
