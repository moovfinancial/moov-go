package testtools

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync/atomic"

	"github.com/google/uuid"
)

var (
	foreignIDPrefix  = uuid.NewString()
	foreignIDCounter atomic.Uint64
)

func RandomBankAccountNumber() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(999999999))
	return fmt.Sprintf("%d", 100000000+n.Int64())
}

// RandomForeignID returns a foreign identifier unique within this process. The
// process-scoped prefix keeps concurrent test processes from colliding.
func RandomForeignID() string {
	return fmt.Sprintf("%s-%d", foreignIDPrefix, foreignIDCounter.Add(1))
}
