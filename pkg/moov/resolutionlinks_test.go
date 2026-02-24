package moov_test

import (
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

func Test_ResolutionLinks(t *testing.T) {
	mc := NewTestClient(t)

	account := getLincolnBank(t, mc)

	var resolutionLinkCode string

	t.Run("create resolution link", func(t *testing.T) {
		createReq := moov.CreateResolutionLink{
			AccountID: account.AccountID,
			Recipient: moov.Recipient{
				Email: "noreply@moov.io",
			},
			Options: moov.ResolutionLinkOptions{
				MerchantName: "Test Merchant",
				AccountName:  "Test Account",
			},
		}
		created, err := mc.CreateResolutionLink(BgCtx(), account.AccountID, createReq)
		NoResponseError(t, err)
		require.NotNil(t, created)
		require.NotEmpty(t, created.ResolutionLinkCode)
		resolutionLinkCode = created.ResolutionLinkCode
	})

	t.Run("list resolution links", func(t *testing.T) {
		links, err := mc.ListResolutionLinks(BgCtx(), account.AccountID)
		NoResponseError(t, err)
		require.NotEmpty(t, links)
	})

	t.Run("get resolution link", func(t *testing.T) {
		link, err := mc.GetResolutionLink(BgCtx(), account.AccountID, resolutionLinkCode)
		NoResponseError(t, err)
		require.NotNil(t, link)
		require.Equal(t, resolutionLinkCode, link.ResolutionLinkCode)
	})

	t.Run("delete resolution link", func(t *testing.T) {
		err := mc.DeleteResolutionLink(BgCtx(), account.AccountID, resolutionLinkCode)
		NoResponseError(t, err)
	})
}
