package moov

const (
	pathPing = "/ping"

	pathOAuth2Token  = "/oauth2/token" // #nosec G101
	pathOAuth2Revoke = "/oauth2/revoke"

	pathAccounts = "/accounts"
	pathAccount  = "/accounts/%s"

	pathApplications    = "/applications"
	pathApplicationKeys = "/applications/%s/keys"

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

	pathInstitutions = "/institutions"

	pathTransferOptions = "/accounts/%s/transfer-options"

	pathTransfers = "/accounts/%s/transfers"
	pathTransfer  = "/accounts/%s/transfers/%s"

	pathSchedules          = "/accounts/%s/schedules"
	pathSchedule           = "/accounts/%s/schedules/%s"
	pathScheduleOccurrence = "/accounts/%s/schedules/%s/occurrences/%s"

	pathCancellations = "/accounts/%s/transfers/%s/cancellations"
	pathCancellation  = "/accounts/%s/transfers/%s/cancellations/%s"

	pathReversals = "/accounts/%s/transfers/%s/reversals"

	pathRefunds = "/accounts/%s/transfers/%s/refunds"
	pathRefund  = "/accounts/%s/transfers/%s/refunds/%s"

	pathReceipts = "/receipts"
	pathReceipt  = "/receipts/%s"

	pathDisputes              = "/accounts/%s/disputes"
	pathDispute               = "/accounts/%s/disputes/%s"
	pathDisputeAccept         = "/accounts/%s/disputes/%s/accept"
	pathDisputeEvidenceText   = "/accounts/%s/disputes/%s/evidence-text"
	pathDisputeSubmitEvidence = "/accounts/%s/disputes/%s/evidence/submit"
	pathDisputeEvidences      = "/accounts/%s/disputes/%s/evidence"
	pathDisputeEvidence       = "/accounts/%s/disputes/%s/evidence/%s"
	pathDisputeEvidenceFile   = "/accounts/%s/disputes/%s/evidence-file"

	pathEndToEndPublicKey = "/end-to-end-keys"
	pathEndToEndTokenTest = "/debug/end-to-end-token" //nolint:gosec

	pathTerminalApplications        = "/terminal-applications"
	pathTerminalApplication         = "/terminal-applications/%s"
	pathTerminalApplicationVersions = "/terminal-applications/%s/versions"

	pathAccountTerminalApplication              = "/accounts/%s/terminal-applications/%s"
	pathAccountTerminalApplications             = "/accounts/%s/terminal-applications"
	pathAccountTerminalApplicationConfiguration = "/accounts/%s/terminal-applications/%s/configuration"

	pathTickets        = "/accounts/%s/tickets"
	pathTicket         = "/accounts/%s/tickets/%s"
	pathTicketMessages = "/accounts/%s/tickets/%s/messages"
)
