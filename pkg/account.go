package moov

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
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

// CreateAccount creates a new account.
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
			return respAccount, err
		}
		return respAccount, nil
	case http.StatusUnauthorized:
		return respAccount, ErrAuthCredentialsNotSet
	case http.StatusUnprocessableEntity:
		return respAccount, ErrBadRequest
	case http.StatusTooManyRequests:
		return respAccount, ErrRateLimit
	}
	return respAccount, ErrDefault(resp.StatusCode)
}

// GetAccount returns an account based on accountID.
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
		err = json.Unmarshal(body, &respAccount)
		if err != nil {
			return respAccount, ErrDefault(resp.StatusCode)
		}
		return respAccount, nil
	case http.StatusUnauthorized:
		return respAccount, ErrAuthCredentialsNotSet
	case http.StatusUnprocessableEntity:
		return respAccount, ErrBadRequest
	case http.StatusTooManyRequests:
		return respAccount, ErrRateLimit
	}
	return respAccount, ErrDefault(resp.StatusCode)
}

// UpdateAccount updates an account.
func (c Client) UpdateAccount(account Account) (Account, error) {
	jsonValue, _ := json.Marshal(account)
	req, _ := http.NewRequest(http.MethodPatch, "https://api.moov.io/accounts/"+account.AccountID, bytes.NewBuffer(jsonValue))
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
		// Account Updated
		err = json.Unmarshal(body, &respAccount)
		if err != nil {
			return respAccount, ErrDefault(resp.StatusCode)
		}
		return respAccount, nil
	case http.StatusUnauthorized:
		return respAccount, ErrAuthCredentialsNotSet
	case http.StatusUnprocessableEntity:
		return respAccount, ErrBadRequest
	case http.StatusTooManyRequests:
		return respAccount, ErrRateLimit
	}
	return respAccount, ErrDefault(resp.StatusCode)
}

type AccountConfigurable func(c *Client)

type AccountParams struct {
	// Name If provided, this query will attempt to find matches against the following Account and Profile fields: diplayName, firstName, middleName, lastName, legalBusinessName
	Name string `url:"name"`
	// Email Filter connected accounts by email address.
	Email string `url:"email"`
	// Type Possible values: individual, business
	Type string `url:"type"`
	// ForeignID Serves as an optional alias from a foreign/external system which can be used to reference this resource.
	ForeignID string `url:"foreignID"`
	// VerificationStatus Possible values: unverified, pending, resubmit, review, verified, failed
	VerificationStatus string `url:"verification_status"`
	// includeDisconnected If true, the response will include disconnected accounts.
	IncludeDisconnected bool `url:"includeDisconnected"`
	// Count Optional parameter to limit the number of results in the query. Default is 20
	Count int `url:"count"`
	// Skip The number of items to offset before starting to collect the result set
	Skip int `url:"skip"`
}

func WithAccountName(name string) AccountConfigurable {
	return func(c *Client) {
		c.AccountParams.Name = name
	}
}

// ListAccounts returns a list of accounts.
func (c Client) ListAccounts(opts ...AccountConfigurable) ([]Account, error) {
	// Apply all the configurable functions to the client
	for _, opt := range opts {
		opt(&c)
	}

	respAccounts := []Account{}
	req, _ := http.NewRequest(http.MethodGet, "https://api.moov.io/accounts", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.Credentials.PublicKey, c.Credentials.SecretKey)

	q := req.URL.Query()
	// Add a name query if it was provided
	if c.AccountParams.Name != "" {
		q.Add("name", c.AccountParams.Name)
	}
	// Add a email query if it was provided
	if c.AccountParams.Email != "" {
		q.Add("email", c.AccountParams.Email)
	}
	// Add a type query if it was provided
	if c.AccountParams.Type != "" {
		q.Add("type", c.AccountParams.Type)
	}
	// Add a foreignID query if it was provided
	if c.AccountParams.ForeignID != "" {
		q.Add("foreignID", c.AccountParams.ForeignID)
	}
	// Add a verification_status query if it was provided
	if c.AccountParams.VerificationStatus != "" {
		q.Add("verification_status", c.AccountParams.VerificationStatus)
	}
	// Add a includeDisconnected query if it was provided
	if c.AccountParams.IncludeDisconnected {
		q.Add("includeDisconnected", "true")
	}
	// Add a count query if it was provided
	if c.AccountParams.Count != 0 {
		q.Add("count", strconv.Itoa(c.AccountParams.Count))
	}
	// Add a skip query if it was provided
	if c.AccountParams.Skip != 0 {
		q.Add("skip", strconv.Itoa(c.AccountParams.Skip))
	}
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return respAccounts, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	switch resp.StatusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &respAccounts)
		if err != nil {
			return respAccounts, ErrDefault(resp.StatusCode)
		}
		return respAccounts, nil
	case http.StatusUnauthorized:
		return respAccounts, ErrAuthCredentialsNotSet
	case http.StatusUnprocessableEntity:
		return respAccounts, ErrBadRequest
	case http.StatusTooManyRequests:
		return respAccounts, ErrRateLimit
	}
	return respAccounts, ErrDefault(resp.StatusCode)
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

	switch resp.StatusCode {
	case http.StatusOK:
		// Account created
		return nil
	case http.StatusUnauthorized:
		return ErrAuthCredentialsNotSet
	case http.StatusUnprocessableEntity:
		log.Println("UnprocessableEntity")
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

	   	log.Println("response Status:", resp.Status)
	   	log.Println("response Headers:", resp.Header)
	   	log.Println("response Body:", string(body))

	   }
*/
