package moov

// TransferConfig configurable options for a transfer.
type TransferConfig struct {
	TipPresets       *TipPresets       `json:"tipPresets,omitempty"`
	TransferControls []TransferControl `json:"transferControls,omitempty"`
}

// TransferControl holds per-account transfer capability flags.
type TransferControl struct {
	AccountID              string `json:"accountID,omitempty"`
	PartnerAccountID       string `json:"partnerAccountID,omitempty"`
	AllowDynamicDescriptor bool   `json:"allowDynamicDescriptor,omitempty"`
	AllowSurcharge         bool   `json:"allowSurcharge,omitempty"`
	DebtRepayment          bool   `json:"debtRepayment,omitempty"`
}

// TipPresets suggested customer tip values for a transfer.
type TipPresets struct {
	CalculationBasis   *TipCalculationBasis `json:"calculationBasis,omitempty"`
	PercentageOptions  []int                `json:"percentageOptions,omitempty"`
	FixedAmountOptions []AmountDecimal      `json:"fixedAmountOptions,omitempty"`
}

// UpsertTransferConfig request payload for creating or replacing a transfer config.
type UpsertTransferConfig struct {
	TipPresets *UpsertTipPresets `json:"tipPresets,omitempty"`
}

// UpsertTipPresets suggested customer tip values for transfer config creation or replacement.
type UpsertTipPresets struct {
	CalculationBasis   *TipCalculationBasis `json:"calculationBasis,omitempty"`
	PercentageOptions  []int                `json:"percentageOptions,omitempty"`
	FixedAmountOptions []AmountDecimal      `json:"fixedAmountOptions,omitempty"`
}

// TipCalculationBasis which subtotal should be used for percentage tip calculations.
type TipCalculationBasis string

// List of TipCalculationBasis
const (
	TipCalculationBasis_PreTax  TipCalculationBasis = "pre-tax"
	TipCalculationBasis_PostTax TipCalculationBasis = "post-tax"
)

// TransferConfigValidationError validation details for transfer config requests.
type TransferConfigValidationError struct {
	TipPresetsCalculationBasis   *string                                 `json:"TipPresets.CalculationBasis,omitempty"`
	TipPresetsPercentageOptions  map[string]string                       `json:"TipPresets.PercentageOptions,omitempty"`
	TipPresetsFixedAmountOptions map[string]AmountDecimalValidationError `json:"TipPresets.FixedAmountOptions,omitempty"`
}

// AmountDecimalValidationError validation details for AmountDecimal fields.
type AmountDecimalValidationError struct {
	Currency     *string `json:"currency,omitempty"`
	ValueDecimal *string `json:"valueDecimal,omitempty"`
}
