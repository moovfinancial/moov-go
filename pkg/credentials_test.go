package moov_test

import (
	"os"
	"testing"

	moov "github.com/moovfinancial/moov-go/pkg"
	"github.com/stretchr/testify/require"
)

func Test_CredentialsFromEnv(t *testing.T) {
	// If its not set on the local machine just skip for now.
	if _, ok := os.LookupEnv(moov.ENV_MOOV_PUBLIC_KEY); !ok {
		t.SkipNow()
	}

	creds := moov.CredentialsFromEnv()

	// This test is constructed this way as to not leak the public key and secret key on failures
	require.Greater(t, len(creds.PublicKey), 0)
	require.Greater(t, len(creds.SecretKey), 0)
}
