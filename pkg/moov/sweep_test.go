package moov_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

// TODO(vince): update IDs here once I get invited to the prod sandbox account
func Test_SweepConfig_CRUD_Endpoints(t *testing.T) {
	var (
		mc = NewTestClient(t)

		ctx        = context.Background()
		accountID  = "ff54f56f-6b7e-439f-9280-bc5edd3f93fc"
		walletID   = "24ff6ab1-01c0-4c50-bd61-d9fe096a0b72"
		minBalance = "100.0"
	)

	sweepConfigs, err := mc.ListSweepConfigs(ctx, accountID)
	require.NoError(t, err)

	var sweepConfig *moov.SweepConfig
	// If no sweep configs found, create one
	if len(sweepConfigs) == 0 {
		sweepConfig, err = mc.CreateSweepConfig(ctx, moov.CreateSweepConfig{
			AccountID: accountID,
			WalletID:  walletID,
			Status:    moov.SweepConfigStatus_Enabled,
			// todo: update these
			PushPaymentMethodID: "df0ac7e9-639b-41ca-ab0e-7180f356ed86",
			PullPaymentMethodID: "155647dd-38ae-4ac6-a9a4-8c02bec9141e",
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
