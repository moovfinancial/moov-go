package mv2610

import (
	"time"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

// IssuedCardActivity is an authorization or a card transaction. Exactly one of AuthorizationID and CardTransactionID is set.
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

// CreateAuthorizationSimulation is the request to simulate an authorization on an issued card in test mode.
type CreateAuthorizationSimulation struct {
	IssuedCardID string `json:"issuedCardID"`
	// Decimal-formatted amount, such as "12.34".
	Amount       string                              `json:"amount"`
	MerchantData *moov.IssuedCardTransactionMerchant `json:"merchantData,omitempty"`
}

// CreateClearingSimulation is the request to simulate a clearing on an issued card authorization in test mode.
type CreateClearingSimulation struct {
	// Decimal-formatted amount, such as "12.34".
	Amount *string `json:"amount,omitempty"`
}

// AuthorizationSimulationAsyncResponse is returned when a simulation did not complete before the request timed out.
type AuthorizationSimulationAsyncResponse struct {
	AuthorizationID string `json:"authorizationID"`
}
