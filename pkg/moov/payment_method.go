package moov

import (
	"context"
	"net/http"
)

type PaymentMethod struct {
	PaymentMethodID   string      `json:"paymentMethodID,omitempty"`
	PaymentMethodType string      `json:"paymentMethodType,omitempty"`
	Wallet            Wallet      `json:"wallet,omitempty"`
	BankAccount       BankAccount `json:"bankAccount,omitempty"`
	Card              Card        `json:"card,omitempty"`
	ApplePay          ApplePay    `json:"applePay,omitempty"`
}

type PaymentMethodListFilter callArg

func WithPaymentMethodSourceID(id string) callArg {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["sourceID"] = id
		return nil
	})
}

// WithPaymentMethodType filters the payment methods by the payment method type. example: moov-wallet, card-payment, ach-debit-collect
func WithPaymentMethodType(t string) callArg {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["paymentMethodType"] = t
		return nil
	})
}

// ListPaymentMethods lists all payment methods that are associated with a Moov account
// https://docs.moov.io/api/index.html#tag/Payment-methods/operation/getPaymentMethods
func (c Client) ListPaymentMethods(ctx context.Context, accountID string, opts ...PaymentMethodListFilter) ([]PaymentMethod, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathPaymentMethods, accountID),
		prependArgs(opts, AcceptJson())...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[PaymentMethod](resp)
}

// GetPaymentMethod retrieves a payment method for the given payment method id
// https://docs.moov.io/api/index.html#tag/Payment-methods/operation/getPaymentMethod
func (c Client) GetPaymentMethod(ctx context.Context, accountID string, paymentMethodID string) (*PaymentMethod, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathPaymentMethods, accountID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[PaymentMethod](resp)
}
