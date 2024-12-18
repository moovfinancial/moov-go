package moov_test

import (
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

func TestCreateApplicationKeys(t *testing.T) {
	t.Skip("skipping until we have an application, requires manual testing against a special setup")

	mc := NewTestClient(t)

	apps, err := mc.ListApplications(BgCtx())
	require.NoError(t, err)
	require.NotEmpty(t, apps)

	createKey := moov.CreateApplicationKey{
		Name:        "moov-go-created-key",
		Description: "moov-go-created-key",
		Origins:     []string{"https://example.com"},
	}

	key, err := mc.CreateApplicationKeys(BgCtx(), apps[0].ApplicationID, createKey)
	require.NoError(t, err)
	require.NotNil(t, key)
}
