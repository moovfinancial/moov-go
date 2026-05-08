package mv2604

import (
	"context"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

type TransferClient struct {
	*moov.Client
}

func NewTransferClient(client *moov.Client) TransferClient {
	return TransferClient{Client: client}
}

func (t TransferClient) PatchTransfer(ctx context.Context, accountID, transferID string, update PatchTransfer) (*moov.Transfer, error) {
	return moov.PatchTransferGeneric(ctx, t.Client, moov.Version2026_04, accountID, transferID, update)
}

type PatchTransfer struct {
	// ForeignID an optional alias from a foreign/external system which can be used to reference this resource.
	ForeignID *moov.Nullable[string] `json:"foreignID,omitempty"`
	// Free-form key-value pair list. Useful for storing information that is not captured elsewhere.
	Metadata *moov.Nullable[map[string]string] `json:"metadata,omitempty"`
}
