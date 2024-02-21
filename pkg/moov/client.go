package moov

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	pathBankAccounts     = "/accounts/%s/bank-accounts"
	pathMicroDeposits    = "/accounts/%s/bank-accounts/%s/microdeposits"
	pathCards            = "/accounts/%s/cards"
	pathApplePay         = "/accounts/%s/apple-pay"
	pathApplePayDomains  = "/accounts/%s/apple-pay/domains"
	pathApplePaySessions = "/accounts/%s/apple-pay/sessions"
	pathApplePayTokens   = "/accounts/%s/apple-pay/tokens" // #nosec G101
	pathPaymentMethods   = "/accounts/%s/payment-methods"
	pathWallets          = "/accounts/%s/wallets"
	pathWallet           = "/accounts/%s/wallets/%s"
	pathWalletTrans      = "/accounts/%s/wallets/%s/transactions"
	pathWalletTran       = "/accounts/%s/wallets/%s/transactions/%s"
	pathTransfers        = "/transfers"
	pathTransfer         = "/transfers/%s"
	pathRefunds          = "/transfers/%s/refunds"
	pathRefund           = "/transfers/%s/refunds/%s"
	pathReversals        = "/transfers/%s/reversals"
	pathTransferOptions  = "/transfer-options"
	pathDisputes         = "/disputes"
	pathDisputeID        = "/disputes/%s"
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

type Client struct {
	Credentials Credentials
	HttpClient  *http.Client
}

// NewClient returns a moov.Client with credentials read from environment variables.
func NewClient(configurables ...ClientConfigurable) (*Client, error) {
	// Default client configuration if no configurables were specificied
	client := &Client{
		Credentials: CredentialsFromEnv(),
		HttpClient:  DefaultHttpClient(),
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

	return client, nil
}

type ClientConfigurable func(c *Client) error

func WithCredentials(credentials Credentials) ClientConfigurable {
	return func(c *Client) error {
		c.Credentials = credentials
		return c.Credentials.Validate()
	}
}

func WithHttpClient(client *http.Client) ClientConfigurable {
	return func(c *Client) error {
		c.HttpClient = client
		return nil
	}
}
