package schedules

import (
	"context"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

func Test_Subscription(t *testing.T) {
	ctx := context.Background()

	// Lets setup an example environment where the client, customer, and merchant already exist.
	env := Setup(t, ctx)

	schedule, err := env.Client.CreateSchedule(ctx, env.PartnerID, moov.CreateSchedule{
		Description: "Streaming Services",

		// Add in a recurring schedule that goes on indefinitely that bills every month at this time.
		RecurTransfer: &moov.RecurTransfer{
			Start:          &env.Now,
			RecurrenceRule: "FREQ=MONTHLY",
			Transfer: moov.ScheduleTransfer{
				Description: "Monthly payment",
				Amount: moov.ScheduleAmount{
					Value:    1,
					Currency: "USD",
				},
				PartnerID: env.PartnerID,
				Source: moov.SchedulePaymentMethod{
					PaymentMethodID: env.CustomerPmId,
				},
				Destination: moov.SchedulePaymentMethod{
					PaymentMethodID: env.MerchantPmId,
				},
			},
		},
	})

	require.NoError(t, err)

	occ, err := env.Client.GetScheduleOccurrence(ctx, env.PartnerID, schedule.ScheduleID, moov.OccurrenceLatestToTime(env.Now))
	require.NoError(t, err)

	if occ.TransferStatus == nil {
		// payment hasn't ran yet
	} else if *occ.TransferStatus == string(moov.TransferStatus_Failed) {
		// last transfer failed, don't allow access to premium features
	} else {
		// payment is in the middle of processing and hasn't fully completed and been deposited in merchants account
	}
}
