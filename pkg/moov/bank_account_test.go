package moov_test

import (
	"fmt"
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
	}), moov.WaitForPaymentMethod())

	t.Cleanup(func() {
		if resp != nil {
			_ = mc.DeleteBankAccount(BgCtx(), account.AccountID, resp.BankAccountID)
		}
	})

	t.Run("add bank account and wait for payment methods", func(t *testing.T) {
		moov.DebugPrintResponse(err, fmt.Printf)
		require.NoError(t, err)
		require.NotNil(t, resp)

		// We told it to wait for the payment methods to be created, so this shouldn't be empty
		require.NotEmpty(t, resp.PaymentMethods)
	})

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
