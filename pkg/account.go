package moov

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	CreatedOn       time.Time         `json:"-"`
	UpdatedOn       time.Time         `json:"-"`
	DisabledOn      time.Time         `json:"-"`
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

func (p Profile) jsonValue() interface{} {
	if p.Individual != (Individual{}) {
		type aliasIndividual struct {
			Individual interface{} `json:"individual,omitempty"`
		}
		return aliasIndividual{Individual: p.Individual.jsonValue()}
	}
	if p.Business != (Business{}) {
		type aliasBusiness struct {
			Business Business `json:"business,omitempty"`
		}
		return aliasBusiness{Business: p.Business}
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

func (c Client) CreateAccount(account Account) (Account, error) {
	jsonValue, _ := json.Marshal(account)
	req, _ := http.NewRequest(http.MethodPost, "https://api.moov.io/accounts", bytes.NewBuffer(jsonValue))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.Credentials.PublicKey, c.Credentials.SecretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return account, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	respAccount := Account{}

	switch resp.StatusCode {
	case http.StatusOK:
		// Account created
		err = json.Unmarshal(body, &respAccount)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
		}
		return respAccount, nil
	case http.StatusUnauthorized:
		return respAccount, ErrAuthCreditionalsNotSet
	case http.StatusUnprocessableEntity:
		fmt.Println("UnprocessableEntity")
		// TODO: error Account sent to server is missing or malformed
	}
	return respAccount, nil
}

func (c Client) GetAccount(accountID string) (Account, error) {
	respAccount := Account{}
	req, _ := http.NewRequest(http.MethodGet, "https://api.moov.io/accounts/"+accountID, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.Credentials.PublicKey, c.Credentials.SecretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return respAccount, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	switch resp.StatusCode {
	case http.StatusOK:
		// Account created
		err = json.Unmarshal(body, &respAccount)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
		}
		return respAccount, nil
	case http.StatusUnauthorized:
		return respAccount, ErrAuthCreditionalsNotSet
	case http.StatusUnprocessableEntity:
		fmt.Println("UnprocessableEntity")
		// TODO: error Account sent to server is missing or malformed
	}
	return respAccount, nil
}

// DeleteAccount deletes an account.
// TODO: Delete is not currently supported by the api
// https://docs.moov.io/guides/dashboard/accounts/#disconnect-accounts
/** func (c Client) DeleteAccount(accountID string) error {
	req, _ := http.NewRequest(http.MethodDelete, "https://api.moov.io/accounts/"+accountID, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.Credentials.PublicKey, c.Credentials.SecretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)

	switch resp.StatusCode {
	case http.StatusOK:
		// Account created
		return nil
	case http.StatusUnauthorized:
		return ErrAuthCreditionalsNotSet
	case http.StatusUnprocessableEntity:
		fmt.Println("UnprocessableEntity")
	}
	return nil
} **/

/*

// HTTP Client template

	   	params := url.Values{}
	   	params.Add("grant_type", "client_credentials")
	   	params.Add("scope", "/accounts.write")

	   	req, err := http.NewRequest("POST", "https://api.moov.io/oauth2/token?"+params.Encode(), nil)

	   	func (c Client) Ping() {
	   	log.Println("ping")
	   	req, _ := http.NewRequest("POST", "https://api.moov.io/ping", nil)
	   	req.Header.Set("Accept", "application/json")
	   	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	   	//req.SetBasicAuth(c.Credentials.PublicKey, c.Credentials.SecretKey)

	   	client := &http.Client{}
	   	resp, err := client.Do(req)
	   	if err != nil {
	   		// Todo: return an error
	   		log.Fatal(err)
	   	}
	   	defer resp.Body.Close()

	   	body, err := io.ReadAll(resp.Body)
	   	if err != nil {
	   		// Todo: return an error
	   		log.Fatal(err)
	   	}

	   	fmt.Println("response Status:", resp.Status)
	   	fmt.Println("response Headers:", resp.Header)
	   	fmt.Println("response Body:", string(body))

	   }
*/
