package moov

// Mode The mode this account is allowed to be used within.
type Mode string

// List of Mode
const (
	MODE_SANDBOX    Mode = "sandbox"
	MODE_PRODUCTION Mode = "production"
)
