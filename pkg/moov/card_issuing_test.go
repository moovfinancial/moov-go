package moov_test

import (
	"bytes"
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
			"createdOn": "2023-11-08T23:06:16Z"
		}`)

	card := new(moov.IssuedCard)

	dec := json.NewDecoder(bytes.NewReader(input))
	dec.DisallowUnknownFields()

	err := dec.Decode(&card)
	require.NoError(t, err)

	require.Equal(t, "ec7e1848-dc80-4ab0-8827-dd7fc0737b43", card.IssuedCardID)
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

	// list issued cards
	requested, err := mc.ListIssuedCards(BgCtx(), MERCHANT_ID,
		moov.WithIssuedCardStates([]moov.IssuedCardState{
			moov.IssuedCardState_Active,
			moov.IssuedCardState_PendingVerification,
		}))
	NoResponseError(t, err)
	require.NotEmpty(t, requested)

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
}
