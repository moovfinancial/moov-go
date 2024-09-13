package sweeps

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

func TestSweepConfigsEndpoints(t *testing.T) {
	mc, err := moov.NewClient()
	require.NoError(t, err)

	var (
		ctx = context.Background()

		accountID  = "ebbf46c6-122a-4367-bc45-7dd555e1d3b9"
		walletID   = "4dbac313-d505-4d51-a0fe-c11787916fcf"
		minBalance = "1000.00"
	)

	sweepConfigs, err := mc.ListSweepConfigs(ctx, accountID)
	require.NoError(t, err)

	var sweepConfig *moov.SweepConfig
	// If no sweep configs found, create one
	if len(sweepConfigs) == 0 {
		// Create a sweep config
		sweepConfig, err := mc.CreateSweepConfig(ctx, moov.CreateSweepConfig{
			AccountID:           accountID,
			WalletID:            walletID,
			Status:              moov.SweepConfigStatus_Enabled,
			PushPaymentMethodID: "b46193d2-6b9b-4a73-afdc-3871779f51e3",
			PullPaymentMethodID: "467dfcdc-463b-4282-83af-db6f47562bf9",
			MinimumBalance:      &minBalance,
		})
		require.NoError(t, err)
		t.Logf("Created sweep config: %+v", sweepConfig)
	} else {
		sweepConfig = &sweepConfigs[0]
	}
	t.Logf("Got sweep config: %+v", sweepConfig)

	// Update the sweep config
	minBalance = "2000.00"
	statementDescriptor := "my-sweeps"

	sweepConfig, err = mc.UpdateSweepConfig(ctx, moov.UpdateSweepConfig{
		AccountID:           accountID,
		SweepConfigID:       sweepConfig.SweepConfigID,
		MinimumBalance:      &minBalance,
		StatementDescriptor: &statementDescriptor,
	})
	require.NoError(t, err)
	t.Logf("updated sweep config: %+v", sweepConfig)

	// Get the sweep config by ID
	sweepConfig, err = mc.GetSweepConfig(ctx, accountID, sweepConfig.SweepConfigID)
	require.NoError(t, err)
	t.Logf("Got sweep config by ID: %+v", sweepConfig)
}

func TestSweepEndpoints(t *testing.T) {
	mc, err := moov.NewClient()
	require.NoError(t, err)

	var (
		accountID = "ebbf46c6-122a-4367-bc45-7dd555e1d3b9"
		walletID  = "4dbac313-d505-4d51-a0fe-c11787916fcf"
		ctx       = context.Background()
	)

	sweeps, err := mc.ListSweeps(ctx, accountID, walletID)
	require.NoError(t, err)

	if len(sweeps) == 0 {
		t.Logf("no sweeps associated with walletID of %v", walletID)
		return
	}
	t.Logf("Listing sweeps returned %d sweeps", len(sweeps))

	sweep := &sweeps[0]
	sweep, err = mc.GetSweep(ctx, accountID, walletID, sweep.SweepID)
	require.NoError(t, err)

	t.Logf("Got first sweep in list: %+v", sweep)
}
