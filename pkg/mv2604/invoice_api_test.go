package mv2604_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2604"
)

func TestCreateInvoiceNotification(t *testing.T) {
	client := mv2604.NewInvoiceClient(newTestClient(t, &expectedRequest{
		method:   http.MethodPost,
		path:     "/accounts/a/invoices/i/notifications",
		jsonBody: `{"channels":["email"]}`,
	}, http.StatusNoContent, ""))

	err := client.CreateInvoiceNotification(t.Context(), "a", "i", mv2604.CreateInvoiceNotification{
		Channels: []mv2604.NotificationChannel{mv2604.NotificationChannelEmail},
	})
	require.NoError(t, err)
}

func TestCreateInvoiceNotificationRequiresInvoiceID(t *testing.T) {
	client := mv2604.NewInvoiceClient(newTestClient(t, nil, http.StatusBadRequest, ""))

	err := client.CreateInvoiceNotification(t.Context(), "a", "", mv2604.CreateInvoiceNotification{})
	require.ErrorIs(t, err, moov.ErrInvoiceIDRequired)
}
