package mhooks

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func ParseEvent(r *http.Request, secret string) (*Event, error) {
	isValid, err := checkSignature(r.Header, secret)
	if err != nil {
		return nil, fmt.Errorf("checking webhook signature: %w", err)
	}

	if !isValid {
		return nil, errors.New("hashed signature using signing secret does not match value from x-signature header")
	}

	var event Event
	err = json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		return nil, fmt.Errorf("decoding event: %w", err)
	}

	return &event, nil
}

type Event struct {
	EventID   string    `json:"eventID"`
	EventType EventType `json:"type"`
	Data      []byte    `json:"data"`
	CreatedOn time.Time `json:"createdOn"`
}

func (p Event) AccountCreated() (*AccountCreated, error) {
	if p.EventType != EventTypeAccountCreated {
		return nil, newInvalidEventTypeError(p.EventType, EventTypeAccountCreated)
	}

	var accountCreated AccountCreated
	err := json.Unmarshal(p.Data, &accountCreated)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling AccountCreated: %w", err)
	}

	return &accountCreated, nil
}

func (p Event) TransferCreated() (*TransferCreated, error) {
	if p.EventType != EventTypeTransferCreated {
		return nil, newInvalidEventTypeError(p.EventType, EventTypeTransferCreated)
	}

	var transferCreated TransferCreated
	err := json.Unmarshal(p.Data, &transferCreated)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling TransferCreated: %w", err)
	}

	return &transferCreated, nil
}

func newInvalidEventTypeError(expected, got EventType) error {
	return fmt.Errorf("invalid event type: expected %v but got %v", expected, got)
}
