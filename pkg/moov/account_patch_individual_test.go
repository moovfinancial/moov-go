package moov_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

// Accounts merges the PATCH body into the stored profile, so "email":"" replaces the
// stored email. An individual patch must leave email out unless the caller sets it.
func TestPatchAccount_IndividualEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected any
	}{
		{name: "empty email is omitted", email: "", expected: nil},
		{name: "supplied email is sent", email: "jane@moov.io", expected: "jane@moov.io"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var body []byte

			// The handler runs on its own goroutine, so it records what it saw rather
			// than asserting.
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, _ = io.ReadAll(r.Body)

				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"accountID":"account-123"}`))
			}))
			t.Cleanup(srv.Close)

			patch := moov.PatchAccount{
				Profile: moov.PatchProfile{
					Individual: &moov.PatchIndividualProfile{
						Name:  moov.Name{FirstName: "Jane", LastName: "Doe"},
						Email: tc.email,
					},
				},
			}

			_, err := newFilesTestClient(t, srv).PatchAccount(context.Background(), "account-123", patch)
			require.NoError(t, err)

			var sent struct {
				Profile struct {
					Individual map[string]any `json:"individual"`
				} `json:"profile"`
			}
			require.NoError(t, json.Unmarshal(body, &sent))
			require.Contains(t, sent.Profile.Individual, "name")

			email, found := sent.Profile.Individual["email"]
			if tc.expected == nil {
				require.False(t, found, "email must not be sent: %s", body)
				return
			}
			require.Equal(t, tc.expected, email)
		})
	}
}
