package moov

import (
	"context"
	"net/http"
	"time"
)

// Guide: https://docs.moov.io/guides/money-movement/scheduling/
// Documentation: https://docs.moov.io/api/money-movement/schedules/create/
func (c Client) CreateSchedule(ctx context.Context, accountID string, schedule CreateSchedule) (*Schedule, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathSchedules, accountID),
		AcceptJson(),
		JsonBody(schedule))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Schedule](resp)
}

// Guide: https://docs.moov.io/guides/money-movement/scheduling/
// Documentation: https://docs.moov.io/api/money-movement/schedules/list/
func (c Client) ListSchedule(ctx context.Context, accountID string, args ...callArg) ([]Schedule, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathSchedules, accountID),
		append(args, AcceptJson())...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[Schedule](resp)
}

// Guide: https://docs.moov.io/guides/money-movement/scheduling/
// Documentation: https://docs.moov.io/api/money-movement/schedules/get/
func (c Client) GetSchedule(ctx context.Context, accountID string, scheduleID string) (*Schedule, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathSchedule, accountID, scheduleID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Schedule](resp)
}

// Guide: https://docs.moov.io/guides/money-movement/scheduling/
// Documentation: https://docs.moov.io/api/money-movement/schedules/update/
func (c Client) UpdateSchedule(ctx context.Context, accountID string, scheduleID string, schedule UpdateSchedule) (*Schedule, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPut, pathSchedule, accountID, scheduleID),
		AcceptJson(),
		JsonBody(schedule))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Schedule](resp)
}

// Guide: https://docs.moov.io/guides/money-movement/scheduling/
// Documentation: https://docs.moov.io/api/money-movement/schedules/delete/
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

// Occurrence with the specific ID
func OccurrenceByID(id string) scheduleOccurrenceFilterArg {
	return func() string { return id }
}

// Occurrence closest to now without going over.
func OccurrenceLatest() scheduleOccurrenceFilterArg {
	return func() string { return "latest" }
}

// Occurrence closest to time `t` without going over.
func OccurrenceLatestToTime(t time.Time) scheduleOccurrenceFilterArg {
	return func() string { return t.UTC().Format(time.RFC3339) }
}

// Guide: https://docs.moov.io/guides/money-movement/scheduling/
func (c Client) GetScheduleOccurrence(ctx context.Context, accountID string, scheduleID string, filter scheduleOccurrenceFilterArg) (*Occurrence, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathScheduleOccurrence, accountID, scheduleID, filter()),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Occurrence](resp)
}
