package moov_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"

	"github.com/joho/godotenv"
)

func BgCtx() context.Context {
	return context.Background()
}

func NewTestClient(t testing.TB, c ...moov.ClientConfigurable) *moov.Client {
	// If we have a 5352b013-ae58-4a63-8a3f-97f316a917cf file written read that and populate the test environment
	secretsPath := filepath.Join("..", "..", "secrets.env")

	if _, err := os.Stat(secretsPath); err == nil {
		secrets, err := godotenv.Read(secretsPath)
		require.NoError(t, err)

		for k, v := range secrets {
			t.Setenv(k, v)
		}
	}

	c = append(c, moov.WithDecoder(strictDecoder))

	mc, err := moov.NewClient(c...)
	require.NoError(t, err)

	require.NoError(t, mc.Ping(BgCtx()), "Unable to ping with credentials")

	return mc
}

func strictDecoder(r io.Reader, contentType string, item any) error {
	if strings.Contains(contentType, "application/json") {
		dec := json.NewDecoder(r)
		dec.DisallowUnknownFields()
		return dec.Decode(item)
	}
	return fmt.Errorf("unknown content-type %s", contentType)
}

func Test_Client(t *testing.T) {
	mc := NewTestClient(t)

	err := mc.Ping(BgCtx())
	require.NoError(t, err)
}

func Test_Client_InvalidCredentials(t *testing.T) {
	_, err := moov.NewClient(moov.WithCredentials(moov.Credentials{}))
	require.Equal(t, moov.ErrCredentialsNotSet, err)
}
