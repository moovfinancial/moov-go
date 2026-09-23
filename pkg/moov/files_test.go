package moov_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

// newFilesTestClient points a client at srv. The test partner's connection is not
// granted any files scopes, so the live suite cannot reach these endpoints.
func newFilesTestClient(t *testing.T, srv *httptest.Server) *moov.Client {
	t.Helper()

	client, err := moov.NewClient(
		moov.WithCredentials(moov.Credentials{PublicKey: "pk", SecretKey: "sk"}),
		moov.WithMoovURLScheme("http"),
	)
	require.NoError(t, err)
	client.Credentials.Host = strings.TrimPrefix(srv.URL, "http://")

	return client
}

func TestDownloadFile(t *testing.T) {
	contents := []byte("account,balance\nmoov,100\n")

	var (
		method  string
		path    string
		version string
		accept  string
	)

	// The handler runs on its own goroutine, so it records what it saw rather than
	// asserting: require.* calls t.FailNow, which is only valid on the test goroutine.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		path = r.URL.Path
		version = r.Header.Get(moov.VersionHeader)
		accept = r.Header.Get("Accept")

		// A name holding a space and a quote, which the service escapes with
		// mime.FormatMediaType. Splitting on "filename=" would mangle it.
		w.Header().Set("Content-Disposition", `attachment; filename="bank \"statement\".csv"`)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write(contents)
	}))
	t.Cleanup(srv.Close)

	actual, err := newFilesTestClient(t, srv).DownloadFile(BgCtx(), "account-123", "file-456")
	require.NoError(t, err)

	require.Equal(t, http.MethodGet, method)
	require.Equal(t, "/accounts/account-123/files/file-456/contents", path)
	require.Equal(t, "*/*", accept)

	// The service registers this route only at 2026.10 and returns a bare 404
	// without the header, so the version is part of the contract.
	require.Equal(t, moov.Version2026_10.String(), version)

	require.Equal(t, contents, actual.Data)
	require.Equal(t, `bank "statement".csv`, actual.FileName)

	// Detected from the contents on upload rather than taken from the uploader,
	// so a csv reports as text/plain.
	require.Equal(t, "text/plain; charset=utf-8", actual.ContentType)
}

func TestDownloadFile_FileNameUnavailable(t *testing.T) {
	tests := []struct {
		name        string
		disposition string
	}{
		{name: "header missing", disposition: ""},
		{name: "header unparsable", disposition: "attachment; filename"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.disposition != "" {
					w.Header().Set("Content-Disposition", tt.disposition)
				}
				w.Header().Set("Content-Type", "application/pdf")
				_, _ = w.Write([]byte("%PDF-1.4"))
			}))
			t.Cleanup(srv.Close)

			// A name the caller cannot use is worth less than an error here, since
			// the bytes are still what was asked for.
			actual, err := newFilesTestClient(t, srv).DownloadFile(BgCtx(), "account-123", "file-456")
			require.NoError(t, err)
			require.Empty(t, actual.FileName)
			require.Equal(t, "application/pdf", actual.ContentType)
			require.Equal(t, []byte("%PDF-1.4"), actual.Data)
		})
	}
}

func TestDownloadFile_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(srv.Close)

	actual, err := newFilesTestClient(t, srv).DownloadFile(BgCtx(), "account-123", "file-456")
	require.Nil(t, actual)

	var httpErr moov.HttpCallResponse
	require.ErrorAs(t, err, &httpErr)
	require.Equal(t, moov.StatusNotFound, httpErr.Status())
}
