package moov

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/joho/godotenv"
)

func init() {
	// If we don't have the environment variables set due to running in an IDE or directly via the go command, load it up here
	if _, ok := os.LookupEnv(ENV_MOOV_SECRET_KEY); !ok {
		godotenv.Load("../secrets.env")
	}
}

func NewTestClient(t require.TestingT, c ...ClientConfigurable) *Client {
	require.NoError(t, godotenv.Load("../secrets.env"), "Unable to load env")

	mc, err := NewClient(c...)
	require.NoError(t, err)

	require.NoError(t, mc.Ping(), "Unable to ping with credentials")

	return mc
}

func Test_Client(t *testing.T) {
	mc := NewTestClient(t)

	err := mc.Ping()
	require.NoError(t, err)
}

func Test_Client_InvalidCredentials(t *testing.T) {
	_, err := NewClient(WithCredentials(Credentials{}))
	require.Equal(t, ErrAuthCredentialsNotSet, err)
}
