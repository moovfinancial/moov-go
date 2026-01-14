package moov_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/moovfinancial/moov-go/internal/testtools"
	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Schedules(t *testing.T) {
	ctx := BgCtx()
	mc := NewTestClient(t)

	// Just bumping now so we don't end up with a bunch more test transfers
	now := time.Date(2040, time.March, 1, 0, 0, 0, 0, time.UTC)

	// We'll start the monthly payments next month
	start := now.AddDate(0, 1, 0)

	partnerId := FACILITATOR_ID

	merchantAccountId := testtools.MERCHANT_ID
	merchantPmId := testtools.MERCHANT_WALLET_PM_ID

	customer := CreateTemporaryTestAccount(t, mc, createTestIndividualAccount())
	customerCard := createTemporaryCard(t, mc, customer.AccountID)
	customerPmId := customerCard.PaymentMethods[0].PaymentMethodID

	schedule, err := mc.CreateSchedule(ctx, partnerId, moov.CreateSchedule{
		Description: "a simple schedule",

		// Lets add a one time occurrence to handle the setup fee's
		Occurrences: []moov.CreateOccurrence{
			{
				RunOn: now,
				RunTransfer: moov.CreateRunTransfer{
					Description: "setup fee and first month example",
					Amount: moov.ScheduleAmount{
						Value:    200,
						Currency: "USD",
					},
					PartnerAccountID: partnerId,
					Source: moov.SchedulePaymentMethod{
						PaymentMethodID: customerPmId,
					},
					Destination: moov.SchedulePaymentMethod{
						PaymentMethodID: merchantPmId,
					},
				},
			},
		},

		// Setup a recurring transfer to handle repayment of say a loan with 6 periods
		Recur: &moov.CreateRecur{
			Start:          &start,
			RecurrenceRule: "FREQ=MONTHLY;COUNT=6",
			RunTransfer: moov.CreateRunTransfer{
				Description: "example of a description for all of the occurrences",
				Amount: moov.ScheduleAmount{
					Value:    100,
					Currency: "USD",
				},
				PartnerAccountID: partnerId,
				Source: moov.SchedulePaymentMethod{
					PaymentMethodID: customerPmId,
				},
				Destination: moov.SchedulePaymentMethod{
					PaymentMethodID: merchantPmId,
				},
			},
		},
	})

	t.Cleanup(func() {
		if schedule != nil {
			mc.CancelSchedule(ctx, partnerId, schedule.ScheduleID)
		}
	})

	t.Run("Created", func(t *testing.T) {
		require.NoError(t, err)
		require.NotNil(t, schedule)
		require.Len(t, schedule.Occurrences, 7)

		require.Equal(t, now, schedule.Occurrences[0].RunOn)
	})

	t.Run("get & list", func(t *testing.T) {

		getTests := []struct {
			name string
			id   string
		}{
			{"customer", customer.AccountID},
			{"partner", partnerId},
			{"merchant", merchantAccountId},
		}

		for _, party := range getTests {
			t.Run(fmt.Sprintf("Get & list from %s accountID", party.name), func(t *testing.T) {
				s2, err := mc.GetSchedule(ctx, party.id, schedule.ScheduleID)
				require.NoError(t, err)
				require.Len(t, s2.Occurrences, 7)
				require.Equal(t, schedule, s2)

				list, err := mc.ListSchedule(ctx, party.id, moov.Count(10), moov.Skip(0))
				require.NoError(t, err)
				require.Contains(t, list, *s2)
			})
		}

	})

	t.Run("Update", func(t *testing.T) {
		upsert := schedule.ToUpdateSchedule()
		require.Len(t, upsert.Occurrences, 7)

		// price increase for all recurring transfers
		upsert.Recur.RunTransfer.Amount.Value = 200

		// Lets shift an occurrence by 1 day
		upsert.Occurrences[0].RunOn = upsert.Occurrences[0].RunOn.AddDate(0, 0, 1)

		// Lets change the amount of the second
		upsert.Occurrences[1].RunTransfer.Amount.Value = 300

		// Add a new one time occurrence to charge a fee in an hour from now
		upsert.Occurrences = append(upsert.Occurrences, moov.UpdateOccurrence{
			RunOn: now.Add(time.Hour),
			RunTransfer: moov.CreateRunTransfer{
				Description: "late fee fine",
				Amount: moov.ScheduleAmount{
					Value:    1,
					Currency: "USD",
				},
				PartnerAccountID: partnerId,
				Source: moov.SchedulePaymentMethod{
					PaymentMethodID: customerPmId,
				},
				Destination: moov.SchedulePaymentMethod{
					PaymentMethodID: merchantPmId,
				},
			},
		})

		// Lets cancel the third occurrence
		upsert.Occurrences[2].Canceled = moov.PtrOf(true)

		// Do the update
		updated, err := mc.UpdateSchedule(ctx, partnerId, schedule.ScheduleID, upsert)
		require.NoError(t, err)

		PrettyDebug(t, updated)

		// occurrences are always ordered so our new one for the fee will be moved up.
		require.Equal(t, int64(1), updated.Occurrences[0].RunTransfer.Amount.Value)

		// check that the new date is set.
		require.Equal(t, upsert.Occurrences[0].RunOn, updated.Occurrences[1].RunOn)

		// occurrences are always ordered so our new one for the fee will be moved up.
		require.Equal(t, int64(300), updated.Occurrences[2].RunTransfer.Amount.Value)

		// Lets cancel the third occurrence
		require.NotNil(t, updated.Occurrences[3].CanceledOn)

		schedule = updated
	})

	t.Run("Cancel", func(t *testing.T) {
		err := mc.CancelSchedule(ctx, partnerId, schedule.ScheduleID)
		require.NoError(t, err)
	})
}

func Test_Schedules_LineItems(t *testing.T) {
	ctx := BgCtx()
	mc := NewTestClient(t)

	now := time.Date(2040, time.March, 1, 0, 0, 0, 0, time.UTC)
	start := now.AddDate(0, 1, 0)

	partnerId := FACILITATOR_ID
	merchantAccountId := testtools.MERCHANT_ID
	merchantPmId := testtools.MERCHANT_WALLET_PM_ID

	customer := CreateTemporaryTestAccount(t, mc, createTestIndividualAccount())
	customerCard := createTemporaryCard(t, mc, customer.AccountID)
	customerPmId := customerCard.PaymentMethods[0].PaymentMethodID

	// Upload test images
	_, imgReader1 := randomImage(t, 100, 100, encodePNG)
	image1, err := mc.UploadImage(ctx, merchantAccountId, imgReader1, &moov.ImageMetadataRequest{
		AltText: moov.PtrOf("Latte image"),
	})
	require.NoError(t, err)
	require.NotNil(t, image1)
	t.Cleanup(func() {
		mc.DeleteImage(ctx, merchantAccountId, image1.ImageID)
	})

	runTransfer := moov.CreateRunTransfer{
		Description: "recurring transfer",
		Amount: moov.ScheduleAmount{
			Value:    550,
			Currency: "USD",
		},
		PartnerAccountID: partnerId,
		Source: moov.SchedulePaymentMethod{
			PaymentMethodID: customerPmId,
		},
		Destination: moov.SchedulePaymentMethod{
			PaymentMethodID: merchantPmId,
		},
		LineItems: &moov.CreateScheduledTransferLineItems{
			Items: []moov.CreateScheduledTransferLineItem{
				{
					Name: "Cappuccino",
					BasePrice: moov.AmountDecimal{
						Currency:     "USD",
						ValueDecimal: "4.00",
					},
					Quantity:  1,
					ProductID: moov.PtrOf("11d58aa0-fb14-4aaf-ac04-8b7cfc282ca4"),
					ImageIDs:  []string{image1.ImageID},
					Options: []moov.CreateScheduledTransferLineItemOption{
						{
							Name:     "Oat Milk",
							Quantity: 1,
							PriceModifier: &moov.AmountDecimal{
								Currency:     "USD",
								ValueDecimal: "1.50",
							},
							Group:    moov.PtrOf("Milk"),
							ImageIDs: []string{image1.ImageID},
						},
					},
				},
			},
		},
	}

	// Schedules may take a moment to be aware of an uploaded image
	var schedule *moov.Schedule
	success := assert.Eventually(t, func() bool {
		var err error
		schedule, err = mc.CreateSchedule(ctx, partnerId, moov.CreateSchedule{
			Description: "schedule with line items",
			Occurrences: []moov.CreateOccurrence{
				{
					RunOn:       now,
					RunTransfer: runTransfer,
				},
			},
			Recur: &moov.CreateRecur{
				Start:          &start,
				RecurrenceRule: "FREQ=MONTHLY;COUNT=3",
				RunTransfer:    runTransfer,
			},
		})
		return err == nil
	}, 5*time.Second, 250*time.Millisecond)
	require.True(t, success, "failed to create schedule")
	require.NotNil(t, schedule)

	t.Cleanup(func() {
		if schedule != nil {
			mc.CancelSchedule(ctx, partnerId, schedule.ScheduleID)
		}
	})

	wantLineItems := &moov.ScheduledTransferLineItems{
		Items: []moov.ScheduledTransferLineItem{
			{
				Name: "Cappuccino",
				BasePrice: moov.AmountDecimal{
					Currency:     "USD",
					ValueDecimal: "4.00",
				},
				Quantity:  1,
				ProductID: moov.PtrOf("11d58aa0-fb14-4aaf-ac04-8b7cfc282ca4"),
				Images: []moov.ScheduledTransferImageMetadata{
					{
						ImageID:  image1.ImageID,
						AltText:  image1.AltText,
						Link:     image1.Link,
						PublicID: image1.PublicID,
					},
				},
				Options: []moov.ScheduledTransferLineItemOption{
					{
						Name:     "Oat Milk",
						Quantity: 1,
						PriceModifier: &moov.AmountDecimal{
							Currency:     "USD",
							ValueDecimal: "1.50",
						},
						Group: moov.PtrOf("Milk"),
						Images: []moov.ScheduledTransferImageMetadata{
							{
								ImageID:  image1.ImageID,
								AltText:  image1.AltText,
								Link:     image1.Link,
								PublicID: image1.PublicID,
							},
						},
					},
				},
			},
		},
	}

	t.Run("created with line items", func(t *testing.T) {
		require.NoError(t, err)
		require.NotNil(t, schedule)

		require.Equal(t, wantLineItems, schedule.Recur.RunTransfer.LineItems)
		for _, occur := range schedule.Occurrences {
			require.Equal(t, wantLineItems, occur.RunTransfer.LineItems)
		}
	})

	t.Run("update line items", func(t *testing.T) {
		upsert := schedule.ToUpdateSchedule()
		require.Len(t, upsert.Occurrences, 4) // 1 one-time + 3 recurring

		// Update the recur line items - change price
		require.NotNil(t, upsert.Recur)
		require.NotNil(t, upsert.Recur.RunTransfer.LineItems)
		upsert.Recur.RunTransfer.LineItems.Items[0].Name = "Latte"

		// Update first occurrence line items - add a new option
		require.NotNil(t, upsert.Occurrences[0].RunTransfer.LineItems)
		upsert.Occurrences[0].RunTransfer.LineItems.Items[0].Options = append(
			upsert.Occurrences[0].RunTransfer.LineItems.Items[0].Options,
			moov.CreateScheduledTransferLineItemOption{
				Name:     "Vanilla Syrup",
				Quantity: 1,
				PriceModifier: &moov.AmountDecimal{
					Currency:     "USD",
					ValueDecimal: "0.75",
				},
				Group:    moov.PtrOf("Flavors"),
				ImageIDs: []string{image1.ImageID},
			},
		)
		upsert.Occurrences[0].RunTransfer.Amount.Value = 625

		// Leave out other occurrences from the update - they should update from the change to recur
		upsert.Occurrences = upsert.Occurrences[:1]

		updated, err := mc.UpdateSchedule(ctx, partnerId, schedule.ScheduleID, upsert)
		require.NoError(t, err)
		require.NotNil(t, updated)

		// Verify recur line items were updated
		require.NotNil(t, updated.Recur)
		require.NotNil(t, updated.Recur.RunTransfer.LineItems)
		require.Equal(t, "Latte", updated.Recur.RunTransfer.LineItems.Items[0].Name)

		// Verify first occurrence has new option
		require.Len(t, updated.Occurrences[0].RunTransfer.LineItems.Items[0].Options, 2)
		require.Equal(t, "Vanilla Syrup", updated.Occurrences[0].RunTransfer.LineItems.Items[0].Options[1].Name)
		require.Equal(t, "0.75", updated.Occurrences[0].RunTransfer.LineItems.Items[0].Options[1].PriceModifier.ValueDecimal)
		require.Equal(t, moov.PtrOf("Flavors"), updated.Occurrences[0].RunTransfer.LineItems.Items[0].Options[1].Group)

		// Verify recurring occurrences reflect updated recur line items
		for i := 1; i < 4; i++ {
			require.NotNil(t, updated.Occurrences[i].RunTransfer.LineItems)
			require.Equal(t, "Latte", updated.Occurrences[i].RunTransfer.LineItems.Items[0].Name)
		}

		schedule = updated
	})
}
