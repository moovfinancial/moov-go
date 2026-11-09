package moov_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

func TestFetchTransfers(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   string
		want       moov.BatchGetTransfersResponse
		wantErr    bool
	}{
		{
			name:       "returns transfers keyed by transfer ID",
			statusCode: http.StatusOK,
			response:   `{"transfer-1":{"transferID":"transfer-1"},"transfer-2":{"transferID":"transfer-2"}}`,
			want: moov.BatchGetTransfersResponse{
				"transfer-1": {TransferID: "transfer-1"},
				"transfer-2": {TransferID: "transfer-2"},
			},
		},
		{
			name:       "preserves empty objects and omits IDs absent from the response",
			statusCode: http.StatusOK,
			response:   `{"transfer-1":{}}`,
			want:       moov.BatchGetTransfersResponse{"transfer-1": {}},
		},
		{
			name:       "returns an empty map for an empty response",
			statusCode: http.StatusOK,
			response:   `{}`,
			want:       moov.BatchGetTransfersResponse{},
		},
		{
			name:       "returns an error on API failure",
			statusCode: http.StatusNotFound,
			response:   `{"error":"not found"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				method string
				path   string
				body   moov.BatchGetTransfersRequest
			)

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				method = r.Method
				path = r.URL.Path
				_ = json.NewDecoder(r.Body).Decode(&body)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			}))
			t.Cleanup(srv.Close)

			client, err := moov.NewClient(
				moov.WithCredentials(moov.Credentials{PublicKey: "pk", SecretKey: "sk"}),
				moov.WithMoovURLScheme("http"),
			)
			require.NoError(t, err)
			client.Credentials.Host = strings.TrimPrefix(srv.URL, "http://")

			request := moov.BatchGetTransfersRequest{TransferIDs: []string{"transfer-1", "transfer-2"}}
			got, err := client.FetchTransfers(context.Background(), "account-1", request)

			require.Equal(t, http.MethodPost, method)
			require.Equal(t, "/accounts/account-1/transfers/.fetch", path)
			require.Equal(t, request, body)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
