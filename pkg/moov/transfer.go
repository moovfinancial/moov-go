package moov

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TransferStatus int

const (
	TransferStatusCreated TransferStatus = iota
	TransferStatusPending
	TransferStatusCompleted
	TransferStatusFailed
	TransferStatusReversed
	TransferStatusQueued
	TransferStatusCanceled
)

var TransferStatusStrings = map[TransferStatus]string{
	TransferStatusCreated:   "created",
	TransferStatusPending:   "pending",
	TransferStatusCompleted: "completed",
	TransferStatusFailed:    "failed",
	TransferStatusReversed:  "reversed",
	TransferStatusQueued:    "queued",
	TransferStatusCanceled:  "canceled",
}

type SynchronousTransfer struct {
	TransferID     string            `json:"transferID,omitempty"`
	CreatedOn      time.Time         `json:"createdOn,omitempty"`
	CompletedOn    time.Time         `json:"completedOn,omitempty"`
	Status         string            `json:"status,omitempty"`
	FailureReason  string            `json:"failureReason,omitempty"`
	Amount         Amount            `json:"amount,omitempty"`
	Description    string            `json:"description,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	FacilitatorFee FacilitatorFee    `json:"facilitatorFee,omitempty"`
	MoovFee        int               `json:"moovFee,omitempty"`
	MoovFeeDecimal string            `json:"moovFeeDecimal,omitempty"`
	MoovFeeDetails MoovFeeDetails    `json:"moovFeeDetails,omitempty"`
	GroupID        string            `json:"groupID,omitempty"`
	RefundedAmount Amount            `json:"refundedAmount,omitempty"`
	Refunds        []Refund          `json:"refunds,omitempty"`
	DisputedAmount Amount            `json:"disputedAmount,omitempty"`
	Disputes       []Dispute         `json:"disputes,omitempty"`
	Source         Source            `json:"source,omitempty"`
	Destination    Destination       `json:"destination,omitempty"`
}

type AsynchronousTransfer struct {
	TransferID string    `json:"transferID,omitempty"`
	CreatedOn  time.Time `json:"createdOn,omitempty"`
}

type Amount struct {
	Currency string `json:"currency,omitempty"`
	Value    int    `json:"value,omitempty"`
}

type FacilitatorFee struct {
	Total         int    `json:"total,omitempty"`
	TotalDecimal  string `json:"totalDecimal,omitempty"`
	Markup        int    `json:"markup,omitempty"`
	MarkupDecimal string `json:"markupDecimal,omitempty"`
}

type MoovFeeDetails struct {
	CardScheme     string `json:"cardScheme,omitempty"`
	Interchange    string `json:"interchange,omitempty"`
	MoovProcessing string `json:"moovProcessing,omitempty"`
}

type Refund struct {
	RefundID    string      `json:"refundID,omitempty"`
	CreatedOn   time.Time   `json:"createdOn,omitempty"`
	UpdatedOn   time.Time   `json:"updatedOn,omitempty"`
	Status      string      `json:"status,omitempty"`
	FailureCode string      `json:"failureCode,omitempty"`
	Amount      Amount      `json:"amount,omitempty"`
	CardDetails CardDetails `json:"cardDetails,omitempty"`
}

type Source struct {
	PaymentMethodID   string          `json:"paymentMethodID,omitempty"`
	PaymentMethodType string          `json:"paymentMethodType,omitempty"`
	Account           TransferAccount `json:"account,omitempty"`
	BankAccount       BankAccount     `json:"bankAccount,omitempty"`
	Wallet            Wallet          `json:"wallet,omitempty"`
	Card              Card            `json:"card,omitempty"`
	ApplePay          ApplePay        `json:"applePay,omitempty"`
	AchDetails        AchDetails      `json:"achDetails,omitempty"`
	CardDetails       CardDetails     `json:"cardDetails,omitempty"`
	TransferID        string          `json:"transferID,omitempty"`
}

type TransferAccount struct {
	AccountID   string `json:"accountID,omitempty"`
	Email       string `json:"email,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

type Destination struct {
	PaymentMethodID   string          `json:"paymentMethodID,omitempty"`
	PaymentMethodType string          `json:"paymentMethodType,omitempty"`
	Account           TransferAccount `json:"account,omitempty"`
	BankAccount       BankAccount     `json:"bankAccount,omitempty"`
	Wallet            Wallet          `json:"wallet,omitempty"`
	Card              Card            `json:"card,omitempty"`
	ApplePay          ApplePay        `json:"applePay,omitempty"`
	AchDetails        AchDetails      `json:"achDetails,omitempty"`
	CardDetails       CardDetails     `json:"cardDetails,omitempty"`
}

type MetadataPayload struct {
	Metadata map[string]string `json:"metadata,omitempty"`
}

type RefundPayload struct {
	Amount int `json:"amount,omitempty"`
}

type TransferOptionsSourcePayload struct {
	PaymentMethodID string `json:"paymentMethodID,omitempty"`
	AccountID       string `json:"accountID,omitempty"`
}

type TransferOptionsDestinationPayload struct {
	PaymentMethodID string `json:"paymentMethodID,omitempty"`
	AccountID       string `json:"accountID,omitempty"`
}

type TransferOptionsPayload struct {
	Source      TransferOptionsSourcePayload      `json:"source,omitempty"`
	Destination TransferOptionsDestinationPayload `json:"destination,omitempty"`
	Amount      Amount                            `json:"amount,omitempty"`
}

type CreatedTransferOptions struct {
	SourceOptions      []Source `json:"sourceOptions,omitempty"`
	DestinationOptions []Source `json:"destinationOptions,omitempty"`
}

type RefundStatus struct {
	Status    string    `json:"status,omitempty"`
	CreatedOn time.Time `json:"createdOn,omitempty"`
}

type CanceledTransfer struct {
	Cancellation RefundStatus `json:"cancellation,omitempty"`
	Refund       Refund       `json:"refund,omitempty"`
}

type CreateTransfer struct {
	Source         Source            `json:"source,omitempty"`
	Destination    Destination       `json:"destination,omitempty"`
	Amount         Amount            `json:"amount,omitempty"`
	FacilitatorFee FacilitatorFee    `json:"facilitatorFee,omitempty"`
	Description    string            `json:"description,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

type CreateTransferArgs callArg

func WithTransferWaitForRailResponse() CreateTransferArgs {
	return WaitFor("rail-response")
}

// Can be specified to overwrite a randomly generated one.
func WithTransferIdempotencyKey(key uuid.UUID) CreateTransferArgs {
	return IdempotencyKey(key.String())
}

// CreateTransfer creates a new transfer
// https://docs.moov.io/api/index.html#tag/Transfers/operation/createTransfer
func (c Client) CreateTransfer(ctx context.Context, transfer CreateTransfer, options ...CreateTransferArgs) (*SynchronousTransfer, *AsynchronousTransfer, error) {
	args := prependArgs[CreateTransferArgs](options,
		// default options
		AcceptJson(),
		JsonBody(transfer),
		WithTransferIdempotencyKey(uuid.New()),
	)

	resp, err := c.CallHttp(ctx, Endpoint(http.MethodPost, pathTransfers), args...)
	if err != nil {
		return nil, nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		st, err := UnmarshalObjectResponse[SynchronousTransfer](resp)
		return st, nil, err
	case StatusStarted:
		st, err := UnmarshalObjectResponse[AsynchronousTransfer](resp)
		return nil, st, err
	case StatusStateConflict:
		return nil, nil, errors.Join(ErrXIdempotencyKey, resp)
	default:
		return nil, nil, resp
	}
}

type ListTransferFilter callArg

func WithTransferAccountIDs(accountIDs []string) ListTransferFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["accoundIDs"] = strings.Join(accountIDs, ",")
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

func WithTransferSkip(skip int) ListTransferFilter {
	return Skip(skip)
}

func WithTransferCount(count int) ListTransferFilter {
	return Count(count)
}

// ListTransfers lists all transfers
// https://docs.moov.io/api/index.html#tag/Transfers/operation/listTransfers
func (c Client) ListTransfers(ctx context.Context, filters ...ListTransferFilter) ([]SynchronousTransfer, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathTransfers),
		prependArgs(filters, AcceptJson())...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[SynchronousTransfer](resp)
}

// GetTransfer retrieves a transfer
// https://docs.moov.io/api/index.html#tag/Transfers/operation/getTransfer
func (c Client) GetTransfer(ctx context.Context, transferID string) (*SynchronousTransfer, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathTransfer, transferID),
		AcceptJson(),
	)
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[SynchronousTransfer](resp)
}

// UpdateTransferMetaData updates the metadata for a transfer
// https://docs.moov.io/api/index.html#tag/Transfers/operation/patchTransfer
func (c Client) UpdateTransferMetaData(ctx context.Context, transferID string, metadata MetadataPayload) (*SynchronousTransfer, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPatch, pathTransfer, transferID),
		AcceptJson(),
		JsonBody(metadata))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[SynchronousTransfer](resp)
}

// TransferOptions lists all transfer options between a source and destination
// https://docs.moov.io/api/#tag/Transfers/operation/createTransferOptions
func (c Client) TransferOptions(ctx context.Context, payload TransferOptionsPayload) (*CreatedTransferOptions, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathTransferOptions),
		AcceptJson(),
		JsonBody(payload))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[CreatedTransferOptions](resp)
}

type CreateRefundArgs callArg

func WithRefundWaitForRailResponse() CreateTransferArgs {
	return WaitFor("rail-response")
}

// Can be specified to overwrite a randomly generated one.
func WithRefundIdempotencyKey(key uuid.UUID) CreateTransferArgs {
	return IdempotencyKey(key.String())
}

// RefundTransfer refunds a transfer
// https://docs.moov.io/api/#tag/Transfers/operation/refundTransfer
func (c Client) RefundTransfer(ctx context.Context, transferID string, refund RefundPayload, options ...CreateRefundArgs) (*Refund, *Refund, error) {
	args := prependArgs(options,
		AcceptJson(),
		WithRefundIdempotencyKey(uuid.New()),
		JsonBody(refund),
	)

	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathRefunds, transferID),
		args...)
	if err != nil {
		return nil, nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		r, err := CompletedObjectOrError[Refund](resp)
		return r, nil, err
	case StatusStarted:
		r, err := CompletedObjectOrError[Refund](resp)
		return nil, r, err
	default:
		return nil, nil, resp
	}
}

// ListRefunds lists all refunds for a transfer
// https://docs.moov.io/api/index.html#tag/Transfers/operation/getRefunds
func (c Client) ListRefunds(ctx context.Context, transferID string) ([]Refund, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathRefunds, transferID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[Refund](resp)
}

// GetRefund retrieves a refund for a transfer
// https://docs.moov.io/api/index.html#tag/Transfers/operation/getRefund
func (c Client) GetRefund(ctx context.Context, transferID string, refundID string) (*Refund, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathRefund, transferID, refundID), AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Refund](resp)
}

type CreateReversalArgs callArg

// Can be specified to overwrite a randomly generated one.
func WithReversalsIdempotencyKey(key uuid.UUID) CreateTransferArgs {
	return IdempotencyKey(key.String())
}

// ReverseTransfer reverses a transfer
// https://docs.moov.io/api/index.html#tag/Transfers/operation/reverseTransfer
func (c Client) ReverseTransfer(ctx context.Context, transferID string, refund RefundPayload, options ...CreateReversalArgs) (*CanceledTransfer, error) {
	args := prependArgs(options,
		AcceptJson(),
		WithReversalsIdempotencyKey(uuid.New()),
		JsonBody(refund),
	)

	resp, err := c.CallHttp(ctx, Endpoint(http.MethodPost, pathTransferReversals, transferID), args...)
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[CanceledTransfer](resp)
}
