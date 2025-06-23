package moov

import "time"

type requestCapabilities struct {
	Capabilities []CapabilityName `json:"capabilities"`
}

type CapabilityName string

const (
	CapabilityName_1099             CapabilityName = "1099"
	CapabilityName_CardIssuing      CapabilityName = "card-issuing"
	CapabilityName_CollectFunds     CapabilityName = "collect-funds"
	CapabilityName_DeveloperAccount CapabilityName = "developer-account"
	CapabilityName_ProductionApp    CapabilityName = "production-app"
	CapabilityName_SendFunds        CapabilityName = "send-funds"
	CapabilityName_Transfers        CapabilityName = "transfers"
	CapabilityName_Wallet           CapabilityName = "wallet"

	// Granular Capability Names Platform
	CapabilityName_PlatformProductionApp   CapabilityName = "platform.production-app"
	CapabilityName_PlatformWalletTransfers CapabilityName = "platform.wallet-transfers"

	// Granular Capability Names Wallet
	Capability_NameWalletBalance CapabilityName = "wallet.balance"

	// Granular Capability Names Collect Funds
	CapabilityName_CollectFundsACH          CapabilityName = "collect-funds.ach"
	CapabilityName_CollectFundsCardPayments CapabilityName = "collect-funds.card-payments"

	// Granular Capability Names Money Transfer
	CapabilityName_MoneyTransferPullFromCard CapabilityName = "money-transfer.pull-from-card"
	CapabilityName_MoneyTransferPushToCard   CapabilityName = "money-transfer.push-to-card"

	// Granular Capability Names Send Funds
	CapabilityName_SendFundsACH        CapabilityName = "send-funds.ach"
	CapabilityName_SendFundsRTP        CapabilityName = "send-funds.rtp"
	CapabilityName_SendFundsPushToCard CapabilityName = "send-funds.push-to-card"
)

// Capability Describes an action or set of actions that an account is permitted to perform.
type Capability struct {
	Capability CapabilityName `json:"capability"`
	// ID of account.
	AccountID    string           `json:"accountID,omitempty"`
	Status       CapabilityStatus `json:"status"`
	Requirements Requirement      `json:"requirements,omitempty"`
	// If status is `disabled`, the reason this capability was disabled.
	DisabledReason string     `json:"disabledReason,omitempty"`
	CreatedOn      time.Time  `json:"createdOn"`
	UpdatedOn      time.Time  `json:"updatedOn"`
	DisabledOn     *time.Time `json:"disabledOn,omitempty"`
}

// CapabilityStatus The status of the capability requested for an account.
type CapabilityStatus string

// List of CapabilityStatus
const (
	CapabilityStatus_Enabled  CapabilityStatus = "enabled"
	CapabilityStatus_Disabled CapabilityStatus = "disabled"
	CapabilityStatus_Pending  CapabilityStatus = "pending"
)

// Requirement Represents individual and business data necessary to facilitate the enabling of a capability for an account.
type Requirement struct {
	CurrentlyDue []RequirementId    `json:"currentlyDue,omitempty"`
	Errors       []RequirementError `json:"errors,omitempty"`
}

// RequirementId The unique ID of what the requirement is asking to be filled out.
type RequirementId string

// List of RequirementID
const (
	RequirementId_Account_TosAcceptance                    RequirementId = "account.tos-acceptance"
	RequirementId_Individual_Mobile                        RequirementId = "individual.mobile"
	RequirementId_Individual_Email                         RequirementId = "individual.email"
	RequirementId_Individual_EmailOrMobile                 RequirementId = "individual.email-or-mobile"
	RequirementId_Individual_Firstname                     RequirementId = "individual.firstname"
	RequirementId_Individual_Lastname                      RequirementId = "individual.lastname"
	RequirementId_Individual_Address                       RequirementId = "individual.address"
	RequirementId_Individual_SsnLast4                      RequirementId = "individual.ssn-last4"
	RequirementId_Individual_Ssn                           RequirementId = "individual.ssn"
	RequirementId_Individual_BirthDate                     RequirementId = "individual.birthdate"
	RequirementId_Business_LegalName                       RequirementId = "business.legalname"
	RequirementId_Business_DescriptionOrWebsite            RequirementId = "business.description-or-website"
	RequirementId_Business_EntityType                      RequirementId = "business.entity-type"
	RequirementId_Business_Dba                             RequirementId = "business.dba"
	RequirementId_Business_Ein                             RequirementId = "business.ein"
	RequirementId_Business_Address                         RequirementId = "business.address"
	RequirementId_Business_Phone                           RequirementId = "business.phone"
	RequirementId_Business_Admins                          RequirementId = "business.admins"
	RequirementId_Business_Controllers                     RequirementId = "business.controllers"
	RequirementId_Business_Owners                          RequirementId = "business.owners"
	RequirementId_Business_Classification                  RequirementId = "business.classification"
	RequirementId_Business_IndustryCodeMcc                 RequirementId = "business.industry-code-mcc"
	RequirementId_Business_IndicateOwnersProvided          RequirementId = "business.indicate-owners-provided"
	RequirementId_Business_AverageTransactionSize          RequirementId = "business.average-transaction-size"
	RequirementId_Business_MaxTransactionSize              RequirementId = "business.max-transaction-size"
	RequirementId_Business_AverageMonthlyTransactionVolume RequirementId = "business.average-monthly-transaction-volume"
	RequirementId_Business_Description                     RequirementId = "business.description"
	RequirementId_Business_UnderwritingDocumentsTierOne    RequirementId = "business.underwriting-documents-tier-one"
	RequirementId_BankAccounts_Name                        RequirementId = "bank-accounts.name"
	RequirementId_BankAccounts_RoutingNumber               RequirementId = "bank-accounts.routing-number"
	RequirementId_BankAccounts_AccountNumber               RequirementId = "bank-accounts.account-number"
	RequirementId_Representative_Mobile                    RequirementId = "representative.{rep-uuid}.mobile"
	RequirementId_Representative_Email                     RequirementId = "representative.{rep-uuid}.email"
	RequirementId_Representative_EmailOrMobile             RequirementId = "representative.{rep-uuid}.email-or-mobile"
	RequirementId_Representative_Firstname                 RequirementId = "representative.{rep-uuid}.firstname"
	RequirementId_Representative_Lastname                  RequirementId = "representative.{rep-uuid}.lastname"
	RequirementId_Representative_Address                   RequirementId = "representative.{rep-uuid}.address"
	RequirementId_Representative_SsnLast4                  RequirementId = "representative.{rep-uuid}.ssn-last4"
	RequirementId_Representative_Ssn                       RequirementId = "representative.{rep-uuid}.ssn"
	RequirementId_Representative_BirthDate                 RequirementId = "representative.{rep-uuid}.birthdate"
	RequirementId_Representative_JobTitle                  RequirementId = "representative.{rep-uuid}.job-title"
	RequirementId_Representative_IsController              RequirementId = "representative.{rep-uuid}.is-controller"
	RequirementId_Representative_IsOwner                   RequirementId = "representative.{rep-uuid}.is-owner"
	RequirementId_Representative_Ownership                 RequirementId = "representative.{rep-uuid}.ownership"
	RequirementId_Document                                 RequirementId = "document.{doc-uuid}"
)

// RequirementError Describes an error fulfilling a Requirement
type RequirementError struct {
	Requirement RequirementId        `json:"requirement,omitempty"`
	ErrorCode   RequirementErrorCode `json:"errorCode,omitempty"`
}

// RequirementErrorCode the model 'RequirementErrorCode'
type RequirementErrorCode string

// List of RequirementErrorCode
const (
	RequirementErrorCode_InvalidValue                RequirementErrorCode = "invalid-value"
	RequirementErrorCode_FailedAutomaticVerification RequirementErrorCode = "failed-automatic-verification"
	RequirementErrorCode_FailedOther                 RequirementErrorCode = "failed-other"
	RequirementErrorCode_InvalidAddress              RequirementErrorCode = "invalid-address"
	RequirementErrorCode_AddressRestricted           RequirementErrorCode = "address-restricted"
	RequirementErrorCode_TaxIdMismatch               RequirementErrorCode = "tax-id-mismatch"
	RequirementErrorCode_DocumentIdMismatch          RequirementErrorCode = "document-id-mismatch"
	RequirementErrorCode_DocumentDateOfBirthMismatch RequirementErrorCode = "document-date-of-birth-mismatch"
	RequirementErrorCode_DocumentNameMismatch        RequirementErrorCode = "document-name-mismatch"
	RequirementErrorCode_DocumentAddressMismatch     RequirementErrorCode = "document-address.mismatch"
	RequirementErrorCode_DocumentNumberMismatch      RequirementErrorCode = "document-number-mismatch"
	RequirementErrorCode_DocumentIncomplete          RequirementErrorCode = "document-incomplete"
	RequirementErrorCode_DocumentFailedRisk          RequirementErrorCode = "document-failed-risk"
	RequirementErrorCode_DocumentIllegible           RequirementErrorCode = "document-illegible"
	RequirementErrorCode_DocumentUnsupported         RequirementErrorCode = "document-unsupported"
	RequirementErrorCode_DocumentNotUploaded         RequirementErrorCode = "document-not-uploaded"
	RequirementErrorCode_DocumentCorrupt             RequirementErrorCode = "document-corrupt"
	RequirementErrorCode_DocumentExpired             RequirementErrorCode = "document-expired"
)
