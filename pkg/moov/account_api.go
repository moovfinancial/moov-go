package moov

import (
	"context"
	"net/http"
	"strconv"
)

// CreateAccount creates a new account.
func (c Client) CreateAccount(ctx context.Context, account CreateAccount) (*Account, *Account, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathAccounts),
		AcceptJson(),
		WaitFor("connection"),
		JsonBody(account))
	if err != nil {
		return nil, nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		a, err := UnmarshalObjectResponse[Account](resp)
		return a, nil, err
	case StatusStarted:
		a, err := UnmarshalObjectResponse[Account](resp)
		return nil, a, err
	default:
		return nil, nil, resp
	}
}

// GetAccount returns an account based on accountID.
func (c Client) GetAccount(ctx context.Context, accountID string) (*Account, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathAccount, accountID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Account](resp)
}

// UpdateAccount updates an account.
func (c Client) UpdateAccount(ctx context.Context, account Account) (*Account, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPatch, pathAccount, account.AccountID),
		AcceptJson(),
		JsonBody(account))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Account](resp)
}

// Func that applies a filter and returns an error if validation fails
type ListAccountFilter callArg

// WithAccountName if provided, this query will attempt to find matches against the following Account and Profile fields: diplayName, firstName, middleName, lastName, legalBusinessName
func WithAccountName(name string) ListAccountFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["name"] = name
		return nil
	})
}

// WithAccountEmail filter connected accounts by email address.
func WithAccountEmail(email string) ListAccountFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["email"] = email
		return nil
	})
}

// WithAccountType filter type possible values: individual, business
func WithAccountType(accountType string) ListAccountFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["type"] = accountType
		return nil
	})
}

// WithAccountForeignID filter as an optional alias from a foreign/external system which can be used to reference this resource.
func WithAccountForeignID(foreignID string) ListAccountFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["foreignID"] = foreignID
		return nil
	})
}

// WithAccountVerificationStatus possible values: unverified, pending, resubmit, review, verified, failed
func WithAccountVerificationStatus(verificationStatus string) ListAccountFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["verification_status"] = verificationStatus
		return nil
	})
}

// WithAccountIncludeDisconnected if true, the response will include disconnected accounts.
func WithAccountIncludeDisconnected() ListAccountFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["includeDisconnected"] = "true"
		return nil
	})
}

// WithAccountCount value to limit the number of results in the query. Default is 20
func WithAccountCount(count int) ListAccountFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["count"] = strconv.Itoa(count)
		return nil
	})
}

// WithAccountSkip the number of items to offset before starting to collect the result set
func WithAccountSkip(skip int) ListAccountFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["skip"] = strconv.Itoa(skip)
		return nil
	})
}

// ListAccounts returns a list of accounts.
func (c Client) ListAccounts(ctx context.Context, opts ...ListAccountFilter) ([]Account, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathAccounts),
		prependArgs(opts, AcceptJson())...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[Account](resp)
}

func (c Client) DisconnectAccount(ctx context.Context, accountID string) error {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodDelete, pathAccount, accountID),
		AcceptJson())
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}
