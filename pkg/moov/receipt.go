package moov

import (
	"context"
	"net/http"
)

func (c Client) CreateReceipt(ctx context.Context, receipts ...CreateReceipt) ([]Receipt, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathReceipts),
		AcceptJson(),
		JsonBody(receipts))
	if err != nil {
		return nil, err
	}

	if resp.Status() == StatusStarted {
		return UnmarshalListResponse[Receipt](resp)
	} else {
		return nil, resp
	}
}

type ListReceiptsFilter callArg

func ReceiptByScheduleID(id string) ListReceiptsFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["id"] = id
		return nil
	})
}

func ReceiptByTransferID(id string) ListReceiptsFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["id"] = id
		return nil
	})
}

func (c Client) ListReceipts(ctx context.Context, filters ...ListReceiptsFilter) ([]Receipt, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathReceipts),
		prependArgs(filters, AcceptJson())...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[Receipt](resp)
}

func (c Client) DeleteReceipt(ctx context.Context, receiptID string) error {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodDelete, pathReceipt, receiptID),
		AcceptJson())
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}
