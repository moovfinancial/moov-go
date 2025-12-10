package moov

import (
	"context"
	"net/http"
)

// GetStatement retrieves a statement
// https://docs.moov.io/api/moov-accounts/billing/get-statement/
func (c Client) GetStatement(ctx context.Context, accountID, statementID string) (*Statement, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathStatement, accountID, statementID),
		AcceptJson(),
	)
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Statement](resp)
}

// GetStatementPDF retrieves a statement in PDF format as []byte
// https://docs.moov.io/api/moov-accounts/billing/get-statement/
func (c Client) GetStatementPDF(ctx context.Context, accountID, statementID string) ([]byte, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathStatement, accountID, statementID),
		AcceptContentType("application/pdf"),
	)
	if err != nil {
		return nil, err
	}

	var b []byte
	if err = resp.Unmarshal(&b); err != nil {
		return nil, err
	}
	return b, nil
}

// ListStatements lists statements for a Moov account
// https://docs.moov.io/api/moov-accounts/billing/list-statements/
func (c Client) ListStatements(ctx context.Context, accountID string, filters ...ListStatementFilter) ([]Statement, error) {
	args := prependArgs(filters, AcceptJson())
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathStatements, accountID), args...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[Statement](resp)
}
