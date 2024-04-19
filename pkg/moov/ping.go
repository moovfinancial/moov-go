package moov

import (
	"context"
	"net/http"
)

// Ping calls the ping endpoint to make sure we have valid credentials
func (c Client) Ping(ctx context.Context) error {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, "/ping"))
	if err != nil {
		return err
	}

	switch resp.Status() {
	case StatusCompleted:
		return nil
	default:
		return resp
	}
}
