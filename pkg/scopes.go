package moov

import (
	"fmt"
	"strings"
)

// Setup namespacing so when you use these scopes its `moov.Scopes.AccountsRead()`
var Scopes scopeList = scopeList{}

// Used for listing the accounts the facilitator has been connected to.
func (sl *scopeList) AccountsRead() ScopeBuilder {
	return appendScope("/accounts.read")
}

// Used for creating a new account thats auto-connected to the facilitator.
func (sl *scopeList) AccountsWrite() ScopeBuilder {
	return appendScope("/accounts.write")
}

func (sl *scopeList) BankAccountsRead(accountID string) ScopeBuilder {
	return appendScope("/accounts/%s/bank-accounts.read", accountID)
}

func (sl *scopeList) BankAccountsWrite(accountID string) ScopeBuilder {
	return appendScope("/accounts/%s/bank-accounts.write", accountID)
}

func (sl *scopeList) CapabilitiesRead(accountID string) ScopeBuilder {
	return appendScope("/accounts/%s/capabilities.read", accountID)
}

func (sl *scopeList) CapabilitiesWrite(accountID string) ScopeBuilder {
	return appendScope("/accounts/%s/capabilities.write", accountID)
}

func (sl *scopeList) CardsRead(accountID string) ScopeBuilder {
	return appendScope("/accounts/%s/cards.read", accountID)
}

func (sl *scopeList) CardsWrite(accountID string) ScopeBuilder {
	return appendScope("/accounts/%s/cards.write", accountID)
}

func (sl *scopeList) IssuedCardsRead(accountID string) ScopeBuilder {
	return appendScope("/accounts/%s/issued-cards.read", accountID)
}

// Allows for the creation of a new issued card
func (sl *scopeList) IssuedCardsWrite(accountID string) ScopeBuilder {
	return appendScope("/accounts/%s/issued-cards.write", accountID)
}

// WARNING: Will return PCI data only allowed use case is by the customer directly to moov otherwise your business will require a PCI audit.
func (sl *scopeList) IssuedCardsReadSecure(accountID string) ScopeBuilder {
	return appendScope("/accounts/%s/issued-cards.read-secure", accountID)
}

// Give to the registered ApplePay merchants to allow them manage their ApplePay integration.
func (sl *scopeList) ApplePayMerchantRead(accountID string) ScopeBuilder {
	return appendScope("/accounts/%s/apple-pay-merchant.read", accountID)
}

// Give to the registered ApplePay merchants to allow them manage their ApplePay integration.
func (sl *scopeList) ApplePayMerchantWrite(accountID string) ScopeBuilder {
	return appendScope("/accounts/%s/apple-pay-merchant.write", accountID)
}

func (sl *scopeList) ApplePayWrite(accountID string) ScopeBuilder {
	return appendScope("/accounts/%s/apple-pay.write", accountID)
}

func (sl *scopeList) AccountProfileRead(accountID string) ScopeBuilder {
	return appendScope("/accounts/%s/profile.read", accountID)
}

func (sl *scopeList) AccountProfileWrite(accountID string) ScopeBuilder {
	return appendScope("/accounts/%s/profile.write", accountID)
}

func (sl *scopeList) AccountRepresentativesRead(accountID string) ScopeBuilder {
	return appendScope("/accounts/%s/representatives.read", accountID)
}

func (sl *scopeList) AccountRepresentativesWrite(accountID string) ScopeBuilder {
	return appendScope("/accounts/%s/representatives.write", accountID)
}

func (sl *scopeList) FilesRead(accountID string) ScopeBuilder {
	return appendScope("/accounts/%s/files.read", accountID)
}

func (sl *scopeList) FilesWrite(accountID string) ScopeBuilder {
	return appendScope("/accounts/%s/files.write", accountID)
}

func (sl *scopeList) PaymentMethodsRead(accountID string) ScopeBuilder {
	return appendScope("/accounts/%s/payment-methods.read", accountID)
}

func (sl *scopeList) WalletsRead(accountID string) ScopeBuilder {
	return appendScope("/accounts/%s/wallets.read", accountID)
}

func (sl *scopeList) Ping() ScopeBuilder {
	return appendScope("/ping.read")
}

func (sl *scopeList) Fed() ScopeBuilder {
	return appendScope("/fed.read")
}

func (sl *scopeList) ProfileEnrichment() ScopeBuilder {
	return appendScope("/profile-enrichment.read")
}

// Boilerplate for setting the above.

type scopeList struct{}

type ScopeBuilder func(sb *scopeBuilder) error

type scopeBuilder struct {
	scopes []string
}

func buildScopes(scopes ...ScopeBuilder) (string, error) {
	sb := &scopeBuilder{
		scopes: []string{},
	}

	for _, scp := range scopes {
		if err := scp(sb); err != nil {
			return "", err
		}
	}

	return strings.Join(sb.scopes, " "), nil
}

func appendScope(scope string, args ...any) ScopeBuilder {
	return func(sb *scopeBuilder) error {
		sb.scopes = append(sb.scopes, fmt.Sprintf(scope, args...))
		return nil
	}
}
