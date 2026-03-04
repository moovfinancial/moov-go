package moov_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

func Test_Transfers(t *testing.T) {
	mc := NewTestClient(t)

	account := getLincolnBank(t, mc)

	options, err := mc.TransferOptions(BgCtx(), FACILITATOR_ID, moov.CreateTransferOptions{
		Source: moov.CreateTransferOptionsTarget{
			AccountID: account.AccountID,
		},
		Destination: moov.CreateTransferOptionsTarget{
			AccountID: FACILITATOR_ID,
		},
		Amount: moov.Amount{
			Currency: "USD",
			Value:    1,
		},
	})

	NoResponseError(t, err)
	source, dest := paymentMethodsFromOptions(t, options, moov.PaymentMethodType_AchDebitFund, moov.PaymentMethodType_MoovWallet)

	t.Run("make async transfer", func(t *testing.T) {
		started, err := mc.CreateTransfer(BgCtx(),
			FACILITATOR_ID,
			moov.CreateTransfer{
				Source: moov.CreateTransfer_Source{
					PaymentMethodID: source,
				},
				Destination: moov.CreateTransfer_Destination{
					PaymentMethodID: dest,
				},
				Amount: moov.Amount{
					Currency: "usd",
					Value:    1,
				},
			}).Started()
		NoResponseError(t, err)

		// We made an async transfer, so completed should be nil, while started not nil
		require.NotNil(t, started)
	})

	t.Run("make sync transfer", func(t *testing.T) {
		completed, started, err := mc.CreateTransfer(BgCtx(),
			FACILITATOR_ID,
			moov.CreateTransfer{
				Source: moov.CreateTransfer_Source{
					PaymentMethodID: source,
				},
				Destination: moov.CreateTransfer_Destination{
					PaymentMethodID: dest,
				},
				Amount: moov.Amount{
					Currency: "usd",
					Value:    1,
				},
			}).WaitForRailResponse()
		NoResponseError(t, err)

		// We made a sync transfer, so completed should be not nil, while started is nil
		require.NotNil(t, completed)
		require.Nil(t, started)
	})

	t.Run("list transfers", func(t *testing.T) {
		transfers, err := mc.ListTransfers(BgCtx(),
			account.AccountID,
			moov.WithTransferAccountIDs([]string{
				account.AccountID,
			}),
			// moov.WithTransferCount(10),
			// moov.WithTransferStatus(string(moov.TransferStatus_Completed)),
		)
		require.NoError(t, err)
		require.Greater(t, len(transfers), 0)
	})

	t.Run("with line items", func(t *testing.T) {
		// upload an image to reference in line items
		_, imgReader := randomImage(t, 100, 100, encodePNG)
		metadata := &moov.ImageMetadataRequest{
			AltText: moov.PtrOf("Test image from moov-go SDK"),
		}
		uploaded, err := mc.UploadImage(BgCtx(), account.AccountID, imgReader, metadata)
		require.NoError(t, err)
		require.NotNil(t, uploaded)
		t.Cleanup(func() {
			_ = mc.DeleteImage(BgCtx(), account.AccountID, uploaded.ImageID)
		})

		completed, _, err := mc.CreateTransfer(BgCtx(),
			FACILITATOR_ID,
			moov.CreateTransfer{
				Source: moov.CreateTransfer_Source{
					PaymentMethodID: source,
				},
				Destination: moov.CreateTransfer_Destination{
					PaymentMethodID: dest,
				},
				Amount: moov.Amount{
					Currency: "USD",
					Value:    11,
				},
				LineItems: &moov.CreateTransferLineItems{
					Items: []moov.CreateTransferLineItem{
						{
							Name:      "Bagel",
							ProductID: moov.PtrOf("1e262367-de3e-4acb-ae02-7f56e83632ee"),
							BasePrice: moov.AmountDecimal{
								Currency:     "USD",
								ValueDecimal: "0.03",
							},
							ImageIDs: []string{uploaded.ImageID},
							Options: []moov.CreateTransferLineItemOption{
								{
									Name:     "Everything Bagel",
									Quantity: 1,
									ImageIDs: []string{uploaded.ImageID},
								},
								{
									Group: moov.PtrOf("Toppings"),
									Name:  "Cream Cheese",
									PriceModifier: &moov.AmountDecimal{
										Currency:     "USD",
										ValueDecimal: "0.01",
									},
									Quantity: 2,
								},
							},
							Quantity: 2,
						},
						{
							Name: "Water",
							BasePrice: moov.AmountDecimal{
								Currency:     "USD",
								ValueDecimal: "0.01",
							},
							Quantity: 1,
						},
					},
				},
			}).WaitForRailResponse()
		NoResponseError(t, err)

		require.NotNil(t, completed)
		wantLineItems := &moov.TransferLineItems{
			Items: []moov.TransferLineItem{
				{
					Name:      "Bagel",
					ProductID: moov.PtrOf("1e262367-de3e-4acb-ae02-7f56e83632ee"),
					BasePrice: moov.AmountDecimal{
						Currency:     "USD",
						ValueDecimal: "0.03",
					},
					Images: []moov.TransferLineItemImageMetadata{
						{
							ImageID:  uploaded.ImageID,
							AltText:  uploaded.AltText,
							Link:     uploaded.Link,
							PublicID: uploaded.PublicID,
						},
					},
					Options: []moov.TransferLineItemOption{
						{
							Name:     "Everything Bagel",
							Quantity: 1,
							Images: []moov.TransferLineItemImageMetadata{
								{
									ImageID:  uploaded.ImageID,
									AltText:  uploaded.AltText,
									Link:     uploaded.Link,
									PublicID: uploaded.PublicID,
								},
							},
						},
						{
							Group: moov.PtrOf("Toppings"),
							Name:  "Cream Cheese",
							PriceModifier: &moov.AmountDecimal{
								Currency:     "USD",
								ValueDecimal: "0.01",
							},
							Quantity: 2,
						},
					},
					Quantity: 2,
				},
				{
					Name: "Water",
					BasePrice: moov.AmountDecimal{
						Currency:     "USD",
						ValueDecimal: "0.01",
					},
					Quantity: 1,
				},
			},
		}
		require.Equal(t, wantLineItems, completed.LineItems)
	})
}

func Test_Cancellations(t *testing.T) {
	mc := NewTestClient(t)
	account := getLincolnBank(t, mc)

	options, err := mc.TransferOptions(BgCtx(), FACILITATOR_ID, moov.CreateTransferOptions{
		Source: moov.CreateTransferOptionsTarget{
			AccountID: account.AccountID,
		},
		Destination: moov.CreateTransferOptionsTarget{
			AccountID: FACILITATOR_ID,
		},
		Amount: moov.Amount{
			Currency: "USD",
			Value:    1,
		},
	})

	NoResponseError(t, err)
	source, dest := paymentMethodsFromOptions(t, options, moov.PaymentMethodType_AchDebitFund, moov.PaymentMethodType_MoovWallet)

	var transferID string
	t.Run("make sync transfer", func(t *testing.T) {
		completed, _, err := mc.CreateTransfer(BgCtx(), FACILITATOR_ID, moov.CreateTransfer{
			Source: moov.CreateTransfer_Source{
				PaymentMethodID: source,
			},
			Destination: moov.CreateTransfer_Destination{
				PaymentMethodID: dest,
			},
			Amount: moov.Amount{
				Currency: "usd",
				Value:    1,
			},
		}).WaitForRailResponse()
		NoResponseError(t, err)

		// We made a sync transfer, so completed should be not nil, while started is nil
		require.NotNil(t, completed)
		transferID = completed.TransferID
	})

	t.Run("create and get cancellation", func(t *testing.T) {
		createdCancellation, err := mc.CancelTransfer(BgCtx(), FACILITATOR_ID, transferID)
		NoResponseError(t, err)
		require.NotEmpty(t, createdCancellation)

		fetchedCancellation, err := mc.GetCancellation(BgCtx(), FACILITATOR_ID, transferID, createdCancellation.CancellationID)
		NoResponseError(t, err)
		require.Equal(t, createdCancellation.CancellationID, fetchedCancellation.CancellationID)
	})
}

// newTransferFilterClient creates a local test server that returns an empty
// transfer list and a moov.Client pointed at it.  The returned function
// retrieves the URL query values captured from the most recent request.
func newTransferFilterClient(t *testing.T) (*moov.Client, func() url.Values) {
	t.Helper()
	var lastQuery url.Values

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lastQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("[]"))
	}))
	t.Cleanup(srv.Close)

	u, err := url.Parse(srv.URL)
	require.NoError(t, err)

	client, err := moov.NewClient(
		moov.WithCredentials(moov.Credentials{
			PublicKey: "test-public",
			SecretKey: "test-secret",
			Host:      u.Host,
		}),
		moov.WithMoovURLScheme("http"),
	)
	require.NoError(t, err)

	return client, func() url.Values { return lastQuery }
}

func Test_ListTransferFilters(t *testing.T) {
	const accountID = "test-account-id"

	t.Run("WithTransferCustomerID", func(t *testing.T) {
		client, query := newTransferFilterClient(t)
		_, _ = client.ListTransfers(BgCtx(), accountID, moov.WithTransferCustomerID("cust-123"))
		require.Equal(t, "cust-123", query().Get("customerId"))
	})

	t.Run("WithTransferAccountIDs single", func(t *testing.T) {
		client, query := newTransferFilterClient(t)
		_, _ = client.ListTransfers(BgCtx(), accountID, moov.WithTransferAccountIDs([]string{"acc-1"}))
		require.Equal(t, "acc-1", query().Get("accountIDs"))
	})

	t.Run("WithTransferAccountIDs multiple", func(t *testing.T) {
		client, query := newTransferFilterClient(t)
		_, _ = client.ListTransfers(BgCtx(), accountID, moov.WithTransferAccountIDs([]string{"acc-1", "acc-2", "acc-3"}))
		require.Equal(t, "acc-1,acc-2,acc-3", query().Get("accountIDs"))
	})

	t.Run("WithTransferStatus", func(t *testing.T) {
		client, query := newTransferFilterClient(t)
		_, _ = client.ListTransfers(BgCtx(), accountID, moov.WithTransferStatus("completed"))
		require.Equal(t, "completed", query().Get("status"))
	})

	t.Run("WithTransferStartDate", func(t *testing.T) {
		client, query := newTransferFilterClient(t)
		ts := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
		_, _ = client.ListTransfers(BgCtx(), accountID, moov.WithTransferStartDate(ts))
		require.Equal(t, ts.Format(time.RFC3339), query().Get("startDateTime"))
	})

	t.Run("WithTransferEndDate", func(t *testing.T) {
		client, query := newTransferFilterClient(t)
		ts := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)
		_, _ = client.ListTransfers(BgCtx(), accountID, moov.WithTransferEndDate(ts))
		require.Equal(t, ts.Format(time.RFC3339), query().Get("endDateTime"))
	})

	t.Run("WithTransferGroup", func(t *testing.T) {
		client, query := newTransferFilterClient(t)
		_, _ = client.ListTransfers(BgCtx(), accountID, moov.WithTransferGroup("group-abc"))
		require.Equal(t, "group-abc", query().Get("groupID"))
	})

	t.Run("WithTransferSchedule", func(t *testing.T) {
		client, query := newTransferFilterClient(t)
		_, _ = client.ListTransfers(BgCtx(), accountID, moov.WithTransferSchedule("sched-xyz"))
		require.Equal(t, "sched-xyz", query().Get("scheduleID"))
	})

	t.Run("WithTransferPaymentLinkCode", func(t *testing.T) {
		client, query := newTransferFilterClient(t)
		_, _ = client.ListTransfers(BgCtx(), accountID, moov.WithTransferPaymentLinkCode("LINK42"))
		require.Equal(t, "LINK42", query().Get("paymentLinkCode"))
	})

	t.Run("WithTransferRefunded", func(t *testing.T) {
		client, query := newTransferFilterClient(t)
		_, _ = client.ListTransfers(BgCtx(), accountID, moov.WithTransferRefunded())
		require.Equal(t, "true", query().Get("refunded"))
	})

	t.Run("WithTransferDisputed", func(t *testing.T) {
		client, query := newTransferFilterClient(t)
		_, _ = client.ListTransfers(BgCtx(), accountID, moov.WithTransferDisputed())
		require.Equal(t, "true", query().Get("disputed"))
	})

	t.Run("WithTransferSkip", func(t *testing.T) {
		client, query := newTransferFilterClient(t)
		_, _ = client.ListTransfers(BgCtx(), accountID, moov.WithTransferSkip(25))
		require.Equal(t, "25", query().Get("skip"))
	})

	t.Run("WithTransferCount", func(t *testing.T) {
		client, query := newTransferFilterClient(t)
		_, _ = client.ListTransfers(BgCtx(), accountID, moov.WithTransferCount(100))
		require.Equal(t, "100", query().Get("count"))
	})
}
