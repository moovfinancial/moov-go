package moov

import (
	"os"
)

const ENV_MOOV_URL = "MOOV_URL"
const ENV_MOOV_PUBLIC_KEY = "MOOV_PUBLIC_KEY"
const ENV_MOOV_SECRET_KEY = "MOOV_SECRET_KEY" // #nosec G101

func CredentialsDefault() Credentials {
	return Credentials{}
}

func CredentialsFromEnv() Credentials {
	creds := CredentialsDefault()
	creds.PublicKey = os.Getenv(ENV_MOOV_PUBLIC_KEY)
	creds.SecretKey = os.Getenv(ENV_MOOV_SECRET_KEY)

	creds.URL = os.Getenv(ENV_MOOV_URL)
	if creds.URL == "" {
		creds.URL = "https://api.moov.io"
	}

	return creds
}

type Credentials struct {
	PublicKey string `yaml:"public_key,omitempty"`
	SecretKey string `yaml:"secret_key,omitempty"`
	URL       string `yaml:"url,omitempty"`
}

func (c *Credentials) Validate() error {
	if c.PublicKey == "" || c.SecretKey == "" {
		return ErrCredentialsNotSet
	}

	return nil
}
