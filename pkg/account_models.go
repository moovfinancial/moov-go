package moov

import (
	"encoding/json"
	"time"
)

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
	CAPABILITIES_TRANSFERS                   = "transfers"
	CAPABILITIES_WALLET                      = "wallet"
	CAPABILITIES_SEND_FUNDS                  = "send-funds"
	CAPABILITIES_COLLECT_FUNDS               = "collect-funds"
	CAPABILITIES_CARD_ISSUING                = "card-issuing"
)

// Accounts represent a legal entity (either a business or an individual) in Moov.
type Account struct {
	Mode            string            `json:"mode,omitempty"`
	AccountID       string            `json:"accountID,omitempty"`
	AccountType     string            `json:"accountType,omitempty"`
	DisaplayName    string            `json:"displayName,omitempty"`
	Profile         Profile           `json:"profile,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
	TermsOfService  TermsOfService    `json:"termsOfService,omitempty"`
	Verification    Verification      `json:"verification,omitempty"`
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

func (n Name) jsonValue() interface{} {
	if n.FirstName != "" {
		return n
	}
	return nil
}

type Phone struct {
	Number      string `json:"number,omitempty"`
	CountryCode string `json:"countryCode,omitempty"`
}

func (p Phone) jsonValue() interface{} {
	if p.Number != "" {
		return p
	}
	return nil
}

type Address struct {
	AddressLine1    string `json:"addressLine1,omitempty"`
	AddressLine2    string `json:"addressLine2,omitempty"`
	City            string `json:"city,omitempty"`
	StateOrProvince string `json:"stateOrProvince,omitempty"`
	PostalCode      string `json:"postalCode,omitempty"`
	Country         string `json:"country,omitempty"`
}

func (a Address) jsonValue() interface{} {
	if a.AddressLine1 != "" {
		return a
	}
	return nil
}

type BirthDate struct {
	Day   int `json:"day,omitempty"`
	Month int `json:"month,omitempty"`
	Year  int `json:"year,omitempty"`
}

func (b BirthDate) jsonValue() interface{} {
	if b.Day != 0 && b.Month != 0 && b.Year != 0 {
		return b
	}
	return nil
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

func (g GovernmentID) jsonValue() interface{} {
	if g.Ssn != (Ssn{}) {
		return g.Ssn
	}
	if g.Itin != (Itin{}) {
		return g.Itin
	}
	return nil
}

type Individual struct {
	Name         Name         `json:"name,omitempty"`
	Phone        Phone        `json:"phone,omitempty"`
	Email        string       `json:"email,omitempty"`
	Address      Address      `json:"address,omitempty"`
	BirthDate    BirthDate    `json:"birthDate,omitempty"`
	GovernmentID GovernmentID `json:"governmentID,omitempty"`

	BirthDateProvided    bool `json:"birthDateProvided,omitempty"`
	GovernmentIDProvided bool `json:"governmentIDProvided,omitempty"`
}

func (i Individual) jsonValue() interface{} {
	type Alias Individual

	type AliasWithInterface struct {
		Alias
		Name         interface{} `json:"name,omitempty"`
		Phone        interface{} `json:"phone,omitempty"`
		Address      interface{} `json:"address,omitempty"`
		BirthDate    interface{} `json:"birthDate,omitempty"`
		GovernmentID interface{} `json:"governmentID,omitempty"`
	}

	return (AliasWithInterface{
		Alias:        Alias(i),
		Name:         i.Name.jsonValue(),
		Phone:        i.Phone.jsonValue(),
		Address:      i.Address.jsonValue(),
		BirthDate:    i.BirthDate.jsonValue(),
		GovernmentID: i.GovernmentID.jsonValue(),
	})
}

type Ein struct {
	Number string `json:"number,omitempty"`
}

func (e Ein) jsonValue() interface{} {
	if e.Number != "" {
		return e
	}
	return nil
}

type TaxID struct {
	Ein Ein `json:"ein,omitempty"`
}

func (t TaxID) jsonValue() interface{} {
	if t.Ein != (Ein{}) {
		type Alias struct {
			Ein interface{} `json:"ein,omitempty"`
		}

		alias := Alias{Ein: t.Ein.jsonValue()}
		return alias
	}
	return nil
}

type IndustryCodes struct {
	Naics string `json:"naics,omitempty"`
	Sic   string `json:"sic,omitempty"`
	Mcc   string `json:"mcc,omitempty"`
}

func (i IndustryCodes) jsonValue() interface{} {
	if i.Naics != "" && i.Sic != "" && i.Mcc != "" {
		return i
	}
	return nil
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

func (i Business) jsonValue() interface{} {
	type Alias Business

	type AliasWithInterface struct {
		Alias
		Address       interface{} `json:"address,omitempty"`
		Phone         interface{} `json:"phone,omitempty"`
		TaxID         interface{} `json:"taxID,omitempty"`
		IndustryCodes interface{} `json:"industryCodes,omitempty"`
	}

	return (AliasWithInterface{
		Alias:         Alias(i),
		Address:       i.Address.jsonValue(),
		Phone:         i.Phone.jsonValue(),
		TaxID:         i.TaxID.jsonValue(),
		IndustryCodes: i.IndustryCodes.jsonValue(),
	})
}

type Profile struct {
	Individual Individual `json:"individual,omitempty"`
	Business   Business   `json:"business,omitempty"`
}

func (p Profile) jsonValue() interface{} {
	if p.Individual != (Individual{}) {
		type aliasIndividual struct {
			Individual interface{} `json:"individual,omitempty"`
		}
		return aliasIndividual{Individual: p.Individual.jsonValue()}
	}
	if p.Business != (Business{}) {
		type aliasBusiness struct {
			Business interface{} `json:"business,omitempty"`
		}
		return aliasBusiness{Business: p.Business.jsonValue()}
	}
	return nil
}

type TermsOfService struct {
	TermsOfServiceToken  TermsOfServiceToken  `json:"termsOfServiceToken,omitempty"`
	TermsOfServiceManual TermsOfServiceManual `json:"termsOfServiceManual,omitempty"`
}

func (t TermsOfService) jsonValue() interface{} {
	if t.TermsOfServiceToken != (TermsOfServiceToken{}) {
		return t.TermsOfServiceToken
	}
	if t.TermsOfServiceManual != (TermsOfServiceManual{}) {
		return t.TermsOfServiceManual

	}
	return nil
}

type TermsOfServiceToken struct {
	Token string `json:"token,omitempty"`
}

type TermsOfServiceManual struct {
	AcceptedDate      string `json:"acceptedDate,omitempty"`
	AcceptedIP        string `json:"acceptedIP,omitempty"`
	AcceptedUserAgent string `json:"acceptedUserAgent,omitempty"`
	AcceptedDomain    string `json:"acceptedDomain,omitempty"`
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

func (s Settings) jsonValue() interface{} {
	if s.CardPayment.StatementDescriptor != "" {
		return s.CardPayment
	}
	if s.AchPayment.CompanyName != "" {
		return s.AchPayment
	}
	return nil
}

type Verification struct {
	Status             string      `json:"status,omitempty"`
	VerificationStatus string      `json:"verificationStatus,omitempty"`
	Documents          []Documents `json:"documents,omitempty"`
}

func (v Verification) jsonValue() interface{} {
	if v.Status != "" {
		return v.Status
	}
	if v.VerificationStatus != "" {
		return v.VerificationStatus
	}
	if len(v.Documents) > 0 {
		return v.Documents
	}
	return nil
}

type Documents struct {
	DocumentID  string    `json:"documentID,omitempty"`
	Type        string    `json:"type,omitempty"`
	ContentType string    `json:"contentType,omitempty"`
	ParseErrors []string  `json:"parseErrors,omitempty"`
	UploadedAt  time.Time `json:"uploadedAt,omitempty"`
}

func (a Account) MarshalJSON() ([]byte, error) {
	// Alias is an alias type of Account to avoid recursion.
	type Alias Account

	type AliasWithInterface struct {
		Alias
		TermsOfService  interface{} `json:"termsOfService,omitempty"`
		Verification    interface{} `json:"verification,omitempty"`
		CustomerSupport interface{} `json:"customerSupport,omitempty"`
		Profile         interface{} `json:"profile,omitempty"`
		Settings        interface{} `json:"settings,omitempty"`
	}

	return json.Marshal(AliasWithInterface{
		Alias:          Alias(a),
		Verification:   a.Verification.jsonValue(),
		Profile:        a.Profile.jsonValue(),
		TermsOfService: a.TermsOfService.jsonValue(),
		Settings:       a.Settings.jsonValue(),
	})
}
