package moov_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"

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
	moov.DebugPrintResponse(err, fmt.Printf)

	require.NoError(t, err)
	require.NotNil(t, completed)
	require.Nil(t, started)
}

func TestGetAccount(t *testing.T) {
	mc := NewTestClient(t)

	account, err := mc.GetAccount(context.Background(), "638481a5-5205-406c-84c7-2fc2239105d1")

	moov.DebugPrintResponse(err, fmt.Printf)
	require.NoError(t, err)
	require.Equal(t, "638481a5-5205-406c-84c7-2fc2239105d1", account.AccountID)
}

func TestUpdateAccount(t *testing.T) {
	mc := NewTestClient(t)

	account := moov.Account{
		AccountID:   "aa19c3a7-4c72-4f64-adfa-9069c80d81cf",
		AccountType: moov.ACCOUNTTYPE_INDIVIDUAL,
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

	moov.DebugPrintResponse(err, fmt.Printf)
	require.NoError(t, err)
	require.Equal(t, "David", result.Profile.Individual.Name.FirstName)
}

func TestListAccounts(t *testing.T) {
	mc := NewTestClient(t)

	accounts, err := mc.ListAccounts(context.Background())
	moov.DebugPrintResponse(err, fmt.Printf)

	require.NoError(t, err)
	require.NotNil(t, accounts)
}

func TestDisconnectAccount(t *testing.T) {
	mc := NewTestClient(t)
	ctx := BgCtx()

	accnt, _, err := mc.CreateAccount(ctx, createTestIndividualAccount())
	moov.DebugPrintResponse(err, fmt.Printf)

	require.NoError(t, err)
	require.NotNil(t, accnt)

	_, err = mc.GetAccount(ctx, accnt.AccountID)
	moov.DebugPrintResponse(err, fmt.Printf)

	require.NoError(t, err)

	err = mc.DisconnectAccount(ctx, accnt.AccountID)
	moov.DebugPrintResponse(err, fmt.Printf)

	require.NoError(t, err)
}
