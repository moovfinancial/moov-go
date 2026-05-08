package testtools

import (
	"github.com/moovfinancial/moov-go/pkg/moov"
)

// FindMerchantWalletPaymentMethod returns the moov-wallet payment method whose
// wallet ID matches MERCHANT_WALLET_ID. Returns the zero value and false if not found.
func FindMerchantWalletPaymentMethod(pms []moov.PaymentMethod) (moov.PaymentMethod, bool) {
	for _, pm := range pms {
		if pm.Wallet != nil && pm.Wallet.WalletID == MERCHANT_WALLET_ID {
			return pm, true
		}
	}
	return moov.PaymentMethod{}, false
}
