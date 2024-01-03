package moov

import (
	"bytes"
	"encoding/json"
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

	accounts, err := mc.ListAccounts(WithAccountName("Lincoln National Corporation"))
	s.NoError(err)

	defaultAccountName := "Lincoln National Corporation"
	for _, account := range accounts {
		if account.DisaplayName == defaultAccountName {
			// set the accountID for testing
			s.accountID = account.AccountID
		}
	}
}

func (s *WalletTestSuite) TearDownSuite() {}

func (s *WalletTestSuite) TestListWallets() {
	mc := NewTestClient(s.T())

	wallets, err := mc.ListWallets(s.accountID)
	s.NoError(err)

	s.Require().NotEmpty(wallets)

	s.walletID = wallets[0].WalletID
}

func (s *WalletTestSuite) TestGetWallet() {
	mc := NewTestClient(s.T())

	walletID := s.walletID
	if walletID == "" {
		walletID = "3097f356-f763-4b24-b282-b53b9bb644ef"
	}

	wallet, err := mc.GetWallet(s.accountID, walletID)
	s.NoError(err)

	s.Equal(walletID, wallet.WalletID)
}

func (s *WalletTestSuite) TestListWalletTransactions() {
	mc := NewTestClient(s.T())

	walletID := s.walletID
	if walletID == "" {
		walletID = "3097f356-f763-4b24-b282-b53b9bb644ef"
	}

	walletTrans, err := mc.ListWalletTransactions(s.accountID, walletID)
	s.NoError(err)

	s.NotNil(walletTrans)
}

func (s *WalletTestSuite) TestGetWalletTransaction() {
	mc := NewTestClient(s.T())

	walletID := s.walletID
	if walletID == "" {
		walletID = "3097f356-f763-4b24-b282-b53b9bb644ef"
	}

	walletTransactionID := s.walletTransactionID
	if walletTransactionID == "" {
		walletTransactionID = "3097f356-f763-4b24-b282-b53b9bb644ef"
	}

	walletTran, err := mc.GetWalletTransaction(s.accountID, walletID, walletTransactionID)
	s.NoError(err)

	s.Equal(walletTransactionID, walletTran.TransactionID)
}
