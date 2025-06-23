package mv2507

import "github.com/moovfinancial/moov-go/pkg/moov"

var Capabilities = moov.CapabilityClient[RequestedCapabilities, moov.Capability]{Version: moov.Version2025_07}

type RequestedCapabilities struct {
	Capabilities []moov.CapabilityName `json:"capabilities"`
}
