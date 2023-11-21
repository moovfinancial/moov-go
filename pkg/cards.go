package moov

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Card struct {
	CardID             string             `json:"cardID,omitempty"`
	Fingerprint        string             `json:"fingerprint,omitempty"`
	Brand              string             `json:"brand,omitempty"`
	CardType           string             `json:"cardType,omitempty"`
	LastFourCardNumber string             `json:"lastFourCardNumber,omitempty"`
	Bin                string             `json:"bin,omitempty"`
	Expiration         Expiration         `json:"expiration,omitempty"`
	HolderName         string             `json:"holderName,omitempty"`
	BillingAddress     Address            `json:"billingAddress,omitempty"`
	CardVerification   CardVerification   `json:"cardVerification,omitempty"`
	Issuer             string             `json:"issuer,omitempty"`
	IssuerCountry      string             `json:"issuerCountry,omitempty"`
	CardOnFile         bool               `json:"cardOnFile,omitempty"`
	MerchantAccountID  string             `json:"merchantAccountID,omitempty"`
	CardAccountUpdater CardAccountUpdater `json:"cardAccountUpdater,omitempty"`
	DomesticPushToCard string             `json:"domesticPushToCard,omitempty"`
}

type Expiration struct {
	Month string `json:"month,omitempty"`
	Year  string `json:"year,omitempty"`
}

type CardVerification struct {
	Cvv          string `json:"cvv,omitempty"`
	AddressLine1 string `json:"addressLine1,omitempty"`
	PostalCode   string `json:"postalCode,omitempty"`
}

type CardAccountUpdater struct {
	UpdatedOn  time.Time `json:"updatedOn,omitempty"`
	UpdateType string    `json:"updateType,omitempty"`
}

type CardDetails struct {
	Status                   string            `json:"status,omitempty"`
	FailureCode              string            `json:"failureCode,omitempty"`
	DynamicDescriptor        string            `json:"dynamicDescriptor,omitempty"`
	TransactionSource        string            `json:"transactionSource,omitempty"`
	InterchangeQualification string            `json:"interchangeQualification,omitempty"`
	StatusUpdates            CardStatusUpdates `json:"statusUpdates,omitempty"`
}

type CardStatusUpdates struct {
	Initiated time.Time `json:"initiated,omitempty"`
	Confirmed time.Time `json:"confirmed,omitempty"`
	Settled   time.Time `json:"settled,omitempty"`
	Failed    time.Time `json:"failed,omitempty"`
	Canceled  time.Time `json:"canceled,omitempty"`
	Completed time.Time `json:"completed,omitempty"`
}

type CardPayload struct {
	Card
	CardNumber string `json:"cardNumber,omitempty"`
	CardCvv    string `json:"cardCvv,omitempty"`
}

// CreateCard creates a new card for the given customer linked to their account
// https://docs.moov.io/api/#tag/Cards/operation/card
func (c Client) CreateCard(accountID string, card Card, cardNumber string, cardCvv string) (Card, error) {
	url := fmt.Sprintf("%s/%s", baseURL, fmt.Sprintf(pathCards, accountID))

	payload := CardPayload{
		Card:       card,
		CardNumber: cardNumber,
		CardCvv:    cardCvv,
	}

	header := map[string]string{"X-Wait-For": "payment-method"}

	resCard := Card{}
	body, statusCode, err := GetHTTPResponse(c, http.MethodPost, url, payload, header)
	if err != nil {
		return resCard, err
	}

	switch statusCode {
	case http.StatusOK:
		// card created
		err = json.Unmarshal(body, &resCard)
		if err != nil {
			log.Println("Error unmarshalling JSON:", err)
		}
		return resCard, nil
	case http.StatusUnauthorized:
		return resCard, ErrAuthCreditionalsNotSet
	case http.StatusNotFound:
		return resCard, ErrNoAccount
	case http.StatusConflict:
		return resCard, errors.New("attempted to link card that already exists on the account")
	case http.StatusUnprocessableEntity:
		return resCard, errors.New("the supplied card data appeared invalid or was declined by the issuer")
	case http.StatusTooManyRequests:
		return resCard, errors.New("request was refused due to rate limiting")
	}
	return resCard, nil
}

// ListCards lists all cards for the given customer Moov account
// https://docs.moov.io/api/#tag/Cards/operation/listCards
func (c Client) ListCards(accountID string) ([]Card, error) {
	var resCards []Card
	url := fmt.Sprintf("%s/%s", baseURL, fmt.Sprintf(pathCards, accountID))

	body, statusCode, err := GetHTTPResponse(c, http.MethodGet, url, nil, nil)
	if err != nil {
		return resCards, err
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &resCards)
		if err != nil {
			log.Println("Error unmarshalling JSON:", err)
		}
		return resCards, nil
	case http.StatusUnauthorized:
		return resCards, ErrAuthCreditionalsNotSet
	case http.StatusNotFound:
		return resCards, ErrNoAccount
	case http.StatusUnprocessableEntity:
		log.Println("UnprocessableEntity")
	}
	return resCards, nil
}

// GetCard retrieves a card for the given customer Moov account
// https://docs.moov.io/api/#tag/Cards/operation/getCard
func (c Client) GetCard(accountID string, cardID string) (Card, error) {
	resCard := Card{}
	url := fmt.Sprintf("%s/%s/%s", baseURL, fmt.Sprintf(pathCards, accountID), cardID)

	body, statusCode, err := GetHTTPResponse(c, http.MethodGet, url, nil, nil)
	if err != nil {
		return resCard, err
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &resCard)
		if err != nil {
			log.Println("Error unmarshalling JSON:", err)
		}
		return resCard, nil
	case http.StatusUnauthorized:
		return resCard, ErrAuthCreditionalsNotSet
	case http.StatusNotFound:
		return resCard, ErrNoAccount
	case http.StatusUnprocessableEntity:
		log.Println("UnprocessableEntity")
	}
	return resCard, nil
}

// UpdateCard Update a linked card and/or resubmit it for verification.
// If a value is provided for CVV, a new verification ($0 authorization) will be submitted for the card. Updating the expiration date or address will update the information stored on file for the card but will not be verified
// https://docs.moov.io/api/#tag/Cards/operation/updateCard
func (c Client) UpdateCard(accountID string, cardID string, card Card, cardCvv string) (Card, error) {
	url := fmt.Sprintf("%s/%s/%s", baseURL, fmt.Sprintf(pathCards, accountID), cardID)

	payload := CardPayload{
		Card:    card,
		CardCvv: cardCvv,
	}

	resCard := Card{}
	body, statusCode, err := GetHTTPResponse(c, http.MethodPatch, url, payload, nil)
	if err != nil {
		return resCard, err
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &resCard)
		if err != nil {
			log.Println("Error unmarshalling JSON:", err)
		}
		return resCard, nil
	case http.StatusUnauthorized:
		return resCard, ErrAuthCreditionalsNotSet
	case http.StatusNotFound:
		return resCard, ErrNoAccount
	case http.StatusConflict:
		return resCard, errors.New("attempted to link card that already exists on the account")
	case http.StatusUnprocessableEntity:
		return resCard, errors.New("the supplied card data appeared invalid or was declined by the issuer")
	case http.StatusTooManyRequests:
		return resCard, errors.New("request was refused due to rate limiting")
	}
	return resCard, nil
}

// DisableCard disables a card associated with a Moov account
// https://docs.moov.io/api/#tag/Cards/operation/deleteCard
func (c Client) DisableCard(accountID string, cardID string) error {
	url := fmt.Sprintf("%s/%s/%s", baseURL, fmt.Sprintf(pathCards, accountID), cardID)

	_, statusCode, err := GetHTTPResponse(c, http.MethodDelete, url, nil, nil)
	if err != nil {
		return err
	}

	switch statusCode {
	case http.StatusNoContent:
		// card deleted
		return nil
	case http.StatusUnauthorized:
		return ErrAuthCreditionalsNotSet
	case http.StatusNotFound:
		return ErrNoAccount
	case http.StatusUnprocessableEntity:
		log.Println("UnprocessableEntity")
	}
	return nil
}
