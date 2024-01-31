package moov

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Dispute struct {
	DisputeID                string              `json:"disputeID,omitempty"`
	CreatedOn                time.Time           `json:"createdOn,omitempty"`
	Amount                   Amount              `json:"amount,omitempty"`
	NetworkReasonCode        string              `json:"networkReasonCode,omitempty"`
	NetworkReasonDescription string              `json:"networkReasonDescription,omitempty"`
	RespondBy                time.Time           `json:"respondBy,omitempty"`
	Status                   string              `json:"status,omitempty"`
	Transfer                 SynchronousTransfer `json:"transfer,omitempty"`
}

type DisputeListFilter callArg

func WithDisputeCount(c int) callArg {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["count"] = fmt.Sprintf("%d", c)
		return nil
	})
}

func WithDisputeSkip(c int) callArg {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["skip"] = fmt.Sprintf("%d", c)
		return nil
	})
}

func WithDisputeResponseStartDate(t time.Time) callArg {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["respondStartDateTime"] = t.Format(time.RFC3339)
		return nil
	})
}

func WithDisputeResponseEndDate(t time.Time) callArg {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["respondEndDateTime"] = t.Format(time.RFC3339)
		return nil
	})
}

func WithDisputeStatus(s string) callArg {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["status"] = s
		return nil
	})
}

func WithDisputeMerchantAccountID(id string) callArg {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["merchantAccountID"] = id
		return nil
	})
}

func WithDisputeCardHolderAccountID(id string) callArg {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["cardholderAccountID"] = id
		return nil
	})
}

func WithDisputeStartDate(t time.Time) callArg {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["startDateTime"] = t.Format(time.RFC3339)
		return nil
	})
}

func WithDisputeEndDate(t time.Time) callArg {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["endDateTime"] = t.Format(time.RFC3339)
		return nil
	})
}

func WithDisputeOrderBy(orderBy string) callArg {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["orderBy"] = orderBy
		return nil
	})
}

// ListDisputes lists of Disputes that are associated with a Moov account
// https://docs.moov.io/api/money-movement/disputes/list/
func (c Client) ListDisputes(ctx context.Context, filters ...DisputeListFilter) ([]Dispute, error) {
	args := prependArgs(filters, AcceptJson())
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathDisputes), args...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[Dispute](resp)
}

// GetDispute retrieves a dispute for the given dispute id
// https://docs.moov.io/api/money-movement/disputes/get/
func (c Client) GetDispute(ctx context.Context, disputeID string) (*Dispute, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathDisputeID, disputeID), AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Dispute](resp)
}
