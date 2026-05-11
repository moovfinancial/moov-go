package moov

import (
	"context"
	"net/http"
)

// CardMetadata describes card BIN attributes and push/pull capabilities returned
// by the card metadata lookup endpoint without linking the card.
type CardMetadata struct {
	Bin                  string `json:"bin,omitempty"`
	Brand                string `json:"brand,omitempty"`
	CardCategory         string `json:"cardCategory,omitempty"`
	CardType             string `json:"cardType,omitempty"`
	Commercial           *bool  `json:"commercial,omitempty"`
	Regulated            *bool  `json:"regulated,omitempty"`
	Issuer               string `json:"issuer,omitempty"`
	IssuerCountry        string `json:"issuerCountry,omitempty"`
	IssuerPhone          string `json:"issuerPhone,omitempty"`
	IssuerURL            string `json:"issuerURL,omitempty"`
	DomesticPullFromCard string `json:"domesticPullFromCard,omitempty"`
	DomesticPushToCard   string `json:"domesticPushToCard,omitempty"`
}

// CardMetadataRequest is the request body for GetCardMetadata. Provide either CardNumber
// or EndToEndToken — exactly one must be populated.
type CardMetadataRequest struct {
	CardNumber string `json:"cardNumber,omitempty"`

	EndToEndToken *EndToEndToken `json:"e2ee,omitempty"`
}

// GetCardMetadata returns BIN attributes and push/pull capabilities for a card identified
// by its full PAN, without linking the card. Provide the PAN via CardNumber or via
// EndToEndToken (encrypted) — exactly one must be populated.
// https://docs.moov.io/api/#tag/Cards/operation/getCardMetadata
func (c Client) GetCardMetadata(ctx context.Context, request CardMetadataRequest) (*CardMetadata, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodPost, pathCardMetadata), AcceptJson(), JsonBody(request))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[CardMetadata](resp)
}
