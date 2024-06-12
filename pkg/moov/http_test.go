package moov

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTTPCallResponse(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
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
