package moov

import (
	"bytes"
	"context"
	"encoding/json"
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
	Account   *BankAccount `json:"account,omitempty"`
	Plaid     *Plaid       `json:"plaid,omitempty"`
	PlaidLink *PlaidLink   `json:"plaidLink,omitempty"`
	MX        *MX          `json:"mx,omitempty"`
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

type CreateBankAccountType callArg

func WithBankAccount(bankAccount BankAccount) CreateBankAccountType {
	return callBuilderFn((func(call *callBuilder) error {
		bankAccountJSON, err := json.Marshal(bankAccount)
		if err != nil {
			return err
		}
		call.body = bytes.NewReader(bankAccountJSON)
		return nil
	}))
}

func WithPlaidLink(plaidLink PlaidLink) CreateBankAccountType {
	return callBuilderFn((func(call *callBuilder) error {
		// need to use a map to add the top level "plaidLink" json key. expected format is:
		//	{
		// 		"plaidLink": {
		// 			"publicToken": "insert token"
		// 		}
		// 	}

		plaidLinkMap := map[string]PlaidLink{
			"plaidLink": plaidLink,
		}
		bankAccountJSON, err := json.Marshal(plaidLinkMap)
		if err != nil {
			return err
		}
		call.body = bytes.NewReader(bankAccountJSON)
		return nil
	}))
}

func WithPlaid(plaid Plaid) CreateBankAccountType {
	return callBuilderFn((func(call *callBuilder) error {
		// need to use a map to add the top level "plaid" json key. expected format is:
		//	{
		// 		"plaid": {
		// 			"token": "insert token"
		// 		}
		// 	}

		plaidMap := map[string]Plaid{
			"plaid": plaid,
		}
		bankAccountJSON, err := json.Marshal(plaidMap)
		if err != nil {
			return err
		}
		call.body = bytes.NewReader(bankAccountJSON)
		return nil
	}))
}

func WithMX(mx MX) CreateBankAccountType {
	return callBuilderFn((func(call *callBuilder) error {
		// need to use a map to add the top level "mx" json key. expected format is:
		//	{
		// 		"mx": {
		// 			"authorizationCode": "testing"
		// 		}
		// 	}

		mxMap := map[string]MX{
			"mx": mx,
		}
		bankAccountJSON, err := json.Marshal(mxMap)
		if err != nil {
			return err
		}
		call.body = bytes.NewReader(bankAccountJSON)
		return nil
	}))
}

// CreateBankAccount creates a new bank account for the given customer account
func (c Client) CreateBankAccount(ctx context.Context, accountID string, opts ...CreateBankAccountType) (*BankAccount, error) {
	args := prependArgs(opts, AcceptJson())
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathBankAccounts, accountID),
		AcceptJson(),
		JsonBody(args))
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
		Endpoint(http.MethodGet, pathBankAccount, accountID, bankAccountID),
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
