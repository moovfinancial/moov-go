package moov

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Credentials struct {
	AccountID string `yaml:"account_id,omitempty"`
	PublicKey string `yaml:"public_key,omitempty"`
	SecretKey string `yaml:"secret_key,omitempty"`
	Domain    string `yaml:"domain,omitempty"`
}

func readConfig() (Credentials, error) {
	var cred Credentials

	// Read the config file
	data, err := ioutil.ReadFile("../config.yaml")
	if err != nil {
		return cred, err
	}

	// Unmarshal the YAML data into the cred struct
	if err := yaml.Unmarshal(data, &cred); err != nil {
		return cred, err
	}

	return cred, nil
}
