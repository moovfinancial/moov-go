package mv2607

import (
	"context"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

type DepositViewClient struct {
	*moov.Client
}

func NewDepositViewClient(client *moov.Client) DepositViewClient {
	return DepositViewClient{Client: client}
}

// CreateDepositView submits a deposit view document for the given account. The
// sourceSystem identifies the core banking system the document originated from and
// is sent as the x-source-system header. The document is sent as the raw request
// body and the API responds with the ingest result on success.
func (d DepositViewClient) CreateDepositView(ctx context.Context, accountID string, sourceSystem SourceSystem, document []byte) (*IngestResponse, error) {
	return moov.CreateDepositView[IngestResponse](ctx, d.Client, moov.Version2026_07, accountID, string(sourceSystem), document)
}
