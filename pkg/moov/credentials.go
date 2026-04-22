package moov

import (
	"os"
)

const ENV_MOOV_HOST = "MOOV_HOST"
const ENV_MOOV_PUBLIC_KEY = "MOOV_PUBLIC_KEY"
const ENV_MOOV_SECRET_KEY = "MOOV_SECRET_KEY" // #nosec G101

func CredentialsDefault() Credentials {
	return Credentials{}
}

func CredentialsFromEnv() Credentials {
	creds := CredentialsDefault()
	creds.PublicKey = os.Getenv(ENV_MOOV_PUBLIC_KEY)
	creds.SecretKey = os.Getenv(ENV_MOOV_SECRET_KEY)

	creds.Host = os.Getenv(ENV_MOOV_HOST)
	if creds.Host == "" {
		creds.Host = "api.moov.io"
	}

	return creds
}

type Credentials struct {
	PublicKey string `yaml:"public_key,omitempty"`
	SecretKey string `yaml:"secret_key,omitempty"`
	// Token, when set, is sent as `Authorization: Bearer <token>` in place of
	// Basic auth from PublicKey/SecretKey. Useful for pass-through scenarios
	// where a caller-supplied access token should authenticate outbound calls.
	Token string `yaml:"token,omitempty"`
	Host  string `yaml:"host,omitempty"`
}

func (c *Credentials) Validate() error {
	if c.Token != "" {
		return nil
	}
	if c.PublicKey == "" || c.SecretKey == "" {
		return ErrCredentialsNotSet
	}

	return nil
}
