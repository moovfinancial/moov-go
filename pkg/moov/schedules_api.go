package moov

import (
	"context"
	"net/http"

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

	switch resp.Status() {
	case StatusCompleted:
		return CompletedObjectOrError[schedules.Schedule](resp)
	default:
		return nil, resp
	}
}

func (c Client) ListSchedule(ctx context.Context, accountID string) ([]schedules.Schedule, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathSchedules, accountID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[schedules.Schedule](resp)
}

func (c Client) GetSchedule(ctx context.Context, accountID string, scheduleID string) (*schedules.Schedule, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathSchedule, accountID, scheduleID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[schedules.Schedule](resp)
}

func (c Client) UpdateSchedule(ctx context.Context, accountID string, schedule schedules.UpsertSchedule) (*schedules.Schedule, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathSchedules, accountID),
		AcceptJson(),
		JsonBody(schedule))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[schedules.Schedule](resp)
}
