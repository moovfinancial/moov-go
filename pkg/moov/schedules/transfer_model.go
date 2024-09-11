package schedules

type Transfer struct {
	Description string `json:"description,omitempty"`
	Amount      Amount `json:"amount,omitempty"`

	PartnerID   string        `json:"partnerAccountID,omitempty"`
	Source      PaymentMethod `json:"source,omitempty"`
	Destination PaymentMethod `json:"destination,omitempty"`
}

type Amount struct {
	Value    int64  `json:"value,omitempty"`
	Currency string `json:"currency,omitempty"`
}

type PaymentMethod struct {
	PaymentMethodID string `json:"paymentMethodID,omitempty"`

	AchDetails  *AchDetails  `json:"achDetails,omitempty"`
	CardDetails *CardDetails `json:"cardDetails,omitempty"`
}

type AchDetails struct {
	CompanyEntryDescription *string `json:"companyEntryDescription,omitempty"`
	OriginatingCompanyName  *string `json:"originatingCompanyName,omitempty"`
}

type CardDetails struct {
	DynamicDescriptor *string `json:"dynamicDescriptor,omitempty"`
}
