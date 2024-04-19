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

func Test_CreateBankAccount_WithBankAccount(t *testing.T) {
	mc := NewTestClient(t)

	account := getLincolnBank(t, mc)

	resp, err := mc.CreateBankAccount(BgCtx(), account.AccountID, moov.WithBankAccount(moov.BankAccountRequest{
		HolderName:    "Sir Test Delete ALot",
		HolderType:    moov.HolderType_Individual,
		AccountType:   moov.BankAccountType_Checking,
		AccountNumber: randomBankAccountNumber(),
		RoutingNumber: "273976369",
	}))

	moov.DebugPrintResponse(err, fmt.Printf)
	require.NoError(t, err)
	require.NotNil(t, resp)

	defer func() {
		_ = mc.DeleteBankAccount(BgCtx(), account.AccountID, resp.BankAccountID)
	}()

	t.Run("get bank account", func(t *testing.T) {
		_, err = mc.GetBankAccount(BgCtx(), account.AccountID, resp.BankAccountID)
		moov.DebugPrintResponse(err, fmt.Printf)
		require.NoError(t, err)
	})

	t.Run("list bank accounts", func(t *testing.T) {
		listed, err := mc.ListBankAccounts(BgCtx(), account.AccountID)
		moov.DebugPrintResponse(err, fmt.Printf)
		require.NoError(t, err)
		require.NotEmpty(t, listed)
	})

	t.Run("micro-deposits", func(t *testing.T) {
		err := mc.MicroDepositInitiate(BgCtx(), account.AccountID, resp.BankAccountID)
		moov.DebugPrintResponse(err, fmt.Printf)
		require.NoError(t, err)

		// sample data
		amounts := []int{0, 0}
		err = mc.MicroDepositConfirm(BgCtx(), account.AccountID, resp.BankAccountID, amounts)
		moov.DebugPrintResponse(err, fmt.Printf)
		require.NoError(t, err)
	})

	t.Run("delete bank account", func(t *testing.T) {
		err = mc.DeleteBankAccount(BgCtx(), account.AccountID, resp.BankAccountID)
		moov.DebugPrintResponse(err, fmt.Printf)
		require.NoError(t, err)
	})
}

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

func randomBankAccountNumber() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(999999999))
	return fmt.Sprintf("%d", 100000000+n.Int64())
}
