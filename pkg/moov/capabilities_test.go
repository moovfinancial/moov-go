package moov_test

import (
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

func Test_Capabilities(t *testing.T) {
	mc := NewTestClient(t)

	account := CreateTemporaryTestAccount(t, mc, createTestBusinessAccount())

	t.Run("requested", func(t *testing.T) {
		requested, err := mc.RequestCapabilities(BgCtx(), account.AccountID, []moov.CapabilityName{
			moov.CapabilityName_1099,
			moov.CapabilityName_CollectFunds,
			moov.CapabilityName_SendFunds,
			moov.CapabilityName_Transfers,
			moov.CapabilityName_Wallet,
		})

		NoResponseError(t, err)
		require.NotEmpty(t, requested)
	})

	t.Run("list", func(t *testing.T) {
		requested, err := mc.ListCapabilities(BgCtx(), account.AccountID)
		NoResponseError(t, err)
		require.NotEmpty(t, requested)
	})

	t.Run("get", func(t *testing.T) {
		cap, err := mc.GetCapability(BgCtx(), account.AccountID, moov.CapabilityName_Transfers)
		NoResponseError(t, err)
		require.NotNil(t, cap)
	})

	t.Run("disable", func(t *testing.T) {
		err := mc.DisableCapability(BgCtx(), account.AccountID, moov.CapabilityName_Transfers)
		NoResponseError(t, err)
	})
}
