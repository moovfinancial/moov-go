package moov

import (
	"bytes"
	"context"
	"net/http"
)

// GetMerchantProcessingAgreement downloads the merchant processing agreement as PDF bytes
// https://docs.moov.io/api/moov-accounts/merchant-agreement/get/
func (c Client) GetMerchantProcessingAgreement(ctx context.Context, accountID string) ([]byte, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathMerchantAgreement, accountID),
		AcceptContentType("application/pdf"),
	)
	if err != nil {
		return []byte{}, err
	}

	buf, err := CompletedObjectOrError[bytes.Buffer](resp)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}
