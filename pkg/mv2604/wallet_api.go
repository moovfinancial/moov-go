package mv2604

import (
	"context"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

type WalletClient struct {
	client *moov.Client
}

func NewWalletClient(client *moov.Client) WalletClient {
	return WalletClient{client: client}
}

func (w WalletClient) UpdateWallet(ctx context.Context, accountID, walletID string, update UpdateWallet) (*moov.Wallet, error) {
	return moov.UpdateWalletGeneric(ctx, w.client, moov.Version2026_04, accountID, walletID, update)
}

type UpdateWallet struct {
	moov.UpdateWallet

	// A free-form description of the wallet.
	Description *moov.Nullable[string] `json:"description,omitempty"`
	// Free-form key-value pair list. Useful for storing information that is not captured elsewhere.
	Metadata *moov.Nullable[map[string]string] `json:"metadata,omitempty"`
}
