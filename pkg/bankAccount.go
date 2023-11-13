package moov

type BankAccount struct {
	BankAccountID         string `json:"bankAccountID,omitempty"`
	Fingerprint           string `json:"fingerprint,omitempty"`
	Status                string `json:"status,omitempty"`
	HolderName            string `json:"holderName,omitempty"`
	HolderType            string `json:"holderType,omitempty"`
	BankName              string `json:"bankName,omitempty"`
	BankAccountType       string `json:"bankAccountType,omitempty"`
	RoutingNumber         string `json:"routingNumber,omitempty"`
	LastFourAccountNumber string `json:"lastFourAccountNumber,omitempty"`
}

// CreateBankAccount creates a new bank account for the given customer account

// GetBankAccount retrieves a bank account for the given customer account

// DelteBankAccount deletes a bank account for the given customer account

// ListBankAccounts lists all bank accounts for the given customer account

// MicroDepositInitiate creates a new micro deposit verification for the given bank account

// MicroDepositConfirm confirms a micro deposit verification for the given bank account
