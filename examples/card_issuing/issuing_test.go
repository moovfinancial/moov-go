package cardissuing

import (
	"context"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"

	"github.com/stretchr/testify/require"
)

// Card issuing is currently in a closed beta. Contact Moov for more information.
func TestCardIssuing(t *testing.T) {
	// Step 1: create Moov client and set some variables

	// The following code shows how you can configure the moov client with
	// your credentials, if you don't want to use environment variables.
	// However, it is recommended to load the credentials from the
	// configuration file.

	mc, err := moov.NewClient() // reads credentials from Environmental variables by default
	require.NoError(t, err)

	// Create a new context or use an existing one
	ctx := context.Background()

	// Ping the server to check credentials
	err = mc.Ping(ctx)
	require.NoError(t, err)

	// Step 2: create account for the user

	// For now just using a known existing account that is already allowed in the closed beta for card issuing
	accountID := "ebbf46c6-122a-4367-bc45-7dd555e1d3b9"

	// Step 3: enable the card-issuing capability for the account

	// Can skip this step if the card-issuing capability has already previously been enabled for the account

	_, err = mc.RequestCapabilities(ctx, accountID, []moov.CapabilityName{
		"card-issuing",
	})
	require.NoError(t, err)

	// Step 4: get the source wallet for funding the issued card

	// For now just using a known existing wallet for the account
	walletID := "4dbac313-d505-4d51-a0fe-c11787916fcf"

	// Step 5: create an issued card

	memo := "example"
	create := moov.CreateIssuedCard{
		FundingWalletID: walletID,
		AuthorizedUser: moov.CreateAuthorizedUser{
			FirstName: "John",
			LastName:  "Doe",
		},
		FormFactor: moov.IssuedCardFormFactor_Virtual,
		Memo:       &memo,
	}
	created, err := mc.CreateIssuedCard(ctx, accountID, create)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, closeIssuedCard(ctx, mc, accountID, created.IssuedCardID))
	})
}

func closeIssuedCard(ctx context.Context, mc *moov.Client, accountID, cardID string) error {
	closed := moov.UpdateIssuedCardState_Closed
	update := moov.UpdateIssuedCard{
		State: &closed,
	}
	return mc.UpdateIssuedCard(ctx, accountID, cardID, update)
}
