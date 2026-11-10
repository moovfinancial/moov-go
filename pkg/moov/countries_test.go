package moov_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

func TestAccountCountries(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		assign     *moov.AccountCountries
		wantBody   []string
		wantResult []string
	}{
		{
			name:       "gets the countries",
			method:     http.MethodGet,
			path:       "/accounts/account-123/countries",
			wantResult: []string{"United States", "Canada"},
		},
		{
			name:       "assigns the countries",
			method:     http.MethodPut,
			path:       "/accounts/account-123/countries",
			assign:     &moov.AccountCountries{Countries: []string{"United States", "Canada"}},
			wantBody:   []string{"United States", "Canada"},
			wantResult: []string{"United States", "Canada"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				method    string
				path      string
				version   string
				accept    string
				body      moov.AccountCountries
				decodeErr error
			)

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				method = r.Method
				path = r.URL.Path
				version = r.Header.Get(moov.VersionHeader)
				accept = r.Header.Get("Accept")
				if tt.assign != nil {
					decodeErr = json.NewDecoder(r.Body).Decode(&body)
				}

				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"countries":["United States","Canada"]}`))
			}))
			t.Cleanup(srv.Close)

			client := newFilesTestClient(t, srv)

			var (
				countries *moov.AccountCountries
				err       error
			)
			if tt.assign != nil {
				countries, err = client.AssignAccountCountries(BgCtx(), "account-123", *tt.assign)
				require.NoError(t, decodeErr)
				require.Equal(t, tt.wantBody, body.Countries)
			} else {
				countries, err = client.GetAccountCountries(BgCtx(), "account-123")
			}
			require.NoError(t, err)

			require.Equal(t, tt.method, method)
			require.Equal(t, tt.path, path)
			require.Equal(t, "application/json", accept)
			require.Equal(t, moov.Version2026_07.String(), version)

			require.Equal(t, tt.wantResult, countries.Countries)
		})
	}
}
