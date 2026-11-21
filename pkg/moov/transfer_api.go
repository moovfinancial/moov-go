package moov

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CreateTransferArgs func(t *createTransferBuilder) callArg
type createTransferBuilder struct {
	idempotencyKey string
}

// Can be specified to overwrite a randomly generated one.
func WithTransferIdempotencyKey(key uuid.UUID) CreateTransferArgs {
	return func(t *createTransferBuilder) callArg {
		t.idempotencyKey = key.String()
		return IdempotencyKey(t.idempotencyKey)
	}
}

// CreateTransfer creates a new transfer
// https://docs.moov.io/api/index.html#tag/Transfers/operation/createTransfer
func (c Client) CreateTransfer(ctx context.Context, partnerAccountID string, transfer CreateTransfer, options ...CreateTransferArgs) CreateTransferBuilder {
	builder := &createTransferBuilder{}
	callArgs := []callArg{
		AcceptJson(),
		JsonBody(transfer),
		WithTransferIdempotencyKey(uuid.New())(builder),
	}

	for _, opt := range options {
		callArgs = append(callArgs, opt(builder))
	}

	return CreateTransferBuilder{
		client:   c,
		ctx:      ctx,
		endpoint: Endpoint(http.MethodPost, pathTransfers, partnerAccountID),
		callArgs: callArgs,
	}
}

type CreateTransferBuilder struct {
	client   Client
	ctx      context.Context
	endpoint EndpointArg
	callArgs []callArg
}

// Started initiates the transfers request and doesn't wait beyond creating the transfer
func (r CreateTransferBuilder) Started() (*TransferStarted, error) {
	resp, err := r.client.CallHttp(r.ctx, r.endpoint, r.callArgs...)
	if err != nil {
		return nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		st, err := UnmarshalObjectResponse[TransferStarted](resp)
		return st, err
	case StatusStateConflict:
		return nil, errors.Join(ErrXIdempotencyKey, resp)
	default:
		return nil, resp
	}
}

// Starts a transfer request and waits for a response from the rail (e.g. authorized or declined) before returning the result.
// There are three possible return values (and only one will be not-nil):
// 1) A full transfer with rail-specific details as a result of waiting for the response from the rail.
// 2) A transfer that started but the request timed out waiting for a response from the rail.
// 3) An error attempting to create the transfer.
func (r CreateTransferBuilder) WaitForRailResponse() (*Transfer, *TransferStarted, error) {
	resp, err := r.client.CallHttp(r.ctx, r.endpoint, append(r.callArgs, WaitFor("rail-response"))...)
	if err != nil {
		return nil, nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		transfer, err := UnmarshalObjectResponse[Transfer](resp)
		return transfer, nil, err
	case StatusStarted:
		transferStarted, err := UnmarshalObjectResponse[TransferStarted](resp)
		return nil, transferStarted, err
	case StatusStateConflict:
		return nil, nil, errors.Join(ErrXIdempotencyKey, resp)
	default:
		return nil, nil, resp
	}
}

type ListTransferFilter callArg

func WithTransferAccountIDs(accountIDs []string) ListTransferFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["accountIDs"] = strings.Join(accountIDs, ",")
		return nil
	})
}

func WithTransferStatus(status string) ListTransferFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["status"] = status
		return nil
	})
}

func WithTransferStartDate(start time.Time) ListTransferFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["startDateTime"] = start.Format(time.RFC3339)
		return nil
	})
}

func WithTransferEndDate(end time.Time) ListTransferFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["endDateTime"] = end.Format(time.RFC3339)
		return nil
	})
}

func WithTransferGroup(groupID string) ListTransferFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["groupID"] = groupID
		return nil
	})
}

func WithTransferSchedule(scheduleID string) ListTransferFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["scheduleID"] = scheduleID
		return nil
	})
}

func WithTransferPaymentLinkCode(paymentLinkCode string) ListTransferFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["paymentLinkCode"] = paymentLinkCode
		return nil
	})
}

func WithTransferRefunded() ListTransferFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["refunded"] = "true"
		return nil
	})
}

func WithTransferDisputed() ListTransferFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["disputed"] = "true"
		return nil
	})
}

// WithTransferTypes filters by transfer type, such as "card-payment". The filter starts in v2026.10.00.
func WithTransferTypes(transferTypes []string) ListTransferFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["transferTypes"] = strings.Join(transferTypes, ",")
		return nil
	})
}

func WithTransferSkip(skip int) ListTransferFilter {
	return Skip(skip)
}

func WithTransferCount(count int) ListTransferFilter {
	return Count(count)
}

// ListTransfers lists all transfers
// https://docs.moov.io/api/index.html#tag/Transfers/operation/listTransfers
func (c Client) ListTransfers(ctx context.Context, accountID string, filters ...ListTransferFilter) ([]Transfer, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathTransfers, accountID),
		prependArgs(filters, AcceptJson())...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[Transfer](resp)
}

// GetTransfer retrieves a transfer
// https://docs.moov.io/api/index.html#tag/Transfers/operation/getTransfer
func (c Client) GetTransfer(ctx context.Context, accountID, transferID string) (*Transfer, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathTransfer, accountID, transferID),
		AcceptJson(),
	)
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Transfer](resp)
}

// FetchTransfers retrieves multiple transfers by their IDs
// https://docs.moov.io/api/money-movement/transfers/batch-get-transfers/
func (c Client) FetchTransfers(ctx context.Context, accountID string, request FetchTransfersRequest) (*FetchTransfersResponse, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathTransfersFetch, accountID),
		AcceptJson(),
		JsonBody(request))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[FetchTransfersResponse](resp)
}

type TransferPatcher func(patch *patchTransfer)

func PatchTransferMetadata(metadata map[string]string) TransferPatcher {
	return func(patch *patchTransfer) {
		patch.Metadata = metadata
	}
}

// UpdateTransferMetaData updates the metadata for a transfer
// https://docs.moov.io/api/index.html#tag/Transfers/operation/patchTransfer
func (c Client) PatchTransfer(ctx context.Context, accountID, transferID string, patches ...TransferPatcher) (*Transfer, error) {
	patch := &patchTransfer{}
	for _, p := range patches {
		p(patch)
	}

	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPatch, pathTransfer, accountID, transferID),
		AcceptJson(),
		JsonBody(patch))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Transfer](resp)
}

/* Exported functions used only for making client calls in the versioned packages (e.g. mv2604) */

func CreateTransferGeneric[TRequest any, TTransfer any](ctx context.Context, client *Client, version Version, partnerAccountID string, transfer TRequest, options ...CreateTransferArgs) CreateTransferGenericBuilder[TTransfer] {
	builder := &createTransferBuilder{}
	callArgs := []callArg{
		MoovVersion(version),
		AcceptJson(),
		JsonBody(transfer),
		WithTransferIdempotencyKey(uuid.New())(builder),
	}

	for _, opt := range options {
		callArgs = append(callArgs, opt(builder))
	}

	return CreateTransferGenericBuilder[TTransfer]{
		client:   client,
		ctx:      ctx,
		endpoint: Endpoint(http.MethodPost, pathTransfers, partnerAccountID),
		callArgs: callArgs,
	}
}

type CreateTransferGenericBuilder[TTransfer any] struct {
	client   *Client
	ctx      context.Context
	endpoint EndpointArg
	callArgs []callArg
}

// Started initiates the transfers request and doesn't wait beyond creating the transfer.
func (r CreateTransferGenericBuilder[TTransfer]) Started() (*TransferStarted, error) {
	if r.client == nil {
		return nil, errors.New("client is nil")
	}

	resp, err := r.client.CallHttp(r.ctx, r.endpoint, r.callArgs...)
	if err != nil {
		return nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		st, err := UnmarshalObjectResponse[TransferStarted](resp)
		return st, err
	case StatusStateConflict:
		return nil, errors.Join(ErrXIdempotencyKey, resp)
	default:
		return nil, resp
	}
}

// WaitForRailResponse starts a transfer request and waits for a response from the rail before returning the result.
func (r CreateTransferGenericBuilder[TTransfer]) WaitForRailResponse() (*TTransfer, *TransferStarted, error) {
	if r.client == nil {
		return nil, nil, errors.New("client is nil")
	}

	resp, err := r.client.CallHttp(r.ctx, r.endpoint, append(r.callArgs, WaitFor("rail-response"))...)
	if err != nil {
		return nil, nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		transfer, err := UnmarshalObjectResponse[TTransfer](resp)
		return transfer, nil, err
	case StatusStarted:
		transferStarted, err := UnmarshalObjectResponse[TransferStarted](resp)
		return nil, transferStarted, err
	case StatusStateConflict:
		return nil, nil, errors.Join(ErrXIdempotencyKey, resp)
	default:
		return nil, nil, resp
	}
}

func PatchTransferGeneric[T any](ctx context.Context, client *Client, version Version, accountID, transferID string, update T) (*Transfer, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}

	resp, err := client.CallHttp(
		ctx,
		Endpoint(http.MethodPatch, pathTransfer, accountID, transferID),
		MoovVersion(version),
		AcceptJson(),
		JsonBody(update),
	)
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Transfer](resp)
}

func ListTransfersGeneric[TTransfer any](ctx context.Context, client *Client, version Version, accountID string, filters ...ListTransferFilter) ([]TTransfer, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}

	resp, err := client.CallHttp(ctx,
		Endpoint(http.MethodGet, pathTransfers, accountID),
		prependArgs(filters, MoovVersion(version), AcceptJson())...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[TTransfer](resp)
}

func GetTransferGeneric[TTransfer any](ctx context.Context, client *Client, version Version, accountID, transferID string) (*TTransfer, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}

	resp, err := client.CallHttp(ctx,
		Endpoint(http.MethodGet, pathTransfer, accountID, transferID),
		MoovVersion(version),
		AcceptJson(),
	)
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[TTransfer](resp)
}

// FetchTransfersGeneric retrieves multiple transfers by their IDs, keyed by transfer ID.
func FetchTransfersGeneric[TTransfer any](ctx context.Context, client *Client, version Version, accountID string, request FetchTransfersRequest) (map[string]TTransfer, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}

	resp, err := client.CallHttp(ctx,
		Endpoint(http.MethodPost, pathTransfersFetch, accountID),
		MoovVersion(version),
		AcceptJson(),
		JsonBody(request))
	if err != nil {
		return nil, err
	}

	out, err := CompletedObjectOrError[map[string]TTransfer](resp)
	if err != nil {
		return nil, err
	}
	return *out, nil
}

type CreateCaptureArgs callArg

// Can be specified to overwrite a randomly generated one.
func WithCaptureIdempotencyKey(key uuid.UUID) CreateCaptureArgs {
	return IdempotencyKey(key.String())
}

// CreateCaptureGeneric captures funds from an authorized card payment. The captures endpoints start in v2026.10.00.
func CreateCaptureGeneric[TRequest any, TCapture any](ctx context.Context, client *Client, version Version, accountID, transferID string, capture TRequest, options ...CreateCaptureArgs) (*TCapture, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}

	args := prependArgs(options,
		MoovVersion(version),
		AcceptJson(),
		WithCaptureIdempotencyKey(uuid.New()),
		JsonBody(capture),
	)

	resp, err := client.CallHttp(ctx, Endpoint(http.MethodPost, pathCaptures, accountID, transferID), args...)
	if err != nil {
		return nil, err
	}

	// The API answers 201 for every created capture, including a declined one; the capture status reports the outcome.
	switch resp.Status() {
	case StatusCompleted, StatusStarted:
		return UnmarshalObjectResponse[TCapture](resp)
	default:
		return nil, resp
	}
}

func ListCapturesGeneric[TCapture any](ctx context.Context, client *Client, version Version, accountID, transferID string) ([]TCapture, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}

	resp, err := client.CallHttp(ctx,
		Endpoint(http.MethodGet, pathCaptures, accountID, transferID),
		MoovVersion(version),
		AcceptJson(),
	)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[TCapture](resp)
}

func GetCaptureGeneric[TCapture any](ctx context.Context, client *Client, version Version, accountID, transferID, captureID string) (*TCapture, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}

	resp, err := client.CallHttp(ctx,
		Endpoint(http.MethodGet, pathCapture, accountID, transferID, captureID),
		MoovVersion(version),
		AcceptJson(),
	)
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[TCapture](resp)
}

// ListTransferEventsGeneric lists a transfer's event timeline, oldest to newest. The events endpoint starts in v2026.10.00.
func ListTransferEventsGeneric[TEvent any](ctx context.Context, client *Client, version Version, accountID, transferID string) ([]TEvent, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	if accountID == "" || transferID == "" {
		return nil, errors.New("accountID and transferID are required")
	}

	resp, err := client.CallHttp(ctx,
		Endpoint(http.MethodGet, pathTransferEvents, accountID, transferID),
		MoovVersion(version),
		AcceptJson(),
	)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[TEvent](resp)
}

type CreateRefundArgs callArg

func WithRefundWaitForRailResponse() CreateRefundArgs {
	return WaitFor("rail-response")
}

// Can be specified to overwrite a randomly generated one.
func WithRefundIdempotencyKey(key uuid.UUID) CreateRefundArgs {
	return IdempotencyKey(key.String())
}

// RefundTransfer refunds a transfer
// https://docs.moov.io/api/#tag/Transfers/operation/refundTransfer
func (c Client) RefundTransfer(ctx context.Context, partnerAccountID, transferID string, refund CreateRefund, options ...CreateRefundArgs) (*Refund, *RefundStarted, error) {
	args := prependArgs(options,
		AcceptJson(),
		WithRefundIdempotencyKey(uuid.New()),
		JsonBody(refund),
	)

	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathRefunds, partnerAccountID, transferID),
		args...)
	if err != nil {
		return nil, nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		r, err := CompletedObjectOrError[Refund](resp)
		return r, nil, err
	case StatusStarted:
		r, err := UnmarshalObjectResponse[RefundStarted](resp)
		return nil, r, err
	default:
		return nil, nil, resp
	}
}

func RefundTransferGeneric[TRequest any, TRefund any](ctx context.Context, client *Client, version Version, partnerAccountID, transferID string, refund TRequest, options ...CreateRefundArgs) (*TRefund, *RefundStarted, error) {
	if client == nil {
		return nil, nil, errors.New("client is nil")
	}

	args := prependArgs(options,
		MoovVersion(version),
		AcceptJson(),
		WithRefundIdempotencyKey(uuid.New()),
		JsonBody(refund),
	)

	resp, err := client.CallHttp(ctx,
		Endpoint(http.MethodPost, pathRefunds, partnerAccountID, transferID),
		args...)
	if err != nil {
		return nil, nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		r, err := CompletedObjectOrError[TRefund](resp)
		return r, nil, err
	case StatusStarted:
		r, err := UnmarshalObjectResponse[RefundStarted](resp)
		return nil, r, err
	default:
		return nil, nil, resp
	}
}

// ListRefunds lists all refunds for a transfer
// https://docs.moov.io/api/index.html#tag/Transfers/operation/getRefunds
func (c Client) ListRefunds(ctx context.Context, accountID, transferID string) ([]Refund, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathRefunds, accountID, transferID),
		AcceptJson(),
	)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[Refund](resp)
}

// CreateRefundGeneric refunds a transfer and returns the refund for a 200 or a 202 response.
// From v2026.10.00 a 202 means the wait for the rail timed out, and its body is the refund,
// not the RefundStarted shape that RefundTransferGeneric decodes.
func CreateRefundGeneric[TRequest any, TRefund any](ctx context.Context, client *Client, version Version, partnerAccountID, transferID string, refund TRequest, options ...CreateRefundArgs) (*TRefund, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}

	args := prependArgs(options,
		MoovVersion(version),
		AcceptJson(),
		WithRefundIdempotencyKey(uuid.New()),
		JsonBody(refund),
	)

	resp, err := client.CallHttp(ctx, Endpoint(http.MethodPost, pathRefunds, partnerAccountID, transferID), args...)
	if err != nil {
		return nil, err
	}

	switch resp.Status() {
	case StatusCompleted, StatusStarted:
		return UnmarshalObjectResponse[TRefund](resp)
	default:
		return nil, resp
	}
}

func ListRefundsGeneric[TRefund any](ctx context.Context, client *Client, version Version, accountID, transferID string) ([]TRefund, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}

	resp, err := client.CallHttp(ctx,
		Endpoint(http.MethodGet, pathRefunds, accountID, transferID),
		MoovVersion(version),
		AcceptJson(),
	)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[TRefund](resp)
}

// GetRefund retrieves a refund for a transfer
// https://docs.moov.io/api/index.html#tag/Transfers/operation/getRefund
func (c Client) GetRefund(ctx context.Context, accountID, transferID, refundID string) (*Refund, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathRefund, accountID, transferID, refundID),
		AcceptJson(),
	)
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Refund](resp)
}

func GetRefundGeneric[TRefund any](ctx context.Context, client *Client, version Version, accountID, transferID, refundID string) (*TRefund, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}

	resp, err := client.CallHttp(ctx,
		Endpoint(http.MethodGet, pathRefund, accountID, transferID, refundID),
		MoovVersion(version),
		AcceptJson(),
	)
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[TRefund](resp)
}

type CreateReversalArgs callArg

// Can be specified to overwrite a randomly generated one.
func WithReversalsIdempotencyKey(key uuid.UUID) CreateReversalArgs {
	return IdempotencyKey(key.String())
}

// ReverseTransfer reverses a transfer
// https://docs.moov.io/api/index.html#tag/Transfers/operation/reverseTransfer
func (c Client) ReverseTransfer(ctx context.Context, partnerAccountID, transferID string, refund CreateReversal, options ...CreateReversalArgs) (*CreatedReversal, error) {
	args := prependArgs(options,
		AcceptJson(),
		WithReversalsIdempotencyKey(uuid.New()),
		JsonBody(refund),
	)

	resp, err := c.CallHttp(ctx, Endpoint(http.MethodPost, pathReversals, partnerAccountID, transferID), args...)
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[CreatedReversal](resp)
}

func ReverseTransferGeneric[TRequest any, TCreatedReversal any](ctx context.Context, client *Client, version Version, partnerAccountID, transferID string, reversal TRequest, options ...CreateReversalArgs) (*TCreatedReversal, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}

	args := prependArgs(options,
		MoovVersion(version),
		AcceptJson(),
		WithReversalsIdempotencyKey(uuid.New()),
		JsonBody(reversal),
	)

	resp, err := client.CallHttp(ctx, Endpoint(http.MethodPost, pathReversals, partnerAccountID, transferID), args...)
	if err != nil {
		return nil, err
	}

	// A 202 means the reversal was created but the wait for the rail timed out. Its body is the same reversal.
	switch resp.Status() {
	case StatusCompleted, StatusStarted:
		return UnmarshalObjectResponse[TCreatedReversal](resp)
	default:
		return nil, resp
	}
}

// CancelTransfer cancels a transfer
// https://docs.moov.io/api/money-movement/transfers/cancel/
func (c Client) CancelTransfer(ctx context.Context, accountID string, transferID string) (*Cancellation, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodPost, pathCancellations, accountID, transferID))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Cancellation](resp)
}

type CreateCancellationArgs callArg

// Can be specified to overwrite a randomly generated one.
func WithCancellationIdempotencyKey(key uuid.UUID) CreateCancellationArgs {
	return IdempotencyKey(key.String())
}

// CancelTransferGeneric cancels a transfer. From v2026.10.00 the request needs an idempotency key and accepts an optional body.
func CancelTransferGeneric[TRequest any, TCancellation any](ctx context.Context, client *Client, version Version, accountID, transferID string, cancellation TRequest, options ...CreateCancellationArgs) (*TCancellation, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}

	args := prependArgs(options,
		MoovVersion(version),
		AcceptJson(),
		WithCancellationIdempotencyKey(uuid.New()),
		JsonBody(cancellation),
	)

	resp, err := client.CallHttp(ctx, Endpoint(http.MethodPost, pathCancellations, accountID, transferID), args...)
	if err != nil {
		return nil, err
	}

	// The API spec allows 200 and 202 here, and both carry a cancellation body.
	switch resp.Status() {
	case StatusCompleted, StatusStarted:
		return UnmarshalObjectResponse[TCancellation](resp)
	default:
		return nil, resp
	}
}

func ListCancellationsGeneric[TCancellation any](ctx context.Context, client *Client, version Version, accountID, transferID string) ([]TCancellation, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}

	resp, err := client.CallHttp(ctx,
		Endpoint(http.MethodGet, pathCancellations, accountID, transferID),
		MoovVersion(version),
		AcceptJson(),
	)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[TCancellation](resp)
}

func GetCancellationGeneric[TCancellation any](ctx context.Context, client *Client, version Version, accountID, transferID, cancellationID string) (*TCancellation, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}

	resp, err := client.CallHttp(ctx,
		Endpoint(http.MethodGet, pathCancellation, accountID, transferID, cancellationID),
		MoovVersion(version),
		AcceptJson(),
	)
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[TCancellation](resp)
}

// GetCancellation gets a cancellation
// https://docs.moov.io/api/money-movement/transfers/cancel-details/
func (c Client) GetCancellation(ctx context.Context, accountID string, transferID string, cancellationID string) (*Cancellation, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathCancellation, accountID, transferID, cancellationID))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Cancellation](resp)
}

func (c Client) ListCancellations(ctx context.Context, accountID string, transferID string) ([]Cancellation, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathCancellations, accountID, transferID),
		AcceptJson(),
	)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[Cancellation](resp)
}

// TransferOptions lists all transfer options between a source and destination
// https://docs.moov.io/api/#tag/Transfers/operation/createTransferOptions
func (c Client) TransferOptions(ctx context.Context, partnerAccountID string, payload CreateTransferOptions) (*TransferOptions, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathTransferOptions, partnerAccountID),
		AcceptJson(),
		JsonBody(payload))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[TransferOptions](resp)
}
