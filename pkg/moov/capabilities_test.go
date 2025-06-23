package moov_test

import (
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2507"
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

	t.Run("cannot request granular capabilities", func(t *testing.T) {
		requested, err := mc.RequestCapabilities(BgCtx(), account.AccountID, []moov.CapabilityName{
			moov.CapabilityName_CollectFundsACH,
		})

		require.Error(t, err)
		require.Nil(t, requested)
	})
}

func Test_Capabilities_V2507(t *testing.T) {
	mc := NewTestClient(t)

	account := CreateTemporaryTestAccount(t, mc, createTestBusinessAccount())

	requested := mv2507.RequestedCapabilities{
		Capabilities: []moov.CapabilityName{
			moov.CapabilityName_PlatformProductionApp,
			moov.CapabilityName_PlatformWalletTransfers,
			moov.Capability_NameWalletBalance,
			moov.CapabilityName_CollectFundsACH,
			moov.CapabilityName_CollectFundsCardPayments,
			moov.CapabilityName_MoneyTransferPullFromCard,
			moov.CapabilityName_MoneyTransferPushToCard,
			moov.CapabilityName_SendFundsACH,
			moov.CapabilityName_SendFundsRTP,
			moov.CapabilityName_SendFundsPushToCard,
		},
	}

	t.Run("requested", func(t *testing.T) {
		requested, err := mv2507.Capabilities.Request(BgCtx(), *mc, account.AccountID, requested)

		NoResponseError(t, err)
		require.NotEmpty(t, requested)
	})

	t.Run("list", func(t *testing.T) {
		requested, err := mv2507.Capabilities.List(BgCtx(), *mc, account.AccountID)

		NoResponseError(t, err)
		require.NotEmpty(t, requested)
	})

	t.Run("get", func(t *testing.T) {
		cap, err := mv2507.Capabilities.Get(BgCtx(), *mc, account.AccountID, moov.CapabilityName_Transfers)

		NoResponseError(t, err)
		require.NotNil(t, cap)
	})

	t.Run("disable", func(t *testing.T) {
		err := mv2507.Capabilities.Disable(BgCtx(), *mc, account.AccountID, moov.CapabilityName_Transfers)

		NoResponseError(t, err)
	})
}
