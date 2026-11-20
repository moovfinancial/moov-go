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
	"github.com/moovfinancial/moov-go/pkg/mv2507"
)

// Accounts treats "taxID":null as a request to delete the EIN on an unlocked profile,
// so a partial business patch must leave taxID out of the body unless the caller sets it.
func TestPatchAccount_BusinessTaxID(t *testing.T) {
	patchers := []struct {
		name  string
		patch func(client *moov.Client, patch moov.PatchAccount) error
	}{
		{
			name: "Client.PatchAccount",
			patch: func(client *moov.Client, patch moov.PatchAccount) error {
				_, err := client.PatchAccount(context.Background(), "account-123", patch)
				return err
			},
		},
		{
			name: "mv2507.Accounts.Patch",
			patch: func(client *moov.Client, patch moov.PatchAccount) error {
				_, err := mv2507.Accounts.Patch(context.Background(), *client, "account-123", patch)
				return err
			},
		},
	}

	tests := []struct {
		name     string
		taxID    *moov.TaxID
		expected map[string]any
	}{
		{
			name:     "nil tax ID is omitted",
			taxID:    nil,
			expected: nil,
		},
		{
			name:     "supplied tax ID is sent",
			taxID:    &moov.TaxID{EIN: moov.EIN{Number: "123456789"}},
			expected: map[string]any{"ein": map[string]any{"number": "123456789"}},
		},
	}

	for _, p := range patchers {
		for _, tc := range tests {
			t.Run(p.name+"/"+tc.name, func(t *testing.T) {
				var (
					method string
					path   string
					body   []byte
				)

				// The handler runs on its own goroutine, so it records what it saw rather
				// than asserting.
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					method = r.Method
					path = r.URL.Path
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
					Profile: moov.PatchProfile{
						Business: &moov.PatchBusinessProfile{
							Address: &moov.Address{
								AddressLine1:    "123 Main St",
								City:            "Moov City",
								StateOrProvince: "CA",
								PostalCode:      "12345",
								Country:         "US",
							},
							TaxID: tc.taxID,
						},
					},
				}

				err := p.patch(newFilesTestClient(t, srv), patch)
				require.NoError(t, err)

				require.Equal(t, http.MethodPatch, method)
				require.Equal(t, "/accounts/account-123", path)

				var sent struct {
					Profile struct {
						Business map[string]any `json:"business"`
					} `json:"profile"`
				}
				require.NoError(t, json.Unmarshal(body, &sent))
				require.Contains(t, sent.Profile.Business, "address")

				taxID, found := sent.Profile.Business["taxID"]
				if tc.expected == nil {
					require.False(t, found, "taxID must not be sent: %s", body)
					return
				}
				require.Equal(t, tc.expected, taxID)
			})
		}
	}
}
