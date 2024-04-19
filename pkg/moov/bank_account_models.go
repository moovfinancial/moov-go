package moov

import (
	"time"
)

type createBankAccount struct {
	Account   *BankAccountRequest `json:"account,omitempty"`
	Plaid     *PlaidRequest       `json:"plaid,omitempty"`
	PlaidLink *PlaidLinkRequest   `json:"plaidLink,omitempty"`
	MX        *MXRequest          `json:"mx,omitempty"`
}

type BankAccountRequest struct {
	RoutingNumber string          `json:"routingNumber,omitempty"`
	AccountNumber string          `json:"accountNumber,omitempty"`
	AccountType   BankAccountType `json:"bankAccountType,omitempty"`
	HolderName    string          `json:"holderName,omitempty"`
	HolderType    HolderType      `json:"holderType,omitempty"`
}

// BankAccountType Defines the type of the bank account
type BankAccountType string

// List of HolderType
const (
	BankAccountType_Checking BankAccountType = "checking"
	BankAccountType_Savings  BankAccountType = "savings"
)

// HolderType Defines the type of the account holder
type HolderType string

// List of HolderType
const (
	HolderType_Individual HolderType = "individual"
	HolderType_Business   HolderType = "business"
)

type PlaidRequest struct {
	Token string `json:"token"`
}

type PlaidLinkRequest struct {
	PublicToken string `json:"publicToken"`
}

// MxAuthorizationCode The authorization code of a MX account which allows a processor to retrieve a linked payment account. <br><br> `sandbox` - When linking a bank account to a `sandbox` account using a MX authorization code it will utilize MX's sandbox environment. The MX authorization code provided must be generated from MX's sandbox environment.
type MXRequest struct {
	AuthorizationCode string `json:"authorizationCode,omitempty"`
}

// BankAccountResponse Describes a bank account on a Moov account.
type BankAccount struct {
	// UUID v4
	BankAccountID string `json:"bankAccountID,omitempty"`
	// Once the bank account is linked, we don't reveal the full bank account number. The fingerprint acts as a way to identify whether two linked bank accounts are the same.
	Fingerprint           string                  `json:"fingerprint,omitempty"`
	Status                BankAccountStatus       `json:"status,omitempty"`
	HolderName            string                  `json:"holderName,omitempty"`
	HolderType            HolderType              `json:"holderType,omitempty"`
	BankName              string                  `json:"bankName,omitempty"`
	BankAccountType       BankAccountType         `json:"bankAccountType,omitempty"`
	RoutingNumber         string                  `json:"routingNumber,omitempty"`
	LastFourAccountNumber string                  `json:"lastFourAccountNumber,omitempty"`
	UpdatedOn             time.Time               `json:"updatedOn,omitempty"`
	StatusReason          BankAccountStatusReason `json:"statusReason,omitempty"`
	ExceptionDetails      *ExceptionDetails       `json:"exceptionDetails,omitempty"`
}

// BankAccountStatus The bank account status.
type BankAccountStatus string

// List of BankAccountStatus
const (
	BANKACCOUNTSTATUS_NEW                 BankAccountStatus = "new"
	BANKACCOUNTSTATUS_VERIFIED            BankAccountStatus = "verified"
	BANKACCOUNTSTATUS_VERIFICATION_FAILED BankAccountStatus = "verificationFailed"
	BANKACCOUNTSTATUS_PENDING             BankAccountStatus = "pending"
	BANKACCOUNTSTATUS_ERRORED             BankAccountStatus = "errored"
)

// BankAccountStatusReason The reason the bank account status changed to the current value.
type BankAccountStatusReason string

// List of BankAccountStatusReason
const (
	BANKACCOUNTSTATUSREASON_BANK_ACCOUNT_CREATED            BankAccountStatusReason = "bank-account-created"
	BANKACCOUNTSTATUSREASON_VERIFICATION_INITIATED          BankAccountStatusReason = "verification-initiated"
	BANKACCOUNTSTATUSREASON_MICRO_DEPOSIT_ATTEMPTS_EXCEEDED BankAccountStatusReason = "micro-deposit-attempts-exceeded"
	BANKACCOUNTSTATUSREASON_MICRO_DEPOSIT_EXPIRED           BankAccountStatusReason = "micro-deposit-expired"
	BANKACCOUNTSTATUSREASON_MAX_VERIFICATION_FAILURES       BankAccountStatusReason = "max-verification-failures"
	BANKACCOUNTSTATUSREASON_VERIFICATION_SUCCESSFUL         BankAccountStatusReason = "verification-successful"
	BANKACCOUNTSTATUSREASON_ACH_DEBIT_RETURN                BankAccountStatusReason = "ach-debit-return"
	BANKACCOUNTSTATUSREASON_ACH_CREDIT_RETURN               BankAccountStatusReason = "ach-credit-return" // #nosec G101
	BANKACCOUNTSTATUSREASON_MICRO_DEPOSIT_RETURN            BankAccountStatusReason = "micro-deposit-return"
	BANKACCOUNTSTATUSREASON_ADMIN_ACTION                    BankAccountStatusReason = "admin-action"
	BANKACCOUNTSTATUSREASON_OTHER                           BankAccountStatusReason = "other"
)

// ExceptionDetails Reason for, and details related to, an `errored` or `verificationFailed` bank account status.
type ExceptionDetails struct {
	AchReturnCode AchReturnCode `json:"achReturnCode,omitempty"`
	// Details related to an `errored` or `verificationFailed` bank account status.
	Description string `json:"description,omitempty"`
}

// AchReturnCode The return code of an ACH transaction that caused the bank account status to change.
type AchReturnCode string

// List of ACHReturnCode
const (
	ACHRETURNCODE_R02 AchReturnCode = "R02"
	ACHRETURNCODE_R03 AchReturnCode = "R03"
	ACHRETURNCODE_R04 AchReturnCode = "R04"
	ACHRETURNCODE_R05 AchReturnCode = "R05"
	ACHRETURNCODE_R07 AchReturnCode = "R07"
	ACHRETURNCODE_R08 AchReturnCode = "R08"
	ACHRETURNCODE_R10 AchReturnCode = "R10"
	ACHRETURNCODE_R11 AchReturnCode = "R11"
	ACHRETURNCODE_R12 AchReturnCode = "R12"
	ACHRETURNCODE_R13 AchReturnCode = "R13"
	ACHRETURNCODE_R14 AchReturnCode = "R14"
	ACHRETURNCODE_R15 AchReturnCode = "R15"
	ACHRETURNCODE_R16 AchReturnCode = "R16"
	ACHRETURNCODE_R17 AchReturnCode = "R17"
	ACHRETURNCODE_R20 AchReturnCode = "R20"
	ACHRETURNCODE_R23 AchReturnCode = "R23"
	ACHRETURNCODE_R29 AchReturnCode = "R29"
	ACHRETURNCODE_R34 AchReturnCode = "R34"
	ACHRETURNCODE_R38 AchReturnCode = "R38"
	ACHRETURNCODE_R39 AchReturnCode = "R39"
)
