package mv2604_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2604"
)

func TestPatchTransfer(t *testing.T) {
	for _, tc := range []struct {
		name   string
		update mv2604.PatchTransfer
		want   string
	}{
		{name: "omits foreignID", want: `{}`},
		{
			name:   "clears foreignID",
			update: mv2604.PatchTransfer{ForeignID: moov.SetNull[string]()},
			want:   `{"foreignID":null}`,
		},
		{
			name: "sets foreignID and metadata",
			update: mv2604.PatchTransfer{
				ForeignID: moov.Set("f"),
				Metadata:  moov.Set(map[string]string{"order": "7"}),
			},
			want: `{"foreignID":"f","metadata":{"order":"7"}}`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := mv2604.NewTransferClient(newTestClient(t, &expectedRequest{
				method:   http.MethodPatch,
				path:     "/accounts/a/transfers/t",
				jsonBody: tc.want,
			}, http.StatusOK, `{"transferID":"t"}`))

			got, err := client.PatchTransfer(t.Context(), "a", "t", tc.update)
			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, "t", got.TransferID)
		})
	}
}

func TestPatchTransferRejectsTruncatedResponse(t *testing.T) {
	client := mv2604.NewTransferClient(newTestClient(t, &expectedRequest{
		method:   http.MethodPatch,
		path:     "/accounts/a/transfers/t",
		jsonBody: `{}`,
	}, http.StatusOK, `{"transferID":`))

	got, err := client.PatchTransfer(t.Context(), "a", "t", mv2604.PatchTransfer{})
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
	require.Nil(t, got)
}
