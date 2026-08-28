package moov_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

// allCardFailureCodes lists every CardFailureCode constant. Keep it in sync
// with the const block in transfer_models.go so the invariant test below
// guards the whole set against duplicate or empty string values.
var allCardFailureCodes = []moov.CardFailureCode{
	moov.CardFailureCode_CallIssuer,
	moov.CardFailureCode_DoNotHonor,
	moov.CardFailureCode_ProcessingError,
	moov.CardFailureCode_InvalidTransaction,
	moov.CardFailureCode_InvalidAmount,
	moov.CardFailureCode_NoSuchIssuer,
	moov.CardFailureCode_ReenterTransaction,
	moov.CardFailureCode_CVV_Mismatch,
	moov.CardFailureCode_LostOrStolen,
	moov.CardFailureCode_Insufficient_Funds,
	moov.CardFailureCode_InvalidCardNumber,
	moov.CardFailureCode_InvalidMerchant,
	moov.CardFailureCode_ExpiredCard,
	moov.CardFailureCode_IncorrectPin,
	moov.CardFailureCode_TransactionNotAllowed,
	moov.CardFailureCode_SuspectedFraud,
	moov.CardFailureCode_AmountLimitedExceeded,
	moov.CardFailureCode_VelocityLimitExceeded,
	moov.CardFailureCode_RevocationOfAauthorization,
	moov.CardFailureCode_CardNotActivated,
	moov.CardFailureCode_IssuerNotAvailable,
	moov.CardFailureCode_CouldNotRoute,
	moov.CardFailureCode_CardholderAccounterClosed,
	moov.CardFailureCode_DuplicateTransaction,
	moov.CardFailureCode_AccountClosed,
	moov.CardFailureCode_AccountNotActivated,
	moov.CardFailureCode_AuthenticationFailed,
	moov.CardFailureCode_AuthenticationRequired,
	moov.CardFailureCode_CardholderActionRequired,
	moov.CardFailureCode_FormatError,
	moov.CardFailureCode_InvalidPin,
	moov.CardFailureCode_OfflineApproved,
	moov.CardFailureCode_OfflineDeclined,
	moov.CardFailureCode_PartialApproval,
	moov.CardFailureCode_PaymentStopped,
	moov.CardFailureCode_PinRequired,
	moov.CardFailureCode_RecordNotFound,
	moov.CardFailureCode_SurchargeNotPermitted,
	moov.CardFailureCode_TransactionReversed,
	moov.CardFailureCode_VerificationFailed,
	moov.CardFailureCode_UnknownIssue,
}

// TestCardFailureCode_valuesAreUniqueAndNonEmpty guards the hand-maintained
// enum against the realistic defect of a copy-pasted duplicate or empty string
// value when new codes are added.
func TestCardFailureCode_valuesAreUniqueAndNonEmpty(t *testing.T) {
	seen := make(map[moov.CardFailureCode]bool, len(allCardFailureCodes))

	for _, code := range allCardFailureCodes {
		require.NotEmpty(t, string(code), "CardFailureCode value must not be empty")
		assert.Falsef(t, seen[code], "duplicate CardFailureCode value %q", code)
		seen[code] = true
	}
}

// TestCardFailureCode_newValuesPresent asserts the 16 values added for CAR-5168
// carry their expected kebab-case string, catching an accidental value edit.
func TestCardFailureCode_newValuesPresent(t *testing.T) {
	added := map[moov.CardFailureCode]string{
		moov.CardFailureCode_AccountClosed:            "account-closed",
		moov.CardFailureCode_AccountNotActivated:      "account-not-activated",
		moov.CardFailureCode_AuthenticationFailed:     "authentication-failed",
		moov.CardFailureCode_AuthenticationRequired:   "authentication-required",
		moov.CardFailureCode_CardholderActionRequired: "cardholder-action-required",
		moov.CardFailureCode_FormatError:              "format-error",
		moov.CardFailureCode_InvalidPin:               "invalid-pin",
		moov.CardFailureCode_OfflineApproved:          "offline-approved",
		moov.CardFailureCode_OfflineDeclined:          "offline-declined",
		moov.CardFailureCode_PartialApproval:          "partial-approval",
		moov.CardFailureCode_PaymentStopped:           "payment-stopped",
		moov.CardFailureCode_PinRequired:              "pin-required",
		moov.CardFailureCode_RecordNotFound:           "record-not-found",
		moov.CardFailureCode_SurchargeNotPermitted:    "surcharge-not-permitted",
		moov.CardFailureCode_TransactionReversed:      "transaction-reversed",
		moov.CardFailureCode_VerificationFailed:       "verification-failed",
	}

	require.Len(t, added, 16)
	for code, want := range added {
		assert.Equal(t, want, string(code))
	}
}
