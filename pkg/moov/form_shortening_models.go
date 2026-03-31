package moov

// Industry represents a single industry entry from the form-shortening industries list.
type Industry struct {
	// Unique identifier for the industry (e.g. "clothing-accessories")
	Industry string `json:"industry,omitempty"`
	// Human-readable name for the industry (e.g. "Clothing & Accessories")
	DisplayName string `json:"displayName,omitempty"`
	// Category identifier the industry belongs to (e.g. "retail")
	Category string `json:"category,omitempty"`
	// Human-readable name for the category (e.g. "Retail")
	CategoryDisplayName string `json:"categoryDisplayName,omitempty"`
	// Default merchant category code for the industry (e.g. "5651")
	DefaultMcc string `json:"defaultMcc,omitempty"`
}

type Industries struct {
	Industries []Industry `json:"industries"`
}

// EnrichedIndustryCodes holds industry classification codes for an enriched business profile.
type EnrichedIndustryCodes struct {
	// North American Industry Classification System code
	Naics string `json:"naics,omitempty"`
	// Standard Industrial Classification code
	Sic string `json:"sic,omitempty"`
}

// EnrichedBusiness holds publicly available business information returned by the profile enrichment endpoint.
type EnrichedBusiness struct {
	Address           *Address               `json:"address,omitempty"`
	Email             string                 `json:"email,omitempty"`
	IndustryCodes     *EnrichedIndustryCodes `json:"industryCodes,omitempty"`
	LegalBusinessName string                 `json:"legalBusinessName,omitempty"`
	Phone             *Phone                 `json:"phone,omitempty"`
	Website           string                 `json:"website,omitempty"`
}

// EnrichedBusinessProfile is the response from the profile enrichment endpoint.
type EnrichedBusinessProfile struct {
	Business *EnrichedBusiness `json:"business,omitempty"`
}
