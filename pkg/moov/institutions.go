package moov

import (
	"context"
	"net/http"
)

type FinancialInstitutions struct {
	AchParticipants  []AchParticipant  `json:"achParticipants,omitempty"`
	WireParticipants []WireParticipant `json:"wireParticipants,omitempty"`
}

type AchParticipant struct {
	RoutingNumber      string           `json:"routingNumber,omitempty"`
	OfficeCode         string           `json:"officeCode,omitempty"`
	ServicingFRBNumber string           `json:"servicingFRBNumber,omitempty"`
	RecordTypeCode     string           `json:"recordTypeCode,omitempty"`
	Revised            string           `json:"revised,omitempty"`
	NewRoutingNumber   string           `json:"newRoutingNumber,omitempty"`
	CustomerName       string           `json:"customerName,omitempty"`
	PhoneNumber        string           `json:"phoneNumber,omitempty"`
	StatusCode         string           `json:"statusCode,omitempty"`
	ViewCode           string           `json:"viewCode,omitempty"`
	CleanName          string           `json:"cleanName,omitempty"`
	AchLocation        AchLocation      `json:"achLocation,omitempty"`
	Logo               *InstitutionLogo `json:"logo,omitempty"`
}

type InstitutionLogo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type AchLocation struct {
	Address             string `json:"address,omitempty"`
	City                string `json:"city,omitempty"`
	State               string `json:"state,omitempty"`
	PostalCode          string `json:"postalCode,omitempty"`
	PostalCodeExtension string `json:"postalCodeExtension,omitempty"`
}

type WireParticipant struct {
	RoutingNumber                     string       `json:"routingNumber,omitempty"`
	TelegraphicName                   string       `json:"telegraphicName,omitempty"`
	CustomerName                      string       `json:"customerName,omitempty"`
	Location                          WireLocation `json:"location,omitempty"`
	FundsTransferStatus               string       `json:"fundsTransferStatus,omitempty"`
	FundsSettlementOnlyStatus         string       `json:"fundsSettlementOnlyStatus,omitempty"`
	BookEntrySecuritiesTransferStatus string       `json:"bookEntrySecuritiesTransferStatus,omitempty"`
	Date                              string       `json:"date,omitempty"`
}

type WireLocation struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

type Rail string

var (
	RailAch  Rail = "ach"
	RailWire Rail = "wire"
)

type ListInstitutionsFailter callArg

// WithInstitutionName filters institutions by their name
func WithInstitutionName(name string) ListTransactionFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["name"] = name
		return nil
	})
}

// WithInstitutionRoutingNumber filters institutions by their name
func WithInstitutionRoutingNumber(routingNumber string) ListTransactionFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["routingNumber"] = routingNumber
		return nil
	})
}

// WithInstitutionState filters institutions by their state
func WithInstitutionState(state string) ListTransactionFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["state"] = state
		return nil
	})
}

// WithInstitutionLimit filters institutions by their name
func WithInstitutionLimit(limit int) ListTransactionFilter {
	return Limit(limit)
}

func (c Client) ListInstitutions(ctx context.Context, rail Rail, opts ...ListInstitutionsFailter) (*FinancialInstitutions, error) {
	args := prependArgs(opts, AcceptJson())
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathInstitutions, string(rail)), args...)
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[FinancialInstitutions](resp)
}
