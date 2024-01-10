package moov

import (
	"context"
	"net/http"
)

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

	return CompletedNilOrError(resp)
}
