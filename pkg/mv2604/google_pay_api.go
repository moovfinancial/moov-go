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

func (g GooglePayClient) LinkGooglePayToken(ctx context.Context, accountID string, link LinkGooglePay) ([]moov.PaymentMethod, error) {
	return moov.LinkGooglePayTokenGeneric(ctx, g.Client, moov.Version2026_04, accountID, link)
}

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
