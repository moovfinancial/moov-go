package moov

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
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

	card := new(Card)

	dec := json.NewDecoder(bytes.NewReader(input))
	dec.DisallowUnknownFields()

	err := dec.Decode(&card)
	if err != nil {
		require.NoError(t, err)
	}
	assert.Equal(t, "ec7e1848-dc80-4ab0-8827-dd7fc0737b43", card.CardID)
}

func TestCreateCard(t *testing.T) {
	card := Card{
		Expiration: Expiration{
			Month: "01",
			Year:  "28",
		},
		HolderName: "Jules Jackson",
		BillingAddress: Address{
			AddressLine1:    "123 Main Street",
			AddressLine2:    "Apt 302",
			City:            "Boulder",
			StateOrProvince: "CO",
			PostalCode:      "80301",
			Country:         "US",
		},
		CardOnFile:        false,
		MerchantAccountID: "50469144-f859-46dc-bdbd-9587c2fa7b41",
	}

	mc, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}

	card, err = mc.CreateCard(accountID, card, "4151066396992786", "937")
	if err != nil {
		t.Fatalf("Error creating card: %v", err)
	}

	assert.NotNil(t, card.CardID)
}

func TestListCards(t *testing.T) {
	mc, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}

	cards, err := mc.ListCards(accountID)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(cards))
	assert.NotNil(t, cards)
}

func TestGetCard(t *testing.T) {
	mc, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}
	// sample card id
	cardID := "ec7e1848-dc80-4ab0-8827-dd7fc0737b43"

	card, err := mc.GetCard(accountID, cardID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, cardID, card.CardID)
}

func TestUpdateCard(t *testing.T) {
	card := Card{
		Expiration: Expiration{
			Month: "01",
			Year:  "28",
		},
		BillingAddress: Address{
			AddressLine1:    "123 Main Street",
			AddressLine2:    "Apt 302",
			City:            "Boulder",
			StateOrProvince: "CO",
			PostalCode:      "80301",
			Country:         "US",
		},
		CardOnFile: false,
	}

	mc, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}

	// sample card id
	cardID := "ec7e1848-dc80-4ab0-8827-dd7fc0737b43"

	card, err = mc.UpdateCard(accountID, cardID, card, "937")
	if err != nil {
		t.Fatalf("Error creating card: %v", err)
	}

	assert.NotNil(t, card.CardID)
}

func TestDisableCard(t *testing.T) {
	mc, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}

	// sample card id
	cardID := "ec7e1848-dc80-4ab0-8827-dd7fc0737b43"

	err = mc.DisableCard(accountID, cardID)
	if err != nil {
		assert.Error(t, err)
	}
}
