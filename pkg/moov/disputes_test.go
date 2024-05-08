package moov_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDisputesMarshal(t *testing.T) {
	input := []byte(`{
			"amount": {
				"currency": "USD",
				"value": 1204
			},
			"createdOn": null,
			"disputeID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43",
			"networkReasonCode": null,
			"networkReasonDescription": null,
			"respondBy": null,
			"status": "response-needed",
			"transfer": {
				"transferID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43"
			}}`)

	dispute := new(moov.Dispute)

	dec := json.NewDecoder(bytes.NewReader(input))
	dec.DisallowUnknownFields()

	err := dec.Decode(&dispute)
	require.NoError(t, err)

	assert.Equal(t, "ec7e1848-dc80-4ab0-8827-dd7fc0737b43", dispute.DisputeID)
}

func Test_Disputes(t *testing.T) {
	mc := NewTestClient(t)

	disputes, err := mc.ListDisputes(context.Background(), moov.WithDisputeCount(200), moov.WithDisputeSkip(0))
	require.NoError(t, err)
	require.NotNil(t, disputes)
}

func Test_GetDisputes_NotFound(t *testing.T) {
	mc := NewTestClient(t)

	// We don't have any disputes to test against! So we can at least check for not found vs other possible errors
	disputeID := uuid.NewString()
	dispute, err := mc.GetDispute(context.Background(), disputeID)
	require.Nil(t, dispute)

	// find and cast the error into HttpCallError so it can be inspected
	var httpErr moov.HttpCallResponse
	require.ErrorAs(t, err, &httpErr)

	require.Equal(t, moov.StatusNotFound, httpErr.Status())
}

func Test_AcceptDisputes_Unauthorized(t *testing.T) {
	mc := NewTestClient(t)

	// We don't have any disputes to test against! So using a random id to accept will return unauthorized since we dont own the resource
	disputeID := uuid.NewString()
	dispute, err := mc.AcceptDispute(context.Background(), disputeID)
	require.Nil(t, dispute)

	// find and cast the error into HttpCallError so it can be inspected
	var httpErr moov.HttpCallResponse
	require.ErrorAs(t, err, &httpErr)

	require.Equal(t, moov.StatusUnauthorized, httpErr.Status())
}

func Test_UploadDisputeEvidence_Unauthorized(t *testing.T) {
	mc := NewTestClient(t)

	// We don't have any disputes to test against! So using a random id to upload evidence text will return unauthorized since we dont own the resource
	disputeID := uuid.NewString()
	dispute, err := mc.UploadDisputeEvidence(context.Background(), disputeID, moov.DisputesEvidenceText{
		Text:         "Some evidence text",
		EvidenceType: moov.EvidenceType_GenericEvidence,
	})
	require.Nil(t, dispute)

	// find and cast the error into HttpCallError so it can be inspected
	var httpErr moov.HttpCallResponse
	require.ErrorAs(t, err, &httpErr)

	require.Equal(t, moov.StatusUnauthorized, httpErr.Status())
}

func Test_SubmitDisputeEvidence_Unauthorized(t *testing.T) {
	mc := NewTestClient(t)

	// We don't have any disputes to test against! So using a random id to submit evidence text will return unauthorized since we dont own the resource
	disputeID := uuid.NewString()
	dispute, err := mc.SubmitDisputeEvidence(context.Background(), disputeID)
	require.Nil(t, dispute)

	// find and cast the error into HttpCallError so it can be inspected
	var httpErr moov.HttpCallResponse
	require.ErrorAs(t, err, &httpErr)

	require.Equal(t, moov.StatusUnauthorized, httpErr.Status())
}

func Test_UpdateDisputeEvidence_Unauthorized(t *testing.T) {
	mc := NewTestClient(t)

	// We don't have any disputes to test against! So using a random id to update evidence will return unauthorized since we dont own the resource
	disputeID := uuid.NewString()
	evidenceID := uuid.NewString()
	dispute, err := mc.UpdateDisputeEvidence(context.Background(), disputeID, evidenceID, moov.DisputesEvidenceUpdate{
		EvidenceType: moov.EvidenceType_Other,
	})
	require.Nil(t, dispute)

	// find and cast the error into HttpCallError so it can be inspected
	var httpErr moov.HttpCallResponse
	require.ErrorAs(t, err, &httpErr)

	require.Equal(t, moov.StatusUnauthorized, httpErr.Status())
}

func Test_DeleteDisputeEvidence_NotFound(t *testing.T) {
	mc := NewTestClient(t)

	// We don't have any disputes to test against! So using a random id to update evidence will return unauthorized since we dont own the resource
	disputeID := uuid.NewString()
	evidenceID := uuid.NewString()

	err := mc.DeleteDisputeEvidence(context.Background(), disputeID, evidenceID)
	require.Error(t, err)

	// find and cast the error into HttpCallError so it can be inspected
	var httpErr moov.HttpCallResponse
	require.ErrorAs(t, err, &httpErr)

	require.Equal(t, moov.StatusUnauthorized, httpErr.Status())
}

func Test_UploadEvidenceFile_DisputeNotFound(t *testing.T) {
	mc := NewTestClient(t)

	// We don't have any disputes to test against! So using a random id to update evidence will return unauthorized since we dont own the resource
	disputeID := uuid.NewString()
	file := bytes.NewReader([]byte("test"))
	err := mc.UploadEvidenceFile(context.Background(), disputeID, moov.EvidenceType_CoverLetter, "test.txt", file)
	require.Error(t, err)

	// find and cast the error into HttpCallError so it can be inspected
	var httpErr moov.HttpCallResponse
	require.ErrorAs(t, err, &httpErr)

	require.Equal(t, moov.StatusUnauthorized, httpErr.Status())
}

func Test_ListDisputeEvidence_DisputeNotFound(t *testing.T) {
	mc := NewTestClient(t)

	// We don't have any disputes to test against! So using a random id to update evidence will return unauthorized since we dont own the resource
	disputeID := uuid.NewString()
	evidence, err := mc.ListDisputeEvidence(context.Background(), disputeID)
	require.Nil(t, evidence)
	require.Error(t, err)

	// find and cast the error into HttpCallError so it can be inspected
	var httpErr moov.HttpCallResponse
	require.ErrorAs(t, err, &httpErr)

	require.Equal(t, moov.StatusUnauthorized, httpErr.Status())
}

func Test_GetDisputeEvidence_NotFound(t *testing.T) {
	mc := NewTestClient(t)

	// We don't have any disputes to test against! So using a random id to update evidence will return unauthorized since we dont own the resource
	disputeID := uuid.NewString()
	evidenceID := uuid.NewString()
	evidence, err := mc.GetDisputeEvidence(context.Background(), disputeID, evidenceID)
	require.Nil(t, evidence)
	require.Error(t, err)

	// find and cast the error into HttpCallError so it can be inspected
	var httpErr moov.HttpCallResponse
	require.ErrorAs(t, err, &httpErr)

	require.Equal(t, moov.StatusUnauthorized, httpErr.Status())
}
