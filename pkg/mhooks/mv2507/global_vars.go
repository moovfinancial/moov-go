package mv2507

import (
	"github.com/moovfinancial/moov-go/pkg/moov"
)

var Underwriting moov.UnderwritingClient[moov.UpsertUnderwriting, moov.UnderwritingV2507] = moov.UnderwritingClient[moov.UpsertUnderwriting, moov.UnderwritingV2507]{Version: moov.Q3_2025}
