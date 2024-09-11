package moov

import (
	"context"
	"net/http"
	"time"

	"github.com/moovfinancial/moov-go/pkg/moov/schedules"
)

func (c Client) CreateSchedule(ctx context.Context, accountID string, schedule schedules.UpsertSchedule) (*schedules.Schedule, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathSchedules, accountID),
		AcceptJson(),
		JsonBody(schedule))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[schedules.Schedule](resp)
}

func (c Client) ListSchedule(ctx context.Context, accountID string) ([]schedules.Schedule, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathSchedules, accountID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[schedules.Schedule](resp)
}

func (c Client) GetSchedule(ctx context.Context, accountID string, scheduleID string) (*schedules.Schedule, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathSchedule, accountID, scheduleID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[schedules.Schedule](resp)
}

func (c Client) UpdateSchedule(ctx context.Context, accountID string, scheduleID string, schedule schedules.UpsertSchedule) (*schedules.Schedule, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPut, pathSchedule, accountID, scheduleID),
		AcceptJson(),
		JsonBody(schedule))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[schedules.Schedule](resp)
}

func (c Client) CancelSchedule(ctx context.Context, accountID string, scheduleID string) error {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodDelete, pathSchedule, accountID, scheduleID),
		AcceptJson())
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}

type scheduleOccurrenceFilterArg func() string

func OccurrenceByID(id string) scheduleOccurrenceFilterArg {
	return func() string { return id }
}

func OccurrenceLatest() scheduleOccurrenceFilterArg {
	return func() string { return "latest" }
}

func OccurrenceLatestToTime(t time.Time) scheduleOccurrenceFilterArg {
	return func() string { return t.UTC().Format(time.RFC3339) }
}

func (c Client) GetScheduleOccurrence(ctx context.Context, accountID string, scheduleID string, filter scheduleOccurrenceFilterArg) (*schedules.TransferOccurrence, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathScheduleOccurrence, accountID, scheduleID, filter()),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[schedules.TransferOccurrence](resp)
}
