package mhooks

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

func NewEvent(requestBody io.Reader) (*Event, error) {
	var event Event
	err := json.NewDecoder(requestBody).Decode(&event)
	if err != nil {
		return nil, fmt.Errorf("decoding event: %w", err)
	}

	return &event, nil
}

type Event struct {
	EventID   string          `json:"eventID"`
	EventType EventType       `json:"type"`
	Data      json.RawMessage `json:"data"`
	CreatedOn time.Time       `json:"createdOn"`
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
