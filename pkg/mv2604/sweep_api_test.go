package mv2604_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2604"
)

func TestListSweepConfigs(t *testing.T) {
	for _, tc := range []struct {
		name     string
		response string
		ids      []string
	}{
		{"two configs", `[{"sweepConfigID":"c1"},{"sweepConfigID":"c2"}]`, []string{"c1", "c2"}},
		{"empty list", `[]`, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := mv2604.NewSweepClient(newTestClient(t, &expectedRequest{
				method: http.MethodGet,
				path:   "/accounts/a/sweep-configs",
			}, http.StatusOK, tc.response))

			got, err := client.ListSweepConfigs(t.Context(), "a")
			require.NoError(t, err)
			require.Len(t, got, len(tc.ids))
			for i, id := range tc.ids {
				require.Equal(t, id, got[i].SweepConfigID)
			}
		})
	}
}

func TestGetSweepConfig(t *testing.T) {
	client := mv2604.NewSweepClient(newTestClient(t, &expectedRequest{
		method: http.MethodGet,
		path:   "/accounts/a/sweep-configs/c",
	}, http.StatusOK, `{"sweepConfigID":"c","walletID":"w"}`))

	got, err := client.GetSweepConfig(t.Context(), "a", "c")
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, "c", got.SweepConfigID)
	require.Equal(t, "w", got.WalletID)
}

func TestGetSweepConfigNotFound(t *testing.T) {
	client := mv2604.NewSweepClient(newTestClient(t, &expectedRequest{
		method: http.MethodGet,
		path:   "/accounts/a/sweep-configs/c",
	}, http.StatusNotFound, `{"error":"not found"}`))

	got, err := client.GetSweepConfig(t.Context(), "a", "c")
	require.Nil(t, got)
	var response moov.HttpCallResponse
	require.ErrorAs(t, err, &response)
	require.Equal(t, http.StatusNotFound, response.StatusCode())
}

func TestCreateSweepConfig(t *testing.T) {
	client := mv2604.NewSweepClient(newTestClient(t, &expectedRequest{
		method:   http.MethodPost,
		path:     "/accounts/a/sweep-configs",
		jsonBody: `{"walletID":"w","status":"enabled","pushPaymentMethodID":"p","pullPaymentMethodID":"q"}`,
	}, http.StatusOK, `{"sweepConfigID":"c","status":"enabled"}`))

	got, err := client.CreateSweepConfig(t.Context(), moov.CreateSweepConfig{
		AccountID:           "a",
		WalletID:            "w",
		Status:              moov.SweepConfigStatus_Enabled,
		PushPaymentMethodID: "p",
		PullPaymentMethodID: "q",
	})
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, "c", got.SweepConfigID)
	require.Equal(t, moov.SweepConfigStatus("enabled"), got.Status)
}

func TestUpdateSweepConfig(t *testing.T) {
	client := mv2604.NewSweepClient(newTestClient(t, &expectedRequest{
		method:   http.MethodPatch,
		path:     "/accounts/a/sweep-configs/c",
		jsonBody: `{"minimumBalance":"12.34","statementDescriptor":null}`,
	}, http.StatusOK, `{"sweepConfigID":"c","minimumBalance":"12.34"}`))

	minimumBalance := "12.34"
	got, err := client.UpdateSweepConfig(t.Context(), "a", "c", mv2604.UpdateSweepConfig{
		MinimumBalance:      &minimumBalance,
		StatementDescriptor: moov.SetNull[string](),
	})
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, "c", got.SweepConfigID)
	require.NotNil(t, got.MinimumBalance)
	require.Equal(t, "12.34", *got.MinimumBalance)
}
