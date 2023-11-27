package moov

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

var (
	ErrDuplicateBankAccount = errors.New("Duplciate bank account or invalid routing number")
)

type BankAccount struct {
	BankAccountID         string `json:"bankAccountID,omitempty"`
	Fingerprint           string `json:"fingerprint,omitempty"`
	Status                string `json:"status,omitempty"`
	HolderName            string `json:"holderName,omitempty"`
	HolderType            string `json:"holderType,omitempty"`
	BankName              string `json:"bankName,omitempty"`
	BankAccountType       string `json:"bankAccountType,omitempty"`
	AccountNumber         string `json:"accountNumber,omitempty"`
	RoutingNumber         string `json:"routingNumber,omitempty"`
	LastFourAccountNumber string `json:"lastFourAccountNumber,omitempty"`
}

type AchDetails struct {
	Status                  string           `json:"status,omitempty"`
	TraceNumber             string           `json:"traceNumber,omitempty"`
	Return                  Return           `json:"return,omitempty"`
	Correction              Correction       `json:"correction,omitempty"`
	CompanyEntryDescription string           `json:"companyEntryDescription,omitempty"`
	OriginatingCompanyName  string           `json:"originatingCompanyName,omitempty"`
	StatusUpdates           ACHStatusUpdates `json:"statusUpdates,omitempty"`
	DebitHoldPeriod         string           `json:"debitHoldPeriod,omitempty"`
}

type Correction struct {
	Code        string `json:"code,omitempty"`
	Reason      string `json:"reason,omitempty"`
	Description string `json:"description,omitempty"`
}

type Return struct {
	Code        string `json:"code,omitempty"`
	Reason      string `json:"reason,omitempty"`
	Description string `json:"description,omitempty"`
}

type ACHStatusUpdates struct {
	Initiated  time.Time `json:"initiated,omitempty"`
	Originated time.Time `json:"originated,omitempty"`
	Corrected  time.Time `json:"corrected,omitempty"`
	Returned   time.Time `json:"returned,omitempty"`
	Completed  time.Time `json:"completed,omitempty"`
}

type BankAccountPayload struct {
	Account BankAccount `json:"account"`
}

// CreateBankAccount creates a new bank account for the given customer account
func (c Client) CreateBankAccount(accountID string, bankAccount BankAccount) (BankAccount, error) {
	url := fmt.Sprintf("%s/%s", baseURL, fmt.Sprintf(pathBankAccounts, accountID))

	accountPayload := BankAccountPayload{
		Account: bankAccount,
	}

	respBankAccount := BankAccount{}
	body, statusCode, err := GetHTTPResponse(c, http.MethodPost, url, accountPayload, nil)
	if err != nil {
		return respBankAccount, err
	}

	switch statusCode {
	case http.StatusOK:
		// Account created
		err = json.Unmarshal(body, &respBankAccount)
		if err != nil {
			log.Println("Error unmarshalling JSON:", err)
		}
		return respBankAccount, nil
	case http.StatusConflict:
		return respBankAccount, ErrDuplicateBankAccount
	case http.StatusUnauthorized:
		return respBankAccount, ErrAuthCreditionalsNotSet
	case http.StatusUnprocessableEntity:
		log.Println("UnprocessableEntity")
	}
	return respBankAccount, nil
}

// GetBankAccount retrieves a bank account for the given customer account
func (c Client) GetBankAccount(accountID string, bankAccountID string) (BankAccount, error) {
	resAccount := BankAccount{}
	url := fmt.Sprintf("%s/%s/%s", baseURL, fmt.Sprintf(pathBankAccounts, accountID), bankAccountID)

	body, statusCode, err := GetHTTPResponse(c, http.MethodGet, url, nil, nil)
	if err != nil {
		return resAccount, err
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &resAccount)
		if err != nil {
			log.Println("Error unmarshalling JSON:", err)
		}
		return resAccount, nil
	case http.StatusUnauthorized:
		return resAccount, ErrAuthCreditionalsNotSet
	case http.StatusUnprocessableEntity:
		log.Println("UnprocessableEntity")
	}
	return resAccount, nil
}

// DeleteBankAccount deletes a bank account for the given customer account
func (c Client) DeleteBankAccount(accountID string, bankAccountID string) error {
	url := fmt.Sprintf("%s/%s/%s", baseURL, fmt.Sprintf(pathBankAccounts, accountID), bankAccountID)

	_, statusCode, err := GetHTTPResponse(c, http.MethodDelete, url, nil, nil)
	if err != nil {
		return err
	}

	switch statusCode {
	case http.StatusNoContent:
		// Account deleted
		return nil
	case http.StatusUnauthorized:
		return ErrAuthCreditionalsNotSet
	case http.StatusUnprocessableEntity:
		log.Println("UnprocessableEntity")
	}
	return nil
}

// ListBankAccounts lists all bank accounts for the given customer account
func (c Client) ListBankAccounts(accountID string) ([]BankAccount, error) {
	var resAccounts []BankAccount
	url := fmt.Sprintf("%s/%s", baseURL, fmt.Sprintf(pathBankAccounts, accountID))

	body, statusCode, err := GetHTTPResponse(c, http.MethodGet, url, nil, nil)
	if err != nil {
		return resAccounts, err
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &resAccounts)
		if err != nil {
			log.Println("Error unmarshalling JSON:", err)
		}
		return resAccounts, nil
	case http.StatusUnauthorized:
		return resAccounts, ErrAuthCreditionalsNotSet
	case http.StatusUnprocessableEntity:
		log.Println("UnprocessableEntity")
	}
	return resAccounts, nil
}

// MicroDepositInitiate creates a new micro deposit verification for the given bank account
func (c Client) MicroDepositInitiate(accountID string, bankAccountID string) error {
	url := fmt.Sprintf("%s/%s/%s/micro-deposits", baseURL, fmt.Sprintf(pathBankAccounts, accountID), bankAccountID)

	_, statusCode, err := GetHTTPResponse(c, http.MethodPost, url, nil, nil)
	if err != nil {
		return err
	}

	switch statusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusUnauthorized:
		return ErrAuthCreditionalsNotSet
	case http.StatusUnprocessableEntity:
		log.Println("UnprocessableEntity")
	}
	return nil
}

// MicroDepositConfirm confirms a micro deposit verification for the given bank account
func (c Client) MicroDepositConfirm(accountID string, bankAccountID string, amounts []int) error {
	url := fmt.Sprintf("%s/%s/%s/micro-deposits", baseURL, fmt.Sprintf(pathBankAccounts, accountID), bankAccountID)

	_, statusCode, err := GetHTTPResponse(c, http.MethodPut, url, map[string][]int{"amounts": amounts}, nil)
	if err != nil {
		return err
	}

	switch statusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized:
		return ErrAuthCreditionalsNotSet
	case http.StatusUnprocessableEntity:
		log.Println("UnprocessableEntity")
	}
	return nil
}
