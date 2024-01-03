package moov

import "time"

type ApplePay struct {
	Brand           string     `json:"brand,omitempty"`
	CardType        string     `json:"cardType,omitempty"`
	CardDisplayName string     `json:"cardDisplayName,omitempty"`
	Fingerprint     string     `json:"fingerprint,omitempty"`
	Expiration      Expiration `json:"expiration,omitempty"`
	DynamicLastFour string     `json:"dynamicLastFour,omitempty"`
}

type ApplePayDomains struct {
	DisplayName string   `json:"displayName,omitempty"`
	Domains     []string `json:"domains,omitempty"`
}

type ApplePayDomainsResponse struct {
	AccountID   string    `json:"accountID,omitempty"`
	DisplayName string    `json:"displayName,omitempty"`
	Domains     []string  `json:"domains,omitempty"`
	CreatedOn   time.Time `json:"createdOn,omitempty"`
	UpdatedOn   time.Time `json:"updatedOn,omitempty"`
}

// CreateApplePayDomain creates a new Apple Pay domain for the given customer account
// https://docs.moov.io/api/#tag/Cards/operation/registerApplePayMerchantDomains

// UpdateApplePayDomain updates an existing Apple Pay domain for the given customer account
// https://docs.moov.io/api/#tag/Cards/operation/updateApplePayMerchantDomains

// GetApplePayDomain retrieves an Apple Pay domain for the given customer account
// https://docs.moov.io/api/#tag/Cards/operation/getApplePayMerchantDomains

// CreateApplePaySession creates a new Apple Pay session for the given customer account
// https://docs.moov.io/api/#tag/Cards/operation/getApplePayMerchantDomains

// ApplePayToken creates a new Apple Pay token for the given customer account
// https://docs.moov.io/api/#tag/Cards/operation/getApplePayMerchantDomains
