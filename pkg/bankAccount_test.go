package moov

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBankAccountMarshal(t *testing.T) {
	input := []byte(`{
		"bankAccountID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43",
		"fingerprint": "9948962d92a1ce40c9f918cd9ece3a22bde62fb325a2f1fe2e833969de672ba3",
		"status": "new",
		"holderName": "Jules Jackson",
		"holderType": "individual",
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
