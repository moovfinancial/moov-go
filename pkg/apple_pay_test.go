package moov

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestApplePayMarshal(t *testing.T) {
	input := []byte(`{
		"brand": "Discover",
		"cardType": "debit",
		"cardDisplayName": "Visa 1234",
		"fingerprint": "9948962d92a1ce40c9f918cd9ece3a22bde62fb325a2f1fe2e833969de672ba3",
		"expiration": {
		  "month": "01",
		  "year": "21"
		},
		"dynamicLastFour": "1234"
	  }`)

	applePay := new(ApplePay)

	dec := json.NewDecoder(bytes.NewReader(input))
	dec.DisallowUnknownFields()

	err := dec.Decode(&applePay)

	require.NoError(t, err)
	assert.Equal(t, "Discover", applePay.Brand)
	assert.Equal(t, "debit", applePay.CardType)
	assert.Equal(t, "Visa 1234", applePay.CardDisplayName)
}

type ApplePayTestSuite struct {
	suite.Suite
	// values for testing will be set in init()
	accountID string
}

// listen for 'go test' command --> run test methods
func TestApplePaySuite(t *testing.T) {
	suite.Run(t, new(ApplePayTestSuite))
}

func (s *ApplePayTestSuite) SetupSuite() {
	// Sandbox accounts have a "Lincoln National Corporation" moov account added by default. Get it's AccountID so we can test against it
	mc := NewTestClient(s.T())

	accounts, err := mc.ListAccounts(context.Background(), WithAccountName("Lincoln National Corporation"))
	s.NoError(err)

	for _, account := range accounts {
		if account.DisaplayName == "Lincoln National Corporation" {
			// set the accountID for testing
			s.accountID = account.AccountID
		}
	}
}

func (s *ApplePayTestSuite) TearDownSuite() {
}

func (s *ApplePayTestSuite) TestCreateApplePayDomain() {
	mc := NewTestClient(s.T())

	domains := []string{"checkout.classbooker.dev"}
	resp, err := mc.CreateApplePayDomain(s.accountID, "Example Merchant", domains)

	s.NoError(err)
	assert.Equal(s.T(), domains, resp.Domains)
}

func (s *ApplePayTestSuite) TestUpdateApplePayDomain() {
	mc := NewTestClient(s.T())

	addDomains := []string{"pay.classbooker.dev"}
	removeDomains := []string{"checkout.classbooker.dev"}

	err := mc.UpdateApplePayDomain(s.accountID, addDomains, removeDomains)
	assert.Nil(s.T(), err)
}

func (s *ApplePayTestSuite) TestGetApplePayDomain() {
	mc := NewTestClient(s.T())

	resp, err := mc.GetApplePayDomain(s.accountID)

	s.NoError(err)
	assert.NotNil(s.T(), resp.Domains)
}

func (s *ApplePayTestSuite) TestCreateApplePaySession() {
	mc := NewTestClient(s.T())

	err := mc.CreateApplePaySession(s.accountID, "checkout.classbooker.dev", "Example Merchant")
	s.NoError(err)
}

func (s *ApplePayTestSuite) TestApplePayToken() {
	mc := NewTestClient(s.T())

	token := ApplePayToken{
		PaymentData: ApplePaymentData{
			Version:   "EC_v1",
			Data:      "3+f4oOTwPa6f1UZ6tG...CE=",
			Signature: "MIAGCSqGSIb3DQ.AAAA==",
			Header: ApplePaymentDataHeader{
				EphemeralPublicKey: "MFkwEK...Md==",
				PublicKeyHash:      "l0CnXdMv...D1I=",
				TransactionId:      "32b...4f3",
			},
		},
		PaymentMethod: ApplePaymentMethod{
			DisplayName: "Visa 1234",
			Network:     "Visa",
			Type:        "debit",
		},
		TransactionIdentifier: "32b...4f3",
	}

	address := ApplePayBillingContact{
		AddressLines: []string{
			"123 Sesame Street",
		},
		Locality:           "Phoenix",
		PostalCode:         "30345",
		AdministrativeArea: "AZ",
		CountryCode:        "US",
	}

	resp, err := mc.ApplePayToken(s.accountID, token, address)

	s.NoError(err)
	assert.NotNil(s.T(), resp.PaymentMethodID)
}
