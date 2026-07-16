package moov_test

import (
	"encoding/json"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

func TestLinkApplePayTokenResponse(t *testing.T) {
	var response moov.LinkApplePayTokenResponse
	err := json.Unmarshal([]byte(`{
		"paymentMethods": [{
			"paymentMethodID": "pm_123",
			"paymentMethodType": "apple-pay",
			"applePay": {"brand": "Visa"}
		}],
		"paymentMethodErrors": [{
			"paymentMethodType": "push-to-apple-pay",
			"error": "Apple Cash is not supported",
			"code": "apple-cash-not-supported"
		}]
	}`), &response)

	require.NoError(t, err)
	require.Len(t, response.PaymentMethods, 1)
	require.Equal(t, "pm_123", response.PaymentMethods[0].PaymentMethodID)
	require.Len(t, response.PaymentMethodErrors, 1)
	require.Equal(t, "apple-cash-not-supported", response.PaymentMethodErrors[0].Code)
}
