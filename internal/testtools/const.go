package testtools

import "os"

const (
	defaultPartnerID          = "5352b013-ae58-4a63-8a3f-97f316a917cf"
	defaultPartnerWalletPMID  = "041fdc88-c93d-4cb4-80aa-b2dde9a4fe2e"
	defaultMerchantID         = "ebbf46c6-122a-4367-bc45-7dd555e1d3b9"
	defaultMerchantWalletPMID = "67ebda6c-de48-474c-b49d-2cd3aa7d3f92"
	defaultMerchantWalletID   = "4dbac313-d505-4d51-a0fe-c11787916fcf"
)

// Seeded sandbox account identifiers. Each one can be overridden with the
// matching MOOV_-prefixed environment variable; unset or empty keeps the
// default shared sandbox account.
var (
	PARTNER_ID            = seededID("MOOV_PARTNER_ID", defaultPartnerID)
	PARTNER_WALLET_PM_ID  = seededID("MOOV_PARTNER_WALLET_PM_ID", defaultPartnerWalletPMID)
	MERCHANT_ID           = seededID("MOOV_MERCHANT_ID", defaultMerchantID)
	MERCHANT_WALLET_PM_ID = seededID("MOOV_MERCHANT_WALLET_PM_ID", defaultMerchantWalletPMID)
	MERCHANT_WALLET_ID    = seededID("MOOV_MERCHANT_WALLET_ID", defaultMerchantWalletID)
)

func seededID(envName, fallback string) string {
	if v := os.Getenv(envName); v != "" {
		return v
	}
	return fallback
}
