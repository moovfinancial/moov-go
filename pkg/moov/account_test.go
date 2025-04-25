package moov_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountMarshalResponse(t *testing.T) {
	input := []byte(`{"mode":"sandbox","accountID":"638481a5-5205-406c-84c7-2fc2239105d1","accountType":"individual","displayName":"Wade Arnold","profile":{"individual":{"name":{"firstName":"Wade","lastName":"Arnold"},"phone":{"number":"5555555555","countryCode":"1"},"birthDateProvided":false,"governmentIDProvided":false}},"verification":{"verificationStatus":"unverified","status":"unverified"},"foreignID":"your-correlation-id","createdOn":"2023-11-08T23:06:16.168497001Z","updatedOn":"2023-11-08T23:06:16.168497001Z"}`)

	account := new(moov.Account)

	dec := json.NewDecoder(bytes.NewReader(input))
	dec.DisallowUnknownFields()

	err := dec.Decode(&account)
	require.NoError(t, err)

	assert.Equal(t, "Wade", account.Profile.Individual.Name.FirstName)
	t.Logf("%#v", account)
}

func TestCreateAccountIndividual(t *testing.T) {
	account := createTestIndividualAccount()

	mc := NewTestClient(t)

	completed, started, err := mc.CreateAccount(context.Background(), account)
	NoResponseError(t, err)
	require.NotNil(t, completed)
	require.Nil(t, started)
}

func TestGetAccount(t *testing.T) {
	mc := NewTestClient(t)

	account, err := mc.GetAccount(context.Background(), "638481a5-5205-406c-84c7-2fc2239105d1")
	NoResponseError(t, err)
	require.Equal(t, "638481a5-5205-406c-84c7-2fc2239105d1", account.AccountID)
}

func TestUpdateIndividualAccount(t *testing.T) {
	mc := NewTestClient(t)

	account := moov.Account{
		AccountID:   "aa19c3a7-4c72-4f64-adfa-9069c80d81cf",
		AccountType: moov.AccountType_Individual,
		Profile: moov.Profile{
			Individual: &moov.Individual{
				Name: moov.Name{
					FirstName: "David",
					LastName:  "Arnold",
				},
				Email: "Wade@arnold.com",
				Phone: &moov.Phone{
					Number:      "333-333-3333",
					CountryCode: "1",
				},
			},
		},
	}

	result, err := mc.UpdateAccount(context.Background(), account)
	NoResponseError(t, err)
	require.Equal(t, "David", result.Profile.Individual.Name.FirstName)
}

func TestPatchIndividualAccount(t *testing.T) {
	mc := NewTestClient(t)
	ctx := BgCtx()

	account, _, err := mc.CreateAccount(ctx, createTestIndividualAccount())
	NoResponseError(t, err)
	require.NotNil(t, account)

	patchAccount := moov.PatchAccount{
		Profile: moov.PatchProfile{
			Individual: &moov.PatchIndividualProfile{
				Name: moov.Name{
					FirstName: "John",
					LastName:  "Doe",
				},
				Phone: &moov.Phone{
					Number:      "555-555-6666",
					CountryCode: "1",
				},
				Email: "john.doe@moov.io",
				Address: &moov.Address{
					AddressLine1:    "123 Main St",
					City:            "Moov City",
					StateOrProvince: "CA",
					PostalCode:      "12345",
					Country:         "US",
				},
				BirthDate: &moov.Date{
					Year:  1990,
					Month: 1,
					Day:   1,
				},
				GovernmentID: &moov.GovernmentID{
					SSN: &moov.SSN{
						Full: "123-45-6789",
					},
				},
			},
		},
	}

	result, err := mc.PatchAccount(context.Background(), account.AccountID, patchAccount)
	NoResponseError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.Profile.Individual)
	require.Nil(t, result.Profile.Business)
	require.NotNil(t, result.Profile.Individual.Phone)
	require.Equal(t, "John", result.Profile.Individual.Name.FirstName)
	require.Equal(t, "Doe", result.Profile.Individual.Name.LastName)
	require.Equal(t, "5555556666", result.Profile.Individual.Phone.Number)
	require.Equal(t, "1", result.Profile.Individual.Phone.CountryCode)
	require.Equal(t, "john.doe@moov.io", result.Profile.Individual.Email)
	require.Equal(t, "123 Main St", result.Profile.Individual.Address.AddressLine1)
	require.Equal(t, "Moov City", result.Profile.Individual.Address.City)
	require.Equal(t, "CA", result.Profile.Individual.Address.StateOrProvince)
	require.Equal(t, "12345", result.Profile.Individual.Address.PostalCode)
	require.Equal(t, "US", result.Profile.Individual.Address.Country)
	require.True(t, result.Profile.Individual.GovernmentIDProvided)
	require.True(t, result.Profile.Individual.BirthDateProvided)
}

func TestPatchBusinessAccount(t *testing.T) {
	mc := NewTestClient(t)
	ctx := BgCtx()

	account, _, err := mc.CreateAccount(ctx, createTestBusinessAccount())
	NoResponseError(t, err)
	require.NotNil(t, account)

	t.Run("foreignID, customer support and account settings", func(t *testing.T) {
		foreignID := uuid.NewString()
		patchAccount := moov.PatchAccount{
			ForeignID: foreignID,
			CustomerSupport: &moov.CustomerSupport{
				Email: "test@moov.io",
				Phone: &moov.Phone{
					Number:      "555-555-5555",
					CountryCode: "1",
				},
				Website: "https://moov.io",
				Address: &moov.Address{
					AddressLine1:    "123 Main St",
					City:            "Moov City",
					StateOrProvince: "CA",
					PostalCode:      "12345",
					Country:         "US",
				},
			},
			AccountSettings: &moov.AccountSettings{
				CardPayment: &moov.CardPaymentSettings{
					StatementDescriptor: "Moov",
				},
				AchPayment: &moov.AchPaymentSettings{
					CompanyName: "Moov",
				},
			},
		}

		validate := func(result *moov.Account) {
			require.NotNil(t, result)
			require.NotNil(t, result.CustomerSupport)
			require.NotNil(t, result.CustomerSupport.Phone)
			require.NotNil(t, result.CustomerSupport.Address)
			require.NotNil(t, result.Settings)
			require.NotNil(t, result.Settings.CardPayment)
			require.NotNil(t, result.Settings.AchPayment)
			require.Equal(t, foreignID, result.ForeignID)
			require.Equal(t, "5555555555", result.CustomerSupport.Phone.Number)
			require.Equal(t, "1", result.CustomerSupport.Phone.CountryCode)
			require.Equal(t, "https://moov.io", result.CustomerSupport.Website)
			require.Equal(t, "123 Main St", result.CustomerSupport.Address.AddressLine1)
			require.Equal(t, "Moov City", result.CustomerSupport.Address.City)
			require.Equal(t, "CA", result.CustomerSupport.Address.StateOrProvince)
			require.Equal(t, "12345", result.CustomerSupport.Address.PostalCode)
			require.Equal(t, "US", result.CustomerSupport.Address.Country)
			require.Equal(t, "Moov", result.Settings.CardPayment.StatementDescriptor)
			require.Equal(t, "Moov", result.Settings.AchPayment.CompanyName)
		}

		result, err := mc.PatchAccount(context.Background(), account.AccountID, patchAccount)
		NoResponseError(t, err)
		validate(result)

		result, err = mc.GetAccount(ctx, account.AccountID)
		NoResponseError(t, err)
		validate(result)
	})

	t.Run("business profile", func(t *testing.T) {
		patchAccount := moov.PatchAccount{
			Profile: moov.PatchProfile{
				Business: &moov.PatchBusinessProfile{
					Name: "Moov Inc.",
					Phone: &moov.Phone{
						Number:      "555-555-5555",
						CountryCode: "1",
					},
					Address: &moov.Address{
						AddressLine1:    "123 Main St",
						City:            "Moov City",
						StateOrProvince: "CA",
						PostalCode:      "12345",
						Country:         "US",
					},
					Website: "https://moov.io",
					TaxID:   &moov.TaxID{EIN: moov.EIN{Number: "123456789"}},
				},
			},
		}

		validate := func(result *moov.Account) {
			require.NotNil(t, result)
			require.NotNil(t, result.Profile.Business)
			require.NotNil(t, result.Profile.Business.Phone)
			require.NotNil(t, result.Profile.Business.Address)
			require.Equal(t, "123 Main St", result.Profile.Business.Address.AddressLine1)
			require.Equal(t, "Moov City", result.Profile.Business.Address.City)
			require.Equal(t, "CA", result.Profile.Business.Address.StateOrProvince)
			require.Equal(t, "12345", result.Profile.Business.Address.PostalCode)
			require.Equal(t, "US", result.Profile.Business.Address.Country)
			require.Equal(t, "5555555555", result.Profile.Business.Phone.Number)
			require.Equal(t, "1", result.Profile.Business.Phone.CountryCode)
			require.Equal(t, "https://moov.io", result.Profile.Business.Website)
			require.Equal(t, "Moov Inc.", result.Profile.Business.LegalBusinessName)
			require.True(t, result.Profile.Business.TaxIDProvided)
		}

		result, err := mc.PatchAccount(context.Background(), account.AccountID, patchAccount)
		NoResponseError(t, err)
		validate(result)

		result, err = mc.GetAccount(ctx, account.AccountID)
		NoResponseError(t, err)
		validate(result)

	})

	t.Run("MetaData", func(t *testing.T) {
		patchAccount := moov.PatchAccount{
			Metadata: map[string]string{
				"key1": "value1",
			},
		}

		result, err := mc.PatchAccount(context.Background(), account.AccountID, patchAccount)
		NoResponseError(t, err)
		require.NotNil(t, result)
		require.Contains(t, result.Metadata, "key1")

		result, err = mc.GetAccount(ctx, account.AccountID)
		NoResponseError(t, err)
		require.NotNil(t, result)
		require.Contains(t, result.Metadata, "key1")
	})
}

func TestListAccounts(t *testing.T) {
	mc := NewTestClient(t)

	accounts, err := mc.ListAccounts(context.Background())
	NoResponseError(t, err)
	require.NotNil(t, accounts)
}

func TestDisconnectAccount(t *testing.T) {
	mc := NewTestClient(t)
	ctx := BgCtx()

	accnt, _, err := mc.CreateAccount(ctx, createTestIndividualAccount())
	NoResponseError(t, err)
	require.NotNil(t, accnt)

	_, err = mc.GetAccount(ctx, accnt.AccountID)
	NoResponseError(t, err)

	err = mc.DisconnectAccount(ctx, accnt.AccountID)
	NoResponseError(t, err)
}
