package mv2604

import (
	"context"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

type WalletClient struct {
	*moov.Client
}

func NewWalletClient(client *moov.Client) WalletClient {
	return WalletClient{Client: client}
}

func (w WalletClient) UpdateWallet(ctx context.Context, accountID, walletID string, update UpdateWallet) (*moov.Wallet, error) {
	return moov.UpdateWalletGeneric(ctx, w.Client, moov.Version2026_04, accountID, walletID, update)
}

type UpdateWallet struct {
	Name   *string            `json:"name,omitempty"`
	Status *moov.WalletStatus `json:"status,omitempty"`

	// A free-form description of the wallet.
	Description *moov.Nullable[string] `json:"description,omitempty"`
	// Free-form key-value pair list. Useful for storing information that is not captured elsewhere.
	Metadata *moov.Nullable[map[string]string] `json:"metadata,omitempty"`
}
