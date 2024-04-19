package moov_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

func getLincolnBank(t *testing.T, mc *moov.Client) moov.Account {

	accounts, err := mc.ListAccounts(context.Background(), moov.WithAccountName("Lincoln National Corporation"))
	moov.DebugPrintResponse(err, fmt.Printf)
	require.NoError(t, err)

	for _, account := range accounts {
		if account.DisplayName == "Lincoln National Corporation" {
			return account
		}
	}

	require.FailNow(t, "bank account test account not found")
	return moov.Account{}
}

func Test_CreateBankAccount_WithBankAccount(t *testing.T) {
	mc := NewTestClient(t)

	account := getLincolnBank(t, mc)

	_, err := mc.CreateBankAccount(BgCtx(), account.AccountID, moov.WithBankAccount(moov.BankAccountRequest{
		HolderName:    "Sir Test Delete ALot",
		HolderType:    moov.HolderType_Individual,
		AccountType:   moov.BankAccountType_Checking,
		AccountNumber: randomBankAccountNumber(),
		RoutingNumber: "273976369",
	}))

	moov.DebugPrintResponse(err, fmt.Printf)
	require.NoError(t, err)
}

func randomBankAccountNumber() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(999999999))
	return fmt.Sprintf("%d", 100000000+n.Int64())
}

/*
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

	bankAccount := new(moov.BankAccount)

	dec := json.NewDecoder(bytes.NewReader(input))
	dec.DisallowUnknownFields()

	err := dec.Decode(&bankAccount)
	if err != nil {
		require.NoError(t, err)
	}

	require.Equal(t, "Chase Bank", bankAccount.BankName)
}
*/
/*
@TODO fix by getting rid of the suite

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

	accounts, err := mc.ListAccounts(context.Background(), moov.WithAccountName("Lincoln National Corporation"))
	s.NoError(err)

	for _, account := range accounts {
		if account.DisplayName == "Lincoln National Corporation" {
			// set the accountID for testing
			s.accountID = account.AccountID
		}
	}

	s.Assert().NotEmpty(s.accountID)

	// create a bank account for Lincoln National Corporation
	bankAccount := moov.BankAccount{
		HolderName:      "Sir Test ALot",
		HolderType:      "individual",
		BankAccountType: "checking",
		AccountNumber:   randomBankAccountNumber(),
		RoutingNumber:   "273976369",
	}

	result, err := mc.CreateBankAccount(context.Background(), s.accountID, moov.WithBankAccount(bankAccount))
	s.NoError(err)
	s.NotNil(result)
	bankAccount = *result

	s.bankAccountID = bankAccount.BankAccountID
	s.bankAccounts = append(s.bankAccounts, bankAccount.BankAccountID)

	bankAccountDelete := moov.BankAccount{
		HolderName:      "Sir Test Delete ALot",
		HolderType:      "individual",
		BankAccountType: "checking",
		AccountNumber:   randomBankAccountNumber(),
		RoutingNumber:   "273976369",
	}

	result, err = mc.CreateBankAccount(context.Background(), s.accountID, moov.WithBankAccount(bankAccountDelete))
	s.NoError(err)
	s.NotNil(result)
	bankAccount = *result

	s.bankAccountIDDelete = bankAccount.BankAccountID
	s.bankAccounts = append(s.bankAccounts, bankAccount.BankAccountID)
}

func (s *BankAccountTestSuite) TearDownSuite() {
	mc := NewTestClient(s.T())

	// delete the bank accounts we created
	for _, bankAccountID := range s.bankAccounts {
		if bankAccountID != "" {
			err := mc.DeleteBankAccount(context.Background(), s.accountID, bankAccountID)
			s.NoError(err)
		}
	}
}

func (s *BankAccountTestSuite) TestCreateBankAccount() {
	bankAccount := moov.BankAccount{
		HolderName:      "Jules Jackson",
		HolderType:      "individual",
		BankAccountType: "checking",
		AccountNumber:   randomBankAccountNumber(),
		RoutingNumber:   "273976369",
	}

	mc := NewTestClient(s.T())

	result, err := mc.CreateBankAccount(context.Background(), s.accountID, moov.WithBankAccount(bankAccount))
	s.NoError(err)
	s.NotNil(result)

	s.NotEmpty(bankAccount.BankAccountID)
	s.bankAccounts = append(s.bankAccounts, bankAccount.BankAccountID)
}

func (s *BankAccountTestSuite) TestCreateBankAccountWithPlaid() {
	mc := NewTestClient(s.T())
	plaid := moov.Plaid{
		Token: "fake-token",
	}
	result, err := mc.CreateBankAccount(context.Background(), s.accountID, moov.WithPlaid(plaid))
	s.NoError(err)
	s.NotNil(result)
	s.Equal("1111222233330000", result.AccountNumber)
}

func (s *BankAccountTestSuite) TestCreateBankAccountWithPlaidLink() {
	mc := NewTestClient(s.T())
	plaidLink := moov.PlaidLink{
		PublicToken: "fake-public-token",
	}
	result, err := mc.CreateBankAccount(context.Background(), s.accountID, moov.WithPlaidLink(plaidLink))
	s.NoError(err)
	s.NotNil(result)
	s.Equal("1111222233330000", result.AccountNumber)
}

func (s *BankAccountTestSuite) TestCreateBankAccountWithMX() {
	mc := NewTestClient(s.T())
	mxToken := moov.MX{
		AuthorizationCode: "fake-authorization-code",
	}
	result, err := mc.CreateBankAccount(context.Background(), s.accountID, moov.WithMX(mxToken))
	s.NoError(err)
	s.NotNil(result)
	s.Equal("1111222233330000", result.AccountNumber)
}

func (s *BankAccountTestSuite) TestGetBankAccount() {
	mc := NewTestClient(s.T())
	account, err := mc.GetBankAccount(context.Background(), s.accountID, s.bankAccountID)
	s.NoError(err)
	s.Equal("Sir Test ALot", account.HolderName)
}

func (s *BankAccountTestSuite) TestDeleteBankAccount() {
	mc := NewTestClient(s.T())

	err := mc.DeleteBankAccount(context.Background(), s.accountID, s.bankAccountIDDelete)
	s.NoError(err)
}

func (s *BankAccountTestSuite) TestListBankAccounts() {
	mc := NewTestClient(s.T())

	accounts, err := mc.ListBankAccounts(context.Background(), s.accountID)
	s.NoError(err)

	s.NotNil(accounts)
}

func (s *BankAccountTestSuite) TestMicroDepositInitiate() {
	mc := NewTestClient(s.T())

	err := mc.MicroDepositInitiate(context.Background(), s.accountID, s.bankAccountID)
	s.NoError(err)
}

// TODO: test this could run before TestMicroDepositInitiate
func (s *BankAccountTestSuite) TestMicroDepositConfirm() {
	mc := NewTestClient(s.T())

	err := mc.MicroDepositInitiate(context.Background(), s.accountID, s.bankAccountID)
	s.NoError(err)

	// sample data
	amounts := []int{0, 0}
	err = mc.MicroDepositConfirm(context.Background(), s.accountID, s.bankAccountID, amounts)
	s.NoError(err)
}
*/
