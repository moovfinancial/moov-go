package moov_test

import (
	"fmt"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

func Test_ShareConnection(t *testing.T) {
	mc := NewTestClient(t)

	customer := CreateTemporaryTestAccount(t, mc, createTestIndividualAccount())
	merchant := CreateTemporaryTestAccount(t, mc, createTestIndividualAccount())

	fmt.Println(customer.AccountID, merchant.AccountID)

	shared, err := mc.ShareConnection(BgCtx(), customer.AccountID, moov.ShareConnectionRequest{
		PrincipalAccountID: merchant.AccountID,
		AllowScopes: []string{
			"/accounts.read",
		},
	})
	fmt.Println(shared)
	require.NoError(t, err)
}
