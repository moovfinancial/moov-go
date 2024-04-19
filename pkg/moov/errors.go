package moov

import "errors"

func ErrorAsCallResponse(err error) *CallResponse {
	return errorAsA[CallResponse](err)
}

func ErrorAsHttpCallResponse(err error) *HttpCallResponse {
	return errorAsA[HttpCallResponse](err)
}

func errorAsA[A interface{}](err error) *A {
	t := new(A)
	if errors.As(err, t) {
		return t
	}
	return nil
}

var (
	ErrAuthCredentialsNotSet = errors.New("API Keys are not set or invalid credentials")
	// ErrRequestBody              = errors.New("request body could not be parsed")
	// ErrAuthNetwork              = errors.New("network error")
	ErrNoAccount = errors.New("no account with the specified accountID was found")
	// ErrBadRequest               = errors.New("the request body could not be processed")
	// ErrInvalidBankAccount       = errors.New("the bank account is not a bank account or is already pending verification")
	ErrAmountIncorrect   = errors.New("the amounts provided are incorrect or the bank account is in an unexpected state")
	ErrDuplicateLinkCard = errors.New("attempted to link card that already exists on the account")
	ErrCardDataInvalid   = errors.New("the supplied card data appeared invalid or was declined by the issuer")
	// ErrDuplicatedApplePayDomain = errors.New("apple pay domains already registered for this account")
	// ErrDomainsNotVerified       = errors.New("domains not verified with Apple")
	// ErrDomainsNotRegistered     = errors.New("no apple pay domains registered for this account were found")
	// ErrLinkingApplePayToken     = errors.New("an error occurred when linking an apple pay token")
	// ErrRateLimit                = errors.New("request was refused due to rate limiting")
	ErrXIdempotencyKey = errors.New("attempted to create a transfer using a duplicate X-Idempotency-Key header")
	// ErrURL                      = errors.New("invalid url")

	ErrDuplicateBankAccount = errors.New("duplciate bank account or invalid routing number")
	ErrNoMicroDeposit       = errors.New("no account with the specified accountID was found or micro-deposits have not been sent for the source")

	ErrAccount = errors.New("no account with the specified accountID was found")
	// ErrNoCardUpdateFilters = errors.New("no card update filters provided")
	ErrUpdateCardConflict = errors.New("attempting to update an existing disabled card")
)
