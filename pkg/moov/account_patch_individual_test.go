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

// Accounts merges the PATCH body into the stored profile, so "email":"" or
// "postalCode":"" replaces the stored value and an empty name fails validation. An
// individual patch must leave these fields out unless the caller sets them.
func TestPatchAccount_IndividualProfile(t *testing.T) {
	tests := []struct {
		name       string
		individual moov.PatchIndividualProfile
		expected   map[string]any
	}{
		{
			name:       "unset fields are omitted",
			individual: moov.PatchIndividualProfile{Address: &moov.Address{City: "Moov City"}},
			expected:   map[string]any{"address": map[string]any{"city": "Moov City"}},
		},
		{
			name: "supplied name and email are sent",
			individual: moov.PatchIndividualProfile{
				Name:  moov.Name{FirstName: "Jane", LastName: "Doe"},
				Email: "jane@moov.io",
			},
			expected: map[string]any{
				"name":  map[string]any{"firstName": "Jane", "lastName": "Doe"},
				"email": "jane@moov.io",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var body []byte

			// The handler runs on its own goroutine, so it records what it saw rather
			// than asserting.
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var err error
				body, err = io.ReadAll(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"accountID":"account-123"}`))
			}))
			t.Cleanup(srv.Close)

			patch := moov.PatchAccount{
				Profile: moov.PatchProfile{Individual: &tc.individual},
			}

			_, err := newFilesTestClient(t, srv).PatchAccount(context.Background(), "account-123", patch)
			require.NoError(t, err)

			var sent struct {
				Profile struct {
					Individual map[string]any `json:"individual"`
				} `json:"profile"`
			}
			require.NoError(t, json.Unmarshal(body, &sent))

			for _, field := range []string{"name", "email", "address"} {
				value, found := sent.Profile.Individual[field]
				expected, want := tc.expected[field]
				require.Equal(t, want, found, "%s presence in body: %s", field, body)
				require.Equal(t, expected, value)
			}
		})
	}
}
