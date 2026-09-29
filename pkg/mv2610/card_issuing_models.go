package mv2610

import (
	"time"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

type IssuedCardActivity struct {
	AuthorizationID         *string                            `json:"authorizationID,omitempty"`
	CardTransactionID       *string                            `json:"cardTransactionID,omitempty"`
	Status                  moov.IssuedCardAuthorizationStatus `json:"status"`
	IssuedCardID            string                             `json:"issuedCardID"`
	LastFourCardNumber      *string                            `json:"lastFourCardNumber,omitempty"`
	AuthorizedUserAccountID *string                            `json:"authorizedUserAccountID,omitempty"`
	AuthorizedAmount        *string                            `json:"authorizedAmount,omitempty"`
	ClearedAmount           *string                            `json:"clearedAmount,omitempty"`
	DeclineReason           *moov.IssuingDeclineReason         `json:"declineReason,omitempty"`
	CreatedOn               time.Time                          `json:"createdOn"`
	MerchantData            moov.IssuedCardTransactionMerchant `json:"merchantData"`
}
