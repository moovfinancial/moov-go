package moov

import (
	"context"
	"net/http"
	"time"
)

type Wallet struct {
	WalletID         string           `json:"walletID,omitempty"`
	AvailableBalance AvailableBalance `json:"availableBalance,omitempty"`
}

type AvailableBalance struct {
	// A 3-letter ISO 4217 currency code.
	Currency string `json:"currency,omitempty" otel:"currency"`
	// Quantity in the smallest unit of the specified currency. In USD this is cents, so $12.04 is 1204 and $0.99 would be 99.
	Value int64 `json:"value,omitempty" otel:"value"`
	// Same as `value`, but a decimal-formatted numerical string that represents up to 9 decimal place precision.
	ValueDecimal string `json:"valueDecimal,omitempty"`
}

type Transaction struct {
	WalletID                string                  `json:"walletID,omitempty"`
	TransactionID           string                  `json:"transactionID,omitempty"`
	TransactionType         string                  `json:"transactionType,omitempty"`
	SourceType              string                  `json:"sourceType,omitempty"`
	SourceID                string                  `json:"sourceID,omitempty"`
	Status                  WalletTransactionStatus `json:"status,omitempty"`
	Memo                    string                  `json:"memo,omitempty"`
	CreatedOn               time.Time               `json:"createdOn,omitempty"`
	CompletedOn             time.Time               `json:"completedOn,omitempty"`
	Currency                string                  `json:"currency,omitempty"`
	GrossAmount             int                     `json:"grossAmount,omitempty"`
	GrossAmountDecimal      string                  `json:"grossAmountDecimal,omitempty"`
	Fee                     int                     `json:"fee,omitempty"`
	FeeDecimal              string                  `json:"feeDecimal,omitempty"`
	NetAmount               int                     `json:"netAmount,omitempty"`
	NetAmountDecimal        string                  `json:"netAmountDecimal,omitempty"`
	AvailableBalance        int                     `json:"availableBalance,omitempty"`
	AvailableBalanceDecimal string                  `json:"availableBalanceDecimal,omitempty"`
}

type WalletTransactionStatus string

const (
	WalletTransactionStatus_Pending   WalletTransactionStatus = "pending"
	WalletTransactionStatus_Completed WalletTransactionStatus = "completed"
	WalletTransactionStatus_Canceled  WalletTransactionStatus = "canceled"
	WalletTransactionStatus_Failed    WalletTransactionStatus = "failed"
)

// ListWallets lists all wallets that are associated with a Moov account
// https://docs.moov.io/api/index.html#tag/Wallets/operation/listWalletsForAccount
func (c Client) ListWallets(ctx context.Context, accountID string) ([]Wallet, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathWallets, accountID), AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[Wallet](resp)
}

// GetWallet retrieves a wallet for the given wallet id
// https://docs.moov.io/api/index.html#tag/Wallets/operation/getWalletForAccount
func (c Client) GetWallet(ctx context.Context, accountID string, walletID string) (*Wallet, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathWallet, accountID, walletID), AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Wallet](resp)
}

type ListTransactionFilter callArg

// WithTransactionType filters transactions by transaction type
func WithTransactionType(transactionType string) ListTransactionFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["transactionType"] = transactionType
		return nil
	})
}

// WithSourceType filters transactions by source type (transfer, dispute, issuing-transaction).
func WithTransactionSourceType(sourceType string) ListTransactionFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["sourceType"] = sourceType
		return nil
	})
}

// WithSourceID filters transactions by source ID
func WithTransactionSourceID(sourceID string) ListTransactionFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["sourceID"] = sourceID
		return nil
	})
}

// WithTransactionStatus filters transactions by transaction status (pending, completed, canceled, failed)
func WithTransactionStatus(status string) ListTransactionFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["status"] = status
		return nil
	})
}

// WithTransactionCount filters transactions by transaction count
func WithTransactionCount(count int) ListTransactionFilter {
	return Count(count)
}

// WithTransactionSkip filters transactions by transaction skip
func WithTransactionSkip(skip int) ListTransactionFilter {
	return Skip(skip)
}

// WithCreatedStartDateTime filters transactions by created start date time
func WithCreatedStartDateTime(createdStartDateTime time.Time) ListTransactionFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["createdStartDateTime"] = createdStartDateTime.Format(time.RFC3339)
		return nil
	})
}

// WithCreatedEndDateTime filters transactions by created end date time
func WithCreatedEndDateTime(createdEndDateTime time.Time) ListTransactionFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["createdEndDateTime"] = createdEndDateTime.Format(time.RFC3339)
		return nil
	})
}

// WithCompletedStartDateTime filters transactions by completed start date time
func WithCompletedStartDateTime(completedStartDateTime time.Time) ListTransactionFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["completedStartDateTime"] = completedStartDateTime.Format(time.RFC3339)
		return nil
	})
}

// WithCompletedEndDateTime filters transactions by completed end date time
func WithCompletedEndDateTime(completedEndDateTime time.Time) ListTransactionFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["completedEndDateTime"] = completedEndDateTime.Format(time.RFC3339)
		return nil
	})
}

// ListWalletTransactions lists all transactions for the given wallet id
// https://docs.moov.io/api/index.html#tag/Wallet-transactions
func (c Client) ListWalletTransactions(ctx context.Context, accountID string, walletID string, opts ...ListTransactionFilter) ([]Transaction, error) {
	args := prependArgs(opts, AcceptJson())
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathWalletTransactions, accountID, walletID), args...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[Transaction](resp)
}

// GetWalletTransaction retrieves a transaction for the given wallet id and transaction id
// https://docs.moov.io/api/index.html#tag/Wallet-transactions/operation/getWalletTransaction
func (c Client) GetWalletTransaction(ctx context.Context, accountID string, walletID string, transactionID string) (*Transaction, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathWalletTransaction, accountID, walletID, transactionID), AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Transaction](resp)
}
