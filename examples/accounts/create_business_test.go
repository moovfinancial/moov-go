package accounts

import (
	"context"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateAndPatchBusiness(t *testing.T) {
	// Step 1: create Moov client and set some variables

	// The following code shows how you can configure the moov client with
	// your credentials, if you don't want to use environment variables.
	// However, it is recommended to load the credentials from the
	// configuration file.

	mc, err := moov.NewClient() // reads credentials from Environmental variables
	require.NoError(t, err)

	thirdPartyID := uuid.NewString() // identifier from your system (outside of Moov)

	ctx := context.Background()
	account, _, err := mc.CreateAccount(ctx, moov.CreateAccount{
		Type: moov.AccountType_Business,
		Profile: moov.CreateProfile{
			Business: &moov.CreateBusinessProfile{
				Address: &moov.Address{
					AddressLine1:    "123 Main Street",
					AddressLine2:    "Apt 302",
					City:            "Boulder",
					StateOrProvince: "CO",
					PostalCode:      "80301",
					Country:         "US",
				},
				Type:        "llc",
				Description: "Local fitness center paying out instructors",
				DBA:         "Whole Body Fitness",
				Email:       "amanda@classbooker.dev",
				IndustryCodes: &moov.IndustryCodes{
					Naics: "713940",
					Sic:   "7991",
					Mcc:   "7997",
				},
				Name: "Whole Body Fitness LLC",
				Phone: &moov.Phone{
					Number:      "8185551212",
					CountryCode: "1",
				},
				TaxID: &moov.TaxID{
					EIN: moov.EIN{
						Number: "123456789",
					},
				},
				Website: "www.wholebodyfitnessgym.com",
			},
		},
		RequestedCapabilities: []moov.CapabilityName{
			moov.CapabilityName_Transfers,
			moov.CapabilityName_SendFunds,
			moov.CapabilityName_CollectFunds,
			moov.CapabilityName_Wallet,
		},
		CustomerSupport: &moov.CustomerSupport{
			Address: &moov.Address{
				AddressLine1:    "123 Main Street",
				AddressLine2:    "Unit 302",
				City:            "Boulder",
				StateOrProvince: "CO",
				PostalCode:      "80301",
				Country:         "US",
			},
			Email: "amanda@classbooker.dev",
			Phone: &moov.Phone{
				Number:      "8185551212",
				CountryCode: "1",
			},
		},
		ForeignID: thirdPartyID,
		Metadata: map[string]string{
			"Property1": "string",
			"Property2": "string",
		},
		AccountSettings: &moov.AccountSettings{
			CardPayment: &moov.CardPaymentSettings{
				StatementDescriptor: "Amanda Yang",
			},
		},
	})
	require.NoError(t, err)

	t.Logf("Created Account %s", account.AccountID)

	// Step 2: patch the account with new information
	patchAccount := moov.PatchAccount{
		CustomerSupport: &moov.CustomerSupport{
			Address: &moov.Address{
				AddressLine1:    "123 Main Street",
				AddressLine2:    "Apt 302",
				City:            "Boulder",
				StateOrProvince: "CO",
				PostalCode:      "80301",
				Country:         "US",
			},
			Email: "test@moov.io",
			Phone: &moov.Phone{
				Number:      "8185551212",
				CountryCode: "1",
			},
			Website: "https://moov.io",
		},
	}

	account, err = mc.PatchAccount(ctx, account.AccountID, patchAccount)
	require.NoError(t, err)

	t.Logf("Patched Account %s", account.AccountID)
}
