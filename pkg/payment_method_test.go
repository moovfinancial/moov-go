package moov

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestPaymentMethodMarshal(t *testing.T) {
	input := []byte(`{
		  "paymentMethodID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43",
		  "paymentMethodType": "moov-wallet",
		  "wallet": {
			"walletID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43"
		  }
		}`)

	paymentMethod := new(PaymentMethod)

	dec := json.NewDecoder(bytes.NewReader(input))
	dec.DisallowUnknownFields()

	err := dec.Decode(&paymentMethod)
	require.NoError(t, err)

	require.Equal(t, "ec7e1848-dc80-4ab0-8827-dd7fc0737b43", paymentMethod.PaymentMethodID)
}

type PaymentMethodTestSuite struct {
	suite.Suite
	// values for testing will be set in init()
	accountID       string
	paymentMethodID string
}

// listen for 'go test' command --> run test methods
func TestPaymentMethodSuite(t *testing.T) {
	suite.Run(t, new(PaymentMethodTestSuite))
}

func (s *PaymentMethodTestSuite) SetupSuite() {
	// Sandbox accounts have a "Lincoln National Corporation" moov account added by default. Get it's AccountID so we can test against it
	mc := NewTestClient(s.T())

	accounts, err := mc.ListAccounts(WithAccountName("Lincoln National Corporation"))
	s.NoError(err)

	defaultAccountName := "Daniella Singh"
	for _, account := range accounts {
		if account.DisaplayName == defaultAccountName {
			// set the accountID for testing
			s.accountID = account.AccountID
		}
	}
}

func (s *PaymentMethodTestSuite) TearDownSuite() {
}

func (s *PaymentMethodTestSuite) TestListPaymentMethods() {
	mc := NewTestClient(s.T())

	paymentMethods, err := mc.ListPaymentMethods(s.accountID, "")
	s.NoError(err)

	s.Require().NotNil(paymentMethods)

	s.paymentMethodID = paymentMethods[0].PaymentMethodID
}

func (s *PaymentMethodTestSuite) TestGetPaymentMethod() {
	mc := NewTestClient(s.T())

	paymentMethodID := s.paymentMethodID
	if paymentMethodID == "" {
		paymentMethodID = "2ce45e4e-8d96-45e4-8658-5767423e098d"
	}

	paymentMethod, err := mc.GetPaymentMethod(s.accountID, paymentMethodID)
	s.NoError(err)

	s.Equal(paymentMethodID, paymentMethod.PaymentMethodID)
}
