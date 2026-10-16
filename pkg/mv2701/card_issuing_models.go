package mv2701

import "github.com/moovfinancial/moov-go/pkg/moov"

// CreateIssuedCard is the request to issue a card. v2027.01.00 removed expiration because issuing programs set it.
// Use Controls.ExpiresOn to limit the card lifetime.
type CreateIssuedCard struct {
	AuthorizedUserAccountID *string               `json:"authorizedUserAccountID,omitempty"`
	Nickname                *string               `json:"nickname,omitempty"`
	Metadata                map[string]string     `json:"metadata,omitempty"`
	BillingAddress          *moov.Address         `json:"billingAddress,omitempty"`
	Controls                *moov.IssuingControls `json:"controls,omitempty"`
}
