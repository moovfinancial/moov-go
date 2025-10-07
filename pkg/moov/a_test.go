package moov_test

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/moovfinancial/moov-go/internal/testtools"
	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2507"
	"github.com/stretchr/testify/require"
)

const FACILITATOR_ID = testtools.PARTNER_ID
const FACILITATOR_WALLET_PM_ID = testtools.PARTNER_WALLET_PM_ID

const MERCHANT_WALLET_PM_ID = testtools.MERCHANT_WALLET_PM_ID

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

func createTestIndividualAccount_V2507() mv2507.CreateAccount {
	return mv2507.CreateAccount{
		Type: mv2507.CreateAccountType_Individual,
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

func PrettyDebug(t *testing.T, a any) {
	b, err := json.MarshalIndent(a, "  ", "  ")
	require.NoError(t, err)

	t.Logf("\n%s\n", string(b))
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

func createTemporaryCard(t *testing.T, mc *moov.Client, accountID string) *moov.Card {
	exp := time.Now().UTC().AddDate(0, 7, 0)

	// Create card
	card, err := mc.CreateCard(context.Background(), accountID, moov.CreateCard{
		CardNumber: "4111111111111111",
		CardCvv:    "123",
		Expiration: moov.Expiration{
			Month: exp.Format("01"),
			Year:  exp.Format("06"),
		},
		HolderName: "john doe",
		BillingAddress: moov.Address{
			AddressLine1:    "123 Main Street",
			City:            "City",
			StateOrProvince: "CO",
			PostalCode:      "12345",
			Country:         "US",
		},
		CardOnFile: false,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		if card != nil {
			_ = mc.DisableCard(BgCtx(), accountID, card.CardID)
		}
	})

	require.NotEmpty(t, card.PaymentMethods)
	return card
}
