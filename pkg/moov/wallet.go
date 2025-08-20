package moov

import (
	"context"
	"net/http"
	"time"
)

// TODO(vince,08/13/2025): move all of these models to a separate file to match pattern for other domains
type Wallet struct {
	WalletID         string           `json:"walletID"`
	AvailableBalance AvailableBalance `json:"availableBalance"`

	PartnerAccountID string       `json:"partnerAccountID"`
	Name             string       `json:"name"`
	Status           WalletStatus `json:"status"`
	WalletType       WalletType   `json:"walletType"`
	CreatedOn        time.Time    `json:"createdOn"`

	Metadata    map[string]string `json:"metadata,omitempty"`
	Description *string           `json:"description,omitempty"`
	ClosedOn    *time.Time        `json:"closedOn,omitempty"`
}

type AvailableBalance struct {
	// A 3-letter ISO 4217 currency code.
	Currency string `json:"currency,omitempty" otel:"currency"`
	// Quantity in the smallest unit of the specified currency. In USD this is cents, so $12.04 is 1204 and $0.99 would be 99.
	Value int64 `json:"value" otel:"value"`
	// Same as `value`, but a decimal-formatted numerical string that represents up to 9 decimal place precision.
	ValueDecimal string `json:"valueDecimal"`
}

type WalletTransaction struct {
	WalletID                string                      `json:"walletID,omitempty"`
	TransactionID           string                      `json:"transactionID,omitempty"`
	TransactionType         WalletTransactionType       `json:"transactionType,omitempty"`
	SourceType              WalletTransactionSourceType `json:"sourceType,omitempty"`
	SourceID                string                      `json:"sourceID,omitempty"`
	Status                  WalletTransactionStatus     `json:"status,omitempty"`
	Memo                    string                      `json:"memo,omitempty"`
	CreatedOn               time.Time                   `json:"createdOn,omitempty"`
	CompletedOn             time.Time                   `json:"completedOn,omitempty"`
	Currency                string                      `json:"currency,omitempty"`
	GrossAmount             int                         `json:"grossAmount,omitempty"`
	GrossAmountDecimal      string                      `json:"grossAmountDecimal,omitempty"`
	Fee                     int                         `json:"fee,omitempty"`
	FeeDecimal              string                      `json:"feeDecimal,omitempty"`
	NetAmount               int                         `json:"netAmount,omitempty"`
	NetAmountDecimal        string                      `json:"netAmountDecimal,omitempty"`
	AvailableBalance        int                         `json:"availableBalance,omitempty"`
	AvailableBalanceDecimal string                      `json:"availableBalanceDecimal,omitempty"`
	SweepID                 *string                     `json:"sweepID,omitempty"`
}

type WalletStatus string

const (
	WalletStatus_Active WalletStatus = "active"
	WalletStatus_Closed WalletStatus = "closed"
)

type WalletType string

const (
	WalletType_Default WalletType = "default"
	WalletType_General WalletType = "general"
)

type WalletTransactionStatus string

const (
	WalletTransactionStatus_Pending   WalletTransactionStatus = "pending"
	WalletTransactionStatus_Completed WalletTransactionStatus = "completed"
	WalletTransactionStatus_Canceled  WalletTransactionStatus = "canceled"
	WalletTransactionStatus_Failed    WalletTransactionStatus = "failed"
)

type WalletTransactionType string

const (
	WalletTransactionTypeAccountFunding               WalletTransactionType = "account-funding"
	WalletTransactionTypeAchReversal                  WalletTransactionType = "ach-reversal"
	WalletTransactionTypeAutoSweep                    WalletTransactionType = "auto-sweep"
	WalletTransactionTypeCardPayment                  WalletTransactionType = "card-payment"
	WalletTransactionTypeCardDecline                  WalletTransactionType = "card-decline"
	WalletTransactionTypeCardReversal                 WalletTransactionType = "card-reversal"
	WalletTransactionTypeCashOut                      WalletTransactionType = "cash-out"
	WalletTransactionTypeDispute                      WalletTransactionType = "dispute"
	WalletTransactionTypeDisputeReversal              WalletTransactionType = "dispute-reversal"
	WalletTransactionTypeFacilitatorFee               WalletTransactionType = "facilitator-fee"
	WalletTransactionTypeIssuingRefund                WalletTransactionType = "issuing-refund"
	WalletTransactionTypeIssuingTransaction           WalletTransactionType = "issuing-transaction"
	WalletTransactionTypeIssuingTransactionAdjustment WalletTransactionType = "issuing-transaction-adjustment"
	WalletTransactionTypeIssuingAuthHold              WalletTransactionType = "issuing-auth-hold"
	WalletTransactionTypeIssuingAuthRelease           WalletTransactionType = "issuing-auth-release"
	WalletTransactionTypeIssuingDecline               WalletTransactionType = "issuing-decline"
	WalletTransactionTypeMoovFee                      WalletTransactionType = "moov-fee"
	WalletTransactionTypePayment                      WalletTransactionType = "payment"
	WalletTransactionTypePayout                       WalletTransactionType = "payout"
	WalletTransactionTypeRefund                       WalletTransactionType = "refund"
	WalletTransactionTypeRefundFailure                WalletTransactionType = "refund-failure"
	WalletTransactionTypeRtpFailure                   WalletTransactionType = "rtp-failure"
	WalletTransactionTypeTopUp                        WalletTransactionType = "top-up"
	WalletTransactionTypeWalletTransfer               WalletTransactionType = "wallet-transfer"
)

type WalletTransactionSourceType string

const (
	WalletTransactionSourceTypeTransfer               WalletTransactionSourceType = "transfer"
	WalletTransactionSourceTypeDispute                WalletTransactionSourceType = "dispute"
	WalletTransactionSourceTypeIssuingCardTransaction WalletTransactionSourceType = "issuing-card-transaction"
	WalletTransactionSourceTypeIssuingAuthorization   WalletTransactionSourceType = "issuing-authorization"
	WalletTransactionSourceTypeSweep                  WalletTransactionSourceType = "sweep"
	WalletTransactionSourceTypeFee                    WalletTransactionSourceType = "fee"
)

// TODO(vince,08/13/2025): update this with new query params
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

type CreateWallet struct {
	Name string `json:"name"`
	// TODO(vince,08/13/2025): we may decide to remove this as a required field, fix before merging
	PartnerAccountID string            `json:"partnerAccountID"`
	Description      string            `json:"description"`
	Metadata         map[string]string `json:"metadata"`
}

// TODO(vince,08/13/2025): add docs link
// CreateWallet creates a general wallet
func (c Client) CreateWallet(ctx context.Context, accountID string, create CreateWallet) (*Wallet, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodPost, pathWallets, accountID), AcceptJson(), JsonBody(create))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Wallet](resp)
}

type UpdateWallet struct {
	Name        *string       `json:"name,omitempty"`
	Status      *WalletStatus `json:"status,omitempty"`
	Description *string       `json:"description,omitempty"`
	// Free-form key-value pair list. Useful for storing information that is not captured elsewhere.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// TODO(vince,08/13/2025): add docs link
//
// UpdateWallet updates a wallet
func (c Client) UpdateWallet(ctx context.Context, accountID string, walletID string, update UpdateWallet) (*Wallet, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodPatch, pathWallet, accountID, walletID), AcceptJson(), JsonBody(update))
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

// WithTransactionTypes filters transactions by transaction types. The parameter `transactionTypes` is a comma-separated parameter to filter by transaction types.
func WithTransactionTypes(transactionTypes string) ListTransactionFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["transactionTypes"] = transactionTypes
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

// WithSweepID filters for transactions accrued in a sweep
func WithSweepID(sweepID string) ListTransactionFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["sweepID"] = sweepID
		return nil
	})
}

// ListWalletTransactions lists all transactions for the given wallet id
// https://docs.moov.io/api/index.html#tag/Wallet-transactions
func (c Client) ListWalletTransactions(ctx context.Context, accountID string, walletID string, opts ...ListTransactionFilter) ([]WalletTransaction, error) {
	args := prependArgs(opts, AcceptJson())
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathWalletTransactions, accountID, walletID), args...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[WalletTransaction](resp)
}

// GetWalletTransaction retrieves a transaction for the given wallet id and transaction id
// https://docs.moov.io/api/index.html#tag/Wallet-transactions/operation/getWalletTransaction
func (c Client) GetWalletTransaction(ctx context.Context, accountID string, walletID string, transactionID string) (*WalletTransaction, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathWalletTransaction, accountID, walletID, transactionID), AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[WalletTransaction](resp)
}
