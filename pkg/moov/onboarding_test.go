package moov_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

const onboardingInviteJSON = `{
	"code": "invite-123",
	"link": "https://onboarding.moov.io/invite-123",
	"returnURL": "https://moov.io",
	"termsOfServiceURL": "https://moov.io/terms",
	"scopes": ["accounts.read"],
	"grantScopes": ["accounts.read"],
	"capabilities": ["transfers", "wallet"],
	"feePlanCodes": ["merchant-direct"],
	"redeemedAccountID": "account-456",
	"createdOn": "2026-07-01T00:00:00Z",
	"partner": {"accountID": "partner-123", "accountMode": "sandbox", "displayName": "Partner"}
}`

func Test_OnboardingInvite_CreateGetListRevoke(t *testing.T) {
	var (
		got struct {
			method  string
			path    string
			version string
			body    map[string]any
		}
		bodyErr error
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got.method = r.Method
		got.path = r.URL.Path
		got.version = r.Header.Get(moov.VersionHeader)
		got.body = nil
		bodyErr = nil
		if r.Method == http.MethodPost {
			bodyErr = json.NewDecoder(r.Body).Decode(&got.body)
		}

		switch {
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/missing"):
			w.WriteHeader(http.StatusNotFound)
		case r.Method == http.MethodGet && r.URL.Path == "/onboarding-invites":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte("[" + onboardingInviteJSON + "]"))
		case r.Method == http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(onboardingInviteJSON))
		case r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(onboardingInviteJSON))
		}
	}))
	t.Cleanup(srv.Close)

	mc := newFilesTestClient(t, srv)
	ctx := BgCtx()

	created, err := mc.CreateOnboardingInvite(ctx, moov.OnboardingInviteRequest{
		ReturnURL: "https://moov.io",
		Scopes:    []moov.ApplicationScope{"accounts.read"},
		Capabilities: []moov.CapabilityName{
			moov.CapabilityName_Transfers,
			moov.CapabilityName_Wallet,
		},
		FeePlanCodes: []string{"merchant-direct"},
		Prefill:      &moov.CreateAccount{Type: moov.AccountType_Business},
	})
	require.NoError(t, err)
	require.NoError(t, bodyErr)
	require.Equal(t, http.MethodPost, got.method)
	require.Equal(t, "/onboarding-invites", got.path)
	// Onboarding invites keep the default API version, so no version header is sent.
	require.Empty(t, got.version)
	require.Equal(t, "https://moov.io", got.body["returnURL"])
	require.Equal(t, []any{"accounts.read"}, got.body["scopes"])
	require.Equal(t, []any{"transfers", "wallet"}, got.body["capabilities"])
	require.Equal(t, []any{"merchant-direct"}, got.body["feePlanCodes"])
	require.Equal(t, "business", got.body["prefill"].(map[string]any)["accountType"])
	require.Equal(t, moov.OnboardingInviteCode("invite-123"), created.Code)
	require.Equal(t, "https://moov.io", created.ReturnURL)
	require.Equal(t, []moov.CapabilityName{moov.CapabilityName_Transfers, moov.CapabilityName_Wallet}, created.Capabilities)

	fetched, err := mc.GetOnboardingInvite(ctx, "invite-123")
	require.NoError(t, err)
	require.Equal(t, http.MethodGet, got.method)
	require.Equal(t, "/onboarding-invites/invite-123", got.path)
	require.Empty(t, got.version)
	require.Equal(t, created.Code, fetched.Code)

	listed, err := mc.ListOnboardingInvites(ctx)
	require.NoError(t, err)
	require.Equal(t, "/onboarding-invites", got.path)
	require.Len(t, listed, 1)
	require.Equal(t, created.Code, listed[0].Code)
	require.Equal(t, "account-456", listed[0].RedeemedAccountID)

	err = mc.RevokeOnboardingInvite(ctx, "invite-123")
	require.NoError(t, err)
	require.Equal(t, http.MethodDelete, got.method)
	require.Equal(t, "/onboarding-invites/invite-123", got.path)
	require.Empty(t, got.version)

	_, err = mc.GetOnboardingInvite(ctx, "missing")
	require.Error(t, err)
	var httpErr moov.HttpCallResponse
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, moov.StatusNotFound, httpErr.Status())
}
