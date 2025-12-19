package moov

import (
	"context"
	"net/http"
)

// ShareConnection creates a new connection on the `subjectAccountID` with the `principalAccountID` a connection with the provided configuration.
func (c Client) ShareConnection(ctx context.Context, subjectAccountID string, connection ShareConnectionRequest) (*ShareConnectionResponse, error) {
	httpResp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathConnections, subjectAccountID),
		AcceptJson(),
		JsonBody(connection))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[ShareConnectionResponse](httpResp)
}
