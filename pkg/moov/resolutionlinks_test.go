package moov_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ResolutionLinks(t *testing.T) {
	mc := NewTestClient(t)

	account := getLincolnBank(t, mc)
	resolutionLinkCode := "abc123"

	t.Run("list Resolution Links", func(t *testing.T) {
		filtered, err := mc.ListResolutionLinks(BgCtx(), account.AccountID)
		NoResponseError(t, err)
		require.NotEmpty(t, filtered)
	})

	t.Run("get resolution link", func(t *testing.T) {
		cap, err := mc.GetResolutionLink(BgCtx(), account.AccountID, resolutionLinkCode)
		NoResponseError(t, err)
		require.NotNil(t, cap)
	})

	t.Run("create resolution link", func(t *testing.T) {
		cap, err := mc.CreateResolutionLink(BgCtx(), account.AccountID)
		NoResponseError(t, err)
		require.NotNil(t, cap)
	})

	t.Run("delete resolution link", func(t *testing.T) {
		err := mc.DeleteResolutionLink(BgCtx(), account.AccountID, resolutionLinkCode)
		NoResponseError(t, err)
	})
}
