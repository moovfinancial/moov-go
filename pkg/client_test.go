package moov

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClientErr(t *testing.T) {
	creds := Credentials{}
	_, err := NewClient(creds)
	assert.Equal(t, ErrAuthCreditionalsNotSet, err)
}
