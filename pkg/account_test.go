package moov

import (
	"bytes"
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	creds = Credentials{
		AccountID: "638481a5-5205-406c-84c7-2fc2239105d1",
		PublicKey: "Qo0j0ChSalMKmRI_",
		SecretKey: "junk",
		Domain:    "localhost",
	}
)

func TestAccount(t *testing.T) {
	input := []byte(`{"mode":"sandbox","accountID":"638481a5-5205-406c-84c7-2fc2239105d1","accountType":"individual","displayName":"Wade Arnold","profile":{"individual":{"name":{"firstName":"Wade","lastName":"Arnold"},"phone":{"number":"5555555555","countryCode":"1"},"birthDateProvided":false,"governmentIDProvided":false}},"verification":{"verificationStatus":"unverified","status":"unverified"},"foreignID":"your-correlation-id","createdOn":"2023-11-08T23:06:16.168497001Z","updatedOn":"2023-11-08T23:06:16.168497001Z"}`)

	account := new(Account)

	dec := json.NewDecoder(bytes.NewReader(input))
	dec.DisallowUnknownFields()

	err := dec.Decode(&account)
	if err != nil {
		require.NoError(t, err)
	}

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

	mc, err := NewClient(creds)
	if err != nil {
		log.Fatal(err)
	}

	account, err = mc.CreateAccount(account)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, account.AccountID)
}

func TestGetAccount(t *testing.T) {
	mc, err := NewClient(creds)
	if err != nil {
		log.Fatal(err)
	}

	account, err := mc.GetAccount("638481a5-5205-406c-84c7-2fc2239105d1")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "638481a5-5205-406c-84c7-2fc2239105d1", account.AccountID)
}
