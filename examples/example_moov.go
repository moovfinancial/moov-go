package main

import (
	"context"
	"fmt"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

func main() {
	fmt.Println("Moov SDK Examples")

	// Create a new moov client
	mc, err := moov.NewClient()
	if err != nil {
		// Exit with invalid credentials
		fmt.Println(err)
		return
	}

	// Create a new context or use an existing one
	ctx := context.Background()

	// Ping the server to check credentials
	err = mc.Ping(ctx)
	if err != nil {
		fmt.Println(err)
	}

	// Create access token to instantiate a moov.js drop for creating an account
	token, err := mc.AccountCreationToken(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(token)

	// Create a new moov account server side
	// Create a new account
	account := moov.Account{
		AccountType: moov.INDIVIDUAL,
		Profile: moov.Profile{
			Individual: moov.Individual{
				Name: moov.Name{
					FirstName: "Wade",
					LastName:  "Arnold",
				},
				Email: "wade@arnold.com",
				Phone: moov.Phone{
					Number:      "555-555-5555",
					CountryCode: "1",
				},
			},
		},
	}
	completedAccount, _, err := mc.CreateAccount(ctx, account)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(completedAccount.AccountID)

	// List all accounts with a quary parameter

	// Create a scoped access token for a specific account
}
