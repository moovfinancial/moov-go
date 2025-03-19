package schedules

import (
	"context"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

func Test_Run_Tomorrow_Example(t *testing.T) {
	ctx := context.Background()

	// Lets setup an example environment where the client, customer, and merchant already exist.
	env := Setup(t, ctx)

	// Run the payment tomorrow
	runDate := env.Now.AddDate(0, 0, 1)

	_, err := env.Client.CreateSchedule(ctx, env.PartnerID, moov.CreateSchedule{
		Description: "Delayed Payment",

		// One time occurrence to handle say the tax, title, and registration of a new car.
		Occurrences: []moov.CreateOccurrence{
			{
				RunOn: runDate,
				RunTransfer: moov.RunTransfer{
					Description: "Payment for ...",
					Amount: moov.ScheduleAmount{
						Value:    2,
						Currency: "USD",
					},
					Source: moov.SchedulePaymentMethod{
						PaymentMethodID: env.CustomerPmId,
					},
					Destination: moov.SchedulePaymentMethod{
						PaymentMethodID: env.MerchantPmId,
					},
				},
			},
		},
	})

	require.NoError(t, err)
}
