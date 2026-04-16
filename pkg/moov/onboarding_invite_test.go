package moov_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

func TestOnboardingInviteRequest_Serialization(t *testing.T) {
	t.Run("full request serializes correctly", func(t *testing.T) {
		req := moov.OnboardingInviteRequest{
			ReturnURL:         moov.PtrOf("https://example.com/return"),
			TermsOfServiceURL: moov.PtrOf("https://example.com/tos"),
			Scopes:            []string{"accounts.read", "accounts.write"},
			GrantScopes:       []string{"transfers.read"},
			Capabilities:      []string{"transfers", "send-funds"},
			FeePlanCodes:      []string{"standard"},
			Prefill: &moov.CreateAccount{
				Type: moov.AccountType_Business,
				Profile: moov.CreateProfile{
					Business: &moov.CreateBusinessProfile{
						Name: "Test Business",
					},
				},
			},
		}

		data, err := json.Marshal(req)
		require.NoError(t, err)

		var decoded map[string]interface{}
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)

		require.Equal(t, "https://example.com/return", decoded["returnURL"])
		require.Equal(t, "https://example.com/tos", decoded["termsOfServiceURL"])
		require.Len(t, decoded["scopes"], 2)
		require.Len(t, decoded["grantScopes"], 1)
		require.Len(t, decoded["capabilities"], 2)
		require.Len(t, decoded["feePlanCodes"], 1)
		require.NotNil(t, decoded["prefill"])
	})

	t.Run("minimal request omits empty fields", func(t *testing.T) {
		req := moov.OnboardingInviteRequest{}

		data, err := json.Marshal(req)
		require.NoError(t, err)

		var decoded map[string]interface{}
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)

		require.NotContains(t, decoded, "returnURL")
		require.NotContains(t, decoded, "termsOfServiceURL")
		require.NotContains(t, decoded, "prefill")
	})
}

func TestOnboardingInvite_Deserialization(t *testing.T) {
	t.Run("full response deserializes correctly", func(t *testing.T) {
		responseJSON := `{
			"code": "abc123",
			"link": "https://moov.io/invite/abc123",
			"returnURL": "https://example.com/return",
			"termsOfServiceURL": "https://example.com/tos",
			"scopes": ["accounts.read"],
			"grantScopes": ["transfers.read"],
			"capabilities": ["transfers"],
			"feePlanCodes": ["standard"],
			"redeemedAccountID": "acct-123",
			"prefill": {
				"accountType": "business",
				"profile": {
					"business": {
						"legalBusinessName": "Test Business"
					}
				}
			},
			"createdOn": "2026-01-01T00:00:00Z",
			"revokedOn": "2026-01-02T00:00:00Z",
			"redeemedOn": "2026-01-03T00:00:00Z"
		}`

		var invite moov.OnboardingInvite
		err := json.Unmarshal([]byte(responseJSON), &invite)
		require.NoError(t, err)

		require.Equal(t, "abc123", invite.Code)
		require.Equal(t, "https://moov.io/invite/abc123", invite.Link)
		require.NotNil(t, invite.ReturnURL)
		require.Equal(t, "https://example.com/return", *invite.ReturnURL)
		require.NotNil(t, invite.TermsOfServiceURL)
		require.Equal(t, "https://example.com/tos", *invite.TermsOfServiceURL)
		require.Equal(t, []string{"accounts.read"}, invite.Scopes)
		require.Equal(t, []string{"transfers.read"}, invite.GrantScopes)
		require.Equal(t, []string{"transfers"}, invite.Capabilities)
		require.Equal(t, []string{"standard"}, invite.FeePlanCodes)
		require.NotNil(t, invite.RedeemedAccountID)
		require.Equal(t, "acct-123", *invite.RedeemedAccountID)
		require.NotNil(t, invite.Prefill)
		require.Equal(t, moov.AccountType_Business, invite.Prefill.Type)
		require.Equal(t, time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), invite.CreatedOn)
		require.NotNil(t, invite.RevokedOn)
		require.Equal(t, time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC), *invite.RevokedOn)
		require.NotNil(t, invite.RedeemedOn)
		require.Equal(t, time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC), *invite.RedeemedOn)
	})

	t.Run("minimal response with only required fields", func(t *testing.T) {
		responseJSON := `{
			"code": "xyz789",
			"link": "https://moov.io/invite/xyz789",
			"createdOn": "2026-01-01T00:00:00Z"
		}`

		var invite moov.OnboardingInvite
		err := json.Unmarshal([]byte(responseJSON), &invite)
		require.NoError(t, err)

		require.Equal(t, "xyz789", invite.Code)
		require.Equal(t, "https://moov.io/invite/xyz789", invite.Link)
		require.Nil(t, invite.ReturnURL)
		require.Nil(t, invite.TermsOfServiceURL)
		require.Nil(t, invite.RedeemedAccountID)
		require.Nil(t, invite.Prefill)
		require.Nil(t, invite.RevokedOn)
		require.Nil(t, invite.RedeemedOn)
	})

	t.Run("list response deserializes correctly", func(t *testing.T) {
		responseJSON := `[
			{"code": "inv1", "link": "https://moov.io/invite/inv1", "createdOn": "2026-01-01T00:00:00Z"},
			{"code": "inv2", "link": "https://moov.io/invite/inv2", "createdOn": "2026-01-02T00:00:00Z"}
		]`

		var invites []moov.OnboardingInvite
		err := json.Unmarshal([]byte(responseJSON), &invites)
		require.NoError(t, err)

		require.Len(t, invites, 2)
		require.Equal(t, "inv1", invites[0].Code)
		require.Equal(t, "inv2", invites[1].Code)
	})
}

func TestOnboardingInvite_RoundTrip(t *testing.T) {
	t.Run("request round trips through JSON", func(t *testing.T) {
		original := moov.OnboardingInviteRequest{
			ReturnURL:         moov.PtrOf("https://example.com/return"),
			TermsOfServiceURL: moov.PtrOf("https://example.com/tos"),
			Scopes:            []string{"accounts.read"},
			GrantScopes:       []string{"transfers.read"},
			Capabilities:      []string{"transfers"},
			FeePlanCodes:      []string{"standard"},
		}

		data, err := json.Marshal(original)
		require.NoError(t, err)

		var decoded moov.OnboardingInviteRequest
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)

		require.Equal(t, *original.ReturnURL, *decoded.ReturnURL)
		require.Equal(t, *original.TermsOfServiceURL, *decoded.TermsOfServiceURL)
		require.Equal(t, original.Scopes, decoded.Scopes)
		require.Equal(t, original.GrantScopes, decoded.GrantScopes)
		require.Equal(t, original.Capabilities, decoded.Capabilities)
		require.Equal(t, original.FeePlanCodes, decoded.FeePlanCodes)
	})
}
