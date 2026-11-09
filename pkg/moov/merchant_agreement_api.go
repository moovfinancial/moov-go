package moov

import (
	"bytes"
	"context"
	"net/http"
)

// GetMerchantProcessingAgreement downloads the merchant processing agreement as PDF bytes
// https://docs.moov.io/api/moov-accounts/merchant-agreement/get/
func (c Client) GetMerchantProcessingAgreement(ctx context.Context, accountID string) (MerchantProcessingAgreement, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathMerchantAgreement, accountID),
		AcceptContentType("application/pdf"),
	)
	if err != nil {
		return nil, err
	}

	buf, err := CompletedObjectOrError[bytes.Buffer](resp)
	if err != nil {
		return nil, err
	}
	return MerchantProcessingAgreement(buf.Bytes()), nil
}
