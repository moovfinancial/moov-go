package moov_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Statements(t *testing.T) {
	c := NewTestClient(t)

	statements, err := c.ListStatements(t.Context(), FACILITATOR_ID)
	require.NoError(t, err)
	require.NotNil(t, statements)
	if len(statements) > 0 {
		statement, err := c.GetStatement(t.Context(), FACILITATOR_ID, statements[0].StatementID)
		require.NoError(t, err)
		require.NotNil(t, statement)
		require.Equal(t, statement, statements[0])
	}
}
