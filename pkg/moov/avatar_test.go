package moov_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Avatars(t *testing.T) {
	mc := NewTestClient(t)
	ctx := context.Background()
	accountID := MERCHANT_ID

	// ensure any uploaded avatar is removed even if a test below fails
	t.Cleanup(func() {
		_ = mc.DeleteAvatar(ctx, accountID)
	})

	t.Run("upload avatar", func(t *testing.T) {
		_, imgReader := randomImage(t, 100, 100, encodePNG)

		err := mc.UploadAvatar(ctx, accountID, imgReader)
		require.NoError(t, err)
	})

	t.Run("delete avatar", func(t *testing.T) {
		err := mc.DeleteAvatar(ctx, accountID)
		require.NoError(t, err)
	})
}
