package moov_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

func Test_SweepConfig_CRUD_Endpoints(t *testing.T) {
	var (
		mc = NewTestClient(t)

		ctx = context.Background()
		// IDs for Lincoln National Corporation
		accountID  = "ebbf46c6-122a-4367-bc45-7dd555e1d3b9"
		walletID   = "4dbac313-d505-4d51-a0fe-c11787916fcf"
		minBalance = "1000.0"
	)

	sweepConfigs, err := mc.ListSweepConfigs(ctx, accountID)
	require.NoError(t, err)

	var sweepConfig *moov.SweepConfig
	// If no sweep configs found, create one
	if len(sweepConfigs) == 0 {
		sweepConfig, err = mc.CreateSweepConfig(ctx, moov.CreateSweepConfig{
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
	statementDesc := "my-sweepz"
	sweepConfig, err = mc.UpdateSweepConfig(ctx, moov.UpdateSweepConfig{
		AccountID:           accountID,
		SweepConfigID:       sweepConfig.SweepConfigID,
		StatementDescriptor: &statementDesc,
	})
	require.NoError(t, err)

	// Get by ID
	sweepConfig, err = mc.GetSweepConfig(ctx, accountID, sweepConfig.SweepConfigID)
	require.NoError(t, err)
}
