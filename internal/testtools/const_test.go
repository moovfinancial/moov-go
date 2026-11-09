package testtools

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSeededIDsFromEnv(t *testing.T) {
	if os.Getenv("TESTTOOLS_SEEDED_IDS_HELPER") == "1" {
		fmt.Println(PARTNER_ID)
		fmt.Println(PARTNER_WALLET_PM_ID)
		fmt.Println(MERCHANT_ID)
		fmt.Println(MERCHANT_WALLET_PM_ID)
		fmt.Println(MERCHANT_WALLET_ID)
		os.Exit(0)
	}

	keys := []string{
		"MOOV_PARTNER_ID",
		"MOOV_PARTNER_WALLET_PM_ID",
		"MOOV_MERCHANT_ID",
		"MOOV_MERCHANT_WALLET_PM_ID",
		"MOOV_MERCHANT_WALLET_ID",
	}
	defaults := []string{
		defaultPartnerID,
		defaultPartnerWalletPMID,
		defaultMerchantID,
		defaultMerchantWalletPMID,
		defaultMerchantWalletID,
	}
	overrides := []string{
		"synthetic-partner-id",
		"synthetic-partner-wallet-pm-id",
		"synthetic-merchant-id",
		"synthetic-merchant-wallet-pm-id",
		"synthetic-merchant-wallet-id",
	}

	keySet := make(map[string]bool, len(keys))
	for _, key := range keys {
		keySet[key] = true
	}

	emptyEnv := make([]string, len(keys))
	for i, key := range keys {
		emptyEnv[i] = key + "="
	}
	overrideEnv := make([]string, len(keys))
	for i, key := range keys {
		overrideEnv[i] = key + "=" + overrides[i]
	}

	cases := []struct {
		name string
		env  []string
		want []string
	}{
		{name: "absent overrides keep defaults", env: nil, want: defaults},
		{name: "empty overrides keep defaults", env: emptyEnv, want: defaults},
		{name: "overrides win", env: overrideEnv, want: overrides},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			env := make([]string, 0, len(os.Environ())+len(tc.env)+1)
			for _, entry := range os.Environ() {
				name, _, _ := strings.Cut(entry, "=")
				if name == "TESTTOOLS_SEEDED_IDS_HELPER" || keySet[name] {
					continue
				}
				env = append(env, entry)
			}
			env = append(env, "TESTTOOLS_SEEDED_IDS_HELPER=1")
			env = append(env, tc.env...)

			cmd := exec.Command(os.Args[0], "-test.run=^TestSeededIDsFromEnv$")
			cmd.Env = env
			out, err := cmd.Output()
			require.NoError(t, err)

			got := strings.Split(strings.TrimSpace(string(out)), "\n")
			require.Equal(t, tc.want, got)
		})
	}
}
