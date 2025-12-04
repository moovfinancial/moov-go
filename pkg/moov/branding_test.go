package moov_test

// Tests for the Moov Account Branding API
// This file tests all CRUD operations for branding (colors) for accounts:
// - Create: CreateAccountBranding
// - Read: GetAccountBranding
// - Update: PatchAccountBranding
// - Upsert: UpsertAccountBranding
// See: https://docs.moov.io/api/enrichment/branding/

import (
	"context"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CreateAccountBranding(t *testing.T) {
	mc := NewTestClient(t)

	account := CreateTemporaryTestAccount(t, mc, createTestIndividualAccount())

	brand := moov.Brand{
		Colors: moov.BrandColors{
			Dark: moov.BrandColor{
				Accent: "#1a1a1a",
			},
			Light: moov.BrandColor{
				Accent: "#ffffff",
			},
		},
	}

	result, err := mc.CreateAccountBranding(context.Background(), account.AccountID, brand)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "#1a1a1a", result.Colors.Dark.Accent)
	assert.Equal(t, "#ffffff", result.Colors.Light.Accent)
}

func Test_GetAccountBranding(t *testing.T) {
	mc := NewTestClient(t)

	account := CreateTemporaryTestAccount(t, mc, createTestIndividualAccount())

	// First create branding
	brand := moov.Brand{
		Colors: moov.BrandColors{
			Dark: moov.BrandColor{
				Accent: "#2b2b2b",
			},
			Light: moov.BrandColor{
				Accent: "#f0f0f0",
			},
		},
	}

	_, err := mc.CreateAccountBranding(context.Background(), account.AccountID, brand)
	require.NoError(t, err)

	// Now get the branding
	result, err := mc.GetAccountBranding(context.Background(), account.AccountID)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "#2b2b2b", result.Colors.Dark.Accent)
	assert.Equal(t, "#f0f0f0", result.Colors.Light.Accent)
}

func Test_PatchAccountBranding(t *testing.T) {
	mc := NewTestClient(t)

	account := CreateTemporaryTestAccount(t, mc, createTestIndividualAccount())

	// First create branding
	initialBrand := moov.Brand{
		Colors: moov.BrandColors{
			Dark: moov.BrandColor{
				Accent: "#444444",
			},
			Light: moov.BrandColor{
				Accent: "#dddddd",
			},
		},
	}

	_, err := mc.CreateAccountBranding(context.Background(), account.AccountID, initialBrand)
	require.NoError(t, err)

	// Now patch the branding
	patchBrand := moov.Brand{
		Colors: moov.BrandColors{
			Dark: moov.BrandColor{
				Accent: "#333333",
			},
			Light: moov.BrandColor{
				Accent: "#eeeeee",
			},
		},
	}

	result, err := mc.PatchAccountBranding(context.Background(), account.AccountID, patchBrand)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "#333333", result.Colors.Dark.Accent)
	assert.Equal(t, "#eeeeee", result.Colors.Light.Accent)
}

func Test_UpsertAccountBranding(t *testing.T) {
	mc := NewTestClient(t)

	account := CreateTemporaryTestAccount(t, mc, createTestIndividualAccount())

	// Upsert branding (should create since it doesn't exist)
	brand := moov.Brand{
		Colors: moov.BrandColors{
			Dark: moov.BrandColor{
				Accent: "#444444",
			},
			Light: moov.BrandColor{
				Accent: "#dddddd",
			},
		},
	}

	result, err := mc.UpsertAccountBranding(context.Background(), account.AccountID, brand)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "#444444", result.Colors.Dark.Accent)
	assert.Equal(t, "#dddddd", result.Colors.Light.Accent)

	// Upsert again (should update existing branding)
	updatedBrand := moov.Brand{
		Colors: moov.BrandColors{
			Dark: moov.BrandColor{
				Accent: "#555555",
			},
			Light: moov.BrandColor{
				Accent: "#cccccc",
			},
		},
	}

	updatedResult, err := mc.UpsertAccountBranding(context.Background(), account.AccountID, updatedBrand)
	require.NoError(t, err)
	require.NotNil(t, updatedResult)
	assert.Equal(t, "#555555", updatedResult.Colors.Dark.Accent)
	assert.Equal(t, "#cccccc", updatedResult.Colors.Light.Accent)
}
