package moov_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

func Test_Schedules(t *testing.T) {
	ctx := BgCtx()
	mc := NewTestClient(t)

	// Just bumping now so we don't end up with a bunch more test transfers
	now := time.Date(2040, time.March, 1, 0, 0, 0, 0, time.UTC)
	start := now.AddDate(0, 0, 1)

	partnerId := FACILITATOR_ID

	merchant := getLincolnBank(t, mc)
	bankAcct := createTemporaryBankAccount(t, mc, merchant.AccountID)
	merchantPmId := bankAcct.PaymentMethods[0].PaymentMethodID

	customer := CreateTemporaryTestAccount(t, mc, createTestIndividualAccount())
	customerCard := createTemporaryCard(t, mc, customer.AccountID)
	customerPmId := customerCard.PaymentMethods[0].PaymentMethodID

	schedule, err := mc.CreateSchedule(ctx, partnerId, moov.CreateSchedule{
		Description: "a simple schedule",

		// Setup a recurring transfer to handle repayment of say a loan with 6 periods
		Recur: &moov.Recur{
			Start:          &start,
			RecurrenceRule: "FREQ=MONTHLY;COUNT=6",
			RunTransfer: moov.RunTransfer{
				Description: "example of a description for all of the occurrences",
				Amount: moov.ScheduleAmount{
					Value:    100,
					Currency: "USD",
				},
				PartnerAccountID: FACILITATOR_ID,
				Source: moov.SchedulePaymentMethod{
					PaymentMethodID: customerPmId,
				},
				Destination: moov.SchedulePaymentMethod{
					PaymentMethodID: merchantPmId,
				},
			},
		},

		// Lets add a one time occurrence to handle the setup fee's
		Occurrences: []moov.CreateOccurrence{
			{
				RunOn: now,
				RunTransfer: moov.RunTransfer{
					Description: "setup fee example",
					Amount: moov.ScheduleAmount{
						Value:    200,
						Currency: "USD",
					},
					PartnerAccountID: FACILITATOR_ID,
					Source: moov.SchedulePaymentMethod{
						PaymentMethodID: customerPmId,
					},
					Destination: moov.SchedulePaymentMethod{
						PaymentMethodID: merchantPmId,
					},
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
	})

	t.Run("get & list", func(t *testing.T) {

		getTests := []struct {
			name string
			id   string
		}{
			{"customer", customer.AccountID},
			{"partner", partnerId},
			{"merchant", merchant.AccountID},
		}

		for _, party := range getTests {
			t.Run(fmt.Sprintf("Get & list from %s accountID", party.name), func(t *testing.T) {
				s2, err := mc.GetSchedule(ctx, party.id, schedule.ScheduleID)
				require.NoError(t, err)
				require.Len(t, s2.Occurrences, 7)
				require.Equal(t, schedule, s2)

				list, err := mc.ListSchedule(ctx, party.id)
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
			RunTransfer: moov.RunTransfer{
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
