package moov

type ApplePay struct {
	Brand           string     `json:"brand,omitempty"`
	CardType        string     `json:"cardType,omitempty"`
	CardDisplayName string     `json:"cardDisplayName,omitempty"`
	Fingerprint     string     `json:"fingerprint,omitempty"`
	Expiration      Expiration `json:"expiration,omitempty"`
	DynamicLastFour string     `json:"dynamicLastFour,omitempty"`
}
