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

	// Includes any payment methods generated for a newly created bank account, removing the need to  call the List Payment Methods endpoint following a successful Create BankAccount request.
	// **NOTE: This field is only populated for Create BankAccount requests made with the `X-Wait-For` header.**
	PaymentMethods []BasicPaymentMethod `json:"paymentMethods,omitempty"`
}

// BankAccountStatus The bank account status.
type BankAccountStatus string

// List of BankAccountStatus
const (
	BankAccountStatus_New                BankAccountStatus = "new"
	BankAccountStatus_Verified           BankAccountStatus = "verified"
	BankAccountStatus_VerificationFailed BankAccountStatus = "verificationFailed"
	BankAccountStatus_Pending            BankAccountStatus = "pending"
	BankAccountStatus_Errored            BankAccountStatus = "errored"
)

// BankAccountStatusReason The reason the bank account status changed to the current value.
type BankAccountStatusReason string

// List of BankAccountStatusReason
const (
	BankAccountStatusReason_BankAccountCreated           BankAccountStatusReason = "bank-account-created"
	BankAccountStatusReason_VerificationInitiated        BankAccountStatusReason = "verification-initiated"
	BankAccountStatusReason_MicroDepositAttemptsExceeded BankAccountStatusReason = "micro-deposit-attempts-exceeded"
	BankAccountStatusReason_MicroDepositExpired          BankAccountStatusReason = "micro-deposit-expired"
	BankAccountStatusReason_MaxVerificationFailures      BankAccountStatusReason = "max-verification-failures"
	BankAccountStatusReason_VerificationSuccessful       BankAccountStatusReason = "verification-successful"
	BankAccountStatusReason_AchDebitReturn               BankAccountStatusReason = "ach-debit-return"
	BankAccountStatusReason_AchCreditReturn              BankAccountStatusReason = "ach-credit-return" // #nosec G101
	BankAccountStatusReason_MicroDepositReturn           BankAccountStatusReason = "micro-deposit-return"
	BankAccountStatusReason_AdminAction                  BankAccountStatusReason = "admin-action"
	BankAccountStatusReason_Other                        BankAccountStatusReason = "other"
)

// ExceptionDetails Reason for, and details related to, an `errored` or `verificationFailed` bank account status.
type ExceptionDetails struct {
	// AchReturnCode is the return code of an ACH transaction that caused the bank account status to change.
	AchReturnCode *AchReturnCode `json:"achReturnCode,omitempty"`

	// Details related to an `errored` or `verificationFailed` bank account status.
	Description string `json:"description,omitempty"`

	// RTPRejectionCode is a rejection code of an RTP transaction that caused the bank account status to change.
	RTPRejectionCode *RTPRejectionCode `json:"rtpRejectionCode"`
}

// AchReturnCode is the return code of an ACH transaction that caused the bank account status to change.
type AchReturnCode string

// List of ACHReturnCode
const (
	AchReturnCode_R02 AchReturnCode = "R02"
	AchReturnCode_R03 AchReturnCode = "R03"
	AchReturnCode_R04 AchReturnCode = "R04"
	AchReturnCode_R05 AchReturnCode = "R05"
	AchReturnCode_R07 AchReturnCode = "R07"
	AchReturnCode_R08 AchReturnCode = "R08"
	AchReturnCode_R10 AchReturnCode = "R10"
	AchReturnCode_R11 AchReturnCode = "R11"
	AchReturnCode_R12 AchReturnCode = "R12"
	AchReturnCode_R13 AchReturnCode = "R13"
	AchReturnCode_R14 AchReturnCode = "R14"
	AchReturnCode_R15 AchReturnCode = "R15"
	AchReturnCode_R16 AchReturnCode = "R16"
	AchReturnCode_R17 AchReturnCode = "R17"
	AchReturnCode_R20 AchReturnCode = "R20"
	AchReturnCode_R23 AchReturnCode = "R23"
	AchReturnCode_R29 AchReturnCode = "R29"
	AchReturnCode_R34 AchReturnCode = "R34"
	AchReturnCode_R38 AchReturnCode = "R38"
	AchReturnCode_R39 AchReturnCode = "R39"
)

// RTPRejectionCode is a rejection code of an RTP transaction that caused the bank account status to change.
type RTPRejectionCode string

// List of RTPRejectionCode
const (
	RTPRejectionCode_AC03 RTPRejectionCode = "AC03" // Account Invalid
	RTPRejectionCode_AC04 RTPRejectionCode = "AC04" // Account Closed
	RTPRejectionCode_AC06 RTPRejectionCode = "AC06" // Account Blocked
	RTPRejectionCode_AC14 RTPRejectionCode = "AC14" // Creditor Account Type Invalid
	RTPRejectionCode_AG01 RTPRejectionCode = "AG01" // Transactions Forbidden On Account
	RTPRejectionCode_AG03 RTPRejectionCode = "AG03" // Transaction Type Not Supported
	RTPRejectionCode_MD07 RTPRejectionCode = "MD07" // Customer Deceased
)
