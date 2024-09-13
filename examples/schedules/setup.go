package schedules

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

type Env struct {
	Now time.Time

	Client *moov.Client

	PartnerID string

	Merchant     *moov.Account
	MerchantPmId string

	Customer     *moov.Account
	CustomerPmId string
}

func Setup(t *testing.T, ctx context.Context) *Env {

	godotenv.Load("../../secrets.env")

	// The following code shows how you can configure the moov client with
	// your credentials, if you don't want to use environment variables.
	// However, it is recommended to load the credentials from the
	// configuration file.
	mc, err := moov.NewClient(moov.WithCredentials(moov.CredentialsFromEnv()))
	require.NoError(t, err)

	env := Env{
		// Just bumping time to way ahead so we're not accidently tripping on test data
		Now:       time.Date(2040, time.March, 10, 12, 0, 0, 0, time.UTC),
		Client:    mc,
		PartnerID: "5352b013-ae58-4a63-8a3f-97f316a917cf",
	}

	// Merchant to accept the payments
	merchant, _, err := mc.CreateAccount(ctx, moov.CreateAccount{
		Type: moov.AccountType_Business,
		Profile: moov.CreateProfile{
			Business: &moov.CreateBusinessProfile{
				Name:        "John Does Hobbies",
				Type:        moov.BusinessType_Llc,
				Description: "Merchant in moov-go Schedules example",
				IndustryCodes: &moov.IndustryCodes{
					Mcc:   "6012",
					Naics: "522110",
					Sic:   "6021",
				},
			},
		},
	})
	require.NoError(t, err)
	env.Merchant = merchant
	t.Cleanup(func() {
		mc.DisconnectAccount(ctx, merchant.AccountID)
	})

	merchantBa, err := mc.CreateBankAccount(ctx, merchant.AccountID, moov.WithBankAccount(moov.BankAccountRequest{
		HolderName:    "Schedule business deposit target",
		HolderType:    moov.HolderType_Business,
		AccountType:   moov.BankAccountType_Checking,
		AccountNumber: randomBankAccountNumber(),
		RoutingNumber: "273976369",
	}), moov.WaitForPaymentMethod())
	require.NoError(t, err)
	t.Cleanup(func() {
		mc.DeleteBankAccount(ctx, merchant.AccountID, merchantBa.BankAccountID)
	})

	fmt.Printf("\n\n%+v\n\n", merchantBa.PaymentMethods)

	for _, pm := range merchantBa.PaymentMethods {
		if pm.PaymentMethodType == moov.PaymentMethodType_AchCreditStandard {
			env.MerchantPmId = pm.PaymentMethodID
		}
	}
	require.NotEmpty(t, env.MerchantPmId)

	// Setup the customers account
	customer, _, err := mc.CreateAccount(ctx, moov.CreateAccount{
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
	})
	require.NoError(t, err)
	env.Customer = customer
	t.Cleanup(func() {
		mc.DisconnectAccount(ctx, customer.AccountID)
	})

	// Create card
	exp := env.Now.AddDate(0, 1, 0)
	customerCard, err := mc.CreateCard(context.Background(), customer.AccountID, moov.CreateCard{
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
		if customerCard != nil {
			_ = mc.DisableCard(ctx, customer.AccountID, customerCard.CardID)
		}
	})

	fmt.Printf("\n\n%+v\n\n", customerCard.PaymentMethods)

	for _, pm := range customerCard.PaymentMethods {
		if pm.PaymentMethodType == moov.PaymentMethodType_CardPayment {
			env.CustomerPmId = pm.PaymentMethodID
		}
	}
	require.NotEmpty(t, env.CustomerPmId)

	return &env
}

func randomBankAccountNumber() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(999999999))
	return fmt.Sprintf("%d", 100000000+n.Int64())
}
