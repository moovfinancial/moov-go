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

func (s *InstitutionTestSuite) TestListInstitutions() {
	mc := NewTestClient(s.T())

	ctx := context.Background()
	resp, err := mc.ListInstitutions(ctx, moov.RailAch,
		moov.WithInstitutionName("Chase"),
		moov.WithInstitutionState("FL"),
		moov.WithInstitutionLimit(5),
	)
	s.NoError(err)

	s.Greater(len(resp.AchParticipants), 0)
	s.Len(resp.WireParticipants, 0)

	resp, err = mc.ListInstitutions(ctx, moov.RailAch,
		moov.WithInstitutionRoutingNumber("021000021"),
		moov.WithInstitutionLimit(5),
	)
	s.NoError(err)

	s.Greater(len(resp.AchParticipants), 0)
	s.Len(resp.WireParticipants, 0)
}
