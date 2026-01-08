package moov_test

import (
	"testing"
	"time"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2507"
	"github.com/stretchr/testify/require"
)

func Test_ShareConnection(t *testing.T) {
	mc := NewTestClient(t)

	customer := CreateTemporaryTestAccount(t, mc, createTestIndividualAccount())
	merchant := CreateTemporaryTestAccount(t, mc, createTestIndividualAccount())

	t.Run("share connection between customer and merchant from the partner", func(t *testing.T) {
		shared, err := mc.ShareConnection(BgCtx(), customer.AccountID, moov.ShareConnectionRequest{
			PrincipalAccountID: merchant.AccountID,
			AllowScopes:        []string{"profile.read"},
		})
		require.NoError(t, err)
		require.NotNil(t, shared)

		require.Equal(t, merchant.AccountID, shared.PrincipalAccountID)
		require.Equal(t, customer.AccountID, shared.SubjectAccountID)
		require.Equal(t, []string{"profile.read"}, shared.Scopes)
	})

	t.Run("list accounts connected to merchant from the partner perspective", func(t *testing.T) {
		var connections []mv2507.Account
		var err error

		require.Eventually(t, func() bool {
			connections, err = mv2507.Accounts.ListConnected(BgCtx(), *mc, merchant.AccountID)
			return err == nil && len(connections) > 0
		}, time.Second*1, time.Millisecond*100)

		require.NoError(t, err)
		require.NotNil(t, connections)
		require.Len(t, connections, 1)
		require.Equal(t, customer.AccountID, connections[0].AccountID)
	})
}
