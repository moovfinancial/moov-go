package moov_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: Create faililng test for other card brands in test mode
// https://docs.moov.io/guides/get-started/test-mode/#cards

func TestCardMetadataMarshal(t *testing.T) {
	input := []byte(`{
		  "bin": "411111",
		  "brand": "Visa",
		  "cardCategory": "CLASSIC",
		  "cardType": "credit",
		  "commercial": false,
		  "regulated": false,
		  "issuer": "GRINGOTTS BANK",
		  "issuerCountry": "US",
		  "issuerPhone": "8185551212",
		  "issuerURL": "HTTPS://WWW.EXAMPLE.COM/",
		  "domesticPullFromCard": "supported",
		  "domesticPushToCard": "standard"
		}`)

	metadata := new(moov.CardMetadata)

	dec := json.NewDecoder(bytes.NewReader(input))
	dec.DisallowUnknownFields()

	require.NoError(t, dec.Decode(&metadata))

	assert.Equal(t, "411111", metadata.Bin)
	assert.Equal(t, "Visa", metadata.Brand)
	assert.Equal(t, "CLASSIC", metadata.CardCategory)
	assert.Equal(t, "credit", metadata.CardType)
	require.NotNil(t, metadata.Commercial)
	assert.False(t, *metadata.Commercial)
	require.NotNil(t, metadata.Regulated)
	assert.False(t, *metadata.Regulated)
	assert.Equal(t, "GRINGOTTS BANK", metadata.Issuer)
	assert.Equal(t, "US", metadata.IssuerCountry)
	assert.Equal(t, "8185551212", metadata.IssuerPhone)
	assert.Equal(t, "HTTPS://WWW.EXAMPLE.COM/", metadata.IssuerURL)
	assert.Equal(t, "supported", metadata.DomesticPullFromCard)
	assert.Equal(t, "standard", metadata.DomesticPushToCard)
}

func TestGetCardMetadata(t *testing.T) {
	const responseBody = `{
		"bin": "411111",
		"brand": "Visa",
		"cardCategory": "CLASSIC",
		"cardType": "credit",
		"commercial": false,
		"regulated": false,
		"issuer": "GRINGOTTS BANK",
		"issuerCountry": "US",
		"issuerPhone": "8185551212",
		"issuerURL": "HTTPS://WWW.EXAMPLE.COM/",
		"domesticPullFromCard": "supported",
		"domesticPushToCard": "standard"
	}`

	type capture struct {
		method      string
		path        string
		accept      string
		contentType string
		body        []byte
	}

	newClientWithServer := func(t *testing.T, status int, respBody string) (*moov.Client, *capture) {
		t.Helper()
		cap := &capture{}
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cap.method = r.Method
			cap.path = r.URL.Path
			cap.accept = r.Header.Get("Accept")
			cap.contentType = r.Header.Get("Content-Type")
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			cap.body = body

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			if respBody != "" {
				_, _ = w.Write([]byte(respBody))
			}
		}))
		t.Cleanup(srv.Close)

		host := strings.TrimPrefix(srv.URL, "http://")
		c, err := moov.NewClient(
			moov.WithCredentials(moov.Credentials{PublicKey: "pk", SecretKey: "sk", Host: host}),
			moov.WithMoovURLScheme("http"),
		)
		require.NoError(t, err)
		return c, cap
	}

	t.Run("success with raw cardNumber", func(t *testing.T) {
		c, cap := newClientWithServer(t, http.StatusOK, responseBody)

		metadata, err := c.GetCardMetadata(context.Background(), moov.CardMetadataRequest{
			CardNumber: "4111111111111111",
		})
		require.NoError(t, err)
		require.NotNil(t, metadata)

		assert.Equal(t, http.MethodPost, cap.method)
		assert.Equal(t, "/card-metadata", cap.path)
		assert.Equal(t, "application/json", cap.accept)
		assert.Equal(t, "application/json", cap.contentType)

		var sent map[string]any
		require.NoError(t, json.Unmarshal(cap.body, &sent))
		assert.Equal(t, "4111111111111111", sent["cardNumber"])
		_, hasE2EE := sent["e2ee"]
		assert.False(t, hasE2EE, "e2ee should be omitted when not provided")

		assert.Equal(t, "411111", metadata.Bin)
		assert.Equal(t, "Visa", metadata.Brand)
		assert.Equal(t, "credit", metadata.CardType)
		assert.Equal(t, "GRINGOTTS BANK", metadata.Issuer)
		assert.Equal(t, "supported", metadata.DomesticPullFromCard)
		assert.Equal(t, "standard", metadata.DomesticPushToCard)
		require.NotNil(t, metadata.Commercial)
		assert.False(t, *metadata.Commercial)
		require.NotNil(t, metadata.Regulated)
		assert.False(t, *metadata.Regulated)
	})

	t.Run("success with e2ee token", func(t *testing.T) {
		c, cap := newClientWithServer(t, http.StatusOK, responseBody)

		_, err := c.GetCardMetadata(context.Background(), moov.CardMetadataRequest{
			EndToEndToken: &moov.EndToEndToken{Token: "jwe-token"},
		})
		require.NoError(t, err)

		var sent map[string]any
		require.NoError(t, json.Unmarshal(cap.body, &sent))
		_, hasCardNumber := sent["cardNumber"]
		assert.False(t, hasCardNumber, "cardNumber should be omitted when not provided")
		e2ee, ok := sent["e2ee"].(map[string]any)
		require.True(t, ok, "e2ee should be present in body")
		assert.Equal(t, "jwe-token", e2ee["token"])
	})

	t.Run("error response is surfaced", func(t *testing.T) {
		c, _ := newClientWithServer(t, http.StatusUnprocessableEntity, `{"error":"cardNumber or e2ee is required"}`)

		metadata, err := c.GetCardMetadata(context.Background(), moov.CardMetadataRequest{})
		require.Error(t, err)
		assert.Nil(t, metadata)
	})
}
