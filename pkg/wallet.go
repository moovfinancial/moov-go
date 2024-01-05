package moov

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Wallet struct {
	WalletID         string           `json:"walletID,omitempty"`
	AvailableBalance AvailableBalance `json:"availableBalance,omitempty"`
}

type AvailableBalance struct {
	Currency     string `json:"currency,omitempty"`
	Value        int    `json:"value,omitempty"`
	ValueDecimal string `json:"valueDecimal,omitempty"`
}

type Transaction struct {
	WalletID                string    `json:"walletID,omitempty"`
	TransactionID           string    `json:"transactionID,omitempty"`
	TransactionType         string    `json:"transactionType,omitempty"`
	SourceType              string    `json:"sourceType,omitempty"`
	SourceID                string    `json:"sourceID,omitempty"`
	Status                  string    `json:"status,omitempty"`
	Memo                    string    `json:"memo,omitempty"`
	CreatedOn               time.Time `json:"createdOn,omitempty"`
	CompletedOn             time.Time `json:"completedOn,omitempty"`
	Currency                string    `json:"currency,omitempty"`
	GrossAmount             int       `json:"grossAmount,omitempty"`
	GrossAmountDecimal      string    `json:"grossAmountDecimal,omitempty"`
	Fee                     int       `json:"fee,omitempty"`
	FeeDecimal              string    `json:"feeDecimal,omitempty"`
	NetAmount               int       `json:"netAmount,omitempty"`
	NetAmountDecimal        string    `json:"netAmountDecimal,omitempty"`
	AvailableBalance        int       `json:"availableBalance,omitempty"`
	AvailableBalanceDecimal string    `json:"availableBalanceDecimal,omitempty"`
}

// ListWallets lists all wallets that are associated with a Moov account
// https://docs.moov.io/api/index.html#tag/Wallets/operation/listWalletsForAccount
func (c Client) ListWallets(accountID string) ([]Wallet, error) {
	var resWallets []Wallet
	url := fmt.Sprintf("%s/%s", baseURL, fmt.Sprintf(pathWallets, accountID))

	body, statusCode, err := c.GetHTTPResponse(http.MethodGet, url, nil, nil)
	if err != nil {
		return resWallets, err
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &resWallets)
		if err != nil {
			return resWallets, err
		}
		return resWallets, nil
	case http.StatusTooManyRequests:
		return resWallets, ErrRateLimit
	}
	return resWallets, ErrDefault(statusCode)
}

// GetWallet retrieves a wallet for the given wallet id
// https://docs.moov.io/api/index.html#tag/Wallets/operation/getWalletForAccount
func (c Client) GetWallet(accountID string, walletID string) (Wallet, error) {
	resWallet := Wallet{}
	url := fmt.Sprintf("%s/%s/%s", baseURL, fmt.Sprintf(pathWallets, accountID), walletID)

	body, statusCode, err := c.GetHTTPResponse(http.MethodGet, url, nil, nil)
	if err != nil {
		return resWallet, err
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &resWallet)
		if err != nil {
			return resWallet, err
		}
		return resWallet, nil
	case http.StatusUnauthorized:
		return resWallet, ErrAuthCredentialsNotSet
	case http.StatusNotFound:
		return resWallet, ErrNoAccount
	case http.StatusTooManyRequests:
		return resWallet, ErrRateLimit
	}
	return resWallet, ErrDefault(statusCode)
}

type transactionFilter map[string]string

type ListTransactionFilter func(c transactionFilter) error

func applyTransactionFilters(url *url.URL, opts ...ListTransactionFilter) {
	filters := make(transactionFilter)
	for _, opt := range opts {
		opt(filters)
	}
	q := url.Query()
	for k, v := range filters {
		q.Add(k, v)
	}
	url.RawQuery = q.Encode()
}

// WithTransactionType filters transactions by transaction type
func WithTransactionType(transactionType string) ListTransactionFilter {
	return func(c transactionFilter) error {
		c["transactionType"] = transactionType
		return nil
	}
}

// WithSourceType filters transactions by source type (transfer, dispute, issuing-transaction).
func WithSourceType(sourceType string) ListTransactionFilter {
	return func(c transactionFilter) error {
		c["sourceType"] = sourceType
		return nil
	}
}

// WithSourceID filters transactions by source ID
func WithSourceID(sourceID string) ListTransactionFilter {
	return func(c transactionFilter) error {
		c["sourceID"] = sourceID
		return nil
	}
}

// WithTransactionStatus filters transactions by transaction status (pending, completed, canceled, failed)
func WithTransactionStatus(status string) ListTransactionFilter {
	return func(c transactionFilter) error {
		c["status"] = status
		return nil
	}
}

// WithTransactionCount filters transactions by transaction count
func WithTransactionCount(count int) ListTransactionFilter {
	return func(c transactionFilter) error {
		c["count"] = fmt.Sprintf("%d", count)
		return nil
	}
}

// WithTransactionSkip filters transactions by transaction skip
func WithTransactionSkip(skip int) ListTransactionFilter {
	return func(c transactionFilter) error {
		c["skip"] = fmt.Sprintf("%d", skip)
		return nil
	}
}

// WithCreatedStartDateTime filters transactions by created start date time
func WithCreatedStartDateTime(createdStartDateTime time.Time) ListTransactionFilter {
	return func(c transactionFilter) error {
		c["createdStartDateTime"] = createdStartDateTime.Format(time.RFC3339)
		return nil
	}
}

// WithCreatedEndDateTime filters transactions by created end date time
func WithCreatedEndDateTime(createdEndDateTime time.Time) ListTransactionFilter {
	return func(c transactionFilter) error {
		c["createdEndDateTime"] = createdEndDateTime.Format(time.RFC3339)
		return nil
	}
}

// WithCompletedStartDateTime filters transactions by completed start date time
func WithCompletedStartDateTime(completedStartDateTime time.Time) ListTransactionFilter {
	return func(c transactionFilter) error {
		c["completedStartDateTime"] = completedStartDateTime.Format(time.RFC3339)
		return nil
	}
}

// WithCompletedEndDateTime filters transactions by completed end date time
func WithCompletedEndDateTime(completedEndDateTime time.Time) ListTransactionFilter {
	return func(c transactionFilter) error {
		c["completedEndDateTime"] = completedEndDateTime.Format(time.RFC3339)
		return nil
	}
}

// ListWalletTransactions lists all transactions for the given wallet id
// https://docs.moov.io/api/index.html#tag/Wallet-transactions
func (c Client) ListWalletTransactions(accountID string, walletID string, opts ...ListTransactionFilter) ([]Transaction, error) {
	var resTransactions []Transaction
	url := fmt.Sprintf("%s/%s", baseURL, fmt.Sprintf(pathWalletTrans, accountID, walletID))
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.Credentials.PublicKey, c.Credentials.SecretKey)
	applyTransactionFilters(req.URL, opts...)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return resTransactions, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	switch resp.StatusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &resTransactions)
		if err != nil {
			return resTransactions, err
		}
		return resTransactions, nil
	case http.StatusNotFound:
		return resTransactions, ErrNoAccount
	case http.StatusTooManyRequests:
		return resTransactions, ErrRateLimit
	}
	return resTransactions, ErrDefault(resp.StatusCode)
}

// GetWalletTransaction retrieves a transaction for the given wallet id and transaction id
// https://docs.moov.io/api/index.html#tag/Wallet-transactions/operation/getWalletTransaction
func (c Client) GetWalletTransaction(accountID string, walletID string, transactionID string) (Transaction, error) {
	resTransaction := Transaction{}
	url := fmt.Sprintf("%s/%s/%s", baseURL, fmt.Sprintf(pathWalletTrans, accountID, walletID), transactionID)

	body, statusCode, err := c.GetHTTPResponse(http.MethodGet, url, nil, nil)
	if err != nil {
		return resTransaction, err
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &resTransaction)
		if err != nil {
			return resTransaction, err
		}
		return resTransaction, nil
	case http.StatusUnauthorized:
		return resTransaction, ErrAuthCredentialsNotSet
	case http.StatusNotFound:
		return resTransaction, ErrNoAccount
	case http.StatusTooManyRequests:
		return resTransaction, ErrRateLimit
	}
	return resTransaction, ErrDefault(statusCode)
}
