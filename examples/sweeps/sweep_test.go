package sweeps

import (
	"context"
	"fmt"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

func Example_SweepConfigs_Endpoints() {
	mc, err := moov.NewClient()
	if err != nil {
		fmt.Printf("new client: %v", err)
		return
	}

	var (
		accountID = "2f29c185-91d8-40c7-8963-3915aa673c21"
		walletID  = "ce906d91-61ef-47eb-8e15-38f80b863200"

		ctx        = context.Background()
		minBalance = "100.00"
	)

	// Create a sweep config
	sweepConfig, err := mc.CreateSweepConfig(ctx, moov.CreateSweepConfig{
		AccountID:           accountID,
		WalletID:            walletID,
		Status:              moov.SweepConfigStatus_Enabled,
		PushPaymentMethodID: "8e17b62d-cd7f-4d88-8576-608ba7575860",
		PullPaymentMethodID: "a54ddae8-57ab-4941-a723-097292fe30ac",
		MinimumBalance:      &minBalance,
	})
	if err != nil {
		fmt.Printf("creating sweep config: %v", err)
		return
	}
	fmt.Printf("Created sweep config: %+v", sweepConfig)

	// Get a sweep config by ID
	sweepConfig, err = mc.GetSweepConfig(ctx, accountID, sweepConfig.SweepConfigID)
	if err != nil {
		fmt.Printf("getting sweep config: %v", err)
		return
	}
	fmt.Printf("Got sweep config: %+v", sweepConfig)

	// Update the sweep config
	minBalance = "0"
	pullPaymentMethodID := "219b2089-4dec-45e7-a0cb-68063565868c"
	sweepConfig, err = mc.UpdateSweepConfig(ctx, moov.UpdateSweepConfig{
		AccountID:           accountID,
		SweepConfigID:       sweepConfig.SweepConfigID,
		MinimumBalance:      &minBalance,
		PullPaymentMethodID: &pullPaymentMethodID,
	})
	if err != nil {
		fmt.Printf("udpating sweep config: %v", err)
		return
	}
	fmt.Printf("updated sweep config: %+v", sweepConfig)
}

func Example_Sweeps_Endpoints() {
	mc, err := moov.NewClient()
	if err != nil {
		fmt.Printf("new client: %v", err)
		return
	}

	var (
		accountID = "2f29c185-91d8-40c7-8963-3915aa673c21"
		walletID  = "ce906d91-61ef-47eb-8e15-38f80b863200"
		ctx       = context.Background()
	)

	sweeps, err := mc.ListSweeps(ctx, accountID, walletID)
	if err != nil {
		fmt.Printf("listing sweeps: %v", err)
		return
	}

	if len(sweeps) == 0 {
		fmt.Printf("no sweeps associated with walletID of %v", walletID)
		return
	}
	fmt.Printf("Listing sweeps returned %d sweeps", len(sweeps))

	sweep := &sweeps[0]
	sweep, err = mc.GetSweep(ctx, accountID, walletID, sweep.SweepID)
	if err != nil {
		fmt.Printf("getting sweep: %v", err)
		return
	}
	fmt.Printf("Got first sweep in list: %+v", sweep)
}
