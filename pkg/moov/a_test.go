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

func getLincolnBank(t *testing.T, mc *moov.Client) moov.Account {
	accounts, err := mc.ListAccounts(context.Background(), moov.WithAccountName("Lincoln National Corporation"))
	moov.DebugPrintResponse(err, fmt.Printf)
	require.NoError(t, err)

	for _, account := range accounts {
		if account.DisplayName == "Lincoln National Corporation" {
			return account
		}
	}

	require.FailNow(t, "bank account test account not found")
	return moov.Account{}
}

func randomBankAccountNumber() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(999999999))
	return fmt.Sprintf("%d", 100000000+n.Int64())
}

func createTestIndividualAccount() moov.CreateAccount {
	return moov.CreateAccount{
		Type: moov.ACCOUNTTYPE_INDIVIDUAL,
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
		Type: moov.ACCOUNTTYPE_BUSINESS,
		Profile: moov.CreateProfile{
			Business: &moov.CreateBusinessProfile{
				Name: "John Does Hobbies",
				Type: moov.BUSINESSTYPE_LLC,
			},
		},
	}
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
