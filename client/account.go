package moov

import "time"

const (
	INDIVIDUAL                               = "individual"
	BUSINESS                                 = "business"
	BUSINESS_TYPE_SOLEPROPRIETORSHIOP        = "soleProprietorship"
	BUSINESS_TYPE_UNINCORPORATED_ASSOCIATION = "unincorporatedAssociation"
	BUSINESS_TYPE_TRUST                      = "trust"
	BUSINESS_TYPE_PUBLIC_CORPORATION         = "publicCorporation"
	BUSINESS_TYPE_PRIVATE_CORPORATION        = "privateCorporation"
	BUSINESS_TYPE_LLC                        = "llc"
	BUSINESS_TYPE_PARTNERSHIP                = "partnership"
	BUSINESS_TYPE_UNINCORPORATED_NONPROFIT   = "unincorporatedNonProfit"
	BUSINESS_TYPE_INCORPORATED_NONPROFIT     = "incorporatedNonProfit"
)

// Accounts represent a legal entity (either a business or an individual) in Moov.
type Account struct {
	Mode            string            `json:"mode,omitempty"`
	AccountType     string            `json:"accountType,omitempty"`
	Profile         Profile           `json:"profile,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
	TermsOfService  TermsOfService    `json:"termsOfService,omitempty"`
	ForeignID       string            `json:"foreignID,omitempty"`
	CustomerSupport CustomerSupport   `json:"customerSupport,omitempty"`
	Settings        Settings          `json:"settings,omitempty"`
	Capabilities    []string          `json:"capabilities,omitempty"`
	CreatedOn       time.Time         `json:"createdOn,omitempty"`
	UpdatedOn       time.Time         `json:"updatedOn,omitempty"`
	DisabledOn      time.Time         `json:"disabledOn,omitempty"`
}

type Name struct {
	FirstName  string `json:"firstName,omitempty"`
	MiddleName string `json:"middleName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
	Suffix     string `json:"suffix,omitempty"`
}

type Phone struct {
	Number      string `json:"number,omitempty"`
	CountryCode string `json:"countryCode,omitempty"`
}

type Address struct {
	AddressLine1    string `json:"addressLine1,omitempty"`
	AddressLine2    string `json:"addressLine2,omitempty"`
	City            string `json:"city,omitempty"`
	StateOrProvince string `json:"stateOrProvince,omitempty"`
	PostalCode      string `json:"postalCode,omitempty"`
	Country         string `json:"country,omitempty"`
}

type BirthDate struct {
	Day   int `json:"day,omitempty"`
	Month int `json:"month,omitempty"`
	Year  int `json:"year,omitempty"`
}

type Ssn struct {
	Full     string `json:"full,omitempty"`
	LastFour string `json:"lastFour,omitempty"`
}

type Itin struct {
	Full     string `json:"full,omitempty"`
	LastFour string `json:"lastFour,omitempty"`
}

type GovernmentID struct {
	Ssn  Ssn  `json:"ssn,omitempty"`
	Itin Itin `json:"itin,omitempty"`
}

type Individual struct {
	Name         Name         `json:"name,omitempty"`
	Phone        Phone        `json:"phone,omitempty"`
	Email        string       `json:"email,omitempty"`
	Address      Address      `json:"address,omitempty"`
	BirthDate    BirthDate    `json:"birthDate,omitempty"`
	GovernmentID GovernmentID `json:"governmentID,omitempty"`
}

type Ein struct {
	Number string `json:"number,omitempty"`
}

type TaxID struct {
	Ein Ein `json:"ein,omitempty"`
}

type IndustryCodes struct {
	Naics string `json:"naics,omitempty"`
	Sic   string `json:"sic,omitempty"`
	Mcc   string `json:"mcc,omitempty"`
}

type Business struct {
	LegalBusinessName string        `json:"legalBusinessName,omitempty"`
	DoingBusinessAs   string        `json:"doingBusinessAs,omitempty"`
	BusinessType      string        `json:"businessType,omitempty"`
	Address           Address       `json:"address,omitempty"`
	Phone             Phone         `json:"phone,omitempty"`
	Email             string        `json:"email,omitempty"`
	Website           string        `json:"website,omitempty"`
	Description       string        `json:"description,omitempty"`
	TaxID             TaxID         `json:"taxID,omitempty"`
	IndustryCodes     IndustryCodes `json:"industryCodes,omitempty"`
}

type Profile struct {
	Individual Individual `json:"individual,omitempty"`
	Business   Business   `json:"business,omitempty"`
}

type TermsOfService struct {
	TermsOfServiceToken  TermsOfServiceToken  `json:"termsOfServiceToken,omitempty"`
	TermsOfServiceManual TermsOfServiceManual `json:"termsOfServiceManual,omitempty"`
}

type TermsOfServiceToken struct {
	Token string `json:"token,omitempty"`
}

type TermsOfServiceManual struct {
	AcceptedDate      string `json:"acceptedDate,omitempty"`
	AcceptedIP        string `json:"acceptedIP,omitempty"`
	AcceptedUserAgent string `json:"acceptedUserAgent,omitempty"`
	AccepetedDomain   string `json:"acceptedDomain,omitempty"`
}

type CustomerSupport struct {
	Phone   Phone   `json:"phone,omitempty"`
	Email   string  `json:"email,omitempty"`
	Address Address `json:"address,omitempty"`
	Website string  `json:"website,omitempty"`
}

type CardPayment struct {
	StatementDescriptor string `json:"statementDescriptor,omitempty"`
}

type AchPayment struct {
	CompanyName string `json:"companyName,omitempty"`
}

type Settings struct {
	CardPayment CardPayment `json:"cardPayment,omitempty"`
	AchPayment  AchPayment  `json:"achPayment,omitempty"`
}
