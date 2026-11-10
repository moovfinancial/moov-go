package testtools

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync/atomic"
)

func RandomBankAccountNumber() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(999999999))
	return fmt.Sprintf("%d", 100000000+n.Int64())
}

var (
	foreignIDPrefix = rand.Text()
	foreignIDSeq    atomic.Uint64
)

func RandomForeignID() string {
	return fmt.Sprintf("%s-%d", foreignIDPrefix, foreignIDSeq.Add(1))
}
