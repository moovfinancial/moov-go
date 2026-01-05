package moov_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"

	"github.com/stretchr/testify/require"
)

func TestIssuedCardMarshal(t *testing.T) {
	input := []byte(`{
			"issuedCardID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43",
			"brand": "Visa",
			"lastFourCardNumber": "1234",
			"expiration": {
				"month": "01",
				"year": "21"
			},
			"authorizedUser": {
				"firstName": "Jules",
				"lastName": "Jackson"
			},
			"fundingWalletID": "50469144-f859-46dc-bdbd-9587c2fa7b42",
			"state": "active",
			"formFactor": "virtual",
			"createdOn": "2023-11-08T22:06:16Z"
		}`)

	card := new(moov.IssuedCard)

	dec := json.NewDecoder(bytes.NewReader(input))
	dec.DisallowUnknownFields()

	err := dec.Decode(&card)
	require.NoError(t, err)

	require.Equal(t, "ec7e1848-dc80-4ab0-8827-dd7fc0737b43", card.IssuedCardID)
}

func TestIssuedCardAuthorizationMarshal(t *testing.T) {
	input := []byte(`{
			"authorizationID": "220c75d3-fac6-4572-9379-a2a2fb29f8cf",
			"issuedCardID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43",
			"fundingWalletID": "50469144-f859-46dc-bdbd-9587c2fa7b42",
			"createdOn": "2023-11-08T23:06:16Z",
			"network": "visa",
			"authorizedAmount": "-1.23",
			"status": "cleared",
			"merchantData": {
				"networkID": "123456789012345",
				"name": "Moov Financial",
				"city": "Longmont",
				"country": "US",
				"postalCode": "80525",
				"state": "CO",
				"mcc": "5700"
			},
			"cardTransactions": [
				"86d4b88b-eb1b-4640-941c-d8f087256b90",
				"7e6d0485-bd19-4536-bee6-47d1a7c613ba"
			]
		}`)

	authorization := new(moov.IssuedCardAuthorization)

	dec := json.NewDecoder(bytes.NewReader(input))
	dec.DisallowUnknownFields()

	err := dec.Decode(&authorization)
	require.NoError(t, err)

	require.Equal(t, "220c75d3-fac6-4572-9379-a2a2fb29f8cf", authorization.AuthorizationID)
}

func TestIssuedCardAuthorizationEventMarshal(t *testing.T) {
	input := []byte(`{
			"eventID": "afc99714-7611-4136-bd17-a97f687274b7",
			"eventType": "authorization",
			"createdOn": "2023-11-08T23:06:16Z",
			"amount": "-1.23",
			"result": "approved"
		}`)

	event := new(moov.IssuedCardAuthorizationEvent)

	dec := json.NewDecoder(bytes.NewReader(input))
	dec.DisallowUnknownFields()

	err := dec.Decode(&event)
	require.NoError(t, err)

	require.Equal(t, "afc99714-7611-4136-bd17-a97f687274b7", event.EventID)
}

func TestIssuedCardTransactionMarshal(t *testing.T) {
	input := []byte(`{
			"cardTransactionID": "86d4b88b-eb1b-4640-941c-d8f087256b90",
			"issuedCardID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43",
			"fundingWalletID": "50469144-f859-46dc-bdbd-9587c2fa7b42",
			"amount": "-1.23",
			"authorizationID": "220c75d3-fac6-4572-9379-a2a2fb29f8cf",
			"createdOn": "2023-11-08T23:34:21Z",
			"authorizedOn": "2023-11-08T23:06:16Z",
			"merchantData": {
				"networkID": "123456789012345",
				"name": "Moov Financial",
				"city": "Longmont",
				"country": "US",
				"postalCode": "80525",
				"state": "CO",
				"mcc": "5700"
			}
		}`)

	transaction := new(moov.IssuedCardTransaction)

	dec := json.NewDecoder(bytes.NewReader(input))
	dec.DisallowUnknownFields()

	err := dec.Decode(&transaction)
	require.NoError(t, err)

	require.Equal(t, "86d4b88b-eb1b-4640-941c-d8f087256b90", transaction.CardTransactionID)
}

func Test_CardIssuing(t *testing.T) {
	mc := NewTestClient(t)
	memo := "testing"

	// create issued card
	created, err := mc.CreateIssuedCard(BgCtx(), MERCHANT_ID, moov.CreateIssuedCard{
		FundingWalletID: MERCHANT_WALLET_ID,
		AuthorizedUser: moov.CreateAuthorizedUser{
			FirstName: "John",
			LastName:  "Doe",
		},
		FormFactor: moov.IssuedCardFormFactor_Virtual,
		Memo:       &memo,
	})
	NoResponseError(t, err)
	require.NotNil(t, created)
	t.Cleanup(func() {
		require.NoError(t, closeIssuedCard(BgCtx(), mc, MERCHANT_ID, created.IssuedCardID))
	})

	// list issued cards
	cards, err := mc.ListIssuedCards(BgCtx(), MERCHANT_ID,
		moov.WithIssuedCardStates([]moov.IssuedCardState{
			moov.IssuedCardState_Active,
			moov.IssuedCardState_PendingVerification,
		}))
	NoResponseError(t, err)
	require.NotEmpty(t, cards)

	// update issued card
	closed := moov.UpdateIssuedCardState_Closed
	err = mc.UpdateIssuedCard(BgCtx(), MERCHANT_ID, created.IssuedCardID, moov.UpdateIssuedCard{
		State: &closed,
	})
	NoResponseError(t, err)

	// get issued card
	card, err := mc.GetIssuedCard(BgCtx(), MERCHANT_ID, created.IssuedCardID)
	NoResponseError(t, err)
	require.NotNil(t, card)
	require.Equal(t, created.IssuedCardID, card.IssuedCardID)
	require.Equal(t, string(closed), string(card.State))

	// list issued card authorizations
	authorizations, err := mc.ListIssuedCardAuthorizations(BgCtx(), MERCHANT_ID,
		moov.WithIssuedCardAuthorizationCardID(card.IssuedCardID))
	NoResponseError(t, err)
	require.Empty(t, authorizations)

	// list issued card transactions
	transactions, err := mc.ListIssuedCardTransactions(BgCtx(), MERCHANT_ID,
		moov.WithIssuedCardTransactionCardID(card.IssuedCardID))
	NoResponseError(t, err)
	require.Empty(t, transactions)
}

func closeIssuedCard(ctx context.Context, mc *moov.Client, accountID, cardID string) error {
	card, err := mc.GetIssuedCard(BgCtx(), MERCHANT_ID, cardID)
	if err != nil {
		return err
	}

	if card.State != moov.IssuedCardState_Closed {
		closed := moov.UpdateIssuedCardState_Closed
		update := moov.UpdateIssuedCard{
			State: &closed,
		}
		return mc.UpdateIssuedCard(ctx, accountID, cardID, update)
	}
	return nil
}
