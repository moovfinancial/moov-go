package wallets

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/internal/testtools"
	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2604"
)

func ExampleClient_wallet() {
	mc, err := moov.NewClient()
	if err != nil {
		fmt.Printf("new Moov client: %v\n", err)
		return
	}

	var (
		ctx       = context.Background()
		accountID = "00000000-00000000-00000000-00000000"
	)

	create := moov.CreateWallet{
		Name:        "my general wallet",
		Description: "testing",
		Metadata:    map[string]string{"foo": "bar"},
	}
	createdWallet, err := mc.CreateWallet(ctx, accountID, create)
	if err != nil {
		fmt.Printf("creating wallet: %v\n", err)
		return
	}
	fmt.Printf("Created wallet: %+v\n", createdWallet)

	fetchedWallet, err := mc.GetWallet(ctx, accountID, createdWallet.WalletID)
	if err != nil {
		fmt.Printf("getting wallet: %v\n", err)
		return
	}
	fmt.Printf("getting wallet: %+v\n", fetchedWallet)

	update := moov.UpdateWallet{
		Name:        moov.PtrOf("updated name"),
		Description: moov.PtrOf("inactive wallet"),
		Status:      moov.PtrOf(moov.WalletStatus_Closed),
		Metadata:    map[string]string{"foo": "baz"},
	}
	updatedWallet, err := mc.UpdateWallet(ctx, accountID, createdWallet.WalletID, update)
	if err != nil {
		fmt.Printf("updating wallet: %v\n", err)
		return
	}
	fmt.Printf("Updated wallet: %+v\n", updatedWallet)

	listedWallets, err := mc.ListWallets(ctx, accountID, moov.WithWalletStatus(moov.WalletStatus_Active))
	if err != nil {
		fmt.Printf("listing active wallets: %v\n", err)
		return
	}
	fmt.Printf("listed active wallets: %+v\n", listedWallets)
}

func TestWalletEndpoints(t *testing.T) {
	mc, err := moov.NewClient()
	require.NoError(t, err)
	walletClientV2604 := mv2604.NewWalletClient(mc)

	var (
		ctx       = context.Background()
		accountID = testtools.MERCHANT_ID
	)

	wallets, err := walletClientV2604.ListWallets(ctx, accountID, moov.WithWalletStatus(moov.WalletStatus_Active))
	require.NoError(t, err)

	var walletID string
	for _, w := range wallets {
		if w.WalletType == moov.WalletType_General {
			walletID = w.WalletID
			t.Logf("Reusing existing general wallet: %s", walletID)
			break
		}
	}
	if walletID == "" {
		created, err := walletClientV2604.CreateWallet(ctx, accountID, moov.CreateWallet{
			Name:        "general wallet for v2604 test",
			Description: "initial description",
			Metadata:    map[string]string{"foo": "bar"},
		})
		require.NoError(t, err)
		walletID = created.WalletID
		t.Logf("Created wallet: %+v", created)
	}

	t.Run("v2604.UpdateWallet unsets the description and metadata", func(t *testing.T) {
		updatedWallet, err := walletClientV2604.UpdateWallet(ctx, accountID, walletID, mv2604.UpdateWallet{
			Description: moov.SetNull[string](),
			Metadata:    moov.SetNull[map[string]string](),
		})
		require.NoError(t, err)
		require.Empty(t, updatedWallet.Description)
		require.Empty(t, updatedWallet.Metadata)
		t.Logf("unset description and metadata in wallet: %+v", updatedWallet)

		fetchedWallet, err := walletClientV2604.GetWallet(ctx, accountID, walletID)
		require.NoError(t, err)
		require.Empty(t, fetchedWallet.Description)
		require.Empty(t, fetchedWallet.Metadata)
		t.Logf("got wallet with unset description and metadata: %+v", fetchedWallet)
	})
}
