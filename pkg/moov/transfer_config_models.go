package moov

// TransferConfig configurable options for a transfer.
type TransferConfig struct {
	TipPresets *TipPresets `json:"tipPresets,omitempty"`
}

// TipPresets suggested customer tip values for a transfer.
type TipPresets struct {
	CalculationBasis   *TipCalculationBasis `json:"calculationBasis,omitempty"`
	PercentageOptions  []int                `json:"percentageOptions,omitempty"`
	FixedAmountOptions []AmountDecimal      `json:"fixedAmountOptions,omitempty"`
}

// CreateTransferConfig request payload for creating a transfer config.
type CreateTransferConfig struct {
	TipPresets *CreateTipPresets `json:"tipPresets,omitempty"`
}

// CreateTipPresets suggested customer tip values for transfer config creation.
type CreateTipPresets struct {
	CalculationBasis   *TipCalculationBasis `json:"calculationBasis,omitempty"`
	PercentageOptions  []int                `json:"percentageOptions,omitempty"`
	FixedAmountOptions []AmountDecimal      `json:"fixedAmountOptions,omitempty"`
}

// PutTransferConfig request payload for replacing a transfer config.
type PutTransferConfig struct {
	TipPresets PutTipPresets `json:"tipPresets"`
}

// PutTipPresets suggested customer tip values for transfer config replacement.
type PutTipPresets struct {
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
