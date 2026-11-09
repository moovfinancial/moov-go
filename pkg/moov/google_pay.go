package moov

import (
	"context"
	"fmt"
	"net/http"
)

// LinkGooglePay links a Google Pay token to a Moov account.
type LinkGooglePay struct {
	MerchantAccountID string                     `json:"merchantAccountID"`
	PaymentMethodData GooglePayPaymentMethodData `json:"paymentMethodData"`
}

// GooglePayPaymentMethodData is the Google Pay token payload from the Google Pay API.
type GooglePayPaymentMethodData struct {
	Type             string                    `json:"type,omitempty"`
	Description      string                    `json:"description,omitempty"`
	Info             GooglePayCardInfo         `json:"info,omitempty"`
	TokenizationData GooglePayTokenizationData `json:"tokenizationData,omitempty"`
}

// GooglePayCardInfo describes the card a Google Pay token represents.
type GooglePayCardInfo struct {
	CardNetwork      string                     `json:"cardNetwork,omitempty"`
	CardDetails      string                     `json:"cardDetails,omitempty"`
	AssuranceDetails *GooglePayAssuranceDetails `json:"assuranceDetails,omitempty"`
	BillingAddress   *GooglePayBillingAddress   `json:"billingAddress,omitempty"`
}

// GooglePayTokenizationData carries the gateway token for a Google Pay card.
type GooglePayTokenizationData struct {
	Type  string `json:"type,omitempty"`
	Token string `json:"token,omitempty"`
}

// GooglePayBillingAddress is the billing address Google Pay returns.
type GooglePayBillingAddress struct {
	Name               string `json:"name,omitempty"`
	PostalCode         string `json:"postalCode,omitempty"`
	CountryCode        string `json:"countryCode,omitempty"`
	PhoneNumber        string `json:"phoneNumber,omitempty"`
	Address1           string `json:"address1,omitempty"`
	Address2           string `json:"address2,omitempty"`
	Address3           string `json:"address3,omitempty"`
	Locality           string `json:"locality,omitempty"`
	AdministrativeArea string `json:"administrativeArea,omitempty"`
}

// GooglePayAssuranceDetails carries the Google Pay authentication indicators.
type GooglePayAssuranceDetails struct {
	CardHolderAuthenticated bool `json:"cardHolderAuthenticated,omitempty"`
	AccountVerified         bool `json:"accountVerified,omitempty"`
}

// LinkGooglePayTokenGeneric links a Google Pay token for the accountID.
func LinkGooglePayTokenGeneric(ctx context.Context, client *Client, version Version, accountID string, link LinkGooglePay) ([]PaymentMethod, error) {
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

	return CompletedListOrError[PaymentMethod](resp)
}
