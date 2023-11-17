package moov

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
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

const (
	baseURL  = "https://api.moov.io"
	endpoint = "accounts/%s/bank-accounts"
)

// CreateBankAccount creates a new bank account for the given customer account
func (c Client) CreateBankAccount(accountID string, bankAccount BankAccount) (BankAccount, error) {
	url := fmt.Sprintf("%s/%s", baseURL, fmt.Sprintf(endpoint, accountID))

	accountPayload := BankAccountPayload{
		Account: bankAccount,
	}

	payload, err := json.Marshal(accountPayload)
	if err != nil {
		return bankAccount, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return bankAccount, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.Credentials.PublicKey, c.Credentials.SecretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return bankAccount, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	respAccount := BankAccount{}

	switch resp.StatusCode {
	case http.StatusOK:
		// Account created
		err = json.Unmarshal(body, &respAccount)
		if err != nil {
			log.Println("Error unmarshalling JSON:", err)
		}
		return respAccount, nil
	case http.StatusUnauthorized:
		return respAccount, ErrAuthCreditionalsNotSet
	case http.StatusUnprocessableEntity:
		log.Println("UnprocessableEntity")
	}
	return respAccount, nil
}

// GetBankAccount retrieves a bank account for the given customer account
func (c Client) GetBankAccount(accountID string, bankAccountID string) (BankAccount, error) {
	resAccount := BankAccount{}
	url := fmt.Sprintf("%s/%s/%s", baseURL, fmt.Sprintf(endpoint, accountID), bankAccountID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return BankAccount{}, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.Credentials.PublicKey, c.Credentials.SecretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return resAccount, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	switch resp.StatusCode {
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
	url := fmt.Sprintf("%s/%s/%s", baseURL, fmt.Sprintf(endpoint, accountID), bankAccountID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.Credentials.PublicKey, c.Credentials.SecretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
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
	url := fmt.Sprintf("%s/%s", baseURL, fmt.Sprintf(endpoint, accountID))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return resAccounts, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.Credentials.PublicKey, c.Credentials.SecretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return resAccounts, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	switch resp.StatusCode {
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
	url := fmt.Sprintf("%s/%s/%s/micro-deposits", baseURL, fmt.Sprintf(endpoint, accountID), bankAccountID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.Credentials.PublicKey, c.Credentials.SecretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
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
	url := fmt.Sprintf("%s/%s/%s/micro-deposits", baseURL, fmt.Sprintf(endpoint, accountID), bankAccountID)

	payload, err := json.Marshal(map[string][]int{"amounts": amounts})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.Credentials.PublicKey, c.Credentials.SecretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized:
		return ErrAuthCreditionalsNotSet
	case http.StatusUnprocessableEntity:
		log.Println("UnprocessableEntity")
	}
	return nil
}
