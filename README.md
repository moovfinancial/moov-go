# moov-go
Golang client for interacting with Moov API server side 

### Installation 

```bash 
# Go Client 
go get github.com/moovfinancial/moov-go/pkg
```

This SDK requires and API key. To generate an API login to the Moov Dashboard and follow the following instructions on [API Keys](https://docs.moov.io/guides/get-started/api-keys/). If you have not done so already, use the [Moov Dashboard](https://dashboard.moov.io/signup) to create an account. Note that API Keys for Sandbox and Production are different keys. 

## Basic Usage 

```go 
package main

import (
	"log"

	moov "github.com/moovfinancial/moov-go/pkg"
)

func main() {

	// Setup Moov client
	// Generate API keys from the Moov Dashboard. Sandbox keys are not compatible with production.
	creds := moov.Credentials{
		AccountID: "638481a5-FAKE-406c-84c7-2fc2239105d1",
		PublicKey: "Qo0j0ChFAKEKmRI_",
		SecretKey: "nbyFAKEr0MpM2gxkzo2u2G70jr6ZtFdm",
		Domain:    "localhost",
	}

    mc, err := moov.NewClient(creds)
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

```
