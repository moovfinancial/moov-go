package moov

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountMarshalResponse(t *testing.T) {
	input := []byte(`{"mode":"sandbox","accountID":"638481a5-5205-406c-84c7-2fc2239105d1","accountType":"individual","displayName":"Wade Arnold","profile":{"individual":{"name":{"firstName":"Wade","lastName":"Arnold"},"phone":{"number":"5555555555","countryCode":"1"},"birthDateProvided":false,"governmentIDProvided":false}},"verification":{"verificationStatus":"unverified","status":"unverified"},"foreignID":"your-correlation-id","createdOn":"2023-11-08T23:06:16.168497001Z","updatedOn":"2023-11-08T23:06:16.168497001Z"}`)

	account := new(Account)

	dec := json.NewDecoder(bytes.NewReader(input))
	dec.DisallowUnknownFields()

	err := dec.Decode(&account)
	require.NoError(t, err)

	assert.Equal(t, "Wade", account.Profile.Individual.Name.FirstName)
	t.Logf("%#v", account)
}

func TestCreateAccountIndividual(t *testing.T) {
	account := Account{
		AccountType: INDIVIDUAL,
		Profile: Profile{
			Individual: Individual{
				Name: Name{
					FirstName: "Wade",
					LastName:  "Arnold",
				},
				Email: "Wade@wadearnold.com",
				Phone: Phone{
					Number:      "555-555-5555",
					CountryCode: "1",
				},
			},
		},
	}

	mc := NewTestClient(t)

	account, err := mc.CreateAccount(account)
	require.NoError(t, err)

	assert.NotNil(t, account.AccountID)
}

func TestGetAccount(t *testing.T) {
	mc := NewTestClient(t)

	account, err := mc.GetAccount("638481a5-5205-406c-84c7-2fc2239105d1")
	require.NoError(t, err)

	require.Equal(t, "638481a5-5205-406c-84c7-2fc2239105d1", account.AccountID)
}

func TestUpdateAccount(t *testing.T) {
	mc := NewTestClient(t)

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

	account, err := mc.UpdateAccount(account)
	require.NoError(t, err)

	require.Equal(t, "David", account.Profile.Individual.Name.FirstName)
}

func TestListAccounts(t *testing.T) {
	mc := NewTestClient(t)

	accounts, err := mc.ListAccounts()
	require.NoError(t, err)

	require.NotNil(t, accounts)
}

/* func TestDeleateAccount(t *testing.T) {
	mc, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}

	err = mc.DeleteAccount("30942467-cd58-4a97-9d02-91b2555640c6")
	if err != nil {
		t.Fatal(err)
		assert.Error(t, err)
	}

} */
