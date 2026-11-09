package testtools

import (
	"cmp"
	"os"
)

const (
	defaultPartnerID          = "5352b013-ae58-4a63-8a3f-97f316a917cf"
	defaultPartnerWalletPMID  = "041fdc88-c93d-4cb4-80aa-b2dde9a4fe2e"
	defaultMerchantID         = "ebbf46c6-122a-4367-bc45-7dd555e1d3b9"
	defaultMerchantWalletPMID = "67ebda6c-de48-474c-b49d-2cd3aa7d3f92"
	defaultMerchantWalletID   = "4dbac313-d505-4d51-a0fe-c11787916fcf"
)

// The seeded sandbox resources. Each identifier accepts an optional MOOV_
// environment override; an absent or empty value keeps the default.
var (
	PARTNER_ID            = cmp.Or(os.Getenv("MOOV_PARTNER_ID"), defaultPartnerID)
	PARTNER_WALLET_PM_ID  = cmp.Or(os.Getenv("MOOV_PARTNER_WALLET_PM_ID"), defaultPartnerWalletPMID)
	MERCHANT_ID           = cmp.Or(os.Getenv("MOOV_MERCHANT_ID"), defaultMerchantID)
	MERCHANT_WALLET_PM_ID = cmp.Or(os.Getenv("MOOV_MERCHANT_WALLET_PM_ID"), defaultMerchantWalletPMID)
	MERCHANT_WALLET_ID    = cmp.Or(os.Getenv("MOOV_MERCHANT_WALLET_ID"), defaultMerchantWalletID)
)
