package moov_test

import (
	"os"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

func Test_CredentialsFromEnv(t *testing.T) {
	if os.Getenv(moov.ENV_MOOV_PUBLIC_KEY) == "" || os.Getenv(moov.ENV_MOOV_SECRET_KEY) == "" {
		t.Skip("MOOV_PUBLIC_KEY/MOOV_SECRET_KEY not set; skipping integration test")
	}

	creds := moov.CredentialsFromEnv()

	// This test is constructed this way as to not leak the public key and secret key on failures
	require.Greater(t, len(creds.PublicKey), 0)
	require.Greater(t, len(creds.SecretKey), 0)
}
