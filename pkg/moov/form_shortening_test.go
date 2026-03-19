package moov_test

// Tests for the Moov Form Shortening API
// - ListIndustries
// - EnrichBusinessProfile
// See: https://docs.moov.io/api/v2026.01.00/enrichment/form-shortening/

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ListIndustries(t *testing.T) {
	mc := NewTestClient(t)

	result, err := mc.ListIndustries(context.Background())
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotEmpty(t, result.Industries)

	// Spot-check that industry entries are populated
	first := result.Industries[0]
	assert.NotEmpty(t, first.Industry)
	assert.NotEmpty(t, first.DisplayName)
	assert.NotEmpty(t, first.Category)
	assert.NotEmpty(t, first.CategoryDisplayName)
	assert.NotEmpty(t, first.DefaultMcc)
}

func Test_EnrichBusinessProfile(t *testing.T) {
	mc := NewTestClient(t)

	result, err := mc.EnrichBusinessProfile(context.Background(), "test@moov.io")
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.Business)
	assert.NotEmpty(t, result.Business.LegalBusinessName)
}
