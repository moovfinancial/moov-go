package moov_test

// TODO: Create faililng test for other card brands in test mode
// https://docs.moov.io/guides/get-started/test-mode/#cards

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCardMarshal(t *testing.T) {
	input := []byte(`{
		  "cardID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43",
		  "fingerprint": "9948962d92a1ce40c9f918cd9ece3a22bde62fb325a2f1fe2e833969de672ba3",
		  "brand": "Discover",
		  "cardType": "debit",
		  "lastFourCardNumber": "1234",
		  "bin": "123456",
		  "expiration": {
			"month": "01",
			"year": "21"
		  },
		  "holderName": "Jules Jackson",
		  "billingAddress": {
			"addressLine1": "123 Main Street",
			"addressLine2": "Apt 302",
			"city": "Boulder",
			"stateOrProvince": "CO",
			"postalCode": "80301",
			"country": "US"
		  },
		  "cardVerification": {
			"cvv": "match",
			"addressLine1": "match",
			"postalCode": "match"
		  },
		  "issuer": "GRINGOTTS BANK",
		  "issuerCountry": "US",
		  "cardOnFile": true,
		  "merchantAccountID": "50469144-f859-46dc-bdbd-9587c2fa7b42",
		  "cardAccountUpdater": {
			"updatedOn": "2019-08-24T14:15:22Z",
			"updateType": "number-update"
		  },
		  "domesticPushToCard": "fast-funds"
		}`)

	card := new(moov.Card)

	dec := json.NewDecoder(bytes.NewReader(input))
	dec.DisallowUnknownFields()

	err := dec.Decode(&card)
	if err != nil {
		require.NoError(t, err)
	}
	assert.Equal(t, "ec7e1848-dc80-4ab0-8827-dd7fc0737b43", card.CardID)
}

func TestCards(t *testing.T) {
	// Tests full lifecycle of card resource - create, get, update, list, disable

	mc := NewTestClient(t)

	exp := time.Now().UTC().AddDate(0, 7, 0)
	expMonth := exp.Format("01")
	expYear := exp.Format("06")

	account := CreateTemporaryTestAccount(t, mc, createTestIndividualAccount())

	billingAddress := moov.Address{
		AddressLine1:    "123 Main Street",
		City:            "City",
		StateOrProvince: "CO",
		PostalCode:      "12345",
		Country:         "US",
	}

	// Create card
	card, err := mc.CreateCard(context.Background(), account.AccountID, moov.CreateCard{
		CardNumber: "4111111111111111",
		CardCvv:    "123",
		Expiration: moov.Expiration{
			Month: expMonth,
			Year:  expYear,
		},
		HolderName:     "john doe",
		BillingAddress: billingAddress,
		CardOnFile:     false,
	})
	require.NoError(t, err)
	require.False(t, card.CardOnFile)
	require.NotEmpty(t, card.PaymentMethods)

	// Get card
	created, err := mc.GetCard(context.Background(), account.AccountID, card.CardID)
	require.NoError(t, err)
	require.Equal(t, card.CardID, created.CardID)

	// Update card
	updated, err := mc.UpdateCard(context.Background(), account.AccountID, card.CardID,
		moov.WithCardOnFile(true),
	)
	require.NoError(t, err)
	require.True(t, updated.CardOnFile)

	// make sure other updatable fields didn't get updated unexpectedly
	require.Equal(t, billingAddress, updated.BillingAddress)
	require.Equal(t, expMonth, updated.Expiration.Month)
	require.Equal(t, expYear, updated.Expiration.Year)

	// List cards
	cards, err := mc.ListCards(context.Background(), account.AccountID)
	require.NoError(t, err)
	require.Len(t, cards, 1)

	// Disable card
	err = mc.DisableCard(context.Background(), account.AccountID, card.CardID)
	require.NoError(t, err)

	// List cards (omits disabled)
	cards, err = mc.ListCards(context.Background(), account.AccountID)
	require.NoError(t, err)
	require.Len(t, cards, 0)
}

func TestCard_UpdateBillingAddress(t *testing.T) {
	mc := NewTestClient(t)
	tc := newTestCard(t, mc)
	originalAddress := tc.card.BillingAddress

	// update a single billing address field to ensure the partial update works
	updatedCard, err := mc.UpdateCard(context.Background(), tc.account.AccountID, tc.card.CardID,
		moov.WithCardBillingAddress(moov.AddressPatch{AddressLine1: moov.PtrOf("updated street")}),
	)
	require.NoError(t, err)

	// make sure only the expected field was updated
	require.Equal(t, "updated street", updatedCard.BillingAddress.AddressLine1)
	require.Equal(t, originalAddress.AddressLine2, updatedCard.BillingAddress.AddressLine2)
	require.Equal(t, originalAddress.City, updatedCard.BillingAddress.City)
	require.Equal(t, originalAddress.StateOrProvince, updatedCard.BillingAddress.StateOrProvince)
	require.Equal(t, originalAddress.PostalCode, updatedCard.BillingAddress.PostalCode)
	require.Equal(t, originalAddress.Country, updatedCard.BillingAddress.Country)
}

type testCard struct {
	account moov.Account
	card    moov.Card
}

func newTestCard(t *testing.T, mc *moov.Client) testCard {
	account := CreateTemporaryTestAccount(t, mc, createTestIndividualAccount())

	exp := time.Now().UTC().AddDate(0, 7, 0)

	card, err := mc.CreateCard(context.Background(), account.AccountID, moov.CreateCard{
		CardNumber: "5555555555554444", // test mode Mastercard PAN
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
		require.NoError(t, mc.DisableCard(context.Background(), account.AccountID, card.CardID))
	})

	return testCard{
		account: *account,
		card:    *card,
	}
}
