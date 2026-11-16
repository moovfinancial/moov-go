package testtools

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/google/uuid"
)

func RandomBankAccountNumber() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(999999999))
	return fmt.Sprintf("%d", 100000000+n.Int64())
}

func RandomForeignID() string {
	return uuid.NewString()
}
