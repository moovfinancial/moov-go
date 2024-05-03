package moov

import (
	"context"
	"fmt"
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
	Data         string    `json:"data,omitempty"`
	DisputeID    string    `json:"disputeID,omitempty"`
	EvidenceID   string    `json:"evidenceID,omitempty"`
	EvidenceType string    `json:"evidenceType,omitempty"`
	FileName     string    `json:"fileName,omitempty"`
	MimeType     string    `json:"mimeType,omitempty"`
	Size         int       `json:"size,omitempty"`
	Text         string    `json:"text,omitempty"`
	UpdatedOn    time.Time `json:"updatedOn,omitempty"`
}

type DisputeTextEvidenceType string

const (
	DisputeTextEvidenceType_Receipt               DisputeTextEvidenceType = "receipt"
	DisputeTextEvidenceType_ProofOfDelivery       DisputeTextEvidenceType = "proof-of-delivery"
	DisputeTextEvidenceType_CancellationPolicy    DisputeTextEvidenceType = "cancelation-policy"
	DisputeTextEvidenceType_TermsOfService        DisputeTextEvidenceType = "terms-of-service"
	DisputeTextEvidenceType_CustomerCommunication DisputeTextEvidenceType = "customer-communication"
	DisputeTextEvidenceType_GenericEvidence       DisputeTextEvidenceType = "generic-evidence"
	DisputeTextEvidenceType_CoverLetter           DisputeTextEvidenceType = "cover-letter"
	DisputeTextEvidenceType_Other                 DisputeTextEvidenceType = "other"
)

type DisputesEvidenceText struct {
	Text         string                  `json:"text"`
	EvidenceType DisputeTextEvidenceType `json:"evidenceType"`
}

type DisputesEvidenceUpdate struct {
	EvidenceType DisputeTextEvidenceType `json:"evidenceType"`
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
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathDispute, disputeID), AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Dispute](resp)
}

// AcceptDispute sets the status of the dispute to accepted
// https://docs.moov.io/api/money-movement/disputes/post/
func (c Client) AcceptDispute(ctx context.Context, disputeID string) (*Dispute, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodPost, pathDisputeAccept, disputeID), AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Dispute](resp)
}

// UploadDisputeEvidence Uploads text as evidence for a dispute.
// https://docs.moov.io/api/money-movement/disputes/post-text/
func (c Client) UploadDisputeEvidence(ctx context.Context, disputeID string, evidenceText DisputesEvidenceText) ([]DisputeEvidence, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodPost, pathDisputeUploadEvidenceText, disputeID), AcceptJson(), JsonBody(evidenceText))
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[DisputeEvidence](resp)
}

// SubmitDisputeEvidence Submits evidence for a dispute.
// https://docs.moov.io/api/money-movement/disputes/post-evidence/
func (c Client) SubmitDisputeEvidence(ctx context.Context, disputeID string) (*Dispute, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodPost, pathDisputeSubmitEvidence, disputeID), AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Dispute](resp)
}

// UpdateDisputeEvidence Updates dispute evidence by ID.
// https://docs.moov.io/api/money-movement/disputes/patch/
func (c Client) UpdateDisputeEvidence(ctx context.Context, disputeID string, evidenceID string, evidenceUpdate DisputesEvidenceUpdate) (*DisputeEvidence, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodPatch, pathDisputeUpdateEvidence, disputeID, evidenceID), AcceptJson(), JsonBody(evidenceUpdate))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[DisputeEvidence](resp)
}
