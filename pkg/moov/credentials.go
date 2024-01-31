package moov

import (
	"os"
)

const ENV_MOOV_HOST = "MOOV_HOST"
const ENV_MOOV_PUBLIC_KEY = "MOOV_PUBLIC_KEY"
const ENV_MOOV_SECRET_KEY = "MOOV_SECRET_KEY" //nolint:gosec

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
	Host      string `yaml:"host,omitempty"`
}

func (c *Credentials) Validate() error {
	if c.PublicKey == "" || c.SecretKey == "" {
		return ErrAuthCredentialsNotSet
	}
	return nil
}
