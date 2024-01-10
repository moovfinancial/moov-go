package moov

import (
	"errors"
	"time"
)

var (
	ErrDuplicateBankAccount = errors.New("duplciate bank account or invalid routing number")
)

type BankAccount struct {
	BankAccountID         string `json:"bankAccountID,omitempty"`
	Fingerprint           string `json:"fingerprint,omitempty"`
	Status                string `json:"status,omitempty"`
	HolderName            string `json:"holderName,omitempty"`
	HolderType            string `json:"holderType,omitempty"`
	BankName              string `json:"bankName,omitempty"`
	BankAccountType       string `json:"bankAccountType,omitempty"`
	AccountNumber         string `json:"accountNumber,omitempty"`
	RoutingNumber         string `json:"routingNumber,omitempty"`
	LastFourAccountNumber string `json:"lastFourAccountNumber,omitempty"`
}

type AchDetails struct {
	Status                  string           `json:"status,omitempty"`
	TraceNumber             string           `json:"traceNumber,omitempty"`
	Return                  Return           `json:"return,omitempty"`
	Correction              Correction       `json:"correction,omitempty"`
	CompanyEntryDescription string           `json:"companyEntryDescription,omitempty"`
	OriginatingCompanyName  string           `json:"originatingCompanyName,omitempty"`
	StatusUpdates           ACHStatusUpdates `json:"statusUpdates,omitempty"`
	DebitHoldPeriod         string           `json:"debitHoldPeriod,omitempty"`
}

type Correction struct {
	Code        string `json:"code,omitempty"`
	Reason      string `json:"reason,omitempty"`
	Description string `json:"description,omitempty"`
}

type Return struct {
	Code        string `json:"code,omitempty"`
	Reason      string `json:"reason,omitempty"`
	Description string `json:"description,omitempty"`
}

type ACHStatusUpdates struct {
	Initiated  time.Time `json:"initiated,omitempty"`
	Originated time.Time `json:"originated,omitempty"`
	Corrected  time.Time `json:"corrected,omitempty"`
	Returned   time.Time `json:"returned,omitempty"`
	Completed  time.Time `json:"completed,omitempty"`
}

type BankAccountPayload struct {
	Account BankAccount `json:"account"`
}
