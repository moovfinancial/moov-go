package moov

import (
	"context"
	"net/http"
	"strconv"
)

func (c Client) CreateTicket(ctx context.Context, accountID string, create CreateTicket) (*Ticket, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathTickets, accountID),
		MoovVersion(Version2025_07),
		AcceptJson(),
		JsonBody(create))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Ticket](resp)
}

// Func that applies a filter and returns an error if validation fails
type ListTicketFilter callArg

// WithTicketCursor if provided, this query will get the next page of results
func WithTicketCursor(cursor string) ListTicketFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["cursor"] = cursor
		return nil
	})
}

// WithTicketCount value to limit the number of results in the query. Default is 200, maximum is 1000
func WithTicketCount(count int) ListTicketFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["count"] = strconv.Itoa(count)
		return nil
	})
}

// WithTicketStatus possible values: new, in-progress, on-hold, closed
func WithTicketStatus(status TicketStatus) ListTicketFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["status"] = string(status)
		return nil
	})
}

// WithTicketForeignID selects tickets with the specified foreign ID
func WithTicketForeignID(foreignID string) ListTicketFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["foreignID"] = foreignID
		return nil
	})
}

func (c Client) ListTickets(ctx context.Context, accountID string, filters ...ListTicketFilter) (*ListTicket, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathTickets, accountID),
		prependArgs(filters, MoovVersion(Version2025_07), AcceptJson())...)
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[ListTicket](resp)
}

func (c Client) ListTicketMessages(ctx context.Context, accountID, ticketID string) ([]TicketMessage, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathTicketMessages, accountID, ticketID),
		MoovVersion(Version2025_07),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[TicketMessage](resp)
}

func (c Client) GetTicket(ctx context.Context, accountID string, ticketID string) (*Ticket, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathTicket, accountID, ticketID),
		MoovVersion(Version2025_07),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Ticket](resp)
}

func (c Client) UpdateTicket(ctx context.Context, accountID string, ticketID string, update UpdateTicket) (*Ticket, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPatch, pathTicket, accountID, ticketID),
		MoovVersion(Version2025_07),
		AcceptJson(),
		JsonBody(update))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Ticket](resp)
}
