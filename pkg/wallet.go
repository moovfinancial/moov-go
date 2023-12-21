package moov

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	body, statusCode, err := GetHTTPResponse(c, http.MethodGet, url, nil, nil)
	if err != nil {
		return resWallets, err
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &resWallets)
		if err != nil {
			log.Println("Error unmarshalling JSON:", err)
		}
		return resWallets, nil
	case http.StatusTooManyRequests:
		return resWallets, ErrRateLimit
	}
	return resWallets, ErrDefault
}

// GetWallet retrieves a wallet for the given wallet id
// https://docs.moov.io/api/index.html#tag/Wallets/operation/getWalletForAccount
func (c Client) GetWallet(accountID string, walletID string) (Wallet, error) {
	resWallet := Wallet{}
	url := fmt.Sprintf("%s/%s/%s", baseURL, fmt.Sprintf(pathWallets, accountID), walletID)

	body, statusCode, err := GetHTTPResponse(c, http.MethodGet, url, nil, nil)
	if err != nil {
		return resWallet, err
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &resWallet)
		if err != nil {
			log.Println("Error unmarshalling JSON:", err)
		}
		return resWallet, nil
	case http.StatusUnauthorized:
		return resWallet, ErrAuthCredentialsNotSet
	case http.StatusNotFound:
		return resWallet, ErrNoAccount
	case http.StatusTooManyRequests:
		return resWallet, ErrRateLimit
	}
	return resWallet, ErrDefault
}

// ListWalletTransactions lists all transactions for the given wallet id
// https://docs.moov.io/api/index.html#tag/Wallet-transactions
func (c Client) ListWalletTransactions(accountID string, walletID string) ([]Transaction, error) {
	var resTransactions []Transaction
	url := fmt.Sprintf("%s/%s", baseURL, fmt.Sprintf(pathWalletTrans, accountID, walletID))

	body, statusCode, err := GetHTTPResponse(c, http.MethodGet, url, nil, nil)
	if err != nil {
		return resTransactions, err
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &resTransactions)
		if err != nil {
			log.Println("Error unmarshalling JSON:", err)
		}
		return resTransactions, nil
	case http.StatusNotFound:
		return resTransactions, ErrNoAccount
	case http.StatusTooManyRequests:
		return resTransactions, ErrRateLimit
	}
	return resTransactions, ErrDefault
}

// GetWalletTransaction retrieves a transaction for the given wallet id and transaction id
// https://docs.moov.io/api/index.html#tag/Wallet-transactions/operation/getWalletTransaction
func (c Client) GetWalletTransaction(accountID string, walletID string, transactionID string) (Transaction, error) {
	resTransaction := Transaction{}
	url := fmt.Sprintf("%s/%s/%s", baseURL, fmt.Sprintf(pathWalletTrans, accountID, walletID), transactionID)

	body, statusCode, err := GetHTTPResponse(c, http.MethodGet, url, nil, nil)
	if err != nil {
		return resTransaction, err
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &resTransaction)
		if err != nil {
			log.Println("Error unmarshalling JSON:", err)
		}
		return resTransaction, nil
	case http.StatusUnauthorized:
		return resTransaction, ErrAuthCredentialsNotSet
	case http.StatusNotFound:
		return resTransaction, ErrNoAccount
	case http.StatusTooManyRequests:
		return resTransaction, ErrRateLimit
	}
	return resTransaction, ErrDefault
}
