package mv2604

import (
	"context"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

type GooglePayClient struct {
	*moov.Client
}

func NewGooglePayClient(client *moov.Client) GooglePayClient {
	return GooglePayClient{Client: client}
}

func (g GooglePayClient) LinkGooglePayToken(ctx context.Context, accountID string, link moov.LinkGooglePay) ([]moov.PaymentMethod, error) {
	return moov.LinkGooglePayTokenGeneric(ctx, g.Client, moov.Version2026_04, accountID, link)
}
