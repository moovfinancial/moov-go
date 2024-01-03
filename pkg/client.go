package moov

import (
	"errors"
	"fmt"
)

const (
	// SCOPE_ACCOUNTS_READ allows a new Moov account to be created
	SCOPE_ACCOUNTS_READ = "/accounts.read"
	// SCOPE_ACCOUNTS_WRITE View all connected accounts
	SCOPE_ACCOUNTS_WRITE = "/accounts.write"

	// TODO: Add the rest of the scopes https://docs.moov.io/guides/developer-tools/scopes/
	SCOPE_BANK_ACCOUNT_READ       = "/accounts/{accountID}/bank-accounts.read"
	SCOPE_BANK_ACCOUNT_WRITE      = "/accounts/{accountID}/bank-accounts.write"
	SCOPE_CARDS_READ              = "/accounts/{accountID}/cards.read"
	SCOPE_CARDS_WRITE             = "/accounts/{accountID}/cards.write"
	SCOPE_CAPABILITIES_READ       = "/accounts/{accountID}/capabilities.read"
	SCOPE_CAPABILITIES_WRITE      = "/accounts/{accountID}/capabilities.write"
	SCOPE_DOCUMENTS_READ          = "/accounts/{accountID}/documents.read"
	SCOPE_DOCUMENTS_WRITE         = "/accounts/{accountID}/documents.write"
	SCOPE_FILES_READ              = "/accounts/{accountID}/files.read"
	SCOPE_FILES_WRITE             = "/accounts/{accountID}/files.write"
	SCOPE_PAYMENT_METHODS_READ    = "/accounts/{accountID}/payment-methods.read"
	SCOPE_PAYMENT_METHODS_WRITE   = "/accounts/{accountID}/payment-methods.write"
	SCOPE_PROFILE_READ            = "/accounts/{accountID}/profile.read"
	SCOPE_PROFILE_WRITE           = "/accounts/{accountID}/profile.write"
	SCOPE_PREFILE_ENRICHMENT_READ = "/accounts/{accountID}/prefile-enrichment.read"
	SCOPE_REPRESENTATIVES_READ    = "/accounts/{accountID}/representatives.read"
	SCOPE_REPRESENTATIVES_WRITE   = "/accounts/{accountID}/representatives.write"
	SCOPE_TRANSFERS_READ          = "/accounts/{accountID}/transfers.read"
	SCOPE_TRANSFERS_WRITE         = "/accounts/{accountID}/transfers.write"
	SCOPE_WALLETS_READ            = "/accounts/{accountID}/wallets.read"
	SCOPE_FED_READ                = "/fed.read"
	SCOPE_PING_READ               = "/ping.read"
)

const (
	baseURL             = "https://api.moov.io"
	pathBankAccounts    = "accounts/%s/bank-accounts"
	pathCards           = "accounts/%s/cards"
	pathApplePay        = "accounts/%s/apple-pay"
	pathPaymentMethods  = "accounts/%s/payment-methods"
	pathWallets         = "accounts/%s/wallets"
	pathWalletTrans     = "accounts/%s/wallets/%s/transactions"
	pathTransactions    = "accounts/%s/transactions"
	pathTransfers       = "transfers"
	pathTransferOptions = "transfer-options"
	pathDisputes        = "disputes"
)

var (
	ErrAuthCredentialsNotSet    = errors.New("API Keys are not set or invalid credentials")
	ErrRequestBody              = errors.New("request body could not be parsed")
	ErrAuthNetwork              = errors.New("network error")
	ErrNoAccount                = errors.New("no account with the specified accountID was found")
	ErrBadRequest               = errors.New("the request body could not be processed")
	ErrInvalidBankAccount       = errors.New("the bank account is not a bank account or is already pending verification")
	ErrAmountIncorrect          = errors.New("the amounts provided are incorrect or the bank account is in an unexpected state")
	ErrDuplicateLinkCard        = errors.New("attempted to link card that already exists on the account")
	ErrCardDataInvalid          = errors.New("the supplied card data appeared invalid or was declined by the issuer")
	ErrDuplicatedApplePayDomain = errors.New("apple Pay domains already registered for this account")
	ErrDomainsNotVerified       = errors.New("domains not verified with Apple")
	ErrDomainsNotRegistered     = errors.New("no Apple Pay domains registered for this account were found")
	ErrLinkingApplePayToken     = errors.New("an error occurred when linking an Apple Pay token")
	ErrRateLimit                = errors.New("request was refused due to rate limiting")
	ErrXIdempotencyKey          = errors.New("attempted to create a transfer using a duplicate X-Idempotency-Key header")
	ErrURL                      = errors.New("invalid URL")
)

func ErrDefault(code int) error {
	return fmt.Errorf("empty response for unauthorized or any other returned http status code (%d)", code)
}

type TransferStatus int

const (
	TransferStatusCreated TransferStatus = iota
	TransferStatusPending
	TransferStatusCompleted
	TransferStatusFailed
	TransferStatusReversed
	TransferStatusQueued
	TransferStatusCanceled
)

var TransferStatusStrings = map[TransferStatus]string{
	TransferStatusCreated:   "created",
	TransferStatusPending:   "pending",
	TransferStatusCompleted: "completed",
	TransferStatusFailed:    "failed",
	TransferStatusReversed:  "reversed",
	TransferStatusQueued:    "queued",
	TransferStatusCanceled:  "canceled",
}

type Client struct {
	Credentials   Credentials
	AccountParams AccountParams
}

type ClientCredentialsGrantToAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int32  `json:"expires_in"`
	Scope        string `json:"scope"`
}

func NewClient(configurables ...ClientConfigurable) (*Client, error) {
	// Default client configuration if no configurables were specificied
	client := &Client{
		Credentials: CredentialsFromEnv(),
	}

	// Apply all the configurable functions to the client
	for _, configurable := range configurables {
		if err := configurable(client); err != nil {
			return nil, err
		}
	}

	// Lets make sure that whatever they passed in for the credentials is valid.
	if err := client.Credentials.Validate(); err != nil {
		return nil, err
	}

	// Ping the server to make sure we have valid credentials
	// @todo Should let them call this if they want incase that makes 1000's of these...
	// if err := client.Ping(); err != nil {
	// 	return client, err
	// }

	return client, nil
}

type ClientConfigurable func(c *Client) error

func WithCredentials(credentials Credentials) ClientConfigurable {
	return func(c *Client) error {
		c.Credentials = credentials
		return c.Credentials.Validate()
	}
}
