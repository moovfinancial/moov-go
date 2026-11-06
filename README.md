![Moov Logo](https://github.com/moovfinancial/moov-go/assets/120951/3632d9ea-0c64-40e5-8f9e-b13b28b5e197)

[![GoDoc](https://pkg.go.dev/badge/github.com/moovfinancial/moovgo?utm_source=godoc)](https://pkg.go.dev/github.com/moovfinancial/moov-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/moovfinancial/moov-go)](https://goreportcard.com/report/github.com/moovfinancial/moov-go)
[![Apache 2 License](https://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/moovfinancial/moov-go/master/LICENSE)
[![X](https://img.shields.io/twitter/follow/moov?style=social)](https://twitter.com/moov?lang=en)


# Moov - Go Client
The official [Go](http://golang.org) client for the [Moov payments API](https://docs.moov.io/api/).

<!-- Start SDK Installation [installation] -->
## Installation

```go
import (
	"github.com/moovfinancial/moov-go/pkg/moov"
)
```

This SDK requires an API key. To generate an API login to the Moov Dashboard and follow the following instructions on [API Keys](https://docs.moov.io/guides/get-started/api-keys/). If you have not done so already, use the [Moov Dashboard](https://dashboard.moov.io/signup) to create an account.

```bash
export MOOV_PUBLIC_KEY="public key here"
export MOOV_SECRET_KEY="secret key here"
```

In your Go program, create a new Moov client initiated with your public and secret keys.

```go
mc, err := moov.NewClient(
  moov.WithCredentials(moov.CredentialsFromEnv()), // optional, default is to read from environment
)
```
<!-- End SDK Installation [installation] -->

## Examples Usage

Checkout the [examples](./examples/README.md) exist for [ach](./examples/ach/), [card acquiring](./examples/card_acquiring/), debit [push](./examples/debit_card_push/)/[pull](./examples/debit_card_pull/), and [rtp](./examples/rtp/).

## Payment Links

Payment link requests select `X-Moov-Version: v2026.07.00`. Custom amount payment links do not exist in earlier versions.

```go
create := moov.CreatePaymentLink{
	PartnerAccountID:        partnerAccountID,
	MerchantPaymentMethodID: merchantWalletPaymentMethodID,
	Amount:                  &moov.Amount{Currency: "USD", Value: 1050},
	Display: moov.PaymentLinkDisplayOptions{
		Title:        "Invoice #123",
		Description:  "Payment for services",
		CallToAction: moov.CallToAction_Pay,
	},
	Payment: &moov.PaymentLinkPaymentDetails{
		AllowedMethods: []moov.CollectionPaymentMethodType{moov.CollectionPaymentMethodType_CardPayment},
	},
}

link, err := mc.CreatePaymentLink(ctx, accountID, create)

// List filters support count, skip, and comma-separated types.
links, err := mc.ListPaymentLinks(ctx, accountID,
	moov.WithPaymentLinkCount(20),
	moov.WithPaymentLinkSkip(0),
	moov.WithPaymentLinkTypes(moov.PaymentLinkType_Payment),
)

// GetPaymentLinkQRCode returns the raw image bytes and the response content type.
img, contentType, err := mc.GetPaymentLinkQRCode(ctx, accountID, string(link.Code))
```

A custom amount link omits `Amount` and lets the payor choose the amount:

```go
link, err := mc.CreatePaymentLink(ctx, accountID, moov.CreatePaymentLink{
	PartnerAccountID:        partnerAccountID,
	MerchantPaymentMethodID: merchantWalletPaymentMethodID,
	Display:                 display,
	CustomAmountPayment: &moov.PaymentLinkCustomAmountPaymentDetails{
		AllowedMethods: []moov.CollectionPaymentMethodType{moov.CollectionPaymentMethodType_CardPayment},
		AmountRange: &moov.AmountDecimalRange{
			Minimum: &moov.AmountDecimal{Currency: "USD", ValueDecimal: "1.00"},
			Maximum: &moov.AmountDecimal{Currency: "USD", ValueDecimal: "100.00"},
		},
	},
})
```

## Onboarding Invites

Onboarding invites are organization-level and are identified by their invite code.

```go
invite, err := mc.CreateOnboardingInvite(ctx, moov.OnboardingInviteRequest{
	ReturnURL:    "https://example.com/onboarding",
	Scopes:       []moov.ApplicationScope{"accounts.read"},
	Capabilities: []moov.CapabilityName{moov.CapabilityName_Wallet},
	FeePlanCodes: []string{"merchant-direct"},
})

invites, err := mc.ListOnboardingInvites(ctx)
fetched, err := mc.GetOnboardingInvite(ctx, string(invite.Code))
err = mc.RevokeOnboardingInvite(ctx, string(invite.Code))
```

## License

Apache License 2.0 - See [LICENSE](LICENSE) for details.
