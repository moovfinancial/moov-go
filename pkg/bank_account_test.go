package moov

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func randomBankAccountNumber() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(999999999))
	return fmt.Sprintf("%d", 100000000+n.Int64())
}

func TestBankAccountMarshal(t *testing.T) {
	input := []byte(`{
		"bankAccountID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43",
		"fingerprint": "9948962d92a1ce40c9f918cd9ece3a22bde62fb325a2f1fe2e833969de672ba3",
		"status": "new",
		"holderName": "Jules Jackson",
		"holderType": "individual",
		"accountNumber": "0004321567000",
		"bankName": "Chase Bank",
		"bankAccountType": "checking",
		"routingNumber": "string",
		"lastFourAccountNumber": "7000"
	  }`)

	bankAccount := new(BankAccount)

	dec := json.NewDecoder(bytes.NewReader(input))
	dec.DisallowUnknownFields()

	err := dec.Decode(&bankAccount)
	if err != nil {
		require.NoError(t, err)
	}

	require.Equal(t, "Chase Bank", bankAccount.BankName)
}

type BankAccountTestSuite struct {
	suite.Suite
	// values fort testing will be set in init()
	accountID           string
	bankAccountIDDelete string
	bankAccountID       string
	bankAccounts        []string
}

// listen for 'go test' command --> run test methods
func TestBankAccountSuite(t *testing.T) {
	suite.Run(t, new(BankAccountTestSuite))
}

func (s *BankAccountTestSuite) SetupSuite() {
	// Sandbox accounts have a "Lincoln National Corporation" moov account added by default. Get it's AccountID so we can test against it
	mc := NewTestClient(s.T())

	accounts, err := mc.ListAccounts(WithAccountName("Lincoln National Corporation"))
	s.NoError(err)

	for _, account := range accounts {
		if account.DisaplayName == "Lincoln National Corporation" {
			// set the accountID for testing
			s.accountID = account.AccountID
		}
	}

	s.Assert().NotEmpty(s.accountID)

	// create a bank account for Lincoln National Corporation
	bankAccount := BankAccount{
		HolderName:      "Sir Test ALot",
		HolderType:      "individual",
		BankAccountType: "checking",
		AccountNumber:   randomBankAccountNumber(),
		RoutingNumber:   "273976369",
	}

	bankAccount, err = mc.CreateBankAccount(s.accountID, bankAccount)
	s.NoError(err)

	s.bankAccountID = bankAccount.BankAccountID
	s.bankAccounts = append(s.bankAccounts, bankAccount.BankAccountID)

	bankAccountDelete := BankAccount{
		HolderName:      "Sir Test Delete ALot",
		HolderType:      "individual",
		BankAccountType: "checking",
		AccountNumber:   randomBankAccountNumber(),
		RoutingNumber:   "273976369",
	}

	bankAccount, err = mc.CreateBankAccount(s.accountID, bankAccountDelete)
	s.NoError(err)

	s.bankAccountIDDelete = bankAccount.BankAccountID
	s.bankAccounts = append(s.bankAccounts, bankAccount.BankAccountID)
}

func (s *BankAccountTestSuite) TearDownSuite() {
	mc := NewTestClient(s.T())

	// delete the bank accounts we created
	for _, bankAccountID := range s.bankAccounts {
		if bankAccountID != "" {
			err := mc.DeleteBankAccount(s.accountID, bankAccountID)
			s.NoError(err)
		}
	}
}

func (s *BankAccountTestSuite) TestCreateBankAccount() {
	bankAccount := BankAccount{
		HolderName:      "Jules Jackson",
		HolderType:      "individual",
		BankAccountType: "checking",
		AccountNumber:   randomBankAccountNumber(),
		RoutingNumber:   "273976369",
	}

	mc := NewTestClient(s.T())

	bankAccount, err := mc.CreateBankAccount(s.accountID, bankAccount)
	s.NoError(err)

	s.NotEmpty(bankAccount.BankAccountID)
	s.bankAccounts = append(s.bankAccounts, bankAccount.BankAccountID)
}

func (s *BankAccountTestSuite) TestGetBankAccount() {
	mc := NewTestClient(s.T())

	account, err := mc.GetBankAccount(s.accountID, s.bankAccountID)
	s.NoError(err)

	s.Equal(s.bankAccountID, account.BankAccountID)
}

func (s *BankAccountTestSuite) TestDeleteBankAccount() {
	mc := NewTestClient(s.T())

	err := mc.DeleteBankAccount(s.accountID, s.bankAccountIDDelete)
	s.NoError(err)
}

func (s *BankAccountTestSuite) TestListBankAccounts() {
	mc := NewTestClient(s.T())

	accounts, err := mc.ListBankAccounts(s.accountID)
	s.NoError(err)

	s.NotNil(accounts)
}

func (s *BankAccountTestSuite) TestMicroDepositInitiate() {
	mc := NewTestClient(s.T())

	err := mc.MicroDepositInitiate(s.accountID, s.bankAccountID)
	s.NoError(err)
}

// TODO: test this could run before TestMicroDepositInitiate
func (s *BankAccountTestSuite) TestMicroDepositConfirm() {
	mc := NewTestClient(s.T())

	err := mc.MicroDepositInitiate(s.accountID, s.bankAccountID)
	s.NoError(err)
	// sample data
	amounts := []int{0, 0}
	err = mc.MicroDepositConfirm(s.accountID, s.bankAccountID, amounts)
	s.NoError(err)
}
