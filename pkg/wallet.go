package moov

import "time"

type Wallet struct {
	WalletID string `json:"walletID,omitempty"`
	AvailableBalance
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

// GetWallet retrieves a wallet for the given wallet id
// https://docs.moov.io/api/index.html#tag/Wallets/operation/getWalletForAccount

// ListWalletTransactions lists all transactions for the given wallet id
// https://docs.moov.io/api/index.html#tag/Wallet-transactions

// GetWalletTransaction retrieves a transaction for the given wallet id and transaction id
// https://docs.moov.io/api/index.html#tag/Wallet-transactions/operation/getWalletTransaction
