package e2ee_test

import (
	"context"
	"testing"

	"github.com/go-jose/go-jose/v4"
	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

// This is an example of what a client device would be doing to
// pass data via PCI approved end to end encryption.
// This is only needed if you need to send card details from a device
// through an intermediary and then to Moov. If the application
// speaks directly to Moov this isn't needed.
//
// `make examples-e2ee` to run it yourself
//
// Please contact Moov representative to enable
// on production accounts.
func Test_EndToEndEncryption(t *testing.T) {
	ctx := context.Background()

	mc, err := moov.NewClient(moov.WithCredentials(moov.CredentialsFromEnv()))
	require.NoError(t, err)

	// Need to fetch a public key from Moov's API.
	// This is a public endpoint so no credentials required but is aggresively rate limited.
	jwk, err := mc.GenerateEndToEndPublicKey(ctx)
	require.NoError(t, err)

	// Serialize however you want to get it into json
	json := []byte(`
		{
			"test": true
		}
	`)

	// Generate the JWE token to pass to Moov
	// The returned token is what you'll use in the endpoint you're sending the data over for.
	token, err := Encrypt(jwk, json)
	require.NoError(t, err)

	// Send the JWE token to Moov on its debug endpoint to verify it worked.
	// This will be replaced with whatever endpoint you're planning to use it on.
	err = mc.TestEndToEndToken(ctx, token)
	require.NoError(t, err)
}

func Encrypt(pub *jose.JSONWebKey, data []byte) (string, error) {
	crypter, err := jose.NewEncrypter(jose.A256GCM, jose.Recipient{
		Algorithm: jose.ECDH_ES_A256KW,
		Key:       pub,
	}, nil)
	if err != nil {
		return "", err
	}

	obj, err := crypter.Encrypt(data)
	if err != nil {
		return "", err
	}

	json, err := obj.CompactSerialize()
	if err != nil {
		return "", err
	}

	return json, nil
}
