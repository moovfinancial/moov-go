package moov

import (
	"context"
	"fmt"
	"net/http"
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

// MarkInvoicePaid marks an invoice as paid outside of the Moov platform
func (c Client) MarkInvoicePaid(ctx context.Context, accountID, invoiceID string, payment MarkInvoicePaid) (*Invoice, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPut, pathInvoiceMarkPaid, accountID, invoiceID),
		AcceptJson(),
		JsonBody(payment),
	)
	if err != nil {
		return nil, fmt.Errorf("calling http: %w", err)
	}

	return CompletedObjectOrError[Invoice](resp)
}

// TODO(vince,12/15/2025): Uncomment once this is fully supported in production.
// SendInvoice finalizes an invoice, creates and sends a payment link
// func (c Client) SendInvoice(ctx context.Context, accountID, invoiceID string) (*Invoice, error) {
// 	resp, err := c.CallHttp(ctx,
// 		Endpoint(http.MethodPost, pathInvoiceSend, accountID, invoiceID),
// 		AcceptJson(),
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("calling http: %w), err
// 	}

// 	return CompletedObjectOrError[Invoice](resp)
// }

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
