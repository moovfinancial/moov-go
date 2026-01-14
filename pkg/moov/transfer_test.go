package moov_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

		createTrasfer := moov.CreateTransfer{
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
		}

		builder := mc.CreateTransfer(BgCtx(), FACILITATOR_ID, createTrasfer)
		// transfers may take a moment to be aware of uploaded image
		var completed *moov.Transfer
		success := assert.Eventually(t, func() bool {
			var err error
			completed, _, err = builder.WaitForRailResponse()
			return err == nil
		}, 5*time.Second, 250*time.Millisecond)
		require.True(t, success)
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
