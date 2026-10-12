package mhooks

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"net/http"
)

var ErrInvalidSignature = errors.New("calculated signature does not match X-Signature header")

func checkSignature(headers http.Header, secret string) (bool, error) {
	var (
		timestamp = headers.Get("x-timestamp")
		nonce     = headers.Get("x-nonce")
		webhookID = headers.Get("x-webhook-id")
		gotHash   = headers.Get("x-signature")
	)

	concatHeaders := timestamp + "|" + nonce + "|" + webhookID
	wantHash, err := hash([]byte(concatHeaders), []byte(secret))
	if err != nil {
		return false, err
	}

	return hmac.Equal([]byte(*wantHash), []byte(gotHash)), nil
}

// hash generates a SHA512 HMAC hash of p using the secret provided.
func hash(p []byte, secret []byte) (*string, error) {
	h := hmac.New(sha512.New, secret)
	_, err := h.Write(p)
	if err != nil {
		return nil, err
	}
	hash := hex.EncodeToString(h.Sum(nil))
	return &hash, nil
}
