package moov

// TODO: Create faililng test for other card brands in test mode
// https://docs.moov.io/guides/get-started/test-mode/#cards

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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

type CardTestSuite struct {
	suite.Suite
	accountID    string
	cardID       string
	deleteCardID string
	cards        []string
}

func TestCardSuite(t *testing.T) {
	suite.Run(t, new(CardTestSuite))
}

func (s *CardTestSuite) SetupSuite() {
	// Sandbox accounts have a "Lincoln National Corporation" moov account added by default. Get it's AccountID so we can test against it
	mc := NewTestClient(s.T())

	accounts, err := mc.ListAccounts(WithAccountName("Lincoln National Corporation"))
	s.NoError(err)

	for _, account := range accounts {
		if account.DisaplayName == "Lincoln National Corporation" {
			// set the accountID for testing
			s.accountID = account.AccountID
		}
	}

	card := CardPost{
		CardNumber: "371111111111114",
		CardCvv:    "1234",
		Expiration: Expiration{
			Month: "10",
			Year:  "28",
		},
		HolderName: "Wade Arnold",
		BillingAddress: Address{
			AddressLine1:    "123 Main Street",
			City:            "Golden",
			StateOrProvince: "CO",
			PostalCode:      "80401",
			Country:         "US",
		},
	}

	respCard, err := mc.CreateCard(s.accountID, card)
	s.NoError(err, "Error creating card")
	s.cardID = respCard.CardID
	s.deleteCardID = respCard.CardID
	s.cards = append(s.cards, respCard.CardID)
}

func (s *CardTestSuite) TearDownSuite() {
	mc := NewTestClient(s.T())

	// delete the cards we created
	for _, cardID := range s.cards {
		if cardID != "" {
			err := mc.DisableCard(s.accountID, cardID)
			s.NoError(err)
		}
	}
}

func (s *CardTestSuite) TestCreateCard() {
	card := CardPost{
		CardNumber: "4111111111111111",
		CardCvv:    "123",
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
		CardOnFile: false,
	}

	mc := NewTestClient(s.T())

	respCard, err := mc.CreateCard(s.accountID, card)
	s.NoError(err, "Error creating card")

	assert.NotEmpty(s.T(), respCard.CardID)
	s.cardID = respCard.CardID
	s.cards = append(s.cards, respCard.CardID)
}

func (s *CardTestSuite) TestListCards() {
	mc := NewTestClient(s.T())

	cards, err := mc.ListCards(s.accountID)
	s.NoError(err)

	assert.NotNil(s.T(), cards)
}

func (s *CardTestSuite) TestGetCard() {
	mc := NewTestClient(s.T())

	card, err := mc.GetCard(s.accountID, s.cardID)
	s.NoError(err)

	s.Equal(s.cardID, card.CardID)
}
func (s *CardTestSuite) TestUpdateCardBillingAddress() {
	mc := NewTestClient(s.T())
	billingAddress := Address{
		AddressLine1:    "125 Main Street",
		AddressLine2:    "Apt 302",
		City:            "Boulder",
		StateOrProvince: "CO",
		PostalCode:      "80303",
		Country:         "US",
	}

	updatedCard, err := mc.UpdateCard(s.accountID, s.cardID, WithCardBillingAddress(billingAddress))
	s.NoError(err)
	s.Equal(billingAddress, updatedCard.BillingAddress)
	// TODO: This should be "match" but isn't implemented in Moov's test mode and needs a server side fix
	s.Equal("unavailable", updatedCard.CardVerification.AddressLine1)
}

func (s *CardTestSuite) TestUpdateCardExpiration() {
	mc := NewTestClient(s.T())
	exp := Expiration{
		Month: "01",
		Year:  "28",
	}

	updatedCard, err := mc.UpdateCard(s.accountID, s.cardID, WithCardExpiration(exp))
	s.NoError(err)
	s.Equal(exp, updatedCard.Expiration)
}

func (s *CardTestSuite) TestUpdateCardCVV() {
	mc := NewTestClient(s.T())
	updatedCard, err := mc.UpdateCard(s.accountID, s.cardID, WithCardCVV("987"))
	s.NoError(err)
	// TODO: This should be "match" but isn't implemented in Moov's test mode and needs a server side fix
	s.Equal("unavailable", updatedCard.CardVerification.Cvv)
}

func (s *CardTestSuite) TestUpdateMultipleFilters() {
	mc := NewTestClient(s.T())
	updatedCard, err := mc.UpdateCard(s.accountID, s.cardID, WithCardOnFile(true), WithCardCVV("666"))
	s.NoError(err)
	s.True(updatedCard.CardOnFile)
}

func (s *CardTestSuite) TestDisableCard() {
	mc := NewTestClient(s.T())
	err := mc.DisableCard(s.accountID, s.cardID)
	s.NoError(err)
}
