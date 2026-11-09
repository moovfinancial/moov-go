package testtools

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSeededIDsFromEnv(t *testing.T) {
	const envName = "MOOV_TEST_SEEDED_ID"

	t.Run("override wins", func(t *testing.T) {
		t.Setenv(envName, "override")

		require.Equal(t, "override", seededID(envName, "default"))
	})

	t.Run("empty override keeps default", func(t *testing.T) {
		t.Setenv(envName, "")

		require.Equal(t, "default", seededID(envName, "default"))
	})

	t.Run("unset override keeps default", func(t *testing.T) {
		require.Equal(t, "default", seededID(envName, "default"))
	})
}
