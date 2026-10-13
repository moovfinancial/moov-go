package moov

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
)

// CreatePaymentLink creates a new payment link for a Moov account.
// https://docs.moov.io/api/money-movement/payment-links/create/
func (c Client) CreatePaymentLink(ctx context.Context, accountID string, paymentLink CreatePaymentLink) (*PaymentLink, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathPaymentLinks, accountID),
		MoovVersion(Version2026_07),
		AcceptJson(),
		JsonBody(paymentLink),
	)
	if err != nil {
		return nil, fmt.Errorf("calling http: %w", err)
	}

	return CompletedObjectOrError[PaymentLink](resp)
}

// GetPaymentLink retrieves a payment link by its code.
// https://docs.moov.io/api/money-movement/payment-links/get/
func (c Client) GetPaymentLink(ctx context.Context, accountID, paymentLinkCode string) (*PaymentLink, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathPaymentLink, accountID, paymentLinkCode),
		MoovVersion(Version2026_07),
		AcceptJson(),
	)
	if err != nil {
		return nil, fmt.Errorf("calling http: %w", err)
	}

	return CompletedObjectOrError[PaymentLink](resp)
}

// UpdatePaymentLink updates an existing payment link.
// https://docs.moov.io/api/money-movement/payment-links/update/
func (c Client) UpdatePaymentLink(ctx context.Context, accountID, paymentLinkCode string, paymentLink UpdatePaymentLink) (*PaymentLink, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPatch, pathPaymentLink, accountID, paymentLinkCode),
		MoovVersion(Version2026_07),
		AcceptJson(),
		JsonBody(paymentLink),
	)
	if err != nil {
		return nil, fmt.Errorf("calling http: %w", err)
	}

	return CompletedObjectOrError[PaymentLink](resp)
}

// DisablePaymentLink disables a payment link by its code. Once disabled a payment
// link can no longer be used to collect or send funds.
// https://docs.moov.io/api/money-movement/payment-links/disable/
func (c Client) DisablePaymentLink(ctx context.Context, accountID, paymentLinkCode string) error {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodDelete, pathPaymentLink, accountID, paymentLinkCode),
		MoovVersion(Version2026_07),
	)
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}

// GetPaymentLinkQRCode retrieves the QR code image for a payment link and returns
// the raw image bytes along with the content type (e.g. "image/png").
// https://docs.moov.io/api/money-movement/payment-links/get-qr-code/
func (c Client) GetPaymentLinkQRCode(ctx context.Context, accountID, paymentLinkCode string) ([]byte, string, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathPaymentLinkQRCode, accountID, paymentLinkCode),
		MoovVersion(Version2026_07),
		AcceptContentType("image/*"),
	)
	if err != nil {
		return nil, "", fmt.Errorf("calling http: %w", err)
	}

	buf, err := CompletedObjectOrError[bytes.Buffer](resp)
	if err != nil {
		return nil, "", err
	}

	contentType := ""
	if hcr, ok := resp.(*httpCallResponse); ok {
		contentType = hcr.ContentType()
	}
	return buf.Bytes(), contentType, nil
}

// ListPaymentLinkFilter represents a filter option for listing payment links.
type ListPaymentLinkFilter callArg

// WithPaymentLinkCount sets the number of payment links to return.
func WithPaymentLinkCount(count int) ListPaymentLinkFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["count"] = fmt.Sprintf("%d", count)
		return nil
	})
}

// WithPaymentLinkSkip sets the number of payment links to skip.
func WithPaymentLinkSkip(skip int) ListPaymentLinkFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["skip"] = fmt.Sprintf("%d", skip)
		return nil
	})
}

// WithPaymentLinkTypes filters payment links by type.
func WithPaymentLinkTypes(types ...PaymentLinkType) ListPaymentLinkFilter {
	return callBuilderFn(func(call *callBuilder) error {
		if len(types) == 0 {
			return nil
		}
		var strs []string
		for _, t := range types {
			strs = append(strs, string(t))
		}
		call.params["types"] = strings.Join(strs, ",")
		return nil
	})
}

// ListPaymentLinks lists the payment links for a Moov account.
// https://docs.moov.io/api/money-movement/payment-links/list/
func (c Client) ListPaymentLinks(ctx context.Context, accountID string, filters ...ListPaymentLinkFilter) ([]PaymentLink, error) {
	args := prependArgs(filters, MoovVersion(Version2026_07), AcceptJson())
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathPaymentLinks, accountID), args...)
	if err != nil {
		return nil, fmt.Errorf("calling http: %w", err)
	}

	return CompletedListOrError[PaymentLink](resp)
}
