package moov

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Dispute struct {
	DisputeID                string        `json:"disputeID,omitempty"`
	CreatedOn                time.Time     `json:"createdOn,omitempty"`
	Amount                   Amount        `json:"amount,omitempty"`
	NetworkReasonCode        string        `json:"networkReasonCode,omitempty"`
	NetworkReasonDescription string        `json:"networkReasonDescription,omitempty"`
	Phase                    DisputePhase  `json:"phase,omitempty"`
	RespondBy                time.Time     `json:"respondBy,omitempty"`
	Status                   DisputeStatus `json:"status,omitempty"`
	Transfer                 Transfer      `json:"transfer,omitempty"`
}

type DisputeStatus string

const (
	DisputeStatus_ResponseNeeded DisputeStatus = "response-needed"
	DisputeStatus_Resolved       DisputeStatus = "resolved"
	DisputeStatus_UnderReview    DisputeStatus = "under-review"
	DisputeStatus_Closed         DisputeStatus = "closed"
	DisputeStatus_Accepted       DisputeStatus = "accepted"
	DisputeStatus_Expired        DisputeStatus = "expired"
	DisputeStatus_Won            DisputeStatus = "won"
	DisputeStatus_Lost           DisputeStatus = "lost"
)

type DisputePhase string

const (
	DisputePhase_PreDispute DisputePhase = "pre-dispute"
	DisputePhase_Inquiry    DisputePhase = "inquiry"
	DisputePhase_Chargeback DisputePhase = "chargeback"
	DisputePhase_Unknown    DisputePhase = "unknown"
)

type DisputeEvidence struct {
	CreatedOn    time.Time `json:"createdOn,omitempty"`
	DisputeID    string    `json:"disputeID,omitempty"`
	EvidenceID   string    `json:"evidenceID,omitempty"`
	EvidenceType string    `json:"evidenceType,omitempty"`
	FileName     string    `json:"fileName,omitempty"`
	MimeType     string    `json:"mimeType,omitempty"`
	Size         int       `json:"size,omitempty"`
	Text         string    `json:"text,omitempty"`
	UpdatedOn    time.Time `json:"updatedOn,omitempty"`
}

type DisputesEvidenceText struct {
	Text         string       `json:"text"`
	EvidenceType EvidenceType `json:"evidenceType"`
}

type DisputesEvidenceUpdate struct {
	EvidenceType EvidenceType `json:"evidenceType"`
}

type DisputeEvidenceUpload struct {
	DisputeID    string    `json:"disputeID,omitempty"`
	EvidenceID   string    `json:"evidenceID,omitempty"`
	EvidenceType string    `json:"evidenceType,omitempty"`
	FileName     string    `json:"fileName,omitempty"`
	MimeType     string    `json:"mimeType,omitempty"`
	Size         int       `json:"size,omitempty"`
	CreatedOn    time.Time `json:"createdOn,omitempty"`
}

type DisputeListFilter callArg

func WithDisputeCount(c int) DisputeListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["count"] = fmt.Sprintf("%d", c)
		return nil
	})
}

func WithDisputeSkip(c int) DisputeListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["skip"] = fmt.Sprintf("%d", c)
		return nil
	})
}

func WithDisputeResponseStartDate(t time.Time) DisputeListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["respondStartDateTime"] = t.Format(time.RFC3339)
		return nil
	})
}

func WithDisputeResponseEndDate(t time.Time) DisputeListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["respondEndDateTime"] = t.Format(time.RFC3339)
		return nil
	})
}

func WithDisputeStatus(s string) DisputeListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["status"] = s
		return nil
	})
}

func WithDisputeMerchantAccountID(id string) DisputeListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["merchantAccountID"] = id
		return nil
	})
}

func WithDisputeCardHolderAccountID(id string) DisputeListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["cardholderAccountID"] = id
		return nil
	})
}

func WithDisputeStartDate(t time.Time) DisputeListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["startDateTime"] = t.Format(time.RFC3339)
		return nil
	})
}

func WithDisputeEndDate(t time.Time) DisputeListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["endDateTime"] = t.Format(time.RFC3339)
		return nil
	})
}

func WithDisputeOrderBy(orderBy string) DisputeListFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["orderBy"] = orderBy
		return nil
	})
}

type EvidenceType string

const (
	EvidenceType_Receipt               EvidenceType = "receipt"
	EvidenceType_ProofOfDelivery       EvidenceType = "proof-of-delivery"
	EvidenceType_CancelationPolicy     EvidenceType = "cancelation-policy"
	EvidenceType_TermsOfService        EvidenceType = "terms-of-service"
	EvidenceType_CustomerCommunication EvidenceType = "customer-communication"
	EvidenceType_GenericEvidence       EvidenceType = "generic-evidence"
	EvidenceType_CoverLetter           EvidenceType = "cover-letter"
	EvidenceType_Other                 EvidenceType = "other"
)

// ListDisputes lists of Disputes that are associated with a Moov account
// https://docs.moov.io/api/money-movement/disputes/list/
func (c Client) ListDisputes(ctx context.Context, accountID string, filters ...DisputeListFilter) ([]Dispute, error) {
	args := prependArgs(filters, AcceptJson())
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathDisputes, accountID), args...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[Dispute](resp)
}

// GetDispute retrieves a dispute for the given dispute id
// https://docs.moov.io/api/money-movement/disputes/get/
func (c Client) GetDispute(ctx context.Context, accountID string, disputeID string) (*Dispute, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathDispute, accountID, disputeID), AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Dispute](resp)
}

// AcceptDispute sets the status of the dispute to accepted
// https://docs.moov.io/api/money-movement/disputes/post/
func (c Client) AcceptDispute(ctx context.Context, accountID string, disputeID string) (*Dispute, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodPost, pathDisputeAccept, accountID, disputeID), AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Dispute](resp)
}

// UploadDisputeEvidence Uploads text as evidence for a dispute.
// https://docs.moov.io/api/money-movement/disputes/post-text/
func (c Client) UploadDisputeEvidence(ctx context.Context, accountID string, disputeID string, evidenceText DisputesEvidenceText) (*DisputeEvidence, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodPost, pathDisputeEvidenceText, accountID, disputeID), AcceptJson(), JsonBody(evidenceText))

	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[DisputeEvidence](resp)
}

// DeleteDisputeEvidence deletes a piece of dispute evidence for the given dispute and evidence id
// https://docs.moov.io/api/money-movement/disputes/delete
func (c Client) DeleteDisputeEvidence(ctx context.Context, accountID string, disputeID, evidenceID string) error {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodDelete, pathDisputeEvidence, accountID, disputeID, evidenceID), AcceptJson())
	if err != nil {
		return err
	}
	return CompletedNilOrError(resp)
}

// UploadEvidenceFile uploads a new evidence file for the given dispute id
// https://docs.moov.io/api/money-movement/disputes/post-file/
func (c Client) UploadEvidenceFile(ctx context.Context, accountID string, disputeID string, evidenceType EvidenceType, filename string, file io.Reader, mimeType string) (*DisputeEvidenceUpload, error) {
	var multiParts []multipartFn
	multiParts = append(multiParts, MultipartField("evidenceType", string(evidenceType)))
	multiParts = append(multiParts, MultipartFile("file", filename, file, mimeType))

	resp, err := c.CallHttp(ctx, Endpoint(http.MethodPost, pathDisputeEvidenceFile, accountID, disputeID), MultipartBody(multiParts...))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[DisputeEvidenceUpload](resp)
}

// ListDisputeEvidence lists all evidence for the given dispute id
// https://docs.moov.io/api/money-movement/disputes/list-evidence/
func (c Client) ListDisputeEvidence(ctx context.Context, accountID string, disputeID string) ([]DisputeEvidence, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathDisputeEvidences, accountID, disputeID), AcceptJson())

	if err != nil {
		return nil, err
	}

	return CompletedListOrError[DisputeEvidence](resp)
}

// SubmitDisputeEvidence Submits evidence for a dispute.
// https://docs.moov.io/api/money-movement/disputes/post-evidence/
func (c Client) SubmitDisputeEvidence(ctx context.Context, accountID string, disputeID string) (*Dispute, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodPost, pathDisputeSubmitEvidence, accountID, disputeID), AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Dispute](resp)
}

// UpdateDisputeEvidence Updates dispute evidence by ID.
// https://docs.moov.io/api/money-movement/disputes/patch/
func (c Client) UpdateDisputeEvidence(ctx context.Context, accountID string, disputeID string, evidenceID string, evidenceUpdate DisputesEvidenceUpdate) (*DisputeEvidence, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodPatch, pathDisputeEvidence, accountID, disputeID, evidenceID), AcceptJson(), JsonBody(evidenceUpdate))

	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[DisputeEvidence](resp)
}

// GetDisputeEvidence retrieves the piece of dispute evidence for the given dispute and evidence id
// https://docs.moov.io/api/money-movement/disputes/get-evidence/
func (c Client) GetDisputeEvidence(ctx context.Context, accountID string, disputeID, evidenceID string) (*DisputeEvidence, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathDisputeEvidence, accountID, disputeID, evidenceID), AcceptJson())

	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[DisputeEvidence](resp)
}
