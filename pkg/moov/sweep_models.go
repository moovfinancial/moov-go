package moov

import "time"

type SweepConfig struct {
	SweepConfigID string            `json:"sweepConfigID"`
	WalletID      string            `json:"walletID"`
	Status        SweepConfigStatus `json:"status"`
	// The payment method used to push funds to a bank account. `ach-credit-standard` or `ach-credit-same-day` is required.
	PushPaymentMethod SweepConfigPaymentMethod `json:"pushPaymentMethod"`
	// The payment method used to pull funds from a bank account. An `ach-debit-fund` payment method is required.
	PullPaymentMethod SweepConfigPaymentMethod `json:"pullPaymentMethod"`
	// An optional field to specify an amount to maintain in the wallet. This is a decimal-formatted numerical string that represents up to 2 decimal place precision. In USD for example, 12.34 is $12.34 and 0.99 is $0.99. If not supplied, the default is 0.00.
	MinimumBalance *string `json:"minimumBalance"`
	// An optional override of the default NACHA company entry description for sweep transfers.
	StatementDescriptor *string `json:"statementDescriptor,omitempty"`
	// An array of fields that are locked. To request updates, please contact Moov support.
	LockedFields []string   `json:"lockedFields,omitempty"`
	CreatedOn    time.Time  `json:"createdOn"`
	UpdatedOn    time.Time  `json:"updatedOn"`
	DisabledOn   *time.Time `json:"disabledOn,omitempty"`
}

type SweepConfigStatus string

const (
	SweepConfigStatus_Enabled  SweepConfigStatus = "enabled"
	SweepConfigStatus_Disabled SweepConfigStatus = "disabled"
)

type SweepConfigPaymentMethod struct {
	PaymentMethodID string     `json:"paymentMethodID"`
	DisabledOn      *time.Time `json:"disabledOn,omitempty"`
}

type Sweep struct {
	SweepID             string      `json:"sweepID,omitempty"`
	Status              SweepStatus `json:"status,omitempty"`
	AccrualStartedOn    time.Time   `json:"accrualStartedOn,omitempty"`
	AccrualEndedOn      *time.Time  `json:"accrualEndedOn,omitempty"`
	PushPaymentMethodID string      `json:"pushPaymentMethodID,omitempty"`
	PullPaymentMethodID string      `json:"pullPaymentMethodID,omitempty"`
	// The total net Amount of wallet transactions accrued in the sweep.
	// A decimal-formatted numerical string that represents up to 9 decimal place precision.
	AccruedAmount string `json:"accruedAmount,omitempty"`
	// A 3-letter ISO 4217 currency code.
	Currency string `json:"currency,omitempty"`
	// Only set after the sweep transfer has been created.
	TransferID string `json:"transferID,omitempty"`
	// Only set after the sweep transfer has been created.
	// A decimal-formatted numerical string that represents up to 9 decimal place precision.
	TransferAmount string `json:"transferAmount,omitempty"`
	// Amount remaining in the wallet after the sweep transfer completes.
	// A decimal-formatted numerical string that represents up to 9 decimal place precision.
	ResidualBalance string `json:"residualBalance,omitempty"`
}

type SweepStatus string

const (
	SweepStatus_Accruing       SweepStatus = "accruing"
	SweepStatus_ActionRequired SweepStatus = "action-required"
	SweepStatus_Canceled       SweepStatus = "canceled"
	SweepStatus_Closed         SweepStatus = "closed"
	SweepStatus_Failed         SweepStatus = "failed"
	SweepStatus_Paid           SweepStatus = "paid"
)
