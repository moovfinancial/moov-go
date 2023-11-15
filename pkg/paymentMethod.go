package moov

type AutoGenerated struct {
	PaymentMethodID   string `json:"paymentMethodID,omitempty"`
	PaymentMethodType string `json:"paymentMethodType,omitempty"`
	Wallet            Wallet `json:"wallet,omitempty"`
}

// ListPaymentMethods lists all payment methods that are associated with a Moov account
// https://docs.moov.io/api/index.html#tag/Payment-methods/operation/getPaymentMethods

// GetPaymentMethod retrieves a payment method for the given payment method id
// https://docs.moov.io/api/index.html#tag/Payment-methods/operation/getPaymentMethod