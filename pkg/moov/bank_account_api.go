package moov

import (
	"context"
	"errors"
	"net/http"
)

type CreateBankAccountType callArg

func WithBankAccount(bankAccount BankAccountRequest) CreateBankAccountType {
	return JsonBody(createBankAccount{
		Account: &bankAccount,
	})
}

func WithPlaid(plaid PlaidRequest) CreateBankAccountType {
	return JsonBody(createBankAccount{
		Plaid: &plaid,
	})
}

func WithPlaidLink(plaidLink PlaidLinkRequest) CreateBankAccountType {
	return JsonBody(createBankAccount{
		PlaidLink: &plaidLink,
	})
}

func WithMX(mx MXRequest) CreateBankAccountType {
	return JsonBody(createBankAccount{
		MX: &mx,
	})
}

func WaitForPaymentMethod() CreateBankAccountType {
	return WaitFor("payment-method")
}

// CreateBankAccount creates a new bank account for the given customer account
func (c Client) CreateBankAccount(ctx context.Context, accountID string, opts ...CreateBankAccountType) (*BankAccount, error) {
	args := prependArgs(opts, AcceptJson())
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathBankAccounts, accountID),
		args...)
	if err != nil {
		return nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		return CompletedObjectOrError[BankAccount](resp)
	case StatusStateConflict:
		return nil, errors.Join(ErrAlreadyExists, resp)
	default:
		return nil, resp
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
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodDelete, pathBankAccount, accountID, bankAccountID))
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
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodPost, pathBankAccountMicroDeposits, accountID, bankAccountID))
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}

// MicroDepositConfirm confirms a micro deposit verification for the given bank account
func (c Client) MicroDepositConfirm(ctx context.Context, accountID string, bankAccountID string, amounts []int) error {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPut, pathBankAccountMicroDeposits, accountID, bankAccountID),
		AcceptJson(),
		JsonBody(map[string][]int{"amounts": amounts}))
	if err != nil {
		return err
	}

	switch resp.Status() {
	case StatusCompleted:
		return nil
	case StatusStateConflict:
		return errors.Join(ErrMicroDepositAmountsIncorrect, resp)
	default:
		return resp
	}
}

func (c Client) InstantVerificationInitiate(ctx context.Context, accountID, bankAccountID string) error {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodPost, pathBankAccountInstantVerification, accountID, bankAccountID))
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}

func (c Client) InstantVerificationComplete(ctx context.Context, accountID, bankAccountID, code string) error {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPut, pathBankAccountInstantVerification, accountID, bankAccountID),
		AcceptJson(),
		JsonBody(map[string]string{"code": code}))
	if err != nil {
		return err
	}

	switch resp.Status() {
	case StatusCompleted:
		return nil
	case StatusStateConflict:
		return errors.Join(ErrInstantVerificationFailed, resp)
	default:
		return resp
	}
}
