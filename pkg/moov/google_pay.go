package moov

import (
	"context"
	"net/http"
)

// LinkedGooglePayPaymentMethod is a Google Pay payment method linked to a Moov account.
type LinkedGooglePayPaymentMethod struct {
	PaymentMethodID   string                 `json:"paymentMethodID"`
	PaymentMethodType string                 `json:"paymentMethodType"`
	GooglePay         GooglePayPaymentMethod `json:"googlePay"`
}

// GooglePayTokenLink is a Google Pay token link request body accepted by LinkGooglePayTokenGeneric.
type GooglePayTokenLink interface {
	GooglePayTokenLink()
}

// LinkGooglePayTokenGeneric links a Google Pay token for the accountID.
func LinkGooglePayTokenGeneric[T GooglePayTokenLink](ctx context.Context, client *Client, version Version, accountID string, link T) ([]LinkedGooglePayPaymentMethod, error) {
	if client == nil {
		return nil, ErrClientIsNil
	}
	if accountID == "" {
		return nil, ErrAccountIDRequired
	}

	resp, err := client.CallHttp(
		ctx,
		Endpoint(http.MethodPost, pathGooglePayTokens, accountID),
		MoovVersion(version),
		AcceptJson(),
		JsonBody(link),
	)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[LinkedGooglePayPaymentMethod](resp)
}
