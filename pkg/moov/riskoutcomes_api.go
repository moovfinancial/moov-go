package moov

import (
	"context"
	"net/http"
)

// GetTransferRiskOutcomes retrieves the risk rules that contributed to a
// transfer's risk decision. Availability is limited to enrolled partners; the
// account is resolved from the calling credentials.
// https://docs.moov.io/api/money-movement/transfers/get-risk-outcomes/
func (c Client) GetTransferRiskOutcomes(ctx context.Context, transferID string) (*TransferRiskOutcomes, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathTransferRiskOutcomes, transferID),
		AcceptJson(),
		MoovVersion(Version2026_10),
	)
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[TransferRiskOutcomes](resp)
}
