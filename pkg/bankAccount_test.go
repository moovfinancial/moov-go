package moov

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func randomBankAccountNumber() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d", 100000000+rand.Intn(999999999))
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
	assert.Equal(t, "Chase Bank", bankAccount.BankName)
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
	mc, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}
	accounts, err := mc.ListAccounts()
	if err != nil {
		log.Fatal(err)
	}
	for _, account := range accounts {
		if account.DisaplayName == "Lincoln National Corporation" {
			// set the accountID for testing
			s.accountID = account.AccountID
		}
	}

	// create a bank account for Lincoln National Corporation
	bankAccount := BankAccount{
		HolderName:      "Sir Test ALot",
		HolderType:      "individual",
		BankAccountType: "checking",
		AccountNumber:   randomBankAccountNumber(),
		RoutingNumber:   "273976369",
	}

	bankAccount, err = mc.CreateBankAccount(s.accountID, bankAccount)
	if err != nil {
		log.Fatal(err)
	}
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
	if err != nil {
		log.Fatal(err)
	}
	s.bankAccountIDDelete = bankAccount.BankAccountID
	s.bankAccounts = append(s.bankAccounts, bankAccount.BankAccountID)
}

func (s *BankAccountTestSuite) TearDownSuite() {
	mc, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}
	// delete the bank accounts we created
	for _, bankAccountID := range s.bankAccounts {
		if bankAccountID != "" {
			err = mc.DeleteBankAccount(s.accountID, bankAccountID)
			if err != nil {
				log.Fatal(err)
			}
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

	mc, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}

	bankAccount, err = mc.CreateBankAccount(s.accountID, bankAccount)
	if err != nil {
		s.T().Fatalf("Error creating bank account: %v", err)
	}

	assert.NotEmpty(s.T(), bankAccount.BankAccountID)

	s.bankAccounts = append(s.bankAccounts, bankAccount.BankAccountID)
}

func (s *BankAccountTestSuite) TestGetBankAccount() {
	mc, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}
	account, err := mc.GetBankAccount(s.accountID, s.bankAccountID)
	if err != nil {
		s.T().Fatal(err)
	}
	assert.Equal(s.T(), s.bankAccountID, account.BankAccountID)
}

func (s *BankAccountTestSuite) TestDeleteBankAccount() {
	mc, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}
	err = mc.DeleteBankAccount(s.accountID, s.bankAccountIDDelete)
	if err != nil {
		assert.Error(s.T(), err)
	}
}

func (s *BankAccountTestSuite) TestListBankAccounts() {
	mc, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}

	accounts, err := mc.ListBankAccounts(s.accountID)
	if err != nil {
		s.T().Fatal(err)
	}
	fmt.Println(len(accounts))
	assert.NotNil(s.T(), accounts)
}

func (s *BankAccountTestSuite) TestMicroDepositInitiate() {
	mc, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}
	err = mc.MicroDepositInitiate(s.accountID, s.bankAccountID)
	if err != nil {
		assert.Error(s.T(), err)
	}
}

// TODO: test this could run before TestMicroDepositInitiate
func (s *BankAccountTestSuite) TestMicroDepositConfirm() {
	mc, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}
	// sample data
	amounts := []int{22, 21}
	err = mc.MicroDepositConfirm(s.accountID, s.bankAccountID, amounts)
	if err != nil {
		assert.Error(s.T(), err)
	}
}
