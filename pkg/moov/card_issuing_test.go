package moov_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

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
			"authorizedUserAccountID": "a06f9j41-5e3e-4f4a-9b1e-2c3d4e5f6a7b",
			"nickname": "Travel card",
			"fundingWalletID": "50469144-f859-46dc-bdbd-9587c2fa7b42",
			"state": "active",
			"formFactor": "virtual",
			"billingAddress": {
				"addressLine1": "123 Main St",
				"city": "Longmont",
				"stateOrProvince": "CO",
				"postalCode": "80525",
				"country": "US"
			},
			"metadata": {
				"program": "rewards"
			},
			"controls": {
				"singleUse": false,
				"velocityLimits": [
					{
						"amount": 10000,
						"count": 5,
						"interval": "monthly",
						"amountUsed": 2500,
						"amountRemaining": 7500,
						"countUsed": 1,
						"countRemaining": 4,
						"resetsOn": "2023-12-01T00:00:00Z"
					}
				],
				"merchantCategoryRestrictions": {
					"mode": "block",
					"categories": ["gambling", "alcohol-and-bars"],
					"customMCCs": ["5811"],
					"exemptMerchants": [
						{
							"mid": "123456789012345",
							"descriptorPattern": "(?i)moov",
							"name": "Moov Financial"
						}
					]
				},
				"merchantRestrictions": {
					"mode": "allow",
					"merchants": [
						{
							"descriptorPattern": "(?i)amazon"
						}
					]
				},
				"allowedSchedule": {
					"timezone": "America/Denver",
					"windows": [
						{
							"days": ["monday", "tuesday"],
							"startTime": "09:00",
							"endTime": "17:00"
						}
					]
				},
				"expiresOn": "2024-01-01T00:00:00Z"
			},
			"createdOn": "2023-11-08T22:06:16Z"
		}`)

	card := new(moov.IssuedCard)

	dec := json.NewDecoder(bytes.NewReader(input))
	dec.DisallowUnknownFields()

	err := dec.Decode(&card)
	require.NoError(t, err)

	require.Equal(t, "ec7e1848-dc80-4ab0-8827-dd7fc0737b43", card.IssuedCardID)

	require.NotNil(t, card.Controls)
	require.Len(t, card.Controls.VelocityLimits, 1)

	vl := card.Controls.VelocityLimits[0]
	require.Equal(t, int32(10000), *vl.Amount)
	require.Equal(t, int32(5), *vl.Count)
	require.Equal(t, moov.IssuingIntervalLimit_Monthly, *vl.Interval)
	require.Equal(t, int32(2500), *vl.AmountUsed)
	require.Equal(t, int32(7500), *vl.AmountRemaining)
	require.Equal(t, int32(1), *vl.CountUsed)
	require.Equal(t, int32(4), *vl.CountRemaining)
	require.NotNil(t, vl.ResetsOn)

	require.NotNil(t, card.Controls.MerchantCategoryRestrictions)
	require.Equal(t, moov.IssuingControlsRestrictionMode_Block, card.Controls.MerchantCategoryRestrictions.Mode)
	require.Equal(t, []moov.IssuingMerchantCategory{
		moov.IssuingMerchantCategory_Gambling,
		moov.IssuingMerchantCategory_AlcoholAndBars,
	}, card.Controls.MerchantCategoryRestrictions.Categories)
	require.Equal(t, []string{"5811"}, card.Controls.MerchantCategoryRestrictions.CustomMCCs)
	require.Len(t, card.Controls.MerchantCategoryRestrictions.ExemptMerchants, 1)

	require.NotNil(t, card.Controls.MerchantRestrictions)
	require.Equal(t, moov.IssuingControlsRestrictionMode_Allow, card.Controls.MerchantRestrictions.Mode)

	require.NotNil(t, card.Controls.AllowedSchedule)
	require.Equal(t, "America/Denver", card.Controls.AllowedSchedule.Timezone)
	require.Len(t, card.Controls.AllowedSchedule.Windows, 1)
	require.Equal(t, []moov.IssuingScheduleDay{
		moov.IssuingScheduleDay_Monday,
		moov.IssuingScheduleDay_Tuesday,
	}, card.Controls.AllowedSchedule.Windows[0].Days)

	require.NotNil(t, card.Controls.ExpiresOn)
}

func TestCreateIssuedCardMarshal(t *testing.T) {
	create := moov.CreateIssuedCard{
		Controls: &moov.IssuingControls{
			SingleUse: moov.PtrOf(true),
			VelocityLimits: []moov.IssuingVelocityLimit{
				{
					Amount:   moov.PtrOf(int32(10000)),
					Count:    moov.PtrOf(int32(5)),
					Interval: moov.PtrOf(moov.IssuingIntervalLimit_Daily),
				},
			},
			MerchantCategoryRestrictions: &moov.MerchantCategoryRestrictions{
				Mode:            moov.IssuingControlsRestrictionMode_Block,
				Categories:      []moov.IssuingMerchantCategory{moov.IssuingMerchantCategory_Gambling},
				CustomMCCs:      []string{"5811"},
				ExemptMerchants: []moov.MerchantEntry{{Name: moov.PtrOf("Moov Financial")}},
			},
			MerchantRestrictions: &moov.MerchantRestrictions{
				Mode:      moov.IssuingControlsRestrictionMode_Allow,
				Merchants: []moov.MerchantEntry{{DescriptorPattern: moov.PtrOf("(?i)amazon")}},
			},
			AllowedSchedule: &moov.AllowedSchedule{
				Timezone: "America/Denver",
				Windows: []moov.ScheduleWindow{
					{
						Days:      []moov.IssuingScheduleDay{moov.IssuingScheduleDay_Monday},
						StartTime: "09:00",
						EndTime:   "17:00",
					},
				},
			},
			ExpiresOn: moov.PtrOf(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		},
	}

	out, err := json.Marshal(create)
	require.NoError(t, err)

	s := string(out)
	// Verify the request serializes with the API's lowerCamelCase keys (guards against
	// unnamed/PascalCase JSON tags on the new controls types).
	for _, key := range []string{
		`"controls"`,
		`"velocityLimits"`,
		`"count":5`,
		`"merchantCategoryRestrictions"`,
		`"categories"`,
		`"customMCCs"`,
		`"exemptMerchants"`,
		`"merchantRestrictions"`,
		`"allowedSchedule"`,
		`"startTime"`,
		`"expiresOn"`,
	} {
		require.Contains(t, s, key)
	}
}

func TestUpdateIssuingControlsMarshal(t *testing.T) {
	t.Run("unset fields are omitted", func(t *testing.T) {
		out, err := json.Marshal(moov.UpdateIssuedCard{
			Controls: &moov.UpdateIssuingControls{},
		})
		require.NoError(t, err)
		require.JSONEq(t, `{"controls":{}}`, string(out))
	})

	t.Run("clear velocity limits with empty slice", func(t *testing.T) {
		out, err := json.Marshal(moov.UpdateIssuingControls{
			VelocityLimits: &[]moov.IssuingVelocityLimit{},
		})
		require.NoError(t, err)
		require.JSONEq(t, `{"velocityLimits":[]}`, string(out))
	})

	t.Run("explicit null removes a field", func(t *testing.T) {
		out, err := json.Marshal(moov.UpdateIssuingControls{
			MerchantRestrictions: moov.SetNull[moov.MerchantRestrictions](),
			ExpiresOn:            moov.SetNull[time.Time](),
		})
		require.NoError(t, err)
		require.JSONEq(t, `{"merchantRestrictions":null,"expiresOn":null}`, string(out))
	})

	t.Run("set replaces a field with a value", func(t *testing.T) {
		out, err := json.Marshal(moov.UpdateIssuingControls{
			MerchantCategoryRestrictions: moov.Set(moov.MerchantCategoryRestrictions{
				Mode:       moov.IssuingControlsRestrictionMode_Allow,
				Categories: []moov.IssuingMerchantCategory{moov.IssuingMerchantCategory_Groceries},
			}),
		})
		require.NoError(t, err)
		require.JSONEq(t, `{"merchantCategoryRestrictions":{"mode":"allow","categories":["groceries"]}}`, string(out))
	})
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
	nickname := "testing"

	// create issued card
	// The card is issued to (and funded by the default wallet of) the MERCHANT_ID business
	// account. AuthorizedUserAccountID (an optional second cardholder) is omitted.
	created, err := mc.CreateIssuedCard(BgCtx(), MERCHANT_ID, moov.CreateIssuedCard{
		Nickname: &nickname,
		Metadata: map[string]string{"program": "testing"},
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
		}))
	NoResponseError(t, err)
	require.NotEmpty(t, cards)

	// update issued card
	closed := moov.UpdateIssuedCardState_Closed
	card, err := mc.UpdateIssuedCard(BgCtx(), MERCHANT_ID, created.IssuedCardID, moov.UpdateIssuedCard{
		State: &closed,
	})
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
		_, err := mc.UpdateIssuedCard(ctx, accountID, cardID, update)
		return err
	}
	return nil
}
