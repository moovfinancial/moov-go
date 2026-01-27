package moov

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// CreateInvoice creates a new invoice for a Moov account
func (c Client) CreateInvoice(ctx context.Context, accountID string, invoice CreateInvoice) (*Invoice, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathInvoices, accountID),
		AcceptJson(),
		JsonBody(invoice),
	)
	if err != nil {
		return nil, fmt.Errorf("calling http: %w", err)
	}

	return StartedObjectOrError[Invoice](resp)
}

// GetInvoice retrieves an invoice by ID
func (c Client) GetInvoice(ctx context.Context, accountID, invoiceID string) (*Invoice, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathInvoice, accountID, invoiceID),
		AcceptJson(),
	)
	if err != nil {
		return nil, fmt.Errorf("calling http: %w", err)
	}

	return CompletedObjectOrError[Invoice](resp)
}

// UpdateInvoice updates an existing invoice
func (c Client) UpdateInvoice(ctx context.Context, accountID, invoiceID string, invoice UpdateInvoice) (*Invoice, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPatch, pathInvoice, accountID, invoiceID),
		AcceptJson(),
		JsonBody(invoice),
	)
	if err != nil {
		return nil, fmt.Errorf("calling http: %w", err)
	}

	return CompletedObjectOrError[Invoice](resp)
}

// CreateInvoicePayment marks an invoice as paid outside of the Moov platform
func (c Client) CreateInvoicePayment(ctx context.Context, accountID, invoiceID string, payment CreateInvoicePayment) (*InvoicePayment, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathInvoicePayments, accountID, invoiceID),
		AcceptJson(),
		JsonBody(payment),
	)
	if err != nil {
		return nil, fmt.Errorf("calling http: %w", err)
	}

	return StartedObjectOrError[InvoicePayment](resp)
}

// ListInvoiceFilter represents a filter option for listing invoices
type ListInvoiceFilter callArg

// WithInvoiceStatus filters invoices by status
func WithInvoiceStatus(status InvoiceStatus) ListInvoiceFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["status"] = string(status)
		return nil
	})
}

// WithInvoiceCustomerAccountID filters invoices by customer account ID
func WithInvoiceCustomerAccountID(customerAccountID string) ListInvoiceFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["customerAccountID"] = customerAccountID
		return nil
	})
}

// WithInvoiceCount sets the number of invoices to return
func WithInvoiceCount(count int) ListInvoiceFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["count"] = fmt.Sprintf("%d", count)
		return nil
	})
}

// WithInvoiceSkip sets the number of invoices to skip
func WithInvoiceSkip(skip int) ListInvoiceFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["skip"] = fmt.Sprintf("%d", skip)
		return nil
	})
}

// WithInvoiceCreatedStartDateTime filters invoices created on or after the specified date/time
func WithInvoiceCreatedStartDateTime(t time.Time) ListInvoiceFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["createdStartDateTime"] = t.Format(time.RFC3339)
		return nil
	})
}

// WithInvoiceCreatedEndDateTime filters invoices created on or before the specified date/time
func WithInvoiceCreatedEndDateTime(t time.Time) ListInvoiceFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["createdEndDateTime"] = t.Format(time.RFC3339)
		return nil
	})
}

// WithInvoiceDueStartDateTime filters invoices with a due date on or after the specified date/time
func WithInvoiceDueStartDateTime(t time.Time) ListInvoiceFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["dueStartDateTime"] = t.Format(time.RFC3339)
		return nil
	})
}

// WithInvoiceDueEndDateTime filters invoices with a due date on or before the specified date/time
func WithInvoiceDueEndDateTime(t time.Time) ListInvoiceFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["dueEndDateTime"] = t.Format(time.RFC3339)
		return nil
	})
}

// ListInvoices lists all invoices for a Moov account
// https://docs.moov.io/api/
func (c Client) ListInvoices(ctx context.Context, accountID string, filters ...ListInvoiceFilter) ([]Invoice, error) {
	args := prependArgs(filters, AcceptJson())
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathInvoices, accountID), args...)
	if err != nil {
		return nil, fmt.Errorf("calling http: %w", err)
	}

	return CompletedListOrError[Invoice](resp)
}

// ListInvoicePayments lists all payments for an invoice for a Moov account
func (c Client) ListInvoicePayments(ctx context.Context, accountID, invoiceID string) ([]InvoicePayment, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathInvoicePayments, accountID, invoiceID))
	if err != nil {
		return nil, fmt.Errorf("calling http: %w", err)
	}

	return CompletedListOrError[InvoicePayment](resp)
}
