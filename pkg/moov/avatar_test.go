package moov_test

import (
	"context"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

func Test_Avatars(t *testing.T) {
	mc := NewTestClient(t)
	ctx := context.Background()

	// Avatar writes require the /accounts.write scope against an account the
	// facilitator owns, so use the facilitator's own account rather than a
	// connected merchant account.
	accountID := FACILITATOR_ID

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

		// The avatar is a singleton on the shared facilitator account, so a
		// concurrent test run may have already removed it. Tolerate not_found.
		if err != nil {
			var httpErr moov.HttpCallResponse
			require.ErrorAs(t, err, &httpErr)
			require.Equal(t, moov.StatusNotFound, httpErr.Status())
		}
	})
}
