package moov

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	accountID = "e9a255fd-ecdb-4c78-8049-e9b725bdd28c"
)

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

func TestCreateBankAccount(t *testing.T) {
	bankAccount := BankAccount{
		HolderName:      "Jules Jackson",
		HolderType:      "individual",
		BankAccountType: "checking",
		AccountNumber:   "81957272176",
		RoutingNumber:   "273976369",
	}

	mc, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}

	bankAccount, err = mc.CreateBankAccount(accountID, bankAccount)
	if err != nil {
		t.Fatalf("Error creating bank account: %v", err)
	}

	assert.NotNil(t, bankAccount.BankAccountID)
}

func TestGetBankAccount(t *testing.T) {
	mc, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}
	// sample bank account id
	bankAccountID := "fc8cc388-35d0-4da6-bebd-9b9800fe9d4d"

	account, err := mc.GetBankAccount(accountID, bankAccountID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, bankAccountID, account.BankAccountID)
}

func TestDeleteBankAccount(t *testing.T) {
	mc, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}
	// sample bank account id
	bankAccountID := "bdba60ba-3e22-4984-8159-e86ebbf43690"

	err = mc.DeleteBankAccount(accountID, bankAccountID)
	if err != nil {
		assert.Error(t, err)
	}
}

func TestListBankAccounts(t *testing.T) {
	mc, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}

	accounts, err := mc.ListBankAccounts(accountID)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(accounts))
	assert.NotNil(t, accounts)
}

func TestMicroDepositInitiate(t *testing.T) {
	mc, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}
	// sample bank account id
	bankAccountID := "fc8cc388-35d0-4da6-bebd-9b9800fe9d4d"

	err = mc.MicroDepositInitiate(accountID, bankAccountID)
	if err != nil {
		assert.Error(t, err)
	}
}

func TestMicroDepositConfirm(t *testing.T) {
	mc, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}
	// sample data
	bankAccountID := "c9fe9144-dbe5-4d6f-9994-c5be9813718c"
	amounts := []int{22, 21}
	err = mc.MicroDepositConfirm(accountID, bankAccountID, amounts)
	if err != nil {
		assert.Error(t, err)
	}
}
