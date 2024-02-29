package moov

import (
	"context"
	"errors"
	"net/http"
	"time"
)

var (
	ErrDuplicateBankAccount = errors.New("duplciate bank account or invalid routing number")
	ErrNoMicroDeposit       = errors.New("no account with the specified accountID was found or micro-deposits have not been sent for the source")
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

// Plaid is a direct to plaid integration. Use this if you signed an agreement with Plaid directly.
type Plaid struct {
	// Token is the plaid processor_token
	Token string `json:"token"`
}

// PlaidLink is a Moov-managed Plaid integration. Use this if you purchase Plaid through Moov.
type PlaidLink struct {
	// PublicToken is the plaid public_token
	PublicToken string `json:"publicToken"`
}

// MX is authorization code of a MX account which allows a processor to retrieve a linked payment account.
type MX struct {
	AuthorizationCode string `json:"authorizationCode"`
}

// CreateBankAccount creates a new bank account for the given customer account
func (c Client) CreateBankAccount(ctx context.Context, accountID string, bankAccount BankAccount) (*BankAccount, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathBankAccounts, accountID),
		AcceptJson(),
		JsonBody(bankAccount))
	if err != nil {
		return nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		return CompletedObjectOrError[BankAccount](resp)
	case StatusStateConflict:
		return nil, ErrDuplicateBankAccount
	default:
		return nil, resp.Error()
	}
}

// GetBankAccount retrieves a bank account for the given customer account
func (c Client) GetBankAccount(ctx context.Context, accountID string, bankAccountID string) (*BankAccount, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathBankAccounts, accountID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[BankAccount](resp)
}

// DeleteBankAccount deletes a bank account for the given customer account
func (c Client) DeleteBankAccount(ctx context.Context, accountID string, bankAccountID string) error {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodDelete, pathBankAccounts, accountID))
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}

// ListBankAccounts lists all bank accounts for the given customer account
func (c Client) ListBankAccounts(ctx context.Context, accountID string) ([]BankAccount, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathBankAccounts, accountID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[BankAccount](resp)
}

// MicroDepositInitiate creates a new micro deposit verification for the given bank account
func (c Client) MicroDepositInitiate(ctx context.Context, accountID string, bankAccountID string) error {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodPost, pathMicroDeposits, accountID, bankAccountID))
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}

// MicroDepositConfirm confirms a micro deposit verification for the given bank account
func (c Client) MicroDepositConfirm(ctx context.Context, accountID string, bankAccountID string, amounts []int) error {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPut, pathMicroDeposits, accountID, bankAccountID),
		AcceptJson(),
		JsonBody(map[string][]int{"amounts": amounts}))
	if err != nil {
		return err
	}

	switch resp.Status() {
	case StatusCompleted:
		return nil
	case StatusNotFound:
		return ErrNoMicroDeposit
	case StatusStateConflict:
		return ErrAmountIncorrect
	default:
		return resp.Error()
	}
}

// CreatePlaidLink creates a new bank account for the given customer account using Plaid processor_token
func (c Client) CreateBankAccountWithPlaid(ctx context.Context, accountID string, plaid Plaid) (*BankAccount, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathBankAccounts, accountID),
		AcceptJson(),
		JsonBody(plaid))
	if err != nil {
		return nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		return CompletedObjectOrError[BankAccount](resp)
	case StatusStateConflict:
		return nil, ErrDuplicateBankAccount
	default:
		return nil, resp.Error()
	}
}

// CreateBankAccountWithPlaidLink creates a new bank account for the given customer account using Plaid public_token
func (c Client) CreateBankAccountWithPlaidLink(ctx context.Context, accountID string, plaid PlaidLink) (*BankAccount, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathBankAccounts, accountID),
		AcceptJson(),
		JsonBody(plaid))
	if err != nil {
		return nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		return CompletedObjectOrError[BankAccount](resp)
	case StatusStateConflict:
		return nil, ErrDuplicateBankAccount
	default:
		return nil, resp.Error()
	}
}

// CreateBankAccountWithMX creates a new bank account for the given customer account using MX account
func (c Client) CreateBankAccountWithMX(ctx context.Context, accountID string, mx MX) (*BankAccount, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathBankAccounts, accountID),
		AcceptJson(),
		JsonBody(mx))
	if err != nil {
		return nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		return CompletedObjectOrError[BankAccount](resp)
	case StatusStateConflict:
		return nil, ErrDuplicateBankAccount
	default:
		return nil, resp.Error()
	}
}
