package moov

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PaymentMethod struct {
	PaymentMethodID   string `json:"paymentMethodID,omitempty"`
	PaymentMethodType string `json:"paymentMethodType,omitempty"`
	Wallet            Wallet `json:"wallet,omitempty"`
}

// ListPaymentMethods lists all payment methods that are associated with a Moov account
// https://docs.moov.io/api/index.html#tag/Payment-methods/operation/getPaymentMethods
func (c Client) ListPaymentMethods(accountID string, sourceID string) ([]PaymentMethod, error) {
	var resPaymentMethods []PaymentMethod
	url := fmt.Sprintf("%s/%s", baseURL, fmt.Sprintf(pathPaymentMethods, accountID))
	if sourceID != "" {
		url = fmt.Sprintf("%s?sourceID=%s", url, sourceID)
	}

	body, statusCode, err := GetHTTPResponse(c, http.MethodGet, url, nil, nil)
	if err != nil {
		return resPaymentMethods, err
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &resPaymentMethods)
		if err != nil {
			return resPaymentMethods, err
		}
		return resPaymentMethods, nil
	case http.StatusNotFound:
		return resPaymentMethods, ErrNoAccount
	case http.StatusTooManyRequests:
		return resPaymentMethods, ErrRateLimit
	}
	return resPaymentMethods, ErrDefault(statusCode)
}

// GetPaymentMethod retrieves a payment method for the given payment method id
// https://docs.moov.io/api/index.html#tag/Payment-methods/operation/getPaymentMethod
func (c Client) GetPaymentMethod(accountID string, paymentMethodID string) (PaymentMethod, error) {
	resPaymentMethod := PaymentMethod{}
	url := fmt.Sprintf("%s/%s/%s", baseURL, fmt.Sprintf(pathPaymentMethods, accountID), paymentMethodID)

	body, statusCode, err := GetHTTPResponse(c, http.MethodGet, url, nil, nil)
	if err != nil {
		return resPaymentMethod, err
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &resPaymentMethod)
		if err != nil {
			return resPaymentMethod, err
		}
		return resPaymentMethod, nil
	case http.StatusNotFound:
		return resPaymentMethod, ErrNoAccount
	case http.StatusTooManyRequests:
		return resPaymentMethod, ErrRateLimit
	}
	return resPaymentMethod, ErrDefault(statusCode)
}
