package moov

import (
	"context"
	"fmt"
	"net/http"
)

// LinkedGooglePayPaymentMethod is a Google Pay payment method linked to a Moov account.
type LinkedGooglePayPaymentMethod struct {
	PaymentMethodID   string                 `json:"paymentMethodID"`
	PaymentMethodType string                 `json:"paymentMethodType"`
	GooglePay         GooglePayPaymentMethod `json:"googlePay"`
}

// LinkGooglePayTokenGeneric links a Google Pay token for the accountID.
func LinkGooglePayTokenGeneric[T any](ctx context.Context, client *Client, version Version, accountID string, link T) ([]LinkedGooglePayPaymentMethod, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	if accountID == "" {
		return nil, fmt.Errorf("accountID is required")
	}

	resp, err := client.CallHttp(
		ctx,
		Endpoint(http.MethodPost, pathGooglePayTokens, accountID),
		MoovVersion(version),
		AcceptJson(),
		JsonBody(link),
	)
	if err != nil {
		return nil, fmt.Errorf("linking google pay token: %w", err)
	}

	return CompletedListOrError[LinkedGooglePayPaymentMethod](resp)
}
