package moov

import (
	"time"
)

type CreateAccount struct {
	Type                  AccountType            `json:"accountType"`
	Profile               CreateProfile          `json:"profile"`
	Metadata              map[string]string      `json:"metadata,omitempty"`
	TermsOfService        *TermsOfServicePayload `json:"termsOfService,omitempty"`
	ForeignID             string                 `json:"foreignID,omitempty"`
	CustomerSupport       *CustomerSupport       `json:"customerSupport,omitempty"`
	AccountSettings       *AccountSettings       `json:"settings,omitempty"`
	RequestedCapabilities []CapabilityName       `json:"capabilities,omitempty"`
}

type CreateProfile struct {
	Individual *CreateIndividualProfile `json:"individual,omitempty"`
	Business   *CreateBusinessProfile   `json:"business,omitempty"`
}

type CreateIndividualProfile struct {
	Name         Name          `json:"name"`
	Phone        *Phone        `json:"phone,omitempty"`
	Email        string        `json:"email"`
	Address      *Address      `json:"address,omitempty"`
	BirthDate    *Date         `json:"birthDate,omitempty"`
	GovernmentID *GovernmentID `json:"governmentID,omitempty"`
}

type Date struct {
	Year  int `json:"year,omitempty"`
	Month int `json:"month,omitempty"`
	Day   int `json:"day,omitempty"`
}

type GovernmentID struct {
	SSN  *SSN  `json:"ssn,omitempty"`
	ITIN *ITIN `json:"itin,omitempty"`
}

// ITIN is an Individual Taxpayer Identification Number
type ITIN struct {
	Full     string `json:"full,omitempty"`
	LastFour string `json:"lastFour,omitempty"`
}

// SSN is an Individual's Social Security Number
type SSN struct {
	Full     string `json:"full,omitempty"`
	LastFour string `json:"lastFour,omitempty"`
}

type CreateBusinessProfile struct {
	Name             string            `json:"legalBusinessName,omitempty"`
	DBA              string            `json:"doingBusinessAs,omitempty"`
	Type             BusinessType      `json:"businessType,omitempty"`
	Address          *Address          `json:"address,omitempty"`
	Phone            *Phone            `json:"phone,omitempty"`
	Email            string            `json:"email,omitempty"`
	Website          string            `json:"website,omitempty"`
	Description      string            `json:"description,omitempty"`
	TaxID            *TaxID            `json:"taxID,omitempty"`
	IndustryCodes    *IndustryCodes    `json:"industryCodes,omitempty"`
	PrimaryRegulator *PrimaryRegulator `json:"primaryRegulator,omitempty"`
}

type TaxID struct {
	EIN EIN `json:"ein,omitempty"`
}

type EIN struct {
	Number string `json:"number,omitempty"`
}

type TermsOfServicePayload struct {
	Token  string                `json:"token,omitempty"`
	Manual *TermsOfServiceManual `json:"manual,omitempty"`
}

type TermsOfServiceManual struct {
	AcceptanceIP        string    `json:"acceptedIP"`
	AcceptanceDomain    string    `json:"acceptedDomain"`
	AcceptanceUserAgent string    `json:"acceptedUserAgent"`
	AcceptanceDate      time.Time `json:"acceptedDate"`
}

type AccountCapability struct {
	Capability CapabilityName   `json:"capability,omitempty"`
	Status     CapabilityStatus `json:"status,omitempty"`
}

// Account Describes a Moov account.
type Account struct {
	Mode        Mode        `json:"mode,omitempty"`
	AccountID   string      `json:"accountID,omitempty"`
	AccountType AccountType `json:"accountType,omitempty"`
	DisplayName string      `json:"displayName,omitempty"`
	Profile     Profile     `json:"profile,omitempty"`
	// Free-form key-value pair list. Useful for storing information that is not captured elsewhere.
	Metadata       map[string]string   `json:"metadata,omitempty"`
	TermsOfService *TermsOfService     `json:"termsOfService,omitempty"`
	Capabilities   []AccountCapability `json:"capabilities,omitempty"`
	Verification   Verification        `json:"verification,omitempty"`
	// Optional alias from a foreign/external system which can be used to reference this resource.
	ForeignID       string           `json:"foreignID,omitempty"`
	CustomerSupport *CustomerSupport `json:"customerSupport,omitempty"`
	Settings        *AccountSettings `json:"settings,omitempty"`
	CreatedOn       time.Time        `json:"createdOn,omitempty"`
	UpdatedOn       time.Time        `json:"updatedOn,omitempty"`
	DisconnectedOn  *time.Time       `json:"disconnectedOn,omitempty"`
}

// AccountType The type of entity represented by this account.
type AccountType string

// List of AccountType
const (
	AccountType_Individual AccountType = "individual"
	AccountType_Business   AccountType = "business"
)

// Profile Describes a Moov account profile.
type Profile struct {
	Individual *Individual `json:"individual,omitempty"`
	Business   *Business   `json:"business,omitempty"`
}

// Individual Describes an individual.
type Individual struct {
	Name  Name   `json:"name,omitempty"`
	Phone *Phone `json:"phone,omitempty"`
	// Email address.
	Email   string   `json:"email,omitempty"`
	Address *Address `json:"address,omitempty"`
	// Indicates whether this individual's birth date has been provided.
	BirthDateProvided bool `json:"birthDateProvided,omitempty"`
	// Indicates whether a government ID (SSN, ITIN, etc.) has been provided for this individual.
	GovernmentIDProvided bool `json:"governmentIDProvided,omitempty"`
}

// Name An individual's name.
type Name struct {
	// Name this person was given. This is usually the the same as first name.
	FirstName string `json:"firstName"`
	// Name this person was given. This is usually the the same as first name.
	MiddleName string `json:"middleName,omitempty"`
	// Family name of this person. This is usually the the same as last name.
	LastName string `json:"lastName"`
	// Suffix of a given name.
	Suffix string `json:"suffix,omitempty"`
}

// Phone struct for Phone
type Phone struct {
	Number      string `json:"number,omitempty"`
	CountryCode string `json:"countryCode,omitempty"`
}

// Address struct for Address
type Address struct {
	AddressLine1    string `json:"addressLine1"`
	AddressLine2    string `json:"addressLine2,omitempty"`
	City            string `json:"city"`
	StateOrProvince string `json:"stateOrProvince"`
	PostalCode      string `json:"postalCode"`
	Country         string `json:"country"`
}

// AddressPatch is for partial address updates. Only fields that are non-nil will be updated.
type AddressPatch struct {
	AddressLine1    *string `json:"addressLine1"`
	AddressLine2    *string `json:"addressLine2,omitempty"`
	City            *string `json:"city"`
	StateOrProvince *string `json:"stateOrProvince"`
	PostalCode      *string `json:"postalCode"`
	Country         *string `json:"country"`
}

// Business Describes a business.
type Business struct {
	LegalBusinessName string       `json:"legalBusinessName,omitempty"`
	DoingBusinessAs   string       `json:"doingBusinessAs,omitempty"`
	BusinessType      BusinessType `json:"businessType,omitempty"`
	Address           *Address     `json:"address,omitempty"`
	Phone             *Phone       `json:"phone,omitempty"`
	// Email address.
	Email       string `json:"email,omitempty"`
	Website     string `json:"website,omitempty"`
	Description string `json:"description,omitempty"`
	// Indicates whether a tax ID has been provided for this business.
	TaxIDProvided    bool              `json:"taxIDProvided,omitempty"`
	Representatives  []Representative  `json:"representatives,omitempty"`
	OwnersProvided   bool              `json:"ownersProvided"`
	IndustryCodes    *IndustryCodes    `json:"industryCodes,omitempty"`
	PrimaryRegulator *PrimaryRegulator `json:"primaryRegulator,omitempty"`
}

// BusinessType The type of entity represented by this business.
type BusinessType string

// List of BusinessType
const (
	BusinessType_SoleProprietorship        BusinessType = "soleProprietorship"
	BusinessType_UnincorporatedAssociation BusinessType = "unincorporatedAssociation"
	BusinessType_Trust                     BusinessType = "trust"
	BusinessType_PublicCorporation         BusinessType = "publicCorporation"
	BusinessType_PrivateCorporation        BusinessType = "privateCorporation"
	BusinessType_Llc                       BusinessType = "llc"
	BusinessType_Partnership               BusinessType = "partnership"
	BusinessType_UnincorporatedNonProfit   BusinessType = "unincorporatedNonProfit"
	BusinessType_IncorporatedNonProfit     BusinessType = "incorporatedNonProfit"
)

// Responsibilities Describes the job responsibilities of an individual.
type Responsibilities struct {
	// Indicates whether this individual has significant management responsibilities within the business.
	IsController bool `json:"isController"`
	// If `true`, this field indicates that the individual has a business ownership stake of at least 25% in the business. If the representative does not own at least 25% of the business, this field should be `false`.
	IsOwner bool `json:"isOwner"`
	// The percentage of ownership this individual has in the business (required if `isOwner` is `true`).
	OwnershipPercentage int32  `json:"ownershipPercentage"`
	JobTitle            string `json:"jobTitle"`
}

// IndustryCodes Describes industry specific identifiers.
type IndustryCodes struct {
	Naics string `json:"naics,omitempty"`
	Sic   string `json:"sic,omitempty"`
	Mcc   string `json:"mcc,omitempty"`
}

// PrimaryRegulator If the business is a financial institution, this field describes its primary regulator.
type PrimaryRegulator string

// List of PrimaryRegulator
const (
	PrimaryRegulator_OCC  PrimaryRegulator = "OCC"
	PrimaryRegulator_FDIC PrimaryRegulator = "FDIC"
	PrimaryRegulator_NCUA PrimaryRegulator = "NCUA"
	PrimaryRegulator_FRB  PrimaryRegulator = "FRB"
)

// TermsOfService Describes the acceptance of the Terms of Service.
type TermsOfService struct {
	AcceptedDate time.Time `json:"acceptedDate"`
	AcceptedIP   string    `json:"acceptedIP"`
}

// Verification Describes identity verification status and relevant identity verification documents.
type Verification struct {
	VerificationStatus AccountVerificationStatus `json:"verificationStatus,omitempty"`
	Status             VerificationStatus        `json:"status"`
	Details            VerificationStatusDetails `json:"details,omitempty"`
	Documents          []Document                `json:"documents,omitempty"`
}

// AccountVerificationStatus The status of an identity verification for a profile.
type AccountVerificationStatus string

// List of AccountVerificationStatus
const (
	AccountVerificationStatus_Unverified AccountVerificationStatus = "unverified"
	AccountVerificationStatus_Pending    AccountVerificationStatus = "pending"
	AccountVerificationStatus_Resubmit   AccountVerificationStatus = "resubmit"
	AccountVerificationStatus_Review     AccountVerificationStatus = "review"
	AccountVerificationStatus_Verified   AccountVerificationStatus = "verified"
	AccountVerificationStatus_Failed     AccountVerificationStatus = "failed"
)

// VerificationStatus This field is deprecated but available for use until February 2023.
type VerificationStatus string

// List of VerificationStatus
const (
	VerificationStatus_Unverified VerificationStatus = "unverified"
	VerificationStatus_Pending    VerificationStatus = "pending"
	VerificationStatus_Verified   VerificationStatus = "verified"
	VerificationStatus_Errored    VerificationStatus = "errored"
)

// VerificationStatusDetails This field is deprecated but available for use until February 2023.
type VerificationStatusDetails string

// List of VerificationStatusDetails
const (
	VerificationStatusDetails_FailedAutoVerify                      VerificationStatusDetails = "failedAutoVerify"
	VerificationStatusDetails_DocumentDobMismatch                   VerificationStatusDetails = "docDobMismatch"
	VerificationStatusDetails_DocumentNameMismatch                  VerificationStatusDetails = "docNameMismatch"
	VerificationStatusDetails_DocumentAddressMismatch               VerificationStatusDetails = "docAddressMismatch"
	VerificationStatusDetails_DocumentNumberMismatch                VerificationStatusDetails = "docNumberMismatch"
	VerificationStatusDetails_DocumentIncomplete                    VerificationStatusDetails = "docIncomplete"
	VerificationStatusDetails_DocumentFailedRisk                    VerificationStatusDetails = "docFailedRisk"
	VerificationStatusDetails_PotentialAccountSanctionsMatch        VerificationStatusDetails = "potentialAccountSanctionsMatch"
	VerificationStatusDetails_PotentialRepresentativeSanctionsMatch VerificationStatusDetails = "potentialRepresentativeSanctionsMatch"
	VerificationStatusDetails_FailedOther                           VerificationStatusDetails = "failedOther"
)

// Document Describes an uploaded file.
type Document struct {
	// A unique identifier for this document.
	DocumentID  string `json:"documentID"`
	Type        string `json:"type"`
	ContentType string `json:"contentType"`
	// Optional array of errors encountered dring automated parsing.
	ParseErrors []string  `json:"parseErrors,omitempty"`
	UploadedAt  time.Time `json:"uploadedAt"`
}

// CustomerSupport User-provided information that can be displayed on credit card transactions for customers to use when contacting a customer support team. This data is only allowed on a business account.
type CustomerSupport struct {
	Phone *Phone `json:"phone,omitempty"`
	// Email address.
	Email   string   `json:"email,omitempty"`
	Address *Address `json:"address,omitempty"`
	Website string   `json:"website,omitempty"`
}

// AccountSettings User provided settings to manage an account.
type AccountSettings struct {
	CardPayment *CardPaymentSettings `json:"cardPayment,omitempty"`
	AchPayment  *AchPaymentSettings  `json:"achPayment,omitempty"`
}

// CardPaymentSettings User provided settings to manage card payments. This data is only allowed on a business account.
type CardPaymentSettings struct {
	// The description that shows up on credit card transactions. This will default to the accounts display name on account creation.
	StatementDescriptor string `json:"statementDescriptor,omitempty"`
}

// AchPaymentSettings User provided settings to manage ACH payments.
type AchPaymentSettings struct {
	// The description that shows up on ACH transactions. This will default to the account's display name on account creation.
	CompanyName string `json:"companyName,omitempty"`
}
