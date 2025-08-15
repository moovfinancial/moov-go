package mhooks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ParseEvent returns a webhook event with the hydrated payload.
//
// Access the event payload by calling the corresponding getter method.
// For example, if the event.EventType=account.created, call event.AccountCreated().
//
// The webhook signature is verified using the provided signing secret.
func ParseEvent(r *http.Request, secret string) (*Event, error) {
	isValid, err := checkSignature(r.Header, secret)
	if err != nil {
		return nil, fmt.Errorf("checking webhook signature: %w", err)
	}

	if !isValid {
		return nil, ErrInvalidSignature
	}

	var event Event
	err = json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		return nil, fmt.Errorf("decoding event: %w", err)
	}

	var eventData any
	switch event.EventType {
	case EventTypeAccountCreated:
		eventData = &event.accountCreated
	case EventTypeAccountDisconnected:
		eventData = &event.accountDeleted
	case EventTypeAccountUpdated:
		eventData = &event.accountUpdated
	case EventTypeBalanceUpdated:
		eventData = &event.balanceUpdated
	case EventTypeBankAccountCreated:
		eventData = &event.bankAccountCreated
	case EventTypeBankAccountDeleted:
		eventData = &event.bankAccountDeleted
	case EventTypeBankAccountUpdated:
		eventData = &event.bankAccountUpdated
	case EventTypeCancellationCreated:
		eventData = &event.cancellationCreated
	case EventTypeCancellationUpdated:
		eventData = &event.cancellationUpdated
	case EventTypeCardAutoUpdated:
		eventData = &event.cardAutoUpdated
	case EventTypeCapabilityRequested:
		eventData = &event.capabilityRequested
	case EventTypeCapabilityUpdated:
		eventData = &event.capabilityUpdated
	case EventTypeDisputeCreated:
		eventData = &event.disputeCreated
	case EventTypeDisputeUpdated:
		eventData = &event.disputeUpdated
	case EventTypeNetworkIDUpdated:
		eventData = &event.networkIDUpdated
	case EventTypePaymentMethodDisabled:
		eventData = &event.paymentMethodDisabled
	case EventTypePaymentMethodEnabled:
		eventData = &event.paymentMethodEnabled
	case EventTypeRefundCreated:
		eventData = &event.refundCreated
	case EventTypeRefundUpdated:
		eventData = &event.refundUpdated
	case EventTypeRepresentativeCreated:
		eventData = &event.representativeCreated
	case EventTypeRepresentativeDeleted:
		eventData = &event.representativeDeleted
	case EventTypeRepresentativeUpdated:
		eventData = &event.representativeUpdated
	case EventTypeSweepCreated:
		eventData = &event.sweepCreated
	case EventTypeSweepUpdated:
		eventData = &event.sweepUpdated
	case EventTypeTestPing:
		eventData = &event.testPing
	case EventTypeTransferCreated:
		eventData = &event.transferCreated
	case EventTypeTransferUpdated:
		eventData = &event.transferUpdated
	case EventTypeWalletTransactionUpdated:
		eventData = &event.walletTransactionUpdated
	case EventTypeBillingStatementCreated:
		eventData = &event.billingStatementCreated
	default:
		return nil, fmt.Errorf("invalid event type: %v", event.EventType)
	}

	err = json.Unmarshal(event.Data, eventData)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling event data: %w", err)
	}

	return &event, nil
}

type Event struct {
	EventID   string          `json:"eventID"`
	EventType EventType       `json:"type"`
	CreatedOn time.Time       `json:"createdOn"`
	Data      json.RawMessage `json:"data"`

	accountCreated           *AccountCreated
	accountDeleted           *AccountDisconnected
	accountUpdated           *AccountUpdated
	balanceUpdated           *BalanceUpdated
	bankAccountCreated       *BankAccountCreated
	bankAccountDeleted       *BankAccountDeleted
	bankAccountUpdated       *BankAccountUpdated
	cancellationCreated      *CancellationCreated
	cancellationUpdated      *CancellationUpdated
	cardAutoUpdated          *CardAutoUpdated
	capabilityRequested      *CapabilityRequested
	capabilityUpdated        *CapabilityUpdated
	disputeCreated           *DisputeCreated
	disputeUpdated           *DisputeUpdated
	networkIDUpdated         *NetworkIDUpdated
	paymentMethodDisabled    *PaymentMethodDisabled
	paymentMethodEnabled     *PaymentMethodEnabled
	refundCreated            *RefundCreated
	refundUpdated            *RefundUpdated
	representativeCreated    *RepresentativeCreated
	representativeDeleted    *RepresentativeDeleted
	representativeUpdated    *RepresentativeUpdated
	sweepCreated             *SweepCreated
	sweepUpdated             *SweepUpdated
	testPing                 *TestPing
	transferCreated          *TransferCreated
	transferUpdated          *TransferUpdated
	walletTransactionUpdated *WalletTransactionUpdated
	billingStatementCreated  *BillingStatementCreated
}

func (e Event) AccountCreated() (*AccountCreated, error) {
	if e.EventType != EventTypeAccountCreated {
		return nil, newInvalidEventTypeError(EventTypeAccountCreated, e.EventType)
	}

	return e.accountCreated, nil
}

func (e Event) AccountDisconnected() (*AccountDisconnected, error) {
	if e.EventType != EventTypeAccountDisconnected {
		return nil, newInvalidEventTypeError(EventTypeAccountDisconnected, e.EventType)
	}

	return e.accountDeleted, nil
}

func (e Event) AccountUpdated() (*AccountUpdated, error) {
	if e.EventType != EventTypeAccountUpdated {
		return nil, newInvalidEventTypeError(EventTypeAccountUpdated, e.EventType)
	}

	return e.accountUpdated, nil
}

func (e Event) BalanceUpdated() (*BalanceUpdated, error) {
	if e.EventType != EventTypeBalanceUpdated {
		return nil, newInvalidEventTypeError(EventTypeBalanceUpdated, e.EventType)
	}

	return e.balanceUpdated, nil
}

func (e Event) BankAccountCreated() (*BankAccountCreated, error) {
	if e.EventType != EventTypeBankAccountCreated {
		return nil, newInvalidEventTypeError(EventTypeBankAccountCreated, e.EventType)
	}

	return e.bankAccountCreated, nil
}

func (e Event) BankAccountDeleted() (*BankAccountDeleted, error) {
	if e.EventType != EventTypeBankAccountDeleted {
		return nil, newInvalidEventTypeError(EventTypeBankAccountDeleted, e.EventType)
	}

	return e.bankAccountDeleted, nil
}

func (e Event) BankAccountUpdated() (*BankAccountUpdated, error) {
	if e.EventType != EventTypeBankAccountUpdated {
		return nil, newInvalidEventTypeError(EventTypeBankAccountUpdated, e.EventType)
	}

	return e.bankAccountUpdated, nil
}

func (e Event) CancellationCreated() (*CancellationCreated, error) {
	if e.EventType != EventTypeCancellationCreated {
		return nil, newInvalidEventTypeError(EventTypeCancellationCreated, e.EventType)
	}

	return e.cancellationCreated, nil
}

func (e Event) CancellationUpdated() (*CancellationUpdated, error) {
	if e.EventType != EventTypeCancellationUpdated {
		return nil, newInvalidEventTypeError(EventTypeCancellationUpdated, e.EventType)
	}

	return e.cancellationUpdated, nil
}

func (e Event) CardAutoUpdated() (*CardAutoUpdated, error) {
	if e.EventType != EventTypeCardAutoUpdated {
		return nil, newInvalidEventTypeError(EventTypeCardAutoUpdated, e.EventType)
	}

	return e.cardAutoUpdated, nil
}

func (e Event) CapabilityRequested() (*CapabilityRequested, error) {
	if e.EventType != EventTypeCapabilityRequested {
		return nil, newInvalidEventTypeError(EventTypeCapabilityRequested, e.EventType)
	}

	return e.capabilityRequested, nil
}

func (e Event) CapabilityUpdated() (*CapabilityUpdated, error) {
	if e.EventType != EventTypeCapabilityUpdated {
		return nil, newInvalidEventTypeError(EventTypeCapabilityUpdated, e.EventType)
	}

	return e.capabilityUpdated, nil
}

func (e Event) DisputeCreated() (*DisputeCreated, error) {
	if e.EventType != EventTypeDisputeCreated {
		return nil, newInvalidEventTypeError(EventTypeDisputeCreated, e.EventType)
	}

	return e.disputeCreated, nil
}

func (e Event) DisputeUpdated() (*DisputeUpdated, error) {
	if e.EventType != EventTypeDisputeUpdated {
		return nil, newInvalidEventTypeError(EventTypeDisputeUpdated, e.EventType)
	}

	return e.disputeUpdated, nil
}

func (e Event) NetworkIDUpdated() (*NetworkIDUpdated, error) {
	if e.EventType != EventTypeNetworkIDUpdated {
		return nil, newInvalidEventTypeError(EventTypeNetworkIDUpdated, e.EventType)
	}

	return e.networkIDUpdated, nil
}

func (e Event) PaymentMethodDisabled() (*PaymentMethodDisabled, error) {
	if e.EventType != EventTypePaymentMethodDisabled {
		return nil, newInvalidEventTypeError(EventTypePaymentMethodDisabled, e.EventType)
	}

	return e.paymentMethodDisabled, nil
}

func (e Event) PaymentMethodEnabled() (*PaymentMethodEnabled, error) {
	if e.EventType != EventTypePaymentMethodEnabled {
		return nil, newInvalidEventTypeError(EventTypePaymentMethodEnabled, e.EventType)
	}

	return e.paymentMethodEnabled, nil
}

func (e Event) RefundCreated() (*RefundCreated, error) {
	if e.EventType != EventTypeRefundCreated {
		return nil, newInvalidEventTypeError(EventTypeRefundCreated, e.EventType)
	}

	return e.refundCreated, nil
}

func (e Event) RefundUpdated() (*RefundUpdated, error) {
	if e.EventType != EventTypeRefundUpdated {
		return nil, newInvalidEventTypeError(EventTypeRefundUpdated, e.EventType)
	}

	return e.refundUpdated, nil
}

func (e Event) RepresentativeCreated() (*RepresentativeCreated, error) {
	if e.EventType != EventTypeRepresentativeCreated {
		return nil, newInvalidEventTypeError(EventTypeRepresentativeCreated, e.EventType)
	}

	return e.representativeCreated, nil
}

func (e Event) RepresentativeDeleted() (*RepresentativeDeleted, error) {
	if e.EventType != EventTypeRepresentativeDeleted {
		return nil, newInvalidEventTypeError(EventTypeRepresentativeDeleted, e.EventType)
	}

	return e.representativeDeleted, nil
}

func (e Event) RepresentativeUpdated() (*RepresentativeUpdated, error) {
	if e.EventType != EventTypeRepresentativeUpdated {
		return nil, newInvalidEventTypeError(EventTypeRepresentativeUpdated, e.EventType)
	}

	return e.representativeUpdated, nil
}

func (e Event) SweepCreated() (*SweepCreated, error) {
	if e.EventType != EventTypeSweepCreated {
		return nil, newInvalidEventTypeError(EventTypeSweepCreated, e.EventType)
	}

	return e.sweepCreated, nil
}

func (e Event) SweepUpdated() (*SweepUpdated, error) {
	if e.EventType != EventTypeSweepUpdated {
		return nil, newInvalidEventTypeError(EventTypeSweepUpdated, e.EventType)
	}

	return e.sweepUpdated, nil
}

func (e Event) TestPing() (*TestPing, error) {
	if e.EventType != EventTypeTestPing {
		return nil, newInvalidEventTypeError(EventTypeTestPing, e.EventType)
	}

	return e.testPing, nil
}

func (e Event) TransferCreated() (*TransferCreated, error) {
	if e.EventType != EventTypeTransferCreated {
		return nil, newInvalidEventTypeError(EventTypeTransferCreated, e.EventType)
	}

	return e.transferCreated, nil
}

func (e Event) TransferUpdated() (*TransferUpdated, error) {
	if e.EventType != EventTypeTransferUpdated {
		return nil, newInvalidEventTypeError(EventTypeTransferUpdated, e.EventType)
	}

	return e.transferUpdated, nil
}

func (e Event) WalletTransactionUpdated() (*WalletTransactionUpdated, error) {
	if e.EventType != EventTypeWalletTransactionUpdated {
		return nil, newInvalidEventTypeError(EventTypeWalletTransactionUpdated, e.EventType)
	}

	return e.walletTransactionUpdated, nil
}

func (e Event) BillingStatementCreated() (*BillingStatementCreated, error) {
	if e.EventType != EventTypeBillingStatementCreated {
		return nil, newInvalidEventTypeError(EventTypeBillingStatementCreated, e.EventType)
	}

	return e.billingStatementCreated, nil
}

func newInvalidEventTypeError(expected, got EventType) error {
	return fmt.Errorf("invalid event type: expected %v but got %v", expected, got)
}
