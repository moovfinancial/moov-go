package moov_test

import (
	"context"
	"os"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"

	"github.com/joho/godotenv"
)

func init() {
	// If we don't have the environment variables set due to running in an IDE or directly via the go command, load it up here
	if _, ok := os.LookupEnv(moov.ENV_MOOV_SECRET_KEY); !ok {
		godotenv.Load("../secrets.env")
	}
}

func BgCtx() context.Context {
	return context.Background()
}

func NewTestClient(t require.TestingT, c ...moov.ClientConfigurable) *moov.Client {
	mc, err := moov.NewClient(c...)
	require.NoError(t, err)

	require.NoError(t, mc.Ping(BgCtx()), "Unable to ping with credentials")

	return mc
}

func Test_Client(t *testing.T) {
	mc := NewTestClient(t)

	err := mc.Ping(BgCtx())
	require.NoError(t, err)
}

func Test_Client_InvalidCredentials(t *testing.T) {
	_, err := moov.NewClient(moov.WithCredentials(moov.Credentials{}))
	require.Equal(t, moov.ErrAuthCredentialsNotSet, err)
}
