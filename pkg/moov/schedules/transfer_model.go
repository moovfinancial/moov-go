package schedules

type Transfer struct {
	Description string `json:"description,omitempty" spanner:"description"`
	Amount      Amount `json:"amount,omitempty" spanner:"amount"`

	PartnerID   string        `json:"partnerAccountID,omitempty" spanner:"partner_account_id"`
	Source      PaymentMethod `json:"source,omitempty" spanner:"source"`
	Destination PaymentMethod `json:"destination,omitempty" spanner:"destination"`
}

type Amount struct {
	Value    int64  `json:"value,omitempty" spanner:"value"`
	Currency string `json:"currency,omitempty" spanner:"currency"`
}

type PaymentMethod struct {
	PaymentMethodID string `json:"paymentMethodID,omitempty" spanner:"payment_method_id"`

	AchDetails  *AchDetails  `json:"achDetails,omitempty" spanner:"ach_details"`
	CardDetails *CardDetails `json:"cardDetails,omitempty" spanner:"card_details"`
}

type AchDetails struct {
	CompanyEntryDescription *string `json:"companyEntryDescription,omitempty" spanner:"company_entry_description"`
	OriginatingCompanyName  *string `json:"originatingCompanyName,omitempty" spanner:"originating_company_name"`
}

type CardDetails struct {
	DynamicDescriptor *string `json:"dynamicDescriptor,omitempty" spanner:"dynamic_descriptor"`
}
