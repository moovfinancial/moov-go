package moov_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/moov/schedules"
	"github.com/stretchr/testify/require"
)

func Test_Schedules(t *testing.T) {
	ctx := BgCtx()
	mc := NewTestClient(t)

	// Just bumping now so we don't end up with a bunch more test transfers
	now := time.Date(2040, time.March, 1, 0, 0, 0, 0, time.UTC)

	partnerId := FACILITATOR_ID

	merchant := getLincolnBank(t, mc)
	bankAcct := createTemporaryBankAccount(t, mc, merchant.AccountID)
	merchantPmId := bankAcct.PaymentMethods[0].PaymentMethodID

	customer := CreateTemporaryTestAccount(t, mc, createTestIndividualAccount())
	customerCard := createTemporaryCard(t, mc, customer.AccountID)
	customerPmId := customerCard.PaymentMethods[0].PaymentMethodID

	schedule, err := mc.CreateSchedule(ctx, partnerId, schedules.UpsertSchedule{
		Description: "a simple schedule",

		// Setup a recurring transfer to handle repayment of say a loan with 6 periods
		RecurTransfer: &schedules.RecurTransfer{
			RecurrenceRule: "FREQ=DAILY;DTSTART=20460101T150405Z;COUNT=6",
			Transfer: schedules.Transfer{
				Description: "example of a description for all of the occurrences",
				Amount: schedules.Amount{
					Value:    100,
					Currency: "USD",
				},
				PartnerID: FACILITATOR_ID,
				Source: schedules.PaymentMethod{
					PaymentMethodID: customerPmId,
				},
				Destination: schedules.PaymentMethod{
					PaymentMethodID: merchantPmId,
				},
			},
		},

		// Lets add a one time occurrence to handle the setup fee's
		Occurrences: []schedules.UpsertTransferOccurrence{
			{
				RunOn: now,
				Transfer: schedules.Transfer{
					Description: "setup fee example",
					Amount: schedules.Amount{
						Value:    200,
						Currency: "USD",
					},
					PartnerID: FACILITATOR_ID,
					Source: schedules.PaymentMethod{
						PaymentMethodID: customerCard.PaymentMethods[0].PaymentMethodID,
					},
					Destination: schedules.PaymentMethod{
						PaymentMethodID: bankAcct.PaymentMethods[0].PaymentMethodID,
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
		upsert := schedule.ToUpsertSchedule()
		require.Len(t, upsert.Occurrences, 7)

		// price increase for all recurring transfers
		upsert.RecurTransfer.Transfer.Amount.Value = 200

		// Lets shift an occurrence by 1 day
		upsert.Occurrences[0].RunOn = upsert.Occurrences[0].RunOn.AddDate(1, 0, 20)

		// Lets change the amount of the second
		upsert.Occurrences[1].Transfer.Amount.Value = 300

		// Add a new one time occurrence to charge a fee in an hour from now
		upsert.Occurrences = append(upsert.Occurrences, schedules.UpsertTransferOccurrence{
			RunOn: now.Add(time.Hour),
			Transfer: schedules.Transfer{
				Description: "late fee fine",
				Amount: schedules.Amount{
					Value:    1,
					Currency: "USD",
				},
				PartnerID: partnerId,
				Source: schedules.PaymentMethod{
					PaymentMethodID: customerCard.PaymentMethods[0].PaymentMethodID,
				},
				Destination: schedules.PaymentMethod{
					PaymentMethodID: bankAcct.PaymentMethods[0].PaymentMethodID,
				},
			},
		})

		// Lets cancel the third occurrence
		upsert.Occurrences[2].Cancelled = moov.PtrOf(true)

		// Do the update
		updated, err := mc.UpdateSchedule(ctx, partnerId, schedule.ScheduleID, upsert)
		require.NoError(t, err)

		PrettyDebug(t, updated)

		// occurrences are always ordered so our new one for the fee will be moved up.
		require.Equal(t, int64(1), updated.Occurrences[0].Transfer.Amount.Value)

		// check that the new date is set.
		require.Equal(t, upsert.Occurrences[0].RunOn, updated.Occurrences[1].RunOn)

		// occurrences are always ordered so our new one for the fee will be moved up.
		require.Equal(t, int64(300), updated.Occurrences[2].Transfer.Amount.Value)

		// Lets cancel the third occurrence
		require.NotNil(t, updated.Occurrences[3].CancelledOn)

		schedule = updated
	})

	t.Run("Cancel", func(t *testing.T) {
		err := mc.CancelSchedule(ctx, partnerId, schedule.ScheduleID)
		require.NoError(t, err)
	})
}
