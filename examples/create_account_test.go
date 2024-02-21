package main

import (
	"context"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	// Step 1: create Moov client and set some variables

	// The following code shows how you can configure the moov client with
	// your credentials, if you don't want to use environment variables.
	// However, it is recommended to load the credentials from the
	// configuration file.
	mc, err := moov.NewClient(moov.WithCredentials(moov.CredentialsFromEnv()))
	require.NoError(t, err)

	// Create a new context or use an existing one
	ctx := context.Background()

	// Ping the server to check credentials
	err = mc.Ping(ctx)
	require.NoError(t, err)

	// Step 2: create account for the user

	// Add new account
	accountID := "f800d429-0b6d-4962-b681-c16b3d855b40"
	var account *moov.Account

	if accountID != "" {
		account, err = mc.GetAccount(ctx, accountID)
	} else {
		account, _, err = mc.CreateAccount(ctx, moov.Account{
			AccountType: moov.BUSINESS,
			Profile: moov.Profile{
				Business: moov.Business{
					LegalBusinessName: "Sandwich Prep Corp",
					DoingBusinessAs:   "Switches Sandwiches",
					BusinessType:      moov.BUSINESS_TYPE_LLC,
					Address: moov.Address{
						AddressLine1:    "123 Mayo St",
						City:            "Pickleton",
						StateOrProvince: "TX",
						PostalCode:      "43215",
						Country:         "US",
					},
					Phone: moov.Phone{
						Number: "123-456-7890",
						// CountryCode: "+1",
					},
					Email:       "bread@sandwiches.com",
					Website:     "https://sandwich.food",
					Description: "The best spot for sandwiches in the galaxy.",
					TaxID: moov.TaxID{
						Ein: moov.Ein{
							Number: "123213456",
						},
					},
				},
			},
		})
		// map[business:map[address:map[country:{validation_in_invalid only US is currently supported map[]}]
		//
		// validation_in_invalid only US phone numbers are currently supported map[]}]
		//
		// taxID:map[ein:map[number:{validation_match_invalid must be a valid employer identification number map[]}]]]]]
	}
	require.NoError(t, err)

	t.Logf("account: %#v", account)
}
