package moov

import (
	"context"
	"net/http"
)

// sourceSystemHeader identifies the core banking system a deposit view document
// originated from.
const sourceSystemHeader = "X-Source-System"

// CreateDepositView submits a deposit view document for the given account. The
// sourceSystem identifies the core banking system the document originated from and
// is sent as the X-Source-System header. The document is sent as the raw request
// body. On success the API responds with the ingest result, returned as V.
func CreateDepositView[V any](ctx context.Context, client *Client, version Version, accountID, sourceSystem string, document []byte) (*V, error) {
	resp, err := client.CallHttp(ctx,
		Endpoint(http.MethodPost, pathDepositView, accountID),
		MoovVersion(version),
		AcceptJson(),
		Header(sourceSystemHeader, sourceSystem),
		BytesBody("application/json", document))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[V](resp)
}
