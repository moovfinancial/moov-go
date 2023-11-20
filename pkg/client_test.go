package moov

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClientErr(t *testing.T) {
	_, err := NewClient()
	assert.Equal(t, ErrAuthCreditionalsNotSet, err)
}
