package moov_test

import (
	"context"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"

	"github.com/stretchr/testify/suite"
)

type InstitutionTestSuite struct {
	suite.Suite
}

// listen for 'go test' command --> run test methods
func TestInstitutionSuite(t *testing.T) {
	suite.Run(t, new(InstitutionTestSuite))
}

func (s *InstitutionTestSuite) SetupSuite() {}

func (s *InstitutionTestSuite) TearDownSuite() {}

func (s *InstitutionTestSuite) TestSearchInstitutions() {
	mc := NewTestClient(s.T())

	ctx := context.Background()
	resp, err := mc.SearchInstitutions(ctx,
		moov.WithInstitutionName("Chase"),
		moov.WithInstitutionLimit(5),
	)
	s.NoError(err)
	s.NotNil(resp)

	s.Greater(len(resp.Ach), 0)
	s.Greater(len(resp.Rtp), 0)
	s.Greater(len(resp.Wire), 0)

	resp, err = mc.SearchInstitutions(ctx,
		moov.WithInstitutionRoutingNumber("021000021"),
		moov.WithInstitutionLimit(5),
	)
	s.NoError(err)
	s.NotNil(resp)

	s.Greater(len(resp.Ach), 0)
	s.Greater(len(resp.Rtp), 0)
	s.Greater(len(resp.Wire), 0)
}
