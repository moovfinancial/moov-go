![Moov Banner Logo](https://github.com/moovfinancial/moov-go/assets/120951/3632d9ea-0c64-40e5-8f9e-b13b28b5e197)

[![GoDoc](https://godoc.org/github.com/moovfinancial/moovgo?status.svg)](https://godoc./github.com/moovfinancial/moov-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/moovfinancial/moov-go)](https://goreportcard.com/report/github.com/moovfinancial/moov-go)
[![Apache 2 License](https://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/moovfinancial/moov-go/master/LICENSE)
[![Twitter](https://img.shields.io/twitter/follow/moov?style=social)](https://twitter.com/moov?lang=en)

# moov-go

Golang client for interacting with [Moov Financial's API](https://moov.io/product/).

If you believe you have identified a security vulnerability please responsibly report the issue as via email to `security@moov.io`. Please do not post it to a public issue tracker.

### Installation

```bash
# Go Client
go get github.com/moovfinancial/moov-go/pkg
```

This SDK requires and API key. To generate an API login to the Moov Dashboard and follow the following instructions on [API Keys](https://docs.moov.io/guides/get-started/api-keys/). If you have not done so already, use the [Moov Dashboard](https://dashboard.moov.io/signup) to create an account. Note that API Keys for Sandbox and Production are different keys.

After getting API keys, you need to set up on `config.yaml` file
```
public_key: public_key
secret_key: secret_key
```

## Basic Usage

```go
package main

import (
	"log"
	"github.com/moovfinancial/moov-go/pkg/moov"
)

func main() {
	// Setup Moov client
    mc, err := moov.NewClient()
	if err != nil {
        // Network and authentication errors
		log.Fatal(err)
	}
	token, err := mc.SingleUseAccessToken()
	if err != nil {
		log.Fatal(err)
	}
	// Access token for website/mobile app to instantiate Moov.js
	log.Print(token.AccessToken)

	// Create a new account
	account := moov.Account{
		AccountType: moov.INDIVIDUAL,
		Profile: moov.Profile{
			Individual: moov.Individual{
				Name: moov.Name{
					FirstName: "Wade",
					LastName:  "Arnold",
				},
				Email: "Wade@arnold.com",
				Phone: moov.Phone{
					Number:      "555-555-5555",
					CountryCode: "1",
				},
			},
		},
	}
	account, _ = mc.CreateAccount(account)
	log.Print(account.AccountID)

	// List Accounts
	accounts, _ := mc.ListAccounts()
	log.Print(len(accounts)) // 2

	// Get an account
	account, _ = mc.GetAccount("638481a5-5205-406c-84c7-2fc2239105d1")
	log.Print(account.Profile.Individual.Name.FirstName)

	// Update an account
	account := Account{
		AccountID:   "aa19c3a7-4c72-4f64-adfa-9069c80d81cf",
		AccountType: INDIVIDUAL,
		Profile: Profile{
			Individual: Individual{
				Name: Name{
					FirstName: "David",
					LastName:  "Arnold",
				},
				Email: "Wade@arnold.com",
				Phone: Phone{
					Number:      "333-333-3333",
					CountryCode: "1",
				},
			},
		},
	}
	account, _ = mc.UpdateAccount(account)
	log.Print(account.AccountID) // David

}

// scoped access token
scopes := []string{moov.SCOPE_ACCOUNTS_WRITE, moov.SCOPE_BANK_ACCOUNT_WRITE, moov.SCOPE_CARDS_WRITE}
linkPaymentToken, err := mc.ScopedAccessToken(account.AccountID, scopes)
if err != nil {
	log.Fatal(err)
}
log.Print("LinkPaymentToken: ", linkPaymentToken.AccessToken)

```

## License

Apache License 2.0 - See [LICENSE](LICENSE) for details.
