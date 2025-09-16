package moov

import (
	"context"
	"net/http"
	"strconv"
)

type AccountClient[T any, V any] struct {
	Version Version
}

// Create creates a new account with the provided configuration.
//
// It returns:
//   - created (*Account): The fully created account when the server responds with a 200 status code.
//   - started (*Account): The account details when the server responds with a 201 status code.
//   - err (error): Any error encountered during the account creation process.
//
// Only one of created or started will be non-nil, depending on the server's response.
func (ac *AccountClient[T, V]) Create(ctx context.Context, client Client, account T) (created, started *V, err error) {
	resp, err := client.CallHttp(ctx,
		Endpoint(http.MethodPost, pathAccounts),
		MoovVersion(ac.Version),
		AcceptJson(),
		WaitFor("connection"),
		JsonBody(account))
	if err != nil {
		return nil, nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		a, err := UnmarshalObjectResponse[V](resp)
		return a, nil, err
	case StatusStarted:
		a, err := UnmarshalObjectResponse[V](resp)
		return nil, a, err
	default:
		return nil, nil, resp
	}
}

// Get returns an account based on accountID.
func (ac AccountClient[T, V]) Get(ctx context.Context, client Client, accountID string) (*V, error) {
	resp, err := client.CallHttp(ctx,
		Endpoint(http.MethodGet, pathAccount, accountID),
		MoovVersion(ac.Version),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[V](resp)
}

// Patch updates an account.
func (ac AccountClient[T, V]) Patch(ctx context.Context, client Client, accountID string, account PatchAccount) (*V, error) {
	resp, err := client.CallHttp(ctx,
		Endpoint(http.MethodPatch, pathAccount, accountID),
		MoovVersion(ac.Version),
		AcceptJson(),
		JsonBody(account))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[V](resp)
}

// List returns a list of accounts.
func (ac AccountClient[T, V]) List(ctx context.Context, client Client, opts ...ListAccountFilter) ([]V, error) {
	resp, err := client.CallHttp(ctx,
		Endpoint(http.MethodGet, pathAccounts),
		prependArgs(opts, MoovVersion(ac.Version), AcceptJson())...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[V](resp)
}

// Disconnect severs the connection between you and the account specified.
// It will no longer be listed as active in the list of accounts.
// This also means you'll only have read-only access to the account going forward for reporting purposes.
func (ac AccountClient[T, V]) Disconnect(ctx context.Context, client Client, accountID string) error {
	resp, err := client.CallHttp(ctx,
		Endpoint(http.MethodDelete, pathAccount, accountID),
		MoovVersion(ac.Version),
		AcceptJson())
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}

// Legacy

// Only use for Preversioned API calls. Use mvxxxx.Accounts.Create(...) instead.
func (c *Client) CreateAccount(ctx context.Context, account CreateAccount) (created, started *Account, err error) {
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

// Use only for Preversioned API calls. Use mvxxxx.Accounts.Get(...) instead.
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

// Deprecated: use PatchAccount instead
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

// Use only for Preversioned API calls. Use mvxxxx.Accounts.Patch(...) instead.
// PatchAccount updates an account.
func (c Client) PatchAccount(ctx context.Context, accountID string, account PatchAccount) (*Account, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPatch, pathAccount, accountID),
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

// Only use for Preversioned API calls. Use mvxxxx.Accounts.List(...) instead.
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

// Only use for Preversioned API calls. Use mvxxxx.Accounts.Disconnect(...) instead.
func (c Client) DisconnectAccount(ctx context.Context, accountID string) error {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodDelete, pathAccount, accountID),
		AcceptJson())
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}
