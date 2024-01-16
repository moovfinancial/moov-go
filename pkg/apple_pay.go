package moov

import (
	"context"
	"errors"
	"net/http"
	"time"
)

type ApplePayDomains struct {
	DisplayName string   `json:"displayName,omitempty"`
	Domains     []string `json:"domains,omitempty"`
}

type ApplePayDomainsResponse struct {
	AccountID   string    `json:"accountID,omitempty"`
	DisplayName string    `json:"displayName,omitempty"`
	Domains     []string  `json:"domains,omitempty"`
	CreatedOn   time.Time `json:"createdOn,omitempty"`
	UpdatedOn   time.Time `json:"updatedOn,omitempty"`
}

type ApplePaymentDataHeader struct {
	EphemeralPublicKey string `json:"ephemeralPublicKey"`
	PublicKeyHash      string `json:"publicKeyHash"`
	TransactionId      string `json:"transactionId"`
}

type ApplePaymentData struct {
	Version   string                 `json:"version"`
	Data      string                 `json:"data"`
	Signature string                 `json:"signature"`
	Header    ApplePaymentDataHeader `json:"header"`
}

type ApplePaymentMethod struct {
	DisplayName string `json:"displayName"`
	Network     string `json:"network"`
	Type        string `json:"type"`
}

type ApplePayToken struct {
	PaymentData           ApplePaymentData   `json:"paymentData"`
	PaymentMethod         ApplePaymentMethod `json:"paymentMethod"`
	TransactionIdentifier string             `json:"transactionIdentifier"`
}

type ApplePayBillingContact struct {
	AddressLines       []string `json:"addressLines,omitempty"`
	Locality           string   `json:"locality,omitempty"`
	PostalCode         string   `json:"postalCode,omitempty"`
	AdministrativeArea string   `json:"administrativeArea,omitempty"`
	CountryCode        string   `json:"countryCode,omitempty"`
}

type StartApplePaySession struct {
	Domain      string `json:"domain"`
	DisplayName string `json:"displayName"`
}

type LinkApplePay struct {
	Token          ApplePayToken          `json:"token"`
	BillingContact ApplePayBillingContact `json:"billingContact,omitempty"`
}

type LinkedApplePayPaymentMethod struct {
	// ID of the payment method
	PaymentMethodID   string   `json:"paymentMethodID"`
	PaymentMethodType string   `json:"paymentMethodType"`
	ApplePay          ApplePay `json:"applePay"`
}

type ApplePay struct {
	Brand           string     `json:"brand,omitempty"`
	CardType        string     `json:"cardType,omitempty"`
	CardDisplayName string     `json:"cardDisplayName,omitempty"`
	Fingerprint     string     `json:"fingerprint,omitempty"`
	Expiration      Expiration `json:"expiration,omitempty"`
	DynamicLastFour string     `json:"dynamicLastFour,omitempty"`
}

// CreateApplePayDomain creates a new Apple Pay domain for the given customer account
// https://docs.moov.io/api/#tag/Cards/operation/registerApplePayMerchantDomains
func (c Client) CreateApplePayDomain(ctx context.Context, accountID string, domain ApplePayDomains) (*ApplePayDomainsResponse, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathApplePayDomains, accountID),
		AcceptJson(),
		JsonBody(domain))
	if err != nil {
		return nil, err
	}

	errors.Join()

	switch resp.Status() {
	case StatusCompleted:
		return CompletedObjectOrError[ApplePayDomainsResponse](resp)
	default:
		return nil, resp.Error()
	}
}

type PatchApplyPayDomains struct {
	AddDomains    []string `json:"addDomains,omitempty"`
	RemoveDomains []string `json:"removeDomains,omitempty"`
}

// UpdateApplePayDomain updates an existing Apple Pay domain for the given customer account
// https://docs.moov.io/api/#tag/Cards/operation/updateApplePayMerchantDomains
func (c Client) UpdateApplePayDomain(ctx context.Context, accountID string, patch PatchApplyPayDomains) error {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPut, pathApplePayDomains, accountID),
		AcceptJson(),
		JsonBody(patch))
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}

// GetApplePayDomain retrieves an Apple Pay domain for the given customer account
// https://docs.moov.io/api/#tag/Cards/operation/getApplePayMerchantDomains
func (c Client) GetApplePayDomain(ctx context.Context, accountID string) (*ApplePayDomainsResponse, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathApplePayDomains, accountID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[ApplePayDomainsResponse](resp)
}

// CreateApplePaySession creates a new Apple Pay session for the given customer account
// https://docs.moov.io/api/#tag/Cards/operation/getApplePayMerchantDomains
func (c Client) StartApplePaySession(ctx context.Context, accountID string, req StartApplePaySession) (*string, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathApplePaySessions, accountID),
		AcceptJson(),
		JsonBody(req))
	if err != nil {
		return nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		return UnmarshalObjectResponse[string](resp)
	default:
		return nil, resp.Error()
	}
}

// ApplePayToken creates a new Apple Pay token for the given customer account
// https://docs.moov.io/api/#tag/Cards/operation/getApplePayMerchantDomains
func (c Client) LinkApplePayToken(ctx context.Context, accountID string, req LinkApplePay) (*LinkedApplePayPaymentMethod, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathApplePayTokens, accountID),
		AcceptJson(),
		JsonBody(req))
	if err != nil {
		return nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		return UnmarshalObjectResponse[LinkedApplePayPaymentMethod](resp)
	default:
		return nil, resp.Error()
	}
}
