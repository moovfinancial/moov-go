package moov_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

const FACILITATOR_ID = "5352b013-ae58-4a63-8a3f-97f316a917cf"

func getLincolnBank(t *testing.T, mc *moov.Client) *moov.Account {
	accounts, err := mc.ListAccounts(context.Background(), moov.WithAccountName("Lincoln National Corporation"))
	moov.DebugPrintResponse(err, fmt.Printf)
	require.NoError(t, err)

	for _, account := range accounts {
		if account.DisplayName == "Lincoln National Corporation" {
			return &account
		}
	}

	require.FailNow(t, "bank account test account not found")
	return nil
}

func randomBankAccountNumber() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(999999999))
	return fmt.Sprintf("%d", 100000000+n.Int64())
}

func createTestIndividualAccount() moov.CreateAccount {
	return moov.CreateAccount{
		Type: moov.AccountType_Individual,
		Profile: moov.CreateProfile{
			Individual: &moov.CreateIndividualProfile{
				Name: moov.Name{
					FirstName: "John",
					LastName:  "Doe",
				},
				Email: "noreply@moov.io",
				Phone: &moov.Phone{
					Number:      "555-555-5555",
					CountryCode: "1",
				},
			},
		},
	}
}

func createTestBusinessAccount() moov.CreateAccount {
	return moov.CreateAccount{
		Type: moov.AccountType_Business,
		Profile: moov.CreateProfile{
			Business: &moov.CreateBusinessProfile{
				Name:        "John Does Hobbies",
				Type:        moov.BusinessType_Llc,
				Description: "moov-go SDK testing",
				IndustryCodes: &moov.IndustryCodes{
					Mcc:   "6012",
					Naics: "522110",
					Sic:   "6021",
				},
			},
		},
	}
}

func paymentMethodsFromOptions(t *testing.T, options *moov.TransferOptions, sourceType moov.PaymentMethodType, destType moov.PaymentMethodType) (string, string) {
	sourceId := ""
	destId := ""
	for _, pm := range options.SourceOptions {
		if pm.PaymentMethodType == sourceType {
			sourceId = pm.PaymentMethodID
			break
		}
	}
	for _, pm := range options.DestinationOptions {
		if pm.PaymentMethodType == destType {
			destId = pm.PaymentMethodID
			break
		}
	}

	require.NotEmpty(t, sourceId, "unable to find source payment method for type")
	require.NotEmpty(t, destId, "unable to find destination payment method for type")

	return sourceId, destId
}

func NoResponseError(t *testing.T, err error) {
	moov.DebugPrintResponse(err, fmt.Printf)
	require.NoError(t, err)
}

func CreateTemporaryTestAccount(t *testing.T, mc *moov.Client, create moov.CreateAccount) *moov.Account {
	account, started, err := mc.CreateAccount(context.Background(), create)
	moov.DebugPrintResponse(err, fmt.Printf)

	require.NoError(t, err)
	require.NotNil(t, account)
	require.Nil(t, started)

	t.Cleanup(func() {
		if account != nil {
			mc.DisconnectAccount(BgCtx(), account.AccountID)
		}
	})

	return account
}
