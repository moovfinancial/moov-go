package mv2604

import (
	"context"
	"net/http"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

type AccountClient[T any, V any] struct {
	moov.AccountClient[T, V]
}

var Accounts = AccountClient[PatchAccount, moov.Account]{AccountClient: moov.AccountClient[PatchAccount, moov.Account]{Version: moov.Version2026_04}}

type PatchAccount struct {
	ForeignID       string                                `json:"foreignID,omitempty"`
	Profile         PatchProfile                          `json:"profile"`
	Metadata        *moov.Nullable[map[string]string]     `json:"metadata,omitempty"`
	TermsOfService  *moov.Nullable[TermsOfServicePayload] `json:"termsOfService,omitempty"`
	CustomerSupport *moov.Nullable[CustomerSupport]       `json:"customerSupport,omitempty"`
	AccountSettings *moov.Nullable[AccountSettings]       `json:"settings,omitempty"`
}

func (ac AccountClient[T, V]) Patch(ctx context.Context, client moov.Client, accountID string, account PatchAccount) (*V, error) {
	resp, err := client.CallHttp(ctx,
		moov.Endpoint(http.MethodPatch, moov.PathAccount, accountID),
		moov.MoovVersion(ac.Version),
		moov.AcceptJson(),
		moov.JsonBody(account))
	if err != nil {
		return nil, err
	}

	return moov.CompletedObjectOrError[V](resp)
}

type PatchProfile struct {
	Individual *moov.Nullable[PatchIndividualProfile] `json:"individual,omitempty"`
	Business   *moov.Nullable[PatchBusinessProfile]   `json:"business,omitempty"`
}

type PatchIndividualProfile struct {
	Name         moov.Name                         `json:"name"`
	Phone        *moov.Nullable[moov.Phone]        `json:"phone,omitempty"`
	Email        string                            `json:"email"`
	Address      *moov.Nullable[moov.Address]      `json:"address,omitempty"`
	BirthDate    *moov.Nullable[moov.Date]         `json:"birthDate,omitempty"`
	GovernmentID *moov.Nullable[moov.GovernmentID] `json:"governmentID,omitempty"`
}

type PatchBusinessProfile struct {
	Name             string                                `json:"legalBusinessName,omitempty"`
	DBA              string                                `json:"doingBusinessAs,omitempty"`
	Type             *moov.Nullable[moov.BusinessType]     `json:"businessType,omitempty"`
	Address          *moov.Nullable[moov.Address]          `json:"address,omitempty"`
	Phone            *moov.Nullable[moov.Phone]            `json:"phone,omitempty"`
	Email            string                                `json:"email,omitempty"`
	Website          string                                `json:"website,omitempty"`
	Description      string                                `json:"description,omitempty"`
	OwnersProvided   bool                                  `json:"ownersProvided,omitempty"`
	Industry         string                                `json:"industry,omitempty"`
	IndustryCodes    *moov.Nullable[moov.IndustryCodes]    `json:"industryCodes,omitempty"`
	PrimaryRegulator *moov.Nullable[moov.PrimaryRegulator] `json:"primaryRegulator,omitempty"`
	TaxID            *moov.Nullable[moov.TaxID]            `json:"taxID"`
}

type TermsOfServicePayload struct {
	Token  string                                    `json:"token,omitempty"`
	Manual *moov.Nullable[moov.TermsOfServiceManual] `json:"manual,omitempty"`
}

// CustomerSupport User-provided information that can be displayed on credit card transactions for customers to use when contacting a customer support team. This data is only allowed on a business account.
type CustomerSupport struct {
	Phone *moov.Nullable[moov.Phone] `json:"phone,omitempty"`
	// Email address.
	Email   string                       `json:"email,omitempty"`
	Address *moov.Nullable[moov.Address] `json:"address,omitempty"`
	Website string                       `json:"website,omitempty"`
}

// AccountSettings User provided settings to manage an account.
type AccountSettings struct {
	CardPayment *moov.Nullable[moov.CardPaymentSettings] `json:"cardPayment,omitempty"`
	AchPayment  *moov.Nullable[moov.AchPaymentSettings]  `json:"achPayment,omitempty"`
}
