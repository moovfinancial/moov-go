package mhooks

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const (
	testTimestamp = "2024-04-26T21:20:55Z"
	testSecret    = "my-webhook-signing-secret"
	testNonce     = "LwxF1Uk7QOeDq2nzB3theslHbtAo7y3uuncB1PoijwCZZaRZsd8DOtffBT7p"
	testWebhookID = "dff0a709-f982-4475-81e8-214b435c74ab"
	testSignature = "6231d03752de6963087e6aea1c78a27a0617b6df1c071195f30ed85defe34e02fd0bf3995949fe12dafd747c42de9cfae03b8aafcf69cceba5495f4c7b719d82"
)

func TestCheckSignature(t *testing.T) {
	t.Parallel()

	header := func(signature string) http.Header {
		h := make(http.Header)
		h.Set("x-timestamp", testTimestamp)
		h.Set("x-nonce", testNonce)
		h.Set("x-webhook-id", testWebhookID)
		h.Set("x-signature", signature)
		return h
	}

	t.Run("matching tag", func(t *testing.T) {
		t.Parallel()

		ok, err := checkSignature(header(testSignature), testSecret)
		require.NoError(t, err)
		require.True(t, ok)
	})

	t.Run("mismatched tag", func(t *testing.T) {
		t.Parallel()

		bad := testSignature[:len(testSignature)-1] + "0"
		ok, err := checkSignature(header(bad), testSecret)
		require.NoError(t, err)
		require.False(t, ok)
	})

	t.Run("wrong length tag", func(t *testing.T) {
		t.Parallel()

		ok, err := checkSignature(header("abcd"), testSecret)
		require.NoError(t, err)
		require.False(t, ok)
	})

	t.Run("empty tag", func(t *testing.T) {
		t.Parallel()

		ok, err := checkSignature(header(""), testSecret)
		require.NoError(t, err)
		require.False(t, ok)
	})
}

func TestParseEvent_InvalidSignature(t *testing.T) {
	t.Parallel()

	createdOn, err := time.Parse(time.RFC3339, testTimestamp)
	require.NoError(t, err)

	event := Event{
		EventID:   uuid.NewString(),
		EventType: EventTypeTestPing,
		Data:      json.RawMessage(`{"ping":true}`),
		CreatedOn: createdOn,
	}

	var body bytes.Buffer
	err = json.NewEncoder(&body).Encode(event)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/webhook", &body)
	req.Header.Set("x-timestamp", testTimestamp)
	req.Header.Set("x-nonce", testNonce)
	req.Header.Set("x-webhook-id", testWebhookID)
	req.Header.Set("x-signature", "abcd")

	_, err = ParseEvent(req, testSecret)
	require.ErrorIs(t, err, ErrInvalidSignature)
}
