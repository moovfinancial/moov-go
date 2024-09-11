package moov

const (
	pathPing = "/ping"

	pathOAuth2Token  = "/oauth2/token" // #nosec G101
	pathOAuth2Revoke = "/oauth2/revoke"

	pathAccounts = "/accounts"
	pathAccount  = "/accounts/%s"

	pathCapabilities = "/accounts/%s/capabilities"
	pathCapability   = "/accounts/%s/capabilities/%s"

	pathUnderwriting = "/accounts/%s/underwriting"

	pathFiles = "/accounts/%s/files"
	pathFile  = "/accounts/%s/files/%s"

	pathPaymentMethods = "/accounts/%s/payment-methods"
	pathPaymentMethod  = "/accounts/%s/payment-methods/%s"

	pathRepresentatives = "/accounts/%s/representatives"
	pathRepresentative  = "/accounts/%s/representatives/%s"

	pathCards = "/accounts/%s/cards"
	pathCard  = "/accounts/%s/cards/%s"

	pathBankAccounts = "/accounts/%s/bank-accounts"
	pathBankAccount  = "/accounts/%s/bank-accounts/%s"

	pathBankAccountMicroDeposits = "/accounts/%s/bank-accounts/%s/micro-deposits"

	pathBankAccountInstantVerification = "/accounts/%s/bank-accounts/%s/verify"

	pathWallets = "/accounts/%s/wallets"
	pathWallet  = "/accounts/%s/wallets/%s"

	pathWalletTransactions = "/accounts/%s/wallets/%s/transactions"
	pathWalletTransaction  = "/accounts/%s/wallets/%s/transactions/%s"

	pathSweepConfigs = "/accounts/%s/sweep-configs"
	pathSweepConfig  = "/accounts/%s/sweep-configs/%s"

	pathSweeps = "/accounts/%s/wallets/%s/sweeps"
	pathSweep  = "/accounts/%s/wallets/%s/sweeps/%s"

	pathApplePay        = "/accounts/%s/apple-pay"
	pathApplePayDomains = "/accounts/%s/apple-pay/domains"

	pathApplePaySessions = "/accounts/%s/apple-pay/sessions"
	pathApplePayTokens   = "/accounts/%s/apple-pay/tokens" // #nosec G101

	pathInstitutions = "/institutions/%s/search"

	pathTransferOptions = "/transfer-options"

	pathTransfers = "/transfers"
	pathTransfer  = "/transfers/%s"

	pathSchedules          = "/accounts/%s/schedules"
	pathSchedule           = "/accounts/%s/schedules/%s"
	pathScheduleOccurrence = "/accounts/%s/schedules/%s/occurrences/%s"

	pathReversals = "/transfers/%s/reversals"

	pathRefunds = "/transfers/%s/refunds"
	pathRefund  = "/transfers/%s/refunds/%s"

	pathDisputes              = "/disputes"
	pathDispute               = "/disputes/%s"
	pathDisputeAccept         = "/disputes/%s/accept"
	pathDisputeEvidenceText   = "/disputes/%s/evidence-text"
	pathDisputeSubmitEvidence = "/disputes/%s/evidence/submit"
	pathDisputeEvidences      = "/disputes/%s/evidence"
	pathDisputeEvidence       = "/disputes/%s/evidence/%s"
	pathDisputeEvidenceFile   = "/disputes/%s/evidence-file"

	pathEndToEndPublicKey = "/end-to-end-keys"
	pathEndToEndTokenTest = "/debug/end-to-end-token"
)
