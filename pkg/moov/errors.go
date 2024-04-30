package moov

import (
	"errors"
	"strings"
)

func ErrorAsCallResponse(err error) CallResponse {
	if e := errorAsA[CallResponse](err); e != nil {
		return *e
	}
	return nil
}

func ErrorAsHttpCallResponse(err error) HttpCallResponse {
	if e := errorAsA[HttpCallResponse](err); e != nil {
		return *e
	}
	return nil
}

func errorAsA[A interface{}](err error) *A {
	t := new(A)
	if errors.As(err, t) {
		return t
	}
	return nil
}

var (
	ErrCredentialsNotSet            = errors.New("api credentials not set")
	ErrAccountNotFound              = errors.New("no account with the specified accountID was found")
	ErrAlreadyExists                = errors.New("resource already exists")
	ErrMicroDepositAmountsIncorrect = errors.New("the amounts provided are incorrect or the bank account is in an unexpected state")
	ErrXIdempotencyKey              = errors.New("attempted to create a transfer using a duplicate X-Idempotency-Key header")

	// ErrDuplicateBankAccount = errors.New("duplicate bank account or invalid routing number")
	// ErrNoMicroDeposit       = errors.New("no account with the specified accountID was found or micro-deposits have not been sent for the source")
	// ErrAccount              = errors.New("no account with the specified accountID was found")
	// ErrUpdateCardConflict   = errors.New("attempting to update an existing disabled card")

	// ErrCardDataInvalid      = errors.New("the supplied card data appeared invalid or was declined by the issuer")
	// ErrRequestBody              = errors.New("request body could not be parsed")
	// ErrAuthNetwork              = errors.New("network error")
	// ErrBadRequest               = errors.New("the request body could not be processed")
	// ErrInvalidBankAccount       = errors.New("the bank account is not a bank account or is already pending verification")
	// ErrDuplicatedApplePayDomain = errors.New("apple pay domains already registered for this account")
	// ErrDomainsNotVerified       = errors.New("domains not verified with Apple")
	// ErrDomainsNotRegistered     = errors.New("no apple pay domains registered for this account were found")
	// ErrLinkingApplePayToken     = errors.New("an error occurred when linking an apple pay token")
	// ErrRateLimit                = errors.New("request was refused due to rate limiting")
	// ErrURL                      = errors.New("invalid url")
	// ErrNoCardUpdateFilters = errors.New("no card update filters provided")
)

func DebugPrintResponse(err error, f func(format string, a ...any) (n int, err error)) {
	if e := ErrorAsCallResponse(err); e != nil {
		sb := strings.Builder{}
		ErrorAsCallResponse(err).Unmarshal(&sb)
		f("[%s]", sb.String())
	}
}
