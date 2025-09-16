package mv2507

import (
	"time"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

var Accounts = moov.AccountClient[CreateAccount, Account]{Version: moov.Version2025_07}

type CreateAccount struct {
	Type                  CreateAccountType           `json:"accountType"`
	Profile               moov.CreateProfile          `json:"profile"`
	Metadata              map[string]string           `json:"metadata,omitempty"`
	TermsOfService        *moov.TermsOfServicePayload `json:"termsOfService,omitempty"`
	ForeignID             string                      `json:"foreignID,omitempty"`
	CustomerSupport       *moov.CustomerSupport       `json:"customerSupport,omitempty"`
	AccountSettings       *moov.AccountSettings       `json:"settings,omitempty"`
	RequestedCapabilities []moov.CapabilityName       `json:"capabilities,omitempty"`
}

// CreateAccountType The type of entity represented by this account to be created.
type CreateAccountType AccountType

// List of CreateAccountType
const (
	CreateAccountType_Individual CreateAccountType = "individual"
	CreateAccountType_Business   CreateAccountType = "business"
)

// Account Describes a Moov account.
type Account struct {
	Mode        moov.Mode   `json:"mode,omitempty"`
	AccountID   string      `json:"accountID,omitempty"`
	AccountType AccountType `json:"accountType,omitempty"`
	DisplayName string      `json:"displayName,omitempty"`
	Profile     Profile     `json:"profile,omitempty"`
	// Free-form key-value pair list. Useful for storing information that is not captured elsewhere.
	Metadata       map[string]string        `json:"metadata,omitempty"`
	TermsOfService *moov.TermsOfService     `json:"termsOfService,omitempty"`
	Capabilities   []moov.AccountCapability `json:"capabilities,omitempty"`
	Verification   moov.Verification        `json:"verification,omitempty"`
	// Optional alias from a foreign/external system which can be used to reference this resource.
	ForeignID       string                `json:"foreignID,omitempty"`
	CustomerSupport *moov.CustomerSupport `json:"customerSupport,omitempty"`
	Settings        *moov.AccountSettings `json:"settings,omitempty"`
	CreatedOn       time.Time             `json:"createdOn,omitempty"`
	UpdatedOn       time.Time             `json:"updatedOn,omitempty"`
	DisconnectedOn  *time.Time            `json:"disconnectedOn,omitempty"`
}

// AccountType The type of entity represented by this account.
type AccountType string

// List of AccountType
const (
	AccountType_Individual AccountType = "individual"
	AccountType_Business   AccountType = "business"
	AccountType_Guest      AccountType = "guest"
)

// Profile Describes a Moov account profile.
type Profile struct {
	Individual *moov.Individual `json:"individual,omitempty"`
	Business   *moov.Business   `json:"business,omitempty"`
	Guest      *Guest           `json:"guest,omitempty"`
}

// Guest Describes a guest account profile.
type Guest struct {
	Name  string      `json:"name,omitempty"`
	Phone *moov.Phone `json:"phone,omitempty"`
	// Email address.
	Email string `json:"email,omitempty"`
}
