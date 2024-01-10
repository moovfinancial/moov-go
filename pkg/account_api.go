package moov

import (
	"context"
	"net/http"
	"strconv"
)

// CreateAccount creates a new account.
func (c Client) CreateAccount(ctx context.Context, account Account) (*Account, *Account, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, "/accounts"),
		AcceptJson(),
		WaitFor("connected"),
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
		return nil, nil, resp.Error()
	}
}

// GetAccount returns an account based on accountID.
func (c Client) GetAccount(ctx context.Context, accountID string) (*Account, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, "/accounts/%s", accountID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Account](resp)
}

// UpdateAccount updates an account.
func (c Client) UpdateAccount(ctx context.Context, account Account) (*Account, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPatch, "/accounts/%s", account.AccountID),
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
		Endpoint(http.MethodGet, "/accounts"),
		prependArgs(opts, AcceptJson())...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[Account](resp)
}

// DeleteAccount deletes an account.
// TODO: Delete is not currently supported by the api
// https://docs.moov.io/guides/dashboard/accounts/#disconnect-accounts
/** func (c Client) DeleteAccount(accountID string) error {
	req, _ := http.NewRequest(http.MethodDelete, "https://api.moov.io/accounts/"+accountID, nil)
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
		// Account created
		return nil
	case http.StatusUnauthorized:
		return ErrAuthCredentialsNotSet
	case http.StatusUnprocessableEntity:
		log.Println("UnprocessableEntity")
	}
	return nil
} **/

/*

// HTTP Client template

	   	params := url.Values{}
	   	params.Add("grant_type", "client_credentials")
	   	params.Add("scope", "/accounts.write")

	   	req, err := http.NewRequest("POST", "https://api.moov.io/oauth2/token?"+params.Encode(), nil)

	   	func (c Client) Ping() {
	   	log.Println("ping")
	   	req, _ := http.NewRequest("POST", "https://api.moov.io/ping", nil)
	   	req.Header.Set("Accept", "application/json")
	   	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	   	//req.SetBasicAuth(c.Credentials.PublicKey, c.Credentials.SecretKey)

	   	client := &http.Client{}
	   	resp, err := client.Do(req)
	   	if err != nil {
	   		// Todo: return an error
	   		log.Fatal(err)
	   	}
	   	defer resp.Body.Close()

	   	body, err := io.ReadAll(resp.Body)
	   	if err != nil {
	   		// Todo: return an error
	   		log.Fatal(err)
	   	}

	   	log.Println("response Status:", resp.Status)
	   	log.Println("response Headers:", resp.Header)
	   	log.Println("response Body:", string(body))

	   }
*/
