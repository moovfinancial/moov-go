package wallets

import (
	"context"
	"fmt"

	"github.com/moovfinancial/moov-go/pkg/moov"
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
