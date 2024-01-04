package moov

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ApplePay struct {
	Brand           string     `json:"brand,omitempty"`
	CardType        string     `json:"cardType,omitempty"`
	CardDisplayName string     `json:"cardDisplayName,omitempty"`
	Fingerprint     string     `json:"fingerprint,omitempty"`
	Expiration      Expiration `json:"expiration,omitempty"`
	DynamicLastFour string     `json:"dynamicLastFour,omitempty"`
}

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

type LinkApplePay struct {
	Token          ApplePayToken          `json:"token"`
	BillingContact ApplePayBillingContact `json:"billingContact,omitempty"`
}

// CreateApplePayDomain creates a new Apple Pay domain for the given customer account
// https://docs.moov.io/api/#tag/Cards/operation/registerApplePayMerchantDomains
func (c Client) CreateApplePayDomain(accountID string, displayName string, domains []string) (ApplePayDomainsResponse, error) {
	resp := ApplePayDomainsResponse{}

	url := fmt.Sprintf("%s/%s/domains", baseURL, fmt.Sprintf(pathApplePay, accountID))
	payload := ApplePayDomains{
		DisplayName: displayName,
		Domains:     domains,
	}
	body, statusCode, err := c.GetHTTPResponse(http.MethodPost, url, payload, nil)
	if err != nil {
		return resp, err
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &resp)
		if err != nil {
			return resp, ErrDefault(statusCode)
		}
		return resp, nil
	case http.StatusBadRequest:
		return resp, ErrBadRequest
	case http.StatusConflict:
		return resp, ErrDuplicatedApplePayDomain
	case http.StatusExpectationFailed:
		return resp, ErrDomainsNotVerified
	case http.StatusTooManyRequests:
		return resp, ErrRateLimit
	}
	return resp, ErrDefault(statusCode)
}

// UpdateApplePayDomain updates an existing Apple Pay domain for the given customer account
// https://docs.moov.io/api/#tag/Cards/operation/updateApplePayMerchantDomains
func (c Client) UpdateApplePayDomain(accountID string, addDomains []string, removeDomains []string) error {
	url := fmt.Sprintf("%s/%s/domains", baseURL, fmt.Sprintf(pathApplePay, accountID))
	payload := map[string][]string{
		"addDomains":    addDomains,
		"removeDomains": removeDomains,
	}
	_, statusCode, err := c.GetHTTPResponse(http.MethodPatch, url, payload, nil)
	if err != nil {
		return err
	}

	switch statusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusBadRequest:
		return ErrBadRequest
	case http.StatusNotFound:
		return ErrDomainsNotRegistered
	case http.StatusExpectationFailed:
		return ErrDomainsNotVerified
	case http.StatusTooManyRequests:
		return ErrRateLimit
	}
	return ErrDefault(statusCode)
}

// GetApplePayDomain retrieves an Apple Pay domain for the given customer account
// https://docs.moov.io/api/#tag/Cards/operation/getApplePayMerchantDomains
func (c Client) GetApplePayDomain(accountID string) (ApplePayDomainsResponse, error) {
	resp := ApplePayDomainsResponse{}

	url := fmt.Sprintf("%s/%s/domains", baseURL, fmt.Sprintf(pathApplePay, accountID))
	body, statusCode, err := c.GetHTTPResponse(http.MethodGet, url, nil, nil)
	if err != nil {
		return resp, err
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &resp)
		if err != nil {
			return resp, ErrDefault(statusCode)
		}
		return resp, nil
	case http.StatusBadRequest:
		return resp, ErrBadRequest
	case http.StatusConflict:
		return resp, ErrDuplicatedApplePayDomain
	case http.StatusExpectationFailed:
		return resp, ErrDomainsNotVerified
	case http.StatusTooManyRequests:
		return resp, ErrRateLimit
	}
	return resp, ErrDefault(statusCode)
}

// CreateApplePaySession creates a new Apple Pay session for the given customer account
// https://docs.moov.io/api/#tag/Cards/operation/getApplePayMerchantDomains
func (c Client) CreateApplePaySession(accountID string, domain string, displayName string) error {
	url := fmt.Sprintf("%s/%s/sessions", baseURL, fmt.Sprintf(pathApplePay, accountID))
	payload := map[string]string{
		"domains":     domain,
		"displayName": displayName,
	}
	_, statusCode, err := c.GetHTTPResponse(http.MethodPost, url, payload, nil)
	if err != nil {
		return err
	}

	switch statusCode {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		return ErrBadRequest
	case http.StatusConflict:
		return ErrDuplicatedApplePayDomain
	case http.StatusExpectationFailed:
		return ErrDomainsNotVerified
	case http.StatusTooManyRequests:
		return ErrRateLimit
	}
	return ErrDefault(statusCode)
}

// ApplePayToken creates a new Apple Pay token for the given customer account
// https://docs.moov.io/api/#tag/Cards/operation/getApplePayMerchantDomains
func (c Client) ApplePayToken(accountID string, token ApplePayToken, contact ApplePayBillingContact) (PaymentMethod, error) {
	resp := PaymentMethod{}

	url := fmt.Sprintf("%s/%s/tokens", baseURL, fmt.Sprintf(pathApplePay, accountID))
	payload := LinkApplePay{
		Token:          token,
		BillingContact: contact,
	}
	body, statusCode, err := c.GetHTTPResponse(http.MethodPost, url, payload, nil)
	if err != nil {
		return resp, err
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &resp)
		if err != nil {
			return resp, ErrDefault(statusCode)
		}
		return resp, nil
	case http.StatusUnprocessableEntity:
		return resp, ErrLinkingApplePayToken
	case http.StatusTooManyRequests:
		return resp, ErrRateLimit
	}
	return resp, ErrDefault(statusCode)
}
