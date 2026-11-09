package mv2604

import (
	"context"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

type InvoiceClient struct {
	*moov.Client
}

func NewInvoiceClient(client *moov.Client) InvoiceClient {
	return InvoiceClient{Client: client}
}

func (i InvoiceClient) CreateInvoiceNotification(ctx context.Context, accountID, invoiceID string, notification moov.CreateInvoiceNotification) error {
	return moov.CreateInvoiceNotificationGeneric(ctx, i.Client, moov.Version2026_04, accountID, invoiceID, notification)
}
