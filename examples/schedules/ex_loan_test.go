package schedules

import (
	"context"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

func Test_Loan_Example(t *testing.T) {
	ctx := context.Background()

	// Lets setup an example environment where the client, customer, and merchant already exist.
	env := Setup(t, ctx)

	// Start the occurring payments 1 month from today
	paymentsStart := env.Now.AddDate(0, 1, 0)

	_, err := env.Client.CreateSchedule(ctx, env.PartnerID, moov.CreateSchedule{
		Description: "Car Loan",

		// One time occurrence to handle say the tax, title, and registration of a new car.
		Occurrences: []moov.CreateOccurrence{
			{
				RunOn: env.Now,
				RunTransfer: moov.RunTransfer{
					Description: "Tax, Title, Registration",
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

		// Add in a recurring schedule for the remaining 36 payments
		Recur: &moov.Recur{
			Start:          &paymentsStart,
			RecurrenceRule: "FREQ=MONTHLY;BYDAY=+1MO;COUNT=36",
			RunTransfer: moov.RunTransfer{
				Description: "Monthly payment",
				Amount: moov.ScheduleAmount{
					Value:    1,
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
	})

	require.NoError(t, err)
}
