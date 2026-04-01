package moov_test

import (
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

func Test_ResolutionLinks(t *testing.T) {
	mc := NewTestClient(t)

	var resolutionLinkCode string

	t.Run("create resolution link", func(t *testing.T) {
		createReq := moov.CreateResolutionLinkRequest{
			Recipient: moov.Recipient{
				Email: "noreply@moov.io",
			},
		}
		created, err := mc.CreateResolutionLink(BgCtx(), MERCHANT_ID, createReq)
		NoResponseError(t, err)
		require.NotNil(t, created)
		require.NotEmpty(t, created.ResolutionLinkCode)
		resolutionLinkCode = created.ResolutionLinkCode
	})

	t.Run("list resolution links", func(t *testing.T) {
		links, err := mc.ListResolutionLinks(BgCtx(), MERCHANT_ID)
		NoResponseError(t, err)
		require.NotEmpty(t, links)
	})

	t.Run("get resolution link", func(t *testing.T) {
		link, err := mc.GetResolutionLink(BgCtx(), MERCHANT_ID, resolutionLinkCode)
		NoResponseError(t, err)
		require.NotNil(t, link)
		require.Equal(t, resolutionLinkCode, link.ResolutionLinkCode)
	})

	t.Run("delete resolution link", func(t *testing.T) {
		err := mc.DeleteResolutionLink(BgCtx(), MERCHANT_ID, resolutionLinkCode)
		NoResponseError(t, err)
	})
}
