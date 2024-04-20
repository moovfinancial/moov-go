package moov_test

import (
	"testing"
)

func Test_Ping(t *testing.T) {
	mc := NewTestClient(t)
	err := mc.Ping(BgCtx())
	NoResponseError(t, err)
}
