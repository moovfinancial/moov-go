package moov

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestWalletMarshal(t *testing.T) {
	input := []byte(`{
			"walletID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43",
			"availableBalance": {
				"currency": "USD",
				"value": 1204,
				"valueDecimal": "12.987654321"
			}
		}`)

	wallet := new(Wallet)

	dec := json.NewDecoder(bytes.NewReader(input))
	dec.DisallowUnknownFields()

	err := dec.Decode(&wallet)
	require.NoError(t, err)

	require.Equal(t, "ec7e1848-dc80-4ab0-8827-dd7fc0737b43", wallet.WalletID)
}

type WalletTestSuite struct {
	suite.Suite
	// values for testing will be set in init()
	accountID           string
	walletID            string
	walletTransactionID string
}

// listen for 'go test' command --> run test methods
func TestWalletSuite(t *testing.T) {
	suite.Run(t, new(WalletTestSuite))
}

func (s *WalletTestSuite) SetupSuite() {
	// Sandbox accounts have a "Lincoln National Corporation" moov account added by default. Get it's AccountID so we can test against it
	mc := NewTestClient(s.T())

	accounts, err := mc.ListAccounts(context.Background(), WithAccountName("Lincoln National Corporation"))
	s.NoError(err)

	defaultAccountName := "Lincoln National Corporation"
	for _, account := range accounts {
		if account.DisaplayName == defaultAccountName {
			// set the accountID for testing
			s.accountID = account.AccountID
		}
	}
	wallets, err := mc.ListWallets(s.accountID)
	s.NoError(err)

	for _, wallet := range wallets {
		s.walletID = wallet.WalletID
	}
	s.Assert().NotEmpty(s.accountID)

	transactions, err := mc.ListWalletTransactions(s.accountID, s.walletID, WithTransactionCount(1))
	s.NoError(err)
	for _, transaction := range transactions {
		s.walletTransactionID = transaction.TransactionID
	}
	s.Assert().NotEmpty(s.walletTransactionID)
}

func (s *WalletTestSuite) TearDownSuite() {}

func (s *WalletTestSuite) TestListWallets() {
	mc := NewTestClient(s.T())

	wallets, err := mc.ListWallets(s.accountID)
	s.NoError(err)
	// range over wallets and print walletID
	for _, wallet := range wallets {
		log.Printf("Wallet: %s", wallet.WalletID)
	}
	s.Require().NotEmpty(wallets)
}

func (s *WalletTestSuite) TestGetWallet() {
	mc := NewTestClient(s.T())
	wallet, err := mc.GetWallet(s.accountID, s.walletID)
	s.NoError(err)
	s.Equal(s.walletID, wallet.WalletID)
}

func (s *WalletTestSuite) TestListWalletTransactions() {
	mc := NewTestClient(s.T())
	walletTransactions, err := mc.ListWalletTransactions(s.accountID, s.walletID, WithTransactionStatus("completed"), WithTransactionCount(50))
	s.NoError(err)
	s.NotNil(walletTransactions)
	s.Greater(len(walletTransactions), 3)
}

func (s *WalletTestSuite) TestGetWalletTransaction() {
	mc := NewTestClient(s.T())
	walletTran, err := mc.GetWalletTransaction(s.accountID, s.walletID, s.walletTransactionID)
	s.NoError(err)
	s.Equal(s.walletTransactionID, walletTran.TransactionID)
}
