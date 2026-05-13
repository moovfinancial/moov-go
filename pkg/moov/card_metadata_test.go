package moov_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	mc := NewTestClient(t)

	metadata, err := mc.GetCardMetadata(context.Background(), moov.CardMetadataRequest{
		CardNumber: "4111111111111111",
	})
	require.NoError(t, err)
	require.NotNil(t, metadata)

	assert.Equal(t, "411111", metadata.Bin)
	assert.Equal(t, "Visa", metadata.Brand)
	assert.Equal(t, "PERSONAL", metadata.CardCategory)
	assert.Equal(t, "debit", metadata.CardType)
	require.NotNil(t, metadata.Commercial)
	assert.False(t, *metadata.Commercial)
	require.NotNil(t, metadata.Regulated)
	assert.True(t, *metadata.Regulated)
	assert.Equal(t, "Moov Visa Sandbox", metadata.Issuer)
	assert.Equal(t, "US", metadata.IssuerCountry)
	assert.Equal(t, "+ (1) 123-123-1234", metadata.IssuerPhone)
	assert.Equal(t, "HTTPS://MOOV.IO/", metadata.IssuerURL)
	assert.Equal(t, "supported", metadata.DomesticPullFromCard)
	assert.Equal(t, "fast-funds", metadata.DomesticPushToCard)
}
