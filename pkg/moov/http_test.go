package moov

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTTPCallResponse(t *testing.T) {
	t.Run("400", func(t *testing.T) {
		resp := &httpCallResponse{
			resp: &http.Response{
				StatusCode: http.StatusBadRequest,
			},
			body: []byte(`{"error":"X-Idempotency-Key HTTP header must be a valid UUID"}`),
		}
		expected := strings.TrimSpace(`
error from moov - status: bad_request http.request_id:  http.status_code: 400
  X-Idempotency-Key HTTP header must be a valid UUID
`)
		require.Equal(t, expected, resp.Error())
	})

	t.Run("409 and 422", func(t *testing.T) {
		resp := &httpCallResponse{
			resp: &http.Response{
				StatusCode: http.StatusUnprocessableEntity,
			},
			body: []byte(`{"profile":{"business":{"taxID":{"ein":{"number":"must be a valid employer identification number"}}}}}`),
		}
		expected := "error from moov - status: failed_validation http.request_id:  http.status_code: 422 - profile.business.taxID.ein.number: must be a valid employer identification number"
		require.Equal(t, expected, resp.Error())
	})
}
