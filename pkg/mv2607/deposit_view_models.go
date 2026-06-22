package mv2607

import "time"

// IngestResponse is returned when a deposit view document is successfully ingested.
type IngestResponse struct {
	MoovAccountID   string       `json:"moovAccountID"`
	SourceSystem    SourceSystem `json:"sourceSystem"`
	SourceAccountID string       `json:"sourceAccountID"`
	IngestedAt      time.Time    `json:"ingestedAt"`
}

// SourceSystem identifies the core banking system a deposit view document
// originated from. It is sent as the X-Source-System header.
type SourceSystem string

const (
	SourceSystemJHSilverlake   SourceSystem = "jh_silverlake"
	SourceSystemJHCIF2020      SourceSystem = "jh_cif2020"
	SourceSystemJHCoreDirector SourceSystem = "jh_coredirector"
)

// ParseSourceSystem maps a raw string to a known SourceSystem. The second return
// value is false if the string does not match a known source system.
func ParseSourceSystem(s string) (SourceSystem, bool) {
	switch s {
	case string(SourceSystemJHSilverlake):
		return SourceSystemJHSilverlake, true
	case string(SourceSystemJHCIF2020):
		return SourceSystemJHCIF2020, true
	case string(SourceSystemJHCoreDirector):
		return SourceSystemJHCoreDirector, true
	default:
		return "", false
	}
}
