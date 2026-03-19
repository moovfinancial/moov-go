package moov_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/internal/testtools"
	"github.com/moovfinancial/moov-go/pkg/moov"
)

func Test_Invoice_CreateUpdateGet(t *testing.T) {
	var (
		mc                = NewTestClient(t)
		ctx               = t.Context()
		accountID         = testtools.MERCHANT_ID
		customerAccountID = testtools.PARTNER_ID
	)

	create := moov.CreateInvoice{
		CustomerAccountID: customerAccountID,
		Description:       "Testing moov-go",
		LineItems: moov.CreateInvoiceLineItems{
			Items: []moov.CreateInvoiceLineItem{
				{
					Name: "Payment for developing this SDK",
					BasePrice: moov.AmountDecimal{
						Currency:     "USD",
						ValueDecimal: "1.00",
					},
					Quantity: 1,
					Options: []moov.CreateInvoiceLineItemOption{
						{
							Name:     "big TIP",
							Quantity: 1,
							PriceModifier: &moov.AmountDecimal{
								Currency:     "USD",
								ValueDecimal: "3.50",
							},
						},
					},
				},
			},
		},
	}

	// create the invoice
	createdInvoice, err := mc.CreateInvoice(ctx, accountID, create)
	require.NoError(t, err)

	// fetch by ID
	fetchedInvoice, err := mc.GetInvoice(ctx, accountID, createdInvoice.InvoiceID)
	require.NoError(t, err)
	require.Equal(t, createdInvoice, fetchedInvoice)

	// list by status
	listedInvoices, err := mc.ListInvoices(ctx, accountID,
		moov.WithInvoiceStatus(moov.InvoiceStatusDraft),
		moov.WithInvoiceCustomerAccountID(customerAccountID))
	require.NoError(t, err)
	require.Contains(t, listedInvoices, *createdInvoice)

	// Update the dueDate
	now := time.Now().UTC()
	update := moov.UpdateInvoice{DueDate: moov.PtrOf(now)}
	updatedInvoice, err := mc.UpdateInvoice(ctx, accountID, createdInvoice.InvoiceID, update)
	require.NoError(t, err)
	require.Equal(t, now, *updatedInvoice.DueDate)

	// Update invoice status to 'unpaid' to send the invoice to the customer.
	updatedInvoice, err = mc.UpdateInvoice(ctx, accountID, createdInvoice.InvoiceID, moov.UpdateInvoice{Status: moov.PtrOf(moov.InvoiceStatusUnpaid)})
	require.NoError(t, err)
	require.Equal(t, moov.InvoiceStatusUnpaid, updatedInvoice.Status)

	// Create an external payment for the invoice to mark it as paid.
	createdPayment, err := mc.CreateInvoicePayment(ctx, accountID, createdInvoice.InvoiceID, moov.CreateInvoicePayment{
		ForeignID:   moov.PtrOf("abc123"),
		Description: moov.PtrOf("Customer paid with check"),
		Amount:      updatedInvoice.TotalAmount,
	})
	require.NoError(t, err)

	// list payments for the invoice
	payments, err := mc.ListInvoicePayments(ctx, accountID, createdInvoice.InvoiceID)
	require.NoError(t, err)
	require.Len(t, payments, 1)
	latestPayment := payments[0]
	require.Equal(t, *createdPayment, latestPayment)
}

func Test_Invoice_Delete(t *testing.T) {
	var (
		mc                = NewTestClient(t)
		ctx               = t.Context()
		accountID         = testtools.MERCHANT_ID
		customerAccountID = testtools.PARTNER_ID
	)

	create := moov.CreateInvoice{
		CustomerAccountID: customerAccountID,
		Description:       "Invoice to delete",
		LineItems: moov.CreateInvoiceLineItems{
			Items: []moov.CreateInvoiceLineItem{
				{
					Name: "Disposable item",
					BasePrice: moov.AmountDecimal{
						Currency:     "USD",
						ValueDecimal: "10.00",
					},
					Quantity: 1,
				},
			},
		},
	}

	// create a draft invoice
	createdInvoice, err := mc.CreateInvoice(ctx, accountID, create)
	require.NoError(t, err)
	require.Equal(t, moov.InvoiceStatusDraft, createdInvoice.Status)

	// delete the draft invoice
	err = mc.DeleteInvoice(ctx, accountID, createdInvoice.InvoiceID)
	require.NoError(t, err)

	// deleted invoice can still be fetched but has DisabledOn set
	fetchedInvoice, err := mc.GetInvoice(ctx, accountID, createdInvoice.InvoiceID)
	require.NoError(t, err)
	require.NotNil(t, fetchedInvoice.DisabledOn)
}
