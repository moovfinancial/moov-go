package moov

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
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

	dispute := new(Dispute)

	dec := json.NewDecoder(bytes.NewReader(input))
	dec.DisallowUnknownFields()

	err := dec.Decode(&dispute)
	require.NoError(t, err)

	assert.Equal(t, "ec7e1848-dc80-4ab0-8827-dd7fc0737b43", dispute.DisputeID)
}

func Test_Disputes(t *testing.T) {
	mc := NewTestClient(t)

	disputes, err := mc.ListDisputes(context.Background(), WithDisputeCount(200), WithDisputeSkip(0))
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
	var httpErr HttpCallError
	require.ErrorAs(t, err, &httpErr)

	require.Equal(t, StatusNotFound, httpErr.Status())
}
